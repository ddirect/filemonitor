package main

import (
	"filemonitor/pkg/common"
	db "filemonitor/pkg/immudb"
	"flag"
)

type Vault = *db.Vault[common.FileInfo]
type FileInfo = common.FileInfo
type Revision = db.DocumentRevision[FileInfo]

func main() {
	var apiKey, ledgerName, collectionName, id string
	var native, audit bool

	flag.BoolVar(&native, "native", false, "use native instead of index id")
	flag.BoolVar(&audit, "audit", false, "show audit")
	flag.StringVar(&apiKey, "api_key", "", "API key")
	flag.StringVar(&ledgerName, "ledger", "default", "immudb ledger")
	flag.StringVar(&collectionName, "collection", common.FILE_MONITOR_DEFAULT_COLLECTION, "immudb collection")
	flag.StringVar(&id, "id", "", "show specific document by id")

	flag.Parse()

	ledger := db.NewLedger(ledgerName, apiKey)
	vault := db.NewVault[common.FileInfo](ledger, common.FileMonitorCollection(collectionName))

	if id == "" {
		showAll(vault)
	} else {
		if audit {
			if native {
				auditDocumentByNativeId(vault, id)
			} else {
				auditDocumentById(vault, id)
			}
		} else {
			if native {
				showDocumentByNativeId(vault, id)
			} else {
				showDocumentById(vault, id)
			}
		}
	}
}
