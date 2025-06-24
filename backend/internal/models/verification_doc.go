package models

type VerificationDoc struct {
	Id    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	DocId string `json:"docId" db:"doc_id"`
}

type GroupedVerificationDocs struct {
	Groups map[string]*Groups `json:"groups"`
}
type Groups struct {
	Docs []*VerificationDoc `json:"docs"`
}

type GetVerificationDocsDTO struct {
	VerificationId string `json:"verificationId" db:"verification_id"`
}

type GetGroupedVerificationDocsDTO struct {
	InstrumentId string `json:"instrumentId"`
}

type VerificationDocDTO struct {
	Id             string `json:"id" db:"id"`
	VerificationId string `json:"verificationId" db:"verification_id"`
	Name           string `json:"doc" db:"name"`
	DocId          string `json:"docId" db:"doc_id"`
}

type DeleteVerificationDocDTO struct {
	Id string `json:"id" db:"id"`
}
