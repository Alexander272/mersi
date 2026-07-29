package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres/pq_models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SectionRepo struct {
	db *sqlx.DB
}

func NewSectionRepo(db *sqlx.DB) *SectionRepo {
	return &SectionRepo{
		db: db,
	}
}

type Section interface {
	Get(ctx context.Context, req *models.GetSectionsDTO) ([]*models.Section, error)
	GetAll(ctx context.Context, req *models.GetAllSectionsDTO) ([]*models.Section, error)
	GetGrouped(ctx context.Context, req *models.GetGroupedSectionDTO) ([]*models.GroupedSections, error)
	Create(ctx context.Context, dto *models.SectionDTO) error
	Update(ctx context.Context, dto *models.SectionDTO) error
	Delete(ctx context.Context, dto *models.DeleteSectionDTO) error
}

func (r *SectionRepo) Get(ctx context.Context, req *models.GetSectionsDTO) ([]*models.Section, error) {
	query := fmt.Sprintf(`SELECT id, name, realm_id, position, bid_type, verification_day, created_at 
		FROM %s WHERE realm_id=$1 ORDER BY position`,
		SectionTable,
	)
	data := []*models.Section{}

	if err := r.db.SelectContext(ctx, &data, query, req.RealmID); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *SectionRepo) GetAll(ctx context.Context, req *models.GetAllSectionsDTO) ([]*models.Section, error) {
	conditionsParts := []string{}

	if req.IsActive.Defined {
		conditionsParts = append(conditionsParts, fmt.Sprintf("r.is_active=%t", req.IsActive.Value))
	}
	if req.HasReturnNotice.Defined {
		conditionsParts = append(conditionsParts, fmt.Sprintf("r.return_notice=%t", req.HasReturnNotice.Value))
	}

	conditions := ""
	if len(conditionsParts) > 0 {
		conditions = fmt.Sprintf("WHERE %s", strings.Join(conditionsParts, " AND "))
	}

	query := fmt.Sprintf(`SELECT s.id, s.name, realm_id, s.position, s.bid_type, s.verification_day, s.created_at 
		FROM %s AS s JOIN %s AS r ON realm_id=r.id
		%s
		ORDER BY position`,
		SectionTable, RealmTable, conditions,
	)
	data := []*models.Section{}

	if err := r.db.SelectContext(ctx, &data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *SectionRepo) GetGrouped(ctx context.Context, req *models.GetGroupedSectionDTO) ([]*models.GroupedSections, error) {
	query := fmt.Sprintf(`SELECT s.id, s.name, r.name AS title, realm, realm_id, position, bid_type, verification_day, s.created_at 
		FROM %s AS s 
		INNER JOIN %s AS r ON realm_id=r.id ORDER BY r.created_at, position`,
		SectionTable, RealmTable,
	)
	tmp := []*pq_models.Section{}

	if err := r.db.SelectContext(ctx, &tmp, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	data := []*models.GroupedSections{}
	for i, d := range tmp {
		section := &models.Section{
			ID:              d.ID,
			RealmID:         d.RealmID,
			Name:            d.Name,
			Position:        d.Position,
			BidType:         d.BidType,
			VerificationDay: d.VerificationDay,
			CreatedAt:       d.CreatedAt,
		}

		if i == 0 || d.RealmID != data[len(data)-1].ID {
			data = append(data, &models.GroupedSections{
				ID:       d.RealmID,
				Title:    d.RealmTitle,
				Realm:    d.Realm,
				Sections: []*models.Section{section},
			})
		} else {
			data[len(data)-1].Sections = append(data[len(data)-1].Sections, section)
		}
	}
	return data, nil
}

func (r *SectionRepo) Create(ctx context.Context, dto *models.SectionDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, name, realm_id, position, bid_type, verification_day) 
		VALUES (:id, :name, :realm_id, :position, :bid_type, :verification_day)`,
		SectionTable,
	)
	dto.ID = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SectionRepo) Update(ctx context.Context, dto *models.SectionDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET name=:name, realm_id=:realm_id, position=:position, bid_type=:bid_type, verification_day=:verification_day 
		WHERE id=:id`,
		SectionTable,
	)

	res, err := r.db.NamedExecContext(ctx, query, dto)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	if err := checkRowsAffected(res); err != nil {
		return err
	}
	return nil
}

func (r *SectionRepo) Delete(ctx context.Context, dto *models.DeleteSectionDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, SectionTable)

	res, err := r.db.NamedExecContext(ctx, query, dto)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	if err := checkRowsAffected(res); err != nil {
		return err
	}
	return nil
}
