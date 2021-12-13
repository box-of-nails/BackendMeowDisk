package repository

import (
	"database/sql"
	"errors"
	"github.com/box-of-nails/BackendMeowDisk/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(database *sql.DB) UserRepository {
	return UserRepository{db: database}
}

func(userRepo UserRepository) Register(data models.UserData) error{
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

func(userRepo UserRepository) Login (data models.UserData) error {
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
