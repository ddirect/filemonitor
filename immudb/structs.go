package immudb

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Index struct {
	FieldNames []string `json:"fields"`
	IsUnique   bool     `json:"isUnique"`
}

type Collection struct {
	Name        string  `json:"name"`
	IdFieldName string  `json:"idFieldName"`
	Fields      []Field `json:"fields"`
	Indexes     []Index `json:"indexes"`
}

func NewCollectionWithIndex(name string, idFieldName string, field Field, unique bool) Collection {
	return Collection{
		Name:        name,
		IdFieldName: idFieldName,
		Fields:      []Field{field},
		Indexes:     []Index{{FieldNames: []string{field.Name}, IsUnique: unique}},
	}
}

type collectionsResponse struct {
	Items []Collection `json:"collections"`
}

type FieldComparison struct {
	FieldName string `json:"field"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
}

type Expression struct {
	FieldComparisons []FieldComparison `json:"fieldComparisons"`
}

type Query struct {
	Expressions []Expression `json:"expressions"`
	Limit       int          `json:"limit,omitempty"`
}

func byIdQuery(keyName, keyValue string) *Query {
	return &Query{
		Expressions: []Expression{
			{FieldComparisons: []FieldComparison{
				{FieldName: keyName, Operator: "EQ", Value: keyValue}}}},
		Limit: 1,
	}
}

type DocumentRevision[D any] struct {
	Document            D      `json:"document"`
	RevisionString      string `json:"revision"`
	TransactionIdString string `json:"transactionId"`
}

type documentUpdateRequest struct {
	Document any    `json:"document"`
	Query    *Query `json:"query"`
}

type documentResult struct {
	DocumentId          string `json:"documentId"`
	RevisionString      string `json:"revision"`
	TransactionIdString string `json:"transactionId"`
}

type pagedRequest struct {
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
}

type documentSearchRequest struct {
	pagedRequest
	KeepOpen bool   `json:"keepOpen"`
	SearchId string `json:"searchId,omitempty"`
	Query    *Query `json:"query,omitempty"`
}

type documentAuditRequest struct {
	pagedRequest
	OrderByDesc bool `json:"desc"`
}

type documentRevisionsResponse[D any] struct {
	Revisions []DocumentRevision[D] `json:"revisions"`
	SearchId  string                `json:"searchId"`
}
