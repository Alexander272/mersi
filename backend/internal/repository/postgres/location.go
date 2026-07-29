package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type LocationRepo struct {
	db *sqlx.DB
	Transaction
}

func NewLocationRepo(db *sqlx.DB, transaction Transaction) *LocationRepo {
	return &LocationRepo{
		db:          db,
		Transaction: transaction,
	}
}

type Location interface {
	Get(ctx context.Context, dto *models.GetLocationDTO) ([]*models.Location, error)
	GetById(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error)
	GetLast(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error)
	GetSeveralLast(ctx context.Context, req *models.GetSeveralLocationsDTO) ([]*models.Location, error)
	GetUsedByHolder(ctx context.Context, req *models.GetLocationByHolderDTO) ([]*models.Location, error)
	GetUsedByDepartment(ctx context.Context, req *models.GetLocationByDepartmentDTO) ([]*models.Location, error)
	SelectByDepartment(ctx context.Context, dto *models.SelectByDepsDTO) ([]string, error)
	CreateInTx(ctx context.Context, tx Tx, dto *models.LocationDTO) error
	Create(ctx context.Context, dto *models.LocationDTO) error
	CreateSeveral(ctx context.Context, dto []*models.LocationDTO) error
	Update(ctx context.Context, dto *models.LocationDTO) error
	SetPerson(ctx context.Context, id string) error
	SetDepartment(ctx context.Context, id string) error
	Receiving(ctx context.Context, dto *models.ReceivingDTO) error
	ForcedReceipt(ctx context.Context, dto *models.ForcedReceiptDTO) error
	ForcedReceiptAll(ctx context.Context) error
	Delete(ctx context.Context, tx Tx, dto *models.DeleteLocationDTO) error
}

