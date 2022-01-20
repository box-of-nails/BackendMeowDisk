package main

import (
	"database/sql"
	"fmt"
	"github.com/box-of-nails/BackendMeowDisk/user/handlers"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	_ "net/http"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	pass   = "12345"
	dbname = "postgres"
)

type Handlers struct {
	UserHandlers handlers.UserHandlers
}

func InitPostgresql(server *echo.Echo) *sql.DB {
	psqlConn := fmt.Sprintf("host=%s port=%d user= %s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)

	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		server.Logger.Fatal("failed to connect to postgresql", err.Error())
	}
	return db
}

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return client
}

func NewHandlers(db *sql.DB, redis *redis.Client) Handlers {
	userHandlers := handlers.NewUserHandlers(db, redis)
	return Handlers{UserHandlers: userHandlers}
}

func main() {
	server := echo.New()
	db := InitPostgresql(server)
	defer func() {
		if db != nil {
			db.Close()
		}
	}()
	initRedis := InitRedis()
	api := NewHandlers(db, initRedis)
	api.UserHandlers.InitHandlers(server)
	server.Logger.Fatal(server.Start(":8080"))

	fmt.Println("starting server at :8080")
}
