package main

import (
	"encoding/hex"
	"filemonitor/common"
	"fmt"
	"time"
)

func showDocument(doc FileInfo) {
	fmt.Printf("%34s%66s%25s  %s\n",
		doc.Id,
		hex.EncodeToString(doc.Hash),
		time.Unix(0, doc.ModTime).Format("2006-01-02 15:04:05.000"),
		doc.FileName)
}

func showRevision(rev Revision) {
	fmt.Printf("rev %s trans %s\n", rev.RevisionString, rev.TransactionIdString)
	showDocument(rev.Document)
}

func showAll(v Vault) {
	common.ExitIf(v.SearchDocuments(nil, func(rev Revision) {
		showDocument(rev.Document)
	}))
}

func showDocumentById(v Vault, id string) {
	showDocument(loadDocumentById(v, id))
}

func showDocumentByNativeId(v Vault, nativeId string) {
	doc, err := v.LoadDocumentById(nativeId)
	common.ExitIf(err)
	showDocument(doc)
}

func auditDocumentById(v Vault, id string) {
	auditDocumentByNativeId(v, loadDocumentById(v, id).Id)
}

func auditDocumentByNativeId(v Vault, nativeId string) {
	common.ExitIf(v.AuditDocument(nativeId, false, func(rev Revision) {
		showRevision(rev)
	}))
}

func loadDocumentById(v Vault, id string) FileInfo {
	doc, err := v.LoadDocumentBy(common.FILE_INFO_FILE_FIELD_NAME, id)
	common.ExitIf(err)
	return doc
}
