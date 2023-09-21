package immudb

import "fmt"

type Vault[D any] struct {
	ledger     *Ledger
	pathSuffix string
	meta       Collection
}

func NewVault[D any](l *Ledger, meta Collection) *Vault[D] {
	return &Vault[D]{ledger: l, pathSuffix: l.collectionPathSuffix(meta.Name), meta: meta}
}

const MAX_PER_PAGE = 100

func (v *Vault[D]) Name() string {
	return v.meta.Name
}

func (v *Vault[D]) IdFieldName() string {
	return v.meta.IdFieldName
}

func (v *Vault[D]) CreateDocument(doc D) (string, error) {
	var result documentResult
	if err := v.ledger.request("PUT", v.pathSuffix+"/document", doc, &result); err != nil {
		return "", err
	}
	return result.DocumentId, nil
}

func (v *Vault[D]) UpdateDocument(query *Query, doc D) (string, error) {
	var result documentResult
	req := documentUpdateRequest{Document: doc, Query: query}
	if err := v.ledger.request("POST", v.pathSuffix+"/document", &req, &result); err != nil {
		return "", err
	}
	return result.DocumentId, nil
}

func (v *Vault[D]) UpdateDocumentBy(keyName, keyValue string, doc D) (string, error) {
	return v.UpdateDocument(byIdQuery(keyName, keyValue), doc)
}

func (v *Vault[D]) UpdateDocumentById(id string, doc D) error {
	_, err := v.UpdateDocumentBy(v.meta.IdFieldName, id, doc)
	return err
}

func (v *Vault[D]) LoadDocumentBy(keyName, keyValue string) (doc D, err error) {
	found := false
	err = v.SearchDocuments(byIdQuery(keyName, keyValue), func(rev DocumentRevision[D]) {
		doc = rev.Document
		found = true
	})
	if !found && err == nil {
		err = fmt.Errorf("LoadDocumentBy: record not found for %s = %s", keyName, keyValue)
	}
	return
}

func (v *Vault[D]) LoadDocumentById(id string) (doc D, err error) {
	return v.LoadDocumentBy(v.meta.IdFieldName, id)
}

func (v *Vault[D]) SearchDocuments(query *Query, callback func(DocumentRevision[D])) error {
	req := documentSearchRequest{pagedRequest: pagedRequest{Page: 1, PerPage: MAX_PER_PAGE}, KeepOpen: true, Query: query}
	var resp documentRevisionsResponse[D]
	for {
		if err := v.ledger.request("POST", v.pathSuffix+"/documents/search", &req, &resp); err != nil {
			return err
		}
		for _, doc := range resp.Revisions {
			callback(doc)
		}
		if len(resp.Revisions) < req.PerPage {
			return nil
		}
		req.SearchId = resp.SearchId
		req.Page++
	}
}

func (v *Vault[D]) documentPathSuffix(docId, suffix string) string {
	return fmt.Sprintf("%s/document/%s/%s", v.pathSuffix, docId, suffix)
}

func (v *Vault[D]) AuditDocument(id string, orderByDesc bool, callback func(DocumentRevision[D])) error {
	req := documentAuditRequest{pagedRequest: pagedRequest{Page: 1, PerPage: MAX_PER_PAGE}, OrderByDesc: orderByDesc}
	var resp documentRevisionsResponse[D]
	for {
		if err := v.ledger.request("POST", v.documentPathSuffix(id, "audit"), &req, &resp); err != nil {
			return err
		}
		for _, doc := range resp.Revisions {
			callback(doc)
		}
		if len(resp.Revisions) < req.PerPage {
			return nil
		}
		req.Page++
	}
}
