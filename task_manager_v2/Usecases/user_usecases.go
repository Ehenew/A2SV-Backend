package Usecases

import (
	"a2sv-backend/task_manager/Domain"
	"a2sv-backend/task_manager/Infrastructure"
	"a2sv-backend/task_manager/Repositories"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUsecase interface {
	Register(user Domain.User) error
	Login(username, password string) (string, Domain.User, error)
	Promote(userID primitive.ObjectID) error
}

type userUsecase struct {
	userRepo        Repositories.UserRepository
	passwordService Infrastructure.PasswordService
	jwtService      Infrastructure.JWTService
}

func NewUserUsecase(userRepo Repositories.UserRepository, passwordService Infrastructure.PasswordService, jwtService Infrastructure.JWTService) UserUsecase {
	return &userUsecase{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (u *userUsecase) Register(user Domain.User) error {
	// Check if username exists
	existingUser, err := u.userRepo.FindByUsername(user.Username)
	if err == nil && existingUser.Username != "" {
		return errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := u.passwordService.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Check if it's the first user or admin prefix
	count, err := u.userRepo.Count()
	if err != nil {
		return err
	}

	if count == 0 || strings.HasPrefix(user.Username, "admin_") {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}

	_, err = u.userRepo.Create(user)
	return err
}

func (u *userUsecase) Login(username, password string) (string, Domain.User, error) {
	user, err := u.userRepo.FindByUsername(username)
	if err != nil {
		return "", Domain.User{}, errors.New("invalid credentials")
	}

	err = u.passwordService.ComparePassword(user.Password, password)
	if err != nil {
		return "", Domain.User{}, errors.New("invalid credentials")
	}

	token, err := u.jwtService.GenerateToken(user)
	if err != nil {
		return "", Domain.User{}, err
	}

	return token, user, nil
}

func (u *userUsecase) Promote(userID primitive.ObjectID) error {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	user.Role = "admin"
	_, err = u.userRepo.Update(user)
	return err
}
