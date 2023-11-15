package repository

import (
	"WB_game/internal/model"
	"WB_game/pkg/postgres"
	"context"
)

type TaskRepo struct {
	db *postgres.Postgres
}

func NewTaskRepository(db *postgres.Postgres) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) CreateTask(ctx context.Context, task *model.TaskToGen) (*model.Task, error) {
	var newTask model.Task
	query := `INSERT INTO tasks (task_name, weight, status) VALUES ($1, $2, $3) RETURNING task_id, task_name, weight, status;`
	row := r.db.Pool.QueryRow(ctx, query, task.TaskName, task.Weight, task.Status)
	err := row.Scan(&newTask.TaskID, &newTask.TaskName, &newTask.Weight, &newTask.Status)
	if err != nil {
		return nil, err
	}
	return &newTask, nil
}

func (r *TaskRepo) GetTasksForCustomer(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	query := `SELECT t.* FROM tasks t WHERE t.status = false`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var task model.Task
		err = rows.Scan(&task.TaskID, &task.TaskName, &task.Weight, &task.Status, &task.CustomerID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil

}

func (r *TaskRepo) GetTasksForLoader(ctx context.Context, userID int) ([]model.Task, error) {
	query := `SELECT t.* FROM tasks t INNER JOIN assigned_loaders al ON t.task_id = al.task_id WHERE al.loader_id = $1 AND t.status = true `
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
	return tasks, nil
}

func (r *TaskRepo) GetTaskByID(ctx context.Context, taskID int) (*model.Task, error) {
	var task model.Task
	query := `SELECT t.* FROM tasks t WHERE t.task_id = $1`
	err := r.db.Pool.QueryRow(ctx, query, taskID).Scan(&task.TaskID, &task.TaskName, &task.Weight, &task.Status, &task.CustomerID)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepo) UpdateTask(ctx context.Context, taskID int, status bool) error {
	query := `UPDATE tasks t SET status = $1 WHERE t.task_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, status, taskID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepo) AddAssignedLoaders(ctx context.Context, userID int, taskID int) error {
	query := `INSERT INTO assigned_loaders(loader_id, task_id) VALUES ($1, $2)`
	_, err := r.db.Pool.Exec(ctx, query, userID, taskID)
	if err != nil {
		return err
	}
	return nil
}
