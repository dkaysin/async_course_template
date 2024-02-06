package service

import (
	database "async_course/main/internal/database"
	eventbus "async_course/main/internal/eventbus"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

type Service struct {
	config *viper.Viper
	db     *database.Database
	ew     *eventbus.EventWriter
}

func NewService(config *viper.Viper, db *database.Database, ew *eventbus.EventWriter) *Service {
	return &Service{
		config: config,
		db:     db,
		ew:     ew,
	}
}

func (s *Service) AddUser(ctx context.Context, userID string) error {
	return s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		q := `INSERT INTO test_table (user_id) VALUES ($1)`
		_, err := tx.Exec(ctx, q, userID)
		return err
	})
}
