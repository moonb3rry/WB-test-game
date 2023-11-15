package service

import (
	"WB_game/internal/middleware"
	"WB_game/internal/model"
	"WB_game/utils"
	"context"
	"fmt"
	"math/rand"
	"time"
)

type TaskRepo interface {
	CreateTask(ctx context.Context, task model.TaskToGen) (*model.Task, error)
	GetTasksForCustomer(cxt context.Context) (*[]model.Task, error)
	GetTasksForLoader(cxt context.Context, userID int) (*[]model.Task, error)
	GetTaskByID(ctx context.Context, taskID int) (*model.Task, error)
	UpdateTask(ctx context.Context, taskID int, status bool) error
}

type TaskService struct {
	taskRepo TaskRepo
}

func NewTaskService(tr TaskRepo) *TaskService {
	return &TaskService{taskRepo: tr}
}

func (ts TaskService) GenerateTasks(ctx context.Context) (*[]model.Task, error) {
	rand.Seed(time.Now().UnixNano())
	numberOfTasks := 5
	var selectedTasks []model.Task
	for i := 0; i < numberOfTasks; i++ {
		randomIndex := rand.Intn(len(utils.AllTasks))
		TaskToAdd := utils.AllTasks[randomIndex]
		Task, err := ts.taskRepo.CreateTask(ctx, TaskToAdd)
		if err != nil {
			return nil, err
		}
		selectedTasks = append(selectedTasks, *Task)

	}
	return &selectedTasks, nil
}

func (ts TaskService) GetAvailableTasks(ctx context.Context, claims middleware.Claims) (*[]model.Task, error) {
	switch claims.UserRole {

	case "c":
		data, err := ts.taskRepo.GetTasksForCustomer(ctx)
		if err != nil {
			return nil, err
		}
		return data, nil

	case "l":
		data, err := ts.taskRepo.GetTasksForLoader(ctx, claims.UserID)
		if err != nil {
			return nil, err
		}
		return data, nil

	default:
		return nil, fmt.Errorf("недопустимая роль пользователя")
	}
}
