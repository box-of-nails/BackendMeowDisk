package usecase

import (
	"database/sql"
	"github.com/box-of-nails/BackendMeowDisk/models"
	"github.com/box-of-nails/BackendMeowDisk/user/repository"
)

type UserUseCase struct {
	UserRepository repository.UserRepository
}

func NewUserUseCase(db *sql.DB) UserUseCase {
	userRepository := repository.NewUserRepository(db)
	return UserUseCase{UserRepository: userRepository}
}

func(userUcase UserUseCase) Register(users models.UserData) error{
	return userUcase.UserRepository.Register(users)
}

func(userUcase UserUseCase) Login (users models.UserData) error{
	return userUcase.UserRepository.Login(users)
}
