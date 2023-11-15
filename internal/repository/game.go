package repository

import (
	"WB_game/pkg/postgres"
	"context"
)

type GameRepo struct {
	db *postgres.Postgres
}

func NewGameRepository(db *postgres.Postgres) *GameRepo {
	return &GameRepo{db: db}
}

func (r *GameRepo) NewGame(ctx context.Context, userID int, result bool) error {
	query := `INSERT INTO game(user_id, game_result) VALUES ($1, $2)`
	_, err := r.db.Pool.Exec(ctx, query, userID, result)
	if err != nil {
		return err
	}
	return nil
}
