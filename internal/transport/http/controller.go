package http

import (
	"WB_game/internal/middleware"
	"WB_game/internal/model"
	"context"
)

type UserService interface {
	RegisterUser(ctx context.Context, user model.User) error
	LogInUser(ctx context.Context, userLogin model.LoginRequest) (string, error)
	AboutUser(ctx context.Context, claims middleware.Claims) (interface{}, error)
	InfoAboutCustomer(ctx context.Context, claims middleware.Claims) (*model.AboutCustomer, error)
	InfoAboutLoader(ctx context.Context, claims middleware.Claims) (*model.AboutLoader, error)
}

type TaskService interface {
	GenerateTasks(ctx context.Context) (*[]model.Task, error)
	GetAvailableTasks(ctx context.Context, claims middleware.Claims) (*[]model.Task, error)
}

type GameService interface {
	StartGame(ctx context.Context, claims *middleware.Claims, start model.StartTask) (bool, error)
	TryToStart(ctx context.Context, claims *middleware.Claims, start model.StartTask) (bool, error)
}

type controller struct {
	userService UserService
	taskService TaskService
	gameService GameService
}

func newController(us UserService, ts TaskService, gs GameService) *controller {
	return &controller{us, ts, gs}
}
