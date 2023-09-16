package immudb

import (
	"filemonitor/common/log"
	"fmt"
	"os"
)

type Ledger struct {
	baseUrl  string
	name     string
	apiKey   string
	dumpHttp bool
}

const API_URL = "https://vault.immudb.io/ics/api/v1"
const API_KEY_ENV = "IMMUDB_API_KEY"
const HTTP_DUMP_ENV = "IMMUDB_HTTP_DUMP"

func httpDumpFromEnv() bool {
	val := os.Getenv(HTTP_DUMP_ENV)
	if val != "" {
		return val == "1"
	}
	log.Debug("set the %s environment variable to 1 to enable immudb HTTP request and response dumping\n", HTTP_DUMP_ENV)
	return false
}

func NewLedger(name, apiKey string) *Ledger {
	return NewLedgerWithUrl(API_URL, name, apiKey)
}

func NewLedgerWithUrl(baseUrl, name, apiKey string) *Ledger {
	if apiKey == "" {
		apiKey = os.Getenv(API_KEY_ENV)
		if apiKey == "" {
			log.Warning("API key not set - it can be set with the %s environment variable", API_KEY_ENV)
		}
	}
	return &Ledger{
		baseUrl:  fmt.Sprintf("%s/ledger/%s/", baseUrl, name),
		name:     name,
		apiKey:   apiKey,
		dumpHttp: httpDumpFromEnv(),
	}
}

func (l *Ledger) collectionPathSuffix(name string) string {
	return fmt.Sprintf("collection/%s", name)
}

func (l *Ledger) CreateCollection(collection Collection) error {
	return l.request("PUT", l.collectionPathSuffix(collection.Name), &collection, nil)
}

func (l *Ledger) GetCollections() ([]Collection, error) {
	var resp collectionsResponse
	if err := l.request("GET", "collections", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (l *Ledger) GetCollection(name string) (collection Collection, err error) {
	err = l.request("GET", l.collectionPathSuffix(name), nil, &collection)
	return
}

func (l *Ledger) DeleteCollection(name string) error {
	return l.request("DELETE", l.collectionPathSuffix(name), nil, nil)
}
