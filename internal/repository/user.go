package repository

import (
	"WB_game/internal/model"
	"WB_game/pkg/postgres"
	"context"
)

type UserRepo struct {
	db *postgres.Postgres
}

func NewUserRepository(db *postgres.Postgres) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *model.User) (int, error) {
	var userID int
	query := `INSERT INTO users (username, password, user_role) VALUES ($1, $2, $3) RETURNING user_id`
	err := r.db.Pool.QueryRow(ctx, query, user.Username, user.Password, user.UserRole).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *UserRepo) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users u WHERE u.username = $1`
	row := r.db.Pool.QueryRow(ctx, query, name)
	err := row.Scan(&user.UserID, &user.Username, &user.Password, &user.UserRole)
	if err != nil {
		return &model.User{}, err
	}
	return &user, nil
}

func (r *UserRepo) GetInfoAboutLoaderByID(ctx context.Context, ID int) (*model.AboutLoader, error) {
	var loader model.AboutLoader
	query := `SELECT l.max_weight, l.wage, l.alcoholism, l.fatigue FROM loaders l  WHERE l.user_id = $1`
	row := r.db.Pool.QueryRow(ctx, query, ID)
	err := row.Scan(&loader.Weight, &loader.Wage, &loader.Alcoholism, &loader.Fatigue)
	if err != nil {
		return &model.AboutLoader{}, err
	}
	return &loader, nil
}

func (r *UserRepo) GetInfoAboutCustomerByID(ctx context.Context, ID int) (*model.AboutCustomer, error) {
	var customer model.AboutCustomer
	var assignedLoaders []model.Loader
	query1 := `SELECT l.user_id, l.max_weight, l.alcoholism, l.fatigue, l.wage FROM loaders l
    			JOIN assigned_loaders al ON l.user_id = al.loader_id JOIN tasks t ON al.task_id = t.task_id
                WHERE t.customer_id = $1`
	rows, err := r.db.Pool.Query(ctx, query1, ID)
	if err != nil {
		return &model.AboutCustomer{}, nil
	}
	for rows.Next() {
		var loader model.Loader
		err := rows.Scan(&loader.UserID, &loader.MaxWeight, &loader.Alcoholism, &loader.Fatigue, &loader.Wage)
		if err != nil {
			return &model.AboutCustomer{}, err
		}
		assignedLoaders = append(assignedLoaders, loader)
	}
	customer.AssignedLoaders = assignedLoaders

	query2 := `SELECT c.capital FROM customers c WHERE c.user_id = $1`
	row := r.db.Pool.QueryRow(ctx, query2, ID)
	err = row.Scan(&customer.Capital)
	if err != nil {
		return &model.AboutCustomer{}, err

	}

	return &customer, nil
}

func (r *UserRepo) GetLoaderByID(ctx context.Context, ID int) (*model.Loader, error) {
	var loader model.Loader
	query := `SELECT l.* FROM loaders l  WHERE l.user_id = $1`
	row := r.db.Pool.QueryRow(ctx, query, ID)
	err := row.Scan(&loader.UserID, &loader.MaxWeight, &loader.Alcoholism, &loader.Fatigue, &loader.Wage)
	if err != nil {
		return &model.Loader{}, err
	}
	return &loader, nil
}

func (r *UserRepo) CreateCustomer(ctx context.Context, customer *model.Customer) error {
	query := `INSERT INTO customers (user_id, capital) VALUES ($1, $2) RETURNING user_id`
	_, err := r.db.Pool.Exec(ctx, query, customer.UserID, customer.Capital)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) CreateLoader(ctx context.Context, loader *model.Loader) error {
	query := `INSERT INTO loaders (user_id, max_weight, alcoholism, fatigue, wage) VALUES ($1, $2, $3, $4, $5) RETURNING user_id`
	_, err := r.db.Pool.Exec(ctx, query, loader.UserID, loader.MaxWeight, loader.Alcoholism, loader.Fatigue, loader.Wage)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) UpdateCustomer(ctx context.Context, capital int, userID int) (*model.Customer, error) {
	var newCustomer model.Customer
	query := `UPDATE customers c SET capital = $1 WHERE c.user_id = $2 RETURNING c.user_id, c.capital`
	err := r.db.Pool.QueryRow(ctx, query, capital, userID).Scan(&newCustomer.UserID, &newCustomer.Capital)
	if err != nil {
		return &model.Customer{}, err
	}
	return &newCustomer, nil
}

func (r *UserRepo) UpdateLoader(ctx context.Context, loader *model.Loader) (*model.Loader, error) {
	var newLoader model.Loader
	query := `UPDATE loaders l SET fatigue = $1 WHERE l.user_id = $2 RETURNING l.user_id, l.max_weight, l.alcoholism, l.fatigue, l.wage`
	err := r.db.Pool.QueryRow(ctx, query, loader.Fatigue, loader.UserID).Scan(&newLoader.UserID, &newLoader.MaxWeight, &newLoader.Alcoholism, &newLoader.Fatigue, &newLoader.Wage)
	if err != nil {
		return &model.Loader{}, err
	}
	return &newLoader, nil
}