func (r *LocationRepo) Get(ctx context.Context, dto *models.GetLocationDTO) ([]*models.Location, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, status, date_of_receiving, date_of_issue, need_confirmed, has_confirmed,
		COALESCE(e.name, NULLIF(person, ''), '') AS person, COALESCE(d.name, NULLIF(place, ''), '') AS place,
		COALESCE(person_id::text, '') AS person_id, COALESCE(department_id::text, '') AS department_id
		FROM %s AS l
		LEFT JOIN LATERAL (SELECT name FROM %s WHERE l.person_id::uuid=id) AS e ON true
		LEFT JOIN LATERAL (SELECT name FROM %s WHERE l.department_id::uuid=id) AS d ON true
		WHERE instrument_id=$1 ORDER BY date_of_issue DESC, created_at DESC, id`,
		LocationTable, EmployeeTable, DepartmentTable,
	)
	data := []*models.Location{}

	if err := r.db.SelectContext(ctx, &data, query, dto.InstrumentId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return data, nil
}

func (r *LocationRepo) GetById(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date_of_issue, date_of_receiving, status, need_confirmed,
		COALESCE(person_id::text, '') AS person_id, COALESCE(department_id::text, '') AS department_id,
		COALESCE(last_place_id::text, '') AS last_place_id
		FROM %s WHERE id=$1`,
		LocationTable,
	)
	data := &models.Location{}

	if err := r.db.GetContext(ctx, data, query, dto.Id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *LocationRepo) GetLast(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date_of_issue, date_of_receiving, status, need_confirmed,
		COALESCE(person_id::text, '') AS person_id, COALESCE(department_id::text, '') AS department_id,
		COALESCE(last_place_id::text, '') AS last_place_id
		FROM %s WHERE instrument_id=$1 ORDER BY date_of_issue DESC, created_at DESC LIMIT 1`,
		LocationTable,
	)
	data := &models.Location{}

	if err := r.db.GetContext(ctx, data, query, dto.InstrumentId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return data, nil
}

func (r *LocationRepo) GetSeveralLast(ctx context.Context, req *models.GetSeveralLocationsDTO) ([]*models.Location, error) {
	query := fmt.Sprintf(`SELECT * FROM (
			SELECT DISTINCT ON (instrument_id) * FROM %s ORDER BY instrument_id, date_of_issue DESC
		)locs WHERE instrument_id::text = ANY($1) AND status!=%s ORDER BY status, date_of_issue DESC`,
		//) locs WHERE instrument_id::text = ANY($1) ORDER BY status, date_of_issue DESC`,
		LocationTable, constants.LocationStatusMoved,
	)
	data := []*models.Location{}

	if err := r.db.SelectContext(ctx, &data, query, pq.Array(req.InstrumentIds)); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *LocationRepo) GetUsedByHolder(ctx context.Context, req *models.GetLocationByHolderDTO) ([]*models.Location, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, status, date_of_receiving, date_of_issue, need_confirmed, 
		has_confirmed, person_id, department_id, last_place_id FROM (
			SELECT DISTINCT ON (instrument_id) * FROM %s ORDER BY instrument_id, date_of_issue DESC
		) locs WHERE person_id=$1`,
		LocationTable,
	)
	data := []*models.Location{}

	if err := r.db.SelectContext(ctx, &data, query, req.PersonId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}
func (r *LocationRepo) GetUsedByDepartment(ctx context.Context, req *models.GetLocationByDepartmentDTO) ([]*models.Location, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, status, date_of_receiving, date_of_issue, need_confirmed, 
		has_confirmed, person_id, department_id, last_place_id FROM (
			SELECT DISTINCT ON (instrument_id) * FROM %s ORDER BY instrument_id, date_of_issue DESC
		) locs WHERE department_id=$1`,
		LocationTable,
	)
	data := []*models.Location{}

	if err := r.db.SelectContext(ctx, &data, query, req.DepartmentId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *LocationRepo) SelectByDepartment(ctx context.Context, dto *models.SelectByDepsDTO) ([]string, error) {
	// этот запрос некорректно работает (выдает несколько строчек при одном инструменте если в перемещениях есть несколько записей для него)
	// query := fmt.Sprintf(`SELECT s.instrument_id FROM %s AS m
	// 	LEFT JOIN LATERAL (SELECT instrument_id, department_id  FROM %s WHERE instrument_id=m.instrument_id
	// 		ORDER BY date_of_issue DESC, created_at DESC LIMIT 1) AS s ON TRUE
	// 	WHERE s.instrument_id=ANY($1) AND s.department_id=ANY($2) AND status=$3`,
	// 	LocationTable, LocationTable,
	// )
	// этот запрос нормально работает только если передать один инструмент
	// query := fmt.Sprintf(`SELECT instrument_id FROM (
	// 		SELECT instrument_id, department_id, status FROM %s
	// 		WHERE instrument_id::text=ANY($1) AND department_id::text=ANY($2)
	// 		ORDER BY date_of_issue DESC, created_at DESC LIMIT 1
	// 	) AS s WHERE status=$3`,
	// 	LocationTable,
	// )
	query := fmt.Sprintf(`SELECT instrument_id FROM (
			SELECT DISTINCT ON (instrument_id) * FROM %s 
				WHERE instrument_id::text=ANY($1) AND department_id::text=ANY($2) 
				ORDER BY instrument_id, date_of_issue DESC
		) AS s WHERE status=$3`,
		LocationTable,
	)
	if dto.Status == "" {
		dto.Status = constants.LocationStatusUsed
	}
	data := []*models.Location{}

	if err := r.db.SelectContext(ctx, &data, query, pq.Array(dto.InstrumentIds), pq.Array(dto.DepartmentIds), dto.Status); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	instruments := []string{}
	for _, l := range data {
		instruments = append(instruments, l.InstrumentId)
	}
	return instruments, nil
}

