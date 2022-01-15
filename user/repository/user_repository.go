package repository

import (
	"database/sql"
	"errors"
	"github.com/box-of-nails/BackendMeowDisk/models"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"time"
)

type UserRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewUserRepository(database *sql.DB, redis *redis.Client) UserRepository {
	return UserRepository{db: database, redis: redis}
}

func (userRepo UserRepository) SetCoockieinredis(cookie http.Cookie, data models.UserData) error {
	err := userRepo.redis.Set(data.Id, cookie.Value, cookie.Expires.Sub(time.Now())).Err()
	if err != nil {
		return err
	}

	return nil
}

func (userRepo UserRepository) Deletecoockieinredis(data models.UserData) error {
	err := userRepo.redis.Del(data.Id).Err()
	if err != nil {
		return err
	}
	return nil
}

func (userRepo UserRepository) GetCoockieinredis(data models.UserData) string {

	coockieVal, err := userRepo.redis.Get(data.Id).Result()
	if err != nil {
		log.Fatalf("Value not found")
	}

	return coockieVal
}

func (userRepo UserRepository) Register(data models.UserData) error {
	_, err := userRepo.db.Exec(`
	insert into user_data (
	                     id,
	                     login,
	                     password
	                     )
	values ($1,$2,$3)`,
		data.Id,
		data.Login,
		data.Password,
	)
	if err != nil {
		return err
	}
	return nil
}

func (userRepo UserRepository) Login(data models.UserData) error {
	rows, err := userRepo.db.Query(`SELECT "id","login","password" FROM "user_data"`)
	if err != nil {
		panic(err)
	}
	var id, login, password string
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &login, &password)
		if err != nil {
			panic(err)
		}
		if id == data.Id && login == data.Login && password == data.Password {
			return nil
		}
	}
	return errors.New("incorrect login or password; or account does not exist")
}

func (userRepo UserRepository) Logout(models.UserData) error {
	return nil
}
