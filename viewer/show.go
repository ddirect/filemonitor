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
	common.ExitIf(v.SearchDocuments(nil, func(rev Revision) { showDocument(rev.Document) }))
}

func loadDocument(v Vault, id string) FileInfo {
	doc, err := v.LoadDocumentBy(common.FILE_INFO_FILE_FIELD_NAME, id)
	common.ExitIf(err)
	return doc
}

func showById(v Vault, id string) {
	showDocument(loadDocument(v, id))
}

func showByNativeId(v Vault, nativeId string) {
	doc, err := v.LoadDocumentById(nativeId)
	common.ExitIf(err)
	showDocument(doc)
}

func auditById(v Vault, id string) {
	auditByNativeId(v, loadDocument(v, id).Id)
}

func auditByNativeId(v Vault, nativeId string) {
	common.ExitIf(v.AuditDocument(nativeId, false, func(rev Revision) {
		showRevision(rev)
	}))
}
