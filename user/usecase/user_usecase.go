package usecase

import (
	"database/sql"
	"github.com/box-of-nails/BackendMeowDisk/models"
	_ "github.com/box-of-nails/BackendMeowDisk/models"
	"github.com/box-of-nails/BackendMeowDisk/user/repository"
	"github.com/go-redis/redis"
	"net/http"
)

type UserUseCase struct {
	UserRepository repository.UserRepository
}

func NewUserUseCase(db *sql.DB, redis *redis.Client) UserUseCase {
	userRepository := repository.NewUserRepository(db, redis)
	return UserUseCase{UserRepository: userRepository}
}

func (userUcase UserUseCase) SetCoockieinredis(coockie http.Cookie, user models.UserData) error {
	return userUcase.UserRepository.SetCoockieinredis(coockie, user)
}

func (userUcase UserUseCase) GetCoockieinredis(user models.UserData) string {
	return userUcase.UserRepository.GetCoockieinredis(user)
}

func (userUcase UserUseCase) Deletecoockieinredis(user models.UserData) error {
	return userUcase.UserRepository.Deletecoockieinredis(user)
}

func (userUcase UserUseCase) Login(users models.UserData) error {
	return userUcase.UserRepository.Login(users)
}

func (userUcase UserUseCase) Logout(users models.UserData) error {
	return userUcase.UserRepository.Logout(users)
}

func (userUcase UserUseCase) Register(users models.UserData) error {
	return userUcase.UserRepository.Register(users)
}
