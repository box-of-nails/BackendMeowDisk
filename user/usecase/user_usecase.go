package usecase

import (
	"database/sql"
	"github.com/box-of-nails/BackendMeowDisk/models"
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

func (userUcase UserUseCase) Register(users models.UserData) error {
	return userUcase.UserRepository.Register(users)
}

func (userUcase UserUseCase) SetCoockieinredis(http.Cookie, models.UserData) error {
	return userUcase.UserRepository.SetCoockieinredis(http.Cookie{}, models.UserData{})
}

func (userUcase UserUseCase) GetCoockieinredis(data models.UserData) string {
	return userUcase.UserRepository.GetCoockieinredis(models.UserData{})

}

func (userUcase UserUseCase) Login(users models.UserData) error {
	return userUcase.UserRepository.Login(users)
}
