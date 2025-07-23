package pq_models

type VerificationDoc struct {
	Id             string `db:"id"`
	VerificationId string `db:"verification_id"`
	Name           string `db:"name"`
	DocId          string `db:"doc_id"`
	Type           string `db:"type"`
	Path           string `db:"path"`
}
