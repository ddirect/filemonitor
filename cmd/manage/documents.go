package main

import (
	"filemonitor/pkg/common"
	"filemonitor/pkg/immudb"
	"fmt"
)

func showAllDocs(ledger *immudb.Ledger, collectionName string) {
	vault := immudb.NewVault[any](ledger, immudb.Collection{Name: collectionName})
	common.ExitIf(vault.SearchDocuments(nil, func(rev immudb.DocumentRevision[any]) {
		fmt.Println(common.IndentedJson(rev.Document))
	}))
}
