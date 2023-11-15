package service

import (
	"WB_game/internal/middleware"
	"WB_game/internal/model"
	"context"
	"fmt"
)

type GameRepo interface {
	NewGame(ctx context.Context, userID int, result bool) error
}

type GameService struct {
	gameRepo GameRepo
	userRepo UserRepo
	taskRepo TaskRepo
}

func NewGameService(gr GameRepo, ur UserRepo, tr TaskRepo) *GameService {
	return &GameService{gr, ur, tr}
}

func (gs *GameService) StartGame(ctx context.Context, claims *middleware.Claims, start *model.StartTask) (bool, error) {

	if claims.UserRole != "c" {
		return false, fmt.Errorf("недопустимая роль пользователя")
	}

	customer, err := gs.userRepo.GetInfoAboutCustomerByID(ctx, claims.UserID)
	var result bool
	if err != nil {
		return false, err
	}
	task, err := gs.taskRepo.GetTaskByID(ctx, start.TaskID)

	if err != nil {
		return false, err
	}
	var loadersWeight int
	var loadersWage int

	for _, l := range start.AssignedLoaders {
		loader, err := gs.userRepo.GetLoaderByID(ctx, l)
		if err != nil {
			return false, err
		}

		loadersWeight += loader.MaxWeight * (100 - loader.Fatigue/100)
		loadersWage += loader.Wage
		if loader.Alcoholism {
			loader.Fatigue += 50
		} else {
			loader.Fatigue += 20
		}
		if loader.Fatigue > 100 {
			loader.Fatigue = 100
		}
		loader, err = gs.userRepo.UpdateLoader(ctx, loader)
		if err != nil {
			return false, err
		}
		customer.AssignedLoaders = append(customer.AssignedLoaders, *loader)
		err = gs.taskRepo.AddAssignedLoaders(ctx, loader.UserID, task.TaskID)
		if err != nil {
			return false, err
		}
	}

	customer.Capital -= loadersWage
	_, err = gs.userRepo.UpdateCustomer(ctx, customer.Capital, claims.UserID)
	if err != nil {
		return false, err
	}
	if task.Weight <= loadersWeight && customer.Capital >= loadersWage {
		err = gs.taskRepo.UpdateTask(ctx, task.TaskID, !task.Status)
		if err != nil {
			return false, err
		}
		result = true
	}
	print(task.Weight, loadersWeight, customer.Capital, loadersWage, result)
	err = gs.gameRepo.NewGame(ctx, claims.UserID, result)
	if err != nil {
		return false, fmt.Errorf("Не удалось добавить результат" + err.Error())
	}
	return result, nil

}