func (r *LocationRepo) CreateInTx(ctx context.Context, tx Tx, dto *models.LocationDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, instrument_id, date_of_issue, date_of_receiving, status, need_confirmed, 
		person_id, department_id, last_place_id, user_id)
		SELECT id::uuid, instrument_id::uuid, date_of_issue::DATE, date_of_receiving::DATE, 
			status, need_confirmed::boolean, person_id::uuid, s.department_id::uuid, COALESCE(m.department_id::text, ''), s.user_id::uuid
		FROM (VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9))
		AS s(id, instrument_id, date_of_issue, date_of_receiving, status, need_confirmed, person_id, department_id, user_id) 
		LEFT JOIN LATERAL (SELECT department_id FROM %s WHERE instrument_id=s.instrument_id::uuid ORDER BY created_at DESC LIMIT 1) AS m ON true`,
		LocationTable, LocationTable,
	)
	dto.Id = uuid.NewString()

	status := dto.Status
	if status == "" {
		status = constants.LocationStatusUsed
	}

	var personId *string = &dto.PersonId
	if dto.PersonId == "" {
		personId = nil
	}
	var departmentId *string = &dto.DepartmentId
	if dto.DepartmentId == "" {
		departmentId = nil
	}

	args := []interface{}{dto.Id, dto.InstrumentId, dto.DateOfIssue, dto.DateOfReceiving, status, dto.NeedConfirm, personId, departmentId, dto.UserId}
	_, err := r.getExec(tx).ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *LocationRepo) Create(ctx context.Context, dto *models.LocationDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, instrument_id, date_of_issue, date_of_receiving, status, need_confirmed, 
		person_id, department_id, last_place_id, user_id)
		SELECT id::uuid, instrument_id::uuid, date_of_issue::DATE, date_of_receiving::DATE, 
			status, need_confirmed::boolean, person_id::uuid, s.department_id::uuid, COALESCE(m.department_id::text, ''), s.user_id::uuid
		FROM (VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9))
		AS s(id, instrument_id, date_of_issue, date_of_receiving, status, need_confirmed, person_id, department_id, user_id) 
		LEFT JOIN LATERAL (SELECT department_id FROM %s WHERE instrument_id=s.instrument_id::uuid ORDER BY created_at DESC LIMIT 1) AS m ON true`,
		LocationTable, LocationTable,
	)
	dto.Id = uuid.NewString()

	status := dto.Status
	if status == "" {
		status = constants.LocationStatusUsed
	}

	var personId *string = &dto.PersonId
	if dto.PersonId == "" {
		personId = nil
	}
	var departmentId *string = &dto.DepartmentId
	if dto.DepartmentId == "" {
		departmentId = nil
	}

	args := []interface{}{dto.Id, dto.InstrumentId, dto.DateOfIssue, dto.DateOfReceiving, status, dto.NeedConfirm, personId, departmentId, dto.UserId}
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *LocationRepo) CreateSeveral(ctx context.Context, dto []*models.LocationDTO) error {
	args := make([]interface{}, 0)
	values := make([]string, 0, len(dto))
	for i, v := range dto {
		dto[i].Id = uuid.NewString()

		status := v.Status
		if status == "" {
			status = constants.LocationStatusUsed
		}

		var personId *string = &v.PersonId
		if v.PersonId == "" {
			personId = nil
		}
		var departmentId *string = &v.DepartmentId
		if v.DepartmentId == "" {
			departmentId = nil
		}

		tmp := []interface{}{v.Id, v.InstrumentId, v.DateOfIssue, v.DateOfReceiving, status, v.NeedConfirm, personId, departmentId, v.UserId}
		args = append(args, tmp...)
		numbers := []string{}
		for j := range tmp {
			numbers = append(numbers, fmt.Sprintf("$%d", i*len(tmp)+j+1))
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(numbers, ",")))
	}

	query := fmt.Sprintf(`INSERT INTO %s(id, instrument_id, date_of_issue, date_of_receiving, status, need_confirmed, 
		person_id, department_id, last_place_id, user_id)
		SELECT id::uuid, instrument_id::uuid, date_of_issue::DATE, date_of_receiving::DATE, 
			status, need_confirmed::boolean, person_id::uuid, s.department_id::uuid, COALESCE(m.department_id::text, ''), s.user_id::uuid
		FROM (VALUES %s) AS s(id, instrument_id, date_of_issue, date_of_receiving, status, need_confirmed, person_id, department_id, user_id)
		LEFT JOIN LATERAL (SELECT place, department_id FROM %s WHERE instrument_id=s.instrument_id::uuid ORDER BY created_at DESC LIMIT 1) AS m ON true`,
		LocationTable, strings.Join(values, ", "), LocationTable,
	)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *LocationRepo) Update(ctx context.Context, dto *models.LocationDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET date_of_issue=:date_of_issue, date_of_receiving=:date_of_receiving, status=:status, 
		person_id=:person_id, department_id=:department_id WHERE id=:id`,
		LocationTable,
	)

	res, err := r.db.NamedExecContext(ctx, query, dto)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return models.ErrNoRows
	}
	return nil
}

func (r *LocationRepo) SetPerson(ctx context.Context, id string) error {
	query := fmt.Sprintf(`UPDATE %s AS l SET person=e.name FROM %s AS e WHERE l.person_id=e.id AND l.person_id=$1`,
		LocationTable, EmployeeTable,
	)

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
func (r *LocationRepo) SetDepartment(ctx context.Context, id string) error {
	query := fmt.Sprintf(`UPDATE %s AS l SET place=d.name, person=e.name 
		FROM %s AS e INNER JOIN %s AS d ON e.department_id=d.id 
		WHERE l.department_id=d.id AND l.person_id=e.id AND l.department_id=$1`,
		LocationTable, EmployeeTable, DepartmentTable,
	)

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *LocationRepo) Receiving(ctx context.Context, dto *models.ReceivingDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET status=$1, date_of_receiving=$2, has_confirmed=$3 
		WHERE ARRAY[instrument_id] <@ $4 AND date_of_receiving='1970-01-01'::DATE`,
		LocationTable,
	)

	res, err := r.db.ExecContext(ctx, query, dto.Status, time.Now(), dto.HasConfirmed, pq.Array(dto.InstrumentIds))
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return models.ErrNoRows
	}
	return nil
}

