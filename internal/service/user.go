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
	CreateUser(ctx context.Context, user *model.User) (int, error)
	GetUserByName(ctx context.Context, name string) (*model.User, error)
	GetInfoAboutLoaderByID(ctx context.Context, ID int) (*model.AboutLoader, error)
	GetInfoAboutCustomerByID(ctx context.Context, ID int) (*model.AboutCustomer, error)
	GetLoaderByID(ctx context.Context, ID int) (*model.Loader, error)
	CreateCustomer(ctx context.Context, customer *model.Customer) error
	CreateLoader(ctx context.Context, loader *model.Loader) error
	UpdateCustomer(ctx context.Context, capital int, userID int) (*model.Customer, error)
	UpdateLoader(ctx context.Context, loader *model.Loader) (*model.Loader, error)
}

type UserService struct {
	userRepo UserRepo
}

func NewUserService(ur UserRepo) *UserService {
	return &UserService{userRepo: ur}
}

func (us UserService) RegisterUser(ctx context.Context, user *model.User) error {

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

func (us UserService) LogInUser(ctx context.Context, userLogin *model.LoginRequest) (string, error) {

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

func (us UserService) AboutUser(ctx context.Context, claims *middleware.Claims) (interface{}, error) {
	switch claims.UserRole {

	case "c":
		data, err := us.userRepo.GetInfoAboutCustomerByID(ctx, claims.UserID)
		if err != nil {
			return nil, err
		}
		customerData := &model.AboutCustomer{
			data.Capital,
			data.AssignedLoaders,
		}
		return customerData, nil

	case "l":
		data, err := us.userRepo.GetInfoAboutLoaderByID(ctx, claims.UserID)
		if err != nil {
			return nil, err
		}
		loaderData := &model.AboutLoader{
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

func generateCustomer(id int) *model.Customer {
	rand.Seed(time.Now().UnixNano()) // Инициализируйте генератор случайных чисел
	minCapital := 10000
	maxCapital := 100000
	return &model.Customer{
		UserID:  id,
		Capital: rand.Intn(maxCapital-minCapital+1) + minCapital,
	}
}

func generateLoader(id int) *model.Loader {
	var alco bool
	if rand.Intn(2) == 1 {
		alco = true
	} else {
		alco = false
	}

	return &model.Loader{
		UserID:     id,
		MaxWeight:  rand.Intn(30-5+1) + 5,
		Alcoholism: alco,
		Fatigue:    rand.Intn(100 - 0 + 1),
		Wage:       rand.Intn(30000-10000+1) + 10000,
	}
}
