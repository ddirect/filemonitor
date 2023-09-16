package immudb

import (
	cm "filemonitor/common_test"
	"fmt"
	"os"
	"reflect"
	"testing"
)

const LEDGER = "default"

func apiKey(t *testing.T) string {
	key := os.Getenv(API_KEY_ENV)
	if key == "" {
		t.Fatalf("%s environment variable not set", API_KEY_ENV)
	}
	return key
}

type docType struct {
	String   string
	Int      int
	Bool     bool
	CustomId string
}

func randomDoc(customIdValue string) docType {
	return docType{
		String:   cm.RandomString(),
		Int:      cm.Rnd.Intn(1000),
		Bool:     cm.Rnd.Intn(2) == 0,
		CustomId: customIdValue,
	}
}

type slDocType map[string]any

func randomSlDoc(customIdValue string) slDocType {
	return slDocType{
		// avoiding ints since they map to float64 when parsing json
		cm.RandomString(): cm.RandomString(),
		cm.RandomString(): cm.Rnd.Intn(2) == 0,
		"CustomId":        customIdValue,
	}
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func compareCollections(t *testing.T, c1, c2 Collection) {
	if c1.Name != c2.Name || c1.IdFieldName != c2.IdFieldName {
		t.Fatal("collections don't match")
	}
}

func compareDocs(t *testing.T, d1, d2 docType) {
	if !reflect.DeepEqual(d1, d2) {
		t.Fatal("documents don't match")
	}
}

func compareSlDocs(t *testing.T, doc, original slDocType) {
	for key, value := range original {
		if !reflect.DeepEqual(value, doc[key]) {
			t.Fatal("documents don't match")
		}
	}
}

func TestAll(t *testing.T) {
	ledger := NewLedger(LEDGER, apiKey(t))

	collection := NewCollectionWithIndex(cm.RandomString(), "id", Field{Name: "CustomId", Type: "STRING"}, true)
	fmt.Printf("collection name: %s\n", collection.Name)

	checkError(t, ledger.CreateCollection(collection))

	t.Run("collection can be found among all", func(t *testing.T) {
		collections, err := ledger.GetCollections()
		checkError(t, err)
		for _, c := range collections {
			if c.Name == collection.Name {
				compareCollections(t, c, collection)
				return
			}
		}
		t.Fatal("collection not found")
	})

	t.Run("collection can be found", func(t *testing.T) {
		c, err := ledger.GetCollection(collection.Name)
		checkError(t, err)
		compareCollections(t, c, collection)
	})

	vault := NewVault[docType](ledger, collection)
	slVault := NewVault[slDocType](ledger, collection)
	var nativeDocId, nativeSlDocId string
	customDocId := cm.RandomString()
	customSlDocId := cm.RandomString()
	originalDoc := randomDoc(customDocId)
	originalSlDoc := randomSlDoc(customSlDocId)
	t.Run("document can be created", func(t *testing.T) {
		var err error
		nativeDocId, err = vault.CreateDocument(originalDoc)
		checkError(t, err)
	})

	t.Run("schemaless document can be created", func(t *testing.T) {
		var err error
		nativeSlDocId, err = slVault.CreateDocument(originalSlDoc)
		checkError(t, err)
		originalSlDoc[vault.IdFieldName()] = nativeSlDocId
	})

	if nativeDocId == "" || nativeSlDocId == "" {
		t.Fatal("document(s) not created, cannot continue")
	}
	fmt.Printf("doc with schema - native id %s - custom id %s\n", nativeDocId, customDocId)
	fmt.Printf(" schemaless doc - native id %s - custom id %s\n", nativeSlDocId, customSlDocId)

	t.Run("document can be loaded by native id", func(t *testing.T) {
		doc, err := vault.LoadDocumentById(nativeDocId)
		checkError(t, err)
		compareDocs(t, doc, originalDoc)
	})

	t.Run("document can be loaded by custom id", func(t *testing.T) {
		doc, err := vault.LoadDocumentBy("CustomId", customDocId)
		checkError(t, err)
		compareDocs(t, doc, originalDoc)
	})

	t.Run("schemaless document can be loaded by native id", func(t *testing.T) {
		doc, err := slVault.LoadDocumentById(nativeSlDocId)
		checkError(t, err)
		compareSlDocs(t, doc, originalSlDoc)
	})

	t.Run("schemaless document can be loaded by custom id", func(t *testing.T) {
		doc, err := slVault.LoadDocumentBy("CustomId", customSlDocId)
		checkError(t, err)
		compareSlDocs(t, doc, originalSlDoc)
	})

	t.Run("loading a document by native id fails if not found", func(t *testing.T) {
		if _, err := vault.LoadDocumentById("00000000000000000000000000000000"); err == nil {
			t.FailNow()
		}
	})

	t.Run("loading a document by custom id fails if not found", func(t *testing.T) {
		if _, err := vault.LoadDocumentBy("CustomId", "not_an_id"); err == nil {
			t.FailNow()
		}
	})

	expectedAuditCount := 1
	expectedAuditInt := originalDoc.Int
	t.Run("document can be updated by native id", func(t *testing.T) {
		newDoc := originalDoc
		newDoc.Int += 1
		newDoc.String += "_"
		expectedAuditCount++
		expectedAuditInt += newDoc.Int
		checkError(t, vault.UpdateDocumentById(nativeDocId, newDoc))
	})

	t.Run("document can be updated by custom id", func(t *testing.T) {
		newDoc := originalDoc
		newDoc.Int += 2
		newDoc.String += "__"
		expectedAuditCount++
		expectedAuditInt += newDoc.Int
		_, err := vault.UpdateDocumentBy("CustomId", customDocId, newDoc)
		checkError(t, err)
	})

	t.Run("document can be audited", func(t *testing.T) {
		actualCount := 0
		actualInt := 0
		firstString := ""
		vault.AuditDocument(nativeDocId, false, func(rev DocumentRevision[docType]) {
			if firstString == "" {
				firstString = rev.Document.String
			}
			actualCount++
			actualInt += rev.Document.Int
		})
		if actualCount != expectedAuditCount || actualInt != expectedAuditInt || firstString != originalDoc.String {
			t.FailNow()
		}
	})

	t.Run("collection can be deleted", func(t *testing.T) {
		checkError(t, ledger.DeleteCollection(collection.Name))
	})
}
