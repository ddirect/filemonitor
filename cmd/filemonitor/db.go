package main

import (
	"filemonitor/pkg/common"
	db "filemonitor/pkg/immudb"
	"net/http"
)

type Vault = *db.Vault[FileStatus]

func SetupDB(ledgerName, collectionName, apiKey string) (Vault, error) {
	ledger := db.NewLedger(ledgerName, apiKey)
	coll := common.FileMonitorCollection(collectionName)
	err := ledger.CreateCollection(coll)
	if err != nil && db.HttpStatusCode(err) != http.StatusConflict {
		return nil, err
	}
	return db.NewVault[FileStatus](ledger, coll), nil
}
