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
	CreateSeveral(ctx context.Context, dto []*models.InstrumentDTO) error
	Update(ctx context.Context, dto *models.InstrumentDTO) error
	ChangePosition(ctx context.Context, dto *models.ChangePositionDTO) error
	ChangeStatus(ctx context.Context, dto *models.UpdateStatus) error
	ChangeSeveralStatuses(ctx context.Context, dto []*models.UpdateStatus) error
	Delete(ctx context.Context, id string) error
}

func (r *InstrumentRepo) GetById(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
	query := fmt.Sprintf(`SELECT id, position, name, date_of_receipt, type, factory_number, measurement_limits, accuracy, state_register,
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

	enabledFields := map[string]struct{}{
		"name": {}, "type": {}, "measurement_limits": {}, "accuracy": {}, "factory_number": {}, "manufacturer": {}, "country_of_produce": {},
		"responsible": {},
	}

	if _, exist := enabledFields[req.Field]; !exist {
		return nil, fmt.Errorf("field is not allowed")
	}

	query := fmt.Sprintf(`SELECT DISTINCT(%s) AS item FROM %s WHERE %s!='' AND %s IS NOT NULL AND section_id=$1 ORDER BY item`,
		req.Field, InstrumentsTable, req.Field, req.Field,
	)
	tmp := []pq_models.UniqueData{}

	if err := r.db.SelectContext(ctx, &tmp, query, req.SectionId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	data := []string{}
	for _, v := range tmp {
		data = append(data, v.Item)
	}
	return data, nil
}

func (r *InstrumentRepo) Create(ctx context.Context, dto *models.InstrumentDTO) error {
	maxQuery := fmt.Sprintf(`SELECT COALESCE(MAX(position), 0) FROM %s WHERE section_id=$1`, InstrumentsTable)
	var maxPosition int
	if err := r.db.GetContext(ctx, &maxPosition, maxQuery, dto.SectionId); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	query := fmt.Sprintf(`INSERT INTO %s (id, section_id, user_id, position, name, date_of_receipt, type, factory_number, measurement_limits, 
		accuracy, state_register, country_of_produce, manufacturer, responsible, inventory, year_of_issue, inter_verification_interval, 
		act_of_entering, act_of_entering_id, notes, status) 
		VALUES (:id, :section_id, :user_id, :position, :name, :date_of_receipt, :type, :factory_number, :measurement_limits, 
		:accuracy, :state_register, :country_of_produce, :manufacturer, :responsible, :inventory, :year_of_issue, :inter_verification_interval, 
		:act_of_entering, :act_of_entering_id, :notes, :status)`,
		InstrumentsTable,
	)
	dto.Id = uuid.NewString()
	dto.Position = maxPosition + 1
	dto.Status = models.InstrumentStatusWork
	if dto.ActOfEnteringId == "" {
		dto.ActOfEnteringId = uuid.Nil.String()
	}

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *InstrumentRepo) CreateSeveral(ctx context.Context, dto []*models.InstrumentDTO) error {
	maxQuery := fmt.Sprintf(`SELECT COALESCE(MAX(position), 0) FROM %s WHERE section_id=$1`, InstrumentsTable)
	var maxPosition int
	if err := r.db.GetContext(ctx, &maxPosition, maxQuery, dto[0].SectionId); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	query := fmt.Sprintf(`INSERT INTO %s (id, section_id, user_id, position, name, date_of_receipt, type, factory_number, measurement_limits, 
		accuracy, state_register, country_of_produce, manufacturer, responsible, inventory, year_of_issue, inter_verification_interval, 
		act_of_entering, act_of_entering_id, notes, status) 
		VALUES (:id, :section_id, :user_id, :position, :name, :date_of_receipt, :type, :factory_number, :measurement_limits, 
		:accuracy, :state_register, :country_of_produce, :manufacturer, :responsible, :inventory, :year_of_issue, :inter_verification_interval, 
		:act_of_entering, :act_of_entering_id, :notes, :status)`,
		InstrumentsTable,
	)

	for i := range dto {
		dto[i].Id = uuid.NewString()
		dto[i].ActOfEnteringId = uuid.Nil.String()
		dto[i].Position = maxPosition + i + 1
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

func (r *InstrumentRepo) ChangePosition(ctx context.Context, dto *models.ChangePositionDTO) error {
	//TODO фиг знает как это будет отрабатывать с большими объемами данных
	query := fmt.Sprintf(`UPDATE %s AS i SET position=n.position FROM (
			SELECT id, 
				CASE WHEN $1::integer < $2::integer THEN
					(CASE 
						WHEN "position" < $1 OR "position" > $2 THEN "position"
						WHEN "position" >= $1 AND "position" < $2 THEN "position"+1
						WHEN "position" = $2 THEN $1
					END)
				ELSE 
					(CASE
						WHEN "position" > $1 OR "position" < $2 THEN "position"
						WHEN "position" <= $1 AND "position" > $2 THEN "position"-1
						WHEN "position" = $2 THEN $1
					END)
				END AS "position"
			FROM %s WHERE section_id=$3 ORDER BY position
		) AS n WHERE i.id=n.id`,
		InstrumentsTable, InstrumentsTable,
	)

	if _, err := r.db.ExecContext(ctx, query, dto.NewPosition, dto.OldPosition, dto.SectionId); err != nil {
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
func (r *InstrumentRepo) ChangeSeveralStatuses(ctx context.Context, dto []*models.UpdateStatus) error {
	values := []string{}
	args := []interface{}{}
	for i, v := range dto {
		tmp := []interface{}{v.Id, v.Status}
		args = append(args, tmp...)
		numbers := []string{}
		for j := range tmp {
			numbers = append(numbers, fmt.Sprintf("$%d", i*len(tmp)+j+1))
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(numbers, ",")))
	}

	query := fmt.Sprintf(`UPDATE %s AS t SET status=s.status
		FROM (VALUES %s) AS s(id, status)
		WHERE t.id=s.id::uuid`,
		InstrumentsTable, strings.Join(values, ","),
	)

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *InstrumentRepo) Delete(ctx context.Context, id string) error {
	// query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, InstrumentsTable)
	// query := fmt.Sprintf(`UPDATE %s SET status='deleted', deleted=now() WHERE id=$1 AND (
	// 		SELECT status FROM %s WHERE instrument_id=$1 ORDER BY date_of_issue DESC, created_at DESC LIMIT 1
	// 	) = 'reserve'`,
	// 	InstrumentsTable, LocationTable,
	// )
	query := fmt.Sprintf(`UPDATE %s AS i SET status='deleted', deleted=now()
		FROM (SELECT i.id, l.status FROM %s AS i
			LEFT JOIN LATERAL (SELECT status FROM %s WHERE instrument_id=i.id ORDER BY date_of_issue DESC, created_at DESC LIMIT 1) AS l ON TRUE
			WHERE i.id=$1
		) AS s
		WHERE i.id=$1 AND (s.status IS NULL OR s.status='reserve')`,
		InstrumentsTable, InstrumentsTable, LocationTable,
	)

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected. error: %w", err)
	}
	if rows == 0 {
		return models.ErrNoRows
	}

	return nil
}
