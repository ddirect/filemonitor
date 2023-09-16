package common

import (
	"encoding/json"
	db "filemonitor/immudb"
	"fmt"
	"os"
)

func FileMonitorCollection(collectionName string) db.Collection {
	return db.NewCollectionWithIndex(collectionName, FILE_INFO_ID_FIELD_NAME, db.Field{Name: FILE_INFO_FILE_FIELD_NAME, Type: "STRING"}, true)
}

func ExitIf(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func IndentedJson(x any) string {
	data, err := json.MarshalIndent(x, "", "  ")
	ExitIf(err)
	return string(data)
}
