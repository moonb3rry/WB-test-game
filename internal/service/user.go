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

type UserRepo interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	GetUserByName(ctx context.Context, name string) (model.User, error)
	GetInfoAboutLoaderByID(ctx context.Context, ID int) (model.AboutLoader, error)
	GetInfoAboutCustomerByID(ctx context.Context, ID int) (model.AboutCustomer, error)
	GetLoaderByID(ctx context.Context, ID int) (model.Loader, error)
	CreateCustomer(ctx context.Context, customer model.Customer) error
	CreateLoader(ctx context.Context, loader model.Loader) error
	UpdateCustomer(ctx context.Context, capital int, userID int) (model.Customer, error)
	UpdateLoader(ctx context.Context, loader model.Loader) (model.Loader, error)
}

type UserService struct {
	userRepo UserRepo
}

func NewUserService(ur UserRepo) *UserService {
	return &UserService{userRepo: ur}
}

func (us UserService) RegisterUser(ctx context.Context, user model.User) error {

	// валидация
	if user.Username == "" || user.Password == "" {
		return fmt.Errorf("Требуются имя пользователя и пароль")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("ошибка при хэшировании пароля: %w", err)
	}
	user.Password = hashedPassword

	userID, err := us.userRepo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении пользователя: %w", err)
	}

	switch user.UserRole {
	case "c":
		err := us.userRepo.CreateCustomer(ctx, generateCustomer(userID))
		if err != nil {
			return err
		}
		return nil
	case "l":
		err := us.userRepo.CreateLoader(ctx, generateLoader(userID))
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("Требуется роль")
	}
}

func (us UserService) LogInUser(ctx context.Context, userLogin model.LoginRequest) (string, error) {

	// валидация
	if userLogin.Username == "" || userLogin.Password == "" {
		return "", fmt.Errorf("Требуются имя пользователя и пароль")
	}

	// Получение пользователя из базы данных
	user, err := us.userRepo.GetUserByName(ctx, userLogin.Username)
	if err != nil {
		return "", fmt.Errorf("Неверное имя пользователя или пароль")
	}

	// Проверка пароля
	if !utils.CheckPasswordHash(userLogin.Password, user.Password) {
		return "", fmt.Errorf("Неверное имя пользователя или пароль")
	}

	// Генерация JWT-токена
	tokenString, err := middleware.GenerateJWT(user.UserID, user.UserRole)
	if err != nil {
		return "", fmt.Errorf("Ошибка при создании токена")
	}
	return tokenString, nil
}

func (us UserService) AboutUser(ctx context.Context, claims middleware.Claims) (interface{}, error) {
	switch claims.UserRole {

	case "c":
		data, err := us.InfoAboutCustomer(ctx, claims)
		if err != nil {
			return nil, err
		}
		customerData := model.AboutCustomer{
			data.Capital,
			data.AssignedLoaders,
		}
		return customerData, nil

	case "l":
		data, err := us.InfoAboutLoader(ctx, claims)
		if err != nil {
			return nil, err
		}
		loaderData := model.AboutLoader{
			data.Weight,
			data.Wage,
			data.Alcoholism,
			data.Fatigue,
		}
		return loaderData, nil

	default:
		return nil, fmt.Errorf("недопустимая роль пользователя")
	}
}

func (us UserService) InfoAboutCustomer(ctx context.Context, claims middleware.Claims) (*model.AboutCustomer, error) {
	customer, err := us.userRepo.GetInfoAboutCustomerByID(ctx, claims.UserID)
	if err != nil {
		return &model.AboutCustomer{}, err
	}
	return &customer, nil
}

func (us UserService) InfoAboutLoader(ctx context.Context, claims middleware.Claims) (*model.AboutLoader, error) {
	loader, err := us.userRepo.GetInfoAboutLoaderByID(ctx, claims.UserID)
	if err != nil {
		return &model.AboutLoader{}, err
	}
	return &loader, nil
}

func generateCustomer(id int) model.Customer {
	var c model.Customer
	c.UserID = id
	rand.Seed(time.Now().UnixNano()) // Инициализируйте генератор случайных чисел
	min := 10000
	max := 100000
	c.Capital = rand.Intn(max-min+1) + min
	return c
}

func generateLoader(id int) model.Loader {
	var l model.Loader
	l.UserID = id

	min := 10000
	max := 30000
	l.Wage = rand.Intn(max-min+1) + min

	min = 5
	max = 30
	l.MaxWeight = rand.Intn(max-min+1) + min

	min = 0
	max = 100
	l.Fatigue = rand.Intn(max-min+1) + min

	min = 0
	max = 1
	if rand.Intn(max-min+1)+min == 1 {
		l.Alcoholism = true
	} else {
		l.Alcoholism = false
	}

	return l
}
