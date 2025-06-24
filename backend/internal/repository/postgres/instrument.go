package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres/pq_models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type InstrumentRepo struct {
	db *sqlx.DB
}

func NewInstrumentRepo(db *sqlx.DB) *InstrumentRepo {
	return &InstrumentRepo{
		db: db,
	}
}

type Instrument interface {
	GetById(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error)
	GetUniqueData(ctx context.Context, req *models.GetUniqueDTO) ([]string, error)
	Create(ctx context.Context, dto *models.InstrumentDTO) error
	Update(ctx context.Context, dto *models.InstrumentDTO) error
	ChangeStatus(ctx context.Context, dto *models.UpdateStatus) error
	Delete(ctx context.Context, id string) error
}

func (r *InstrumentRepo) GetById(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
	query := fmt.Sprintf(`SELECT name, date_of_receipt, type, factory_number, measurement_limits, accuracy, state_register,
		country_of_produce, manufacturer, responsible, inventory, year_of_issue, inter_verification_interval, act_of_entering, 
		act_of_entering_id, notes, status FROM %s WHERE id=$1`,
		InstrumentsTable,
	)
	data := &models.Instrument{}

	if err := r.db.GetContext(ctx, data, query, req.Id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *InstrumentRepo) GetUniqueData(ctx context.Context, req *models.GetUniqueDTO) ([]string, error) {
	reg := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := reg.ReplaceAllString(req.Field, "${1}_${2}")
	req.Field = strings.ToLower(snake)

	enabledFields := map[string]struct{}{"manufacturer": {}, "country_of_produce": {}, "responsible": {}, "type": {}}

	if _, exist := enabledFields[req.Field]; !exist {
		return nil, fmt.Errorf("field is not allowed")
	}

	query := fmt.Sprintf(`SELECT DISTINCT(%s) AS item FROM %s`, req.Field, InstrumentsTable)
	tmp := []pq_models.UniqueData{}

	if err := r.db.SelectContext(ctx, &tmp, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	data := []string{}
	for _, v := range tmp {
		data = append(data, v.Item)
	}
	return data, nil
}

func (r *InstrumentRepo) Create(ctx context.Context, dto *models.InstrumentDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, section_id, user_id, position, name, date_of_receipt, type, factory_number, measurement_limits, 
		accuracy, state_register, country_of_produce, manufacturer, responsible, inventory, year_of_issue, inter_verification_interval, 
		act_of_entering, act_of_entering_id, notes, status) 
		VALUES (:id, :section_id, :user_id, :position, :name, :date_of_receipt, :type, :factory_number, :measurement_limits, 
		:accuracy, :state_register, :country_of_produce, :manufacturer, :responsible, :inventory, :year_of_issue, :inter_verification_interval, 
		:act_of_entering, :act_of_entering_id, :notes, :status)`,
		InstrumentsTable,
	)
	dto.Id = uuid.NewString()
	dto.Status = models.InstrumentStatusWork
	if dto.ActOfEnteringId == "" {
		dto.ActOfEnteringId = uuid.Nil.String()
	}

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *InstrumentRepo) Update(ctx context.Context, dto *models.InstrumentDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET name=:name, date_of_receipt=:date_of_receipt, type=:type, factory_number=:factory_number, 
		measurement_limits=:measurement_limits, accuracy=:accuracy, state_register=:state_register, country_of_produce=:country_of_produce,
		manufacturer=:manufacturer, responsible=:responsible, inventory=:inventory, year_of_issue=:year_of_issue, 
		inter_verification_interval=:inter_verification_interval, act_of_entering=:act_of_entering, act_of_entering_id=:act_of_entering_id,
		notes=:notes, updated_at=now() WHERE id=:id`,
		InstrumentsTable,
	)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *InstrumentRepo) ChangeStatus(ctx context.Context, dto *models.UpdateStatus) error {
	query := fmt.Sprintf(`UPDATE %s SET status=:status WHERE id=:id`, InstrumentsTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *InstrumentRepo) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, InstrumentsTable)
	// query := fmt.Sprintf(`UPDATE %s SET status='deleted' WHERE id=$1`, InstrumentTable)

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
