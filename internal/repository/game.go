package repository

import (
	"WB_game/pkg/postgres"
	"context"
)

type gameRepo struct {
	db *postgres.Postgres
}

func NewGameRepository(db *postgres.Postgres) *gameRepo {
	return &gameRepo{db: db}
}

func (r *gameRepo) NewGame(ctx context.Context, userID int, result bool) error {
	query := `INSERT INTO game(user_id, game_result) VALUES ($1, $2)`
	_, err := r.db.Pool.Exec(ctx, query, userID, result)
	if err != nil {
		return err
	}
	return nil
}
