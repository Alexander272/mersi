package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type SIRepo struct {
	db *sqlx.DB
}

func NewSIRepo(db *sqlx.DB) *SIRepo {
	return &SIRepo{
		db: db,
	}
}

type SI interface {
	Get(ctx context.Context, req *models.GetSiDTO) ([]*models.SI, error)
}

func (r *SIRepo) formatField(field string) string {
	format := make(map[string]string)

	format["position"] = "position"
	format["name"] = "name"
	format["dateOfReceipt"] = "date_of_receipt"
	format["type"] = "type"
	format["factoryNumber"] = "factory_number"
	format["measurementLimits"] = "measurement_limits"
	format["accuracy"] = "accuracy"
	format["stateRegister"] = "state_register"
	format["countryOfProduce"] = "country_of_produce"
	format["manufacturer"] = "manufacturer"
	format["responsible"] = "responsible"
	format["inventory"] = "inventory"
	format["yearOfIssue"] = "year_of_issue"
	format["interVerificationInterval"] = "inter_verification_interval"
	format["actOfEntering"] = "act_of_entering"
	format["actOfEnteringId"] = "act_of_entering_id"
	format["notes"] = "notes"
	format["verificationDate"] = "date"
	format["nextVerificationDate"] = "next_date"
	format["department"] = "department_id"
	format["place"] = "place"
	format["person"] = "person"
	format["status"] = "m.status"

	return format[field]
}

func (r *SIRepo) Get(ctx context.Context, req *models.GetSiDTO) ([]*models.SI, error) {
	data := []*models.SI{}
	params := []interface{}{req.SectionId, req.Status}
	count := len(params) + 1

	order := " ORDER BY "
	for _, s := range req.Sort {
		order += fmt.Sprintf("%s %s, ", r.formatField(s.Field), s.Type)
	}
	if len(req.Sort) == 0 {
		order += "position, "
	}
	order += "created_at, id"

	filter := ""

	params = append(params, req.Page.Limit, req.Page.Offset)

	query := fmt.Sprintf(`SELECT i.id, position, name, date_of_receipt, type, factory_number, measurement_limits, accuracy, state_register,
		country_of_produce, manufacturer, responsible, inventory, year_of_issue, inter_verification_interval, act_of_entering, act_of_entering_id, notes,
		v.date, v.next_date, COALESCE(cert, '') AS certificate, COALESCE(cert_id, '') AS certificate_id, COUNT(*) OVER() AS total
		FROM %s AS i
		LEFT JOIN LATERAL (SELECT id, date, next_date FROM %s WHERE instrument_id=i.id ORDER BY date DESC, created_at DESC LIMIT 1) AS v ON TRUE
		LEFT JOIN LATERAL (SELECT name AS cert, doc_id::text AS cert_id FROM %s WHERE verification_id=v.id) AS d ON TRUE
		WHERE section_id=$1 AND i.status=$2 %s%s LIMIT $%d OFFSET $%d`,
		InstrumentsTable, VerificationTable, VerificationDocsTable,
		filter, order, count, count+1,
	)
	// logger.Debug("get si", logger.StringAttr("query", query))

	if err := r.db.SelectContext(ctx, &data, query, params...); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}
