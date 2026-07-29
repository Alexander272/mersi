package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ChannelRepo struct {
	db *sqlx.DB
}

func NewChannelRepo(db *sqlx.DB) *ChannelRepo {
	return &ChannelRepo{
		db: db,
	}
}

type Channel interface {
	GetAll(ctx context.Context) ([]*models.Channel, error)
	Create(ctx context.Context, channel *models.Channel) error
	Update(ctx context.Context, channel *models.Channel) error
	Delete(ctx context.Context, id string) error
}

func (r *ChannelRepo) GetAll(ctx context.Context) ([]*models.Channel, error) {
	query := fmt.Sprintf(`SELECT id, name, description, most_channel_id FROM %s ORDER BY name`, ChannelTable)
	channels := []*models.Channel{}

	if err := r.db.SelectContext(ctx, &channels, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return channels, nil
}

func (r *ChannelRepo) Create(ctx context.Context, channel *models.Channel) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, name, description, most_channel_id) VALUES (:id, :name, :description, :most_channel_id)`, ChannelTable)
	channel.ID = uuid.NewString()

	_, err := r.db.NamedExecContext(ctx, query, channel)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *ChannelRepo) Update(ctx context.Context, channel *models.Channel) error {
	query := fmt.Sprintf(`UPDATE %s SET name=:name, description=:description, most_channel_id=:most_channel_id WHERE id=:id`, ChannelTable)

	res, err := r.db.NamedExecContext(ctx, query, channel)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	if err := checkRowsAffected(res); err != nil {
		return err
	}
	return nil
}

func (r *ChannelRepo) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, ChannelTable)

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	if err := checkRowsAffected(res); err != nil {
		return err
	}
	return nil
}
