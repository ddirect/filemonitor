package main

import (
	"filemonitor/common"
	"filemonitor/immudb"
	"fmt"
)

func deleteAllCollections(l *immudb.Ledger) {
	collections, err := l.GetCollections()
	common.ExitIf(err)
	for _, c := range collections {
		fmt.Printf("deleting collection %s...\n", c.Name)
		common.ExitIf(l.DeleteCollection(c.Name))
	}
}

func showAllCollections(l *immudb.Ledger) {
	collections, err := l.GetCollections()
	common.ExitIf(err)
	for _, c := range collections {
		fmt.Println(common.IndentedJson(c))
	}
}