func (r *LocationRepo) ForcedReceipt(ctx context.Context, dto *models.ForcedReceiptDTO) error {
	query := fmt.Sprintf(`UPDATE %s AS m SET date_of_receiving=$1, status='used' 
		WHERE instrument_id=$2 AND date_of_receiving='1970-01-01'::DATE`,
		LocationTable,
	)

	res, err := r.db.ExecContext(ctx, query, time.Now(), dto.InstrumentId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return models.ErrNoRows
	}
	return nil
}

func (r *LocationRepo) ForcedReceiptAll(ctx context.Context) error {
	query := fmt.Sprintf(`UPDATE %s AS m SET date_of_receiving=$1, status=(
			SELECT CASE WHEN status='used' THEN 'reserve' ELSE 'used' END FROM %s
			WHERE instrument_id=m.instrument_id AND date_of_receiving!='1970-01-01'::DATE ORDER BY date_of_issue DESC LIMIT 1
		) WHERE date_of_receiving='1970-01-01'::DATE AND date_of_issue < $2`,
		LocationTable, LocationTable,
	)

	limit := time.Now().Add(-time.Hour * 24 * 20) //20 days ago
	_, err := r.db.ExecContext(ctx, query, time.Now(), limit)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *LocationRepo) Delete(ctx context.Context, tx Tx, dto *models.DeleteLocationDTO) error {
	//? удалить перемещение, если оно не единственное
	query := fmt.Sprintf(`DELETE FROM %s AS m WHERE id=:id
		AND (SELECT COUNT(id) FROM %s WHERE instrument_id=m.instrument_id)>1`,
		LocationTable, LocationTable,
	)

	res, err := r.getExec(tx).NamedExecContext(ctx, query, dto)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return models.ErrNoRows
	}
	return nil
}
