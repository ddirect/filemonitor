package main

import (
	"filemonitor/immudb"
	"flag"
	"fmt"
)

func main() {
	var op, apiKey, ledgerName, collectionName string

	flag.StringVar(&op, "op", "", "operation: delete_all_collections|show_all_collections|show_all_docs")
	flag.StringVar(&apiKey, "api_key", "", "API key")
	flag.StringVar(&ledgerName, "ledger", "default", "immudb ledger")
	flag.StringVar(&collectionName, "collection", "", "immudb collection")

	flag.Parse()

	if op == "" {
		flag.Usage()
		return
	}

	ledger := immudb.NewLedger(ledgerName, apiKey)
	switch op {
	case "show_all_collections":
		showAllCollections(ledger)
	case "delete_all_collections":
		deleteAllCollections(ledger)
	case "show_all_docs":
		if collectionName == "" {
			fmt.Println("please specify the collection name")
			return
		}
		showAllDocs(ledger, collectionName)
	default:
		fmt.Printf("unknown operation '%s'\n", op)
	}
}
