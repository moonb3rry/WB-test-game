package repository

import (
	"WB_game/internal/model"
	"WB_game/pkg/postgres"
	"context"
)

type taskRepo struct {
	db *postgres.Postgres
}

func NewTaskRepository(db *postgres.Postgres) *taskRepo {
	return &taskRepo{db: db}
}

func (r *taskRepo) CreateTask(ctx context.Context, task model.TaskToGen) (*model.Task, error) {
	var newTask model.Task
	query := `INSERT INTO tasks (task_name, weight, status) VALUES ($1, $2, $3) RETURNING task_id, task_name, weight, status;`
	row := r.db.Pool.QueryRow(ctx, query, task.TaskName, task.Weight, task.Status)
	err := row.Scan(&newTask.TaskID, &newTask.TaskName, &newTask.Weight, &newTask.Status)
	if err != nil {
		return nil, err
	}
	return &newTask, nil
}

func (r *taskRepo) GetTasksForCustomer(ctx context.Context) (*[]model.Task, error) {
	query := `SELECT t.* FROM tasks t WHERE t.status = false`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		err = rows.Scan(&task.TaskID, &task.TaskName, &task.Weight, &task.Status, &task.CustomerID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return &tasks, nil

}

func (r *taskRepo) GetTasksForLoader(ctx context.Context, userID int) (*[]model.Task, error) {
	query := `SELECT t.* FROM tasks t LEFT JOIN assigned_loaders al ON t.task_id = al.task_id WHERE al.loader_id = $1 AND t.status = true`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		err = rows.Scan(&task.TaskID, &task.TaskName, &task.Weight, &task.Status, &task.CustomerID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return &tasks, nil
}

func (r *taskRepo) GetTaskByID(ctx context.Context, taskID int) (*model.Task, error) {
	var task model.Task
	query := `SELECT t.* FROM tasks t WHERE t.task_id = $1`
	err := r.db.Pool.QueryRow(ctx, query, taskID).Scan(&task.TaskID, &task.TaskName, &task.Weight, &task.Status, &task.CustomerID)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepo) UpdateTask(ctx context.Context, taskID int, status bool) error {
	query := `UPDATE tasks t SET status = $1 WHERE t.task_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, status, taskID)
	if err != nil {
		return err
	}
	return nil
}
