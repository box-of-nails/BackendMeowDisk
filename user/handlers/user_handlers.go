package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/box-of-nails/BackendMeowDisk/models"
	"github.com/box-of-nails/BackendMeowDisk/user/usecase"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"time"
)

type UserHandlers struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandlers(db *sql.DB, redis *redis.Client) UserHandlers {
	userUseCase := usecase.NewUserUseCase(db, redis)
	return UserHandlers{userUseCase: userUseCase}
}

func (userH UserHandlers) Register(ctx echo.Context) error {
	var user models.UserData

	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		panic(err)
	}
	var users models.UserData
	err = json.Unmarshal(body, &users)
	if err != nil {
		fmt.Print(err)
	}
	err = userH.userUseCase.Register(users)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewV4().String(),
		Expires:  time.Now().AddDate(0, 0, 7),
		HttpOnly: true,
	}
	err = userH.userUseCase.SetCoockieinredis(cookie, user)
	if err != nil {
		return err
	}
	ctx.SetCookie(&cookie)
	return ctx.NoContent(http.StatusOK)
}

func (userH UserHandlers) Login(ctx echo.Context) error {

	var user models.UserData
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		fmt.Print(err)
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Print(err)
	}
	err = userH.userUseCase.Login(user)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, err)
	}
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewV4().String(),
		Expires:  time.Now().AddDate(0, 0, 7),
		HttpOnly: true,
	}
	err = userH.userUseCase.SetCoockieinredis(cookie, user)
	if err != nil {
		return err
	}
	ctx.SetCookie(&cookie)
	return ctx.NoContent(http.StatusNoContent)
}

func (userH UserHandlers) Logout(ctx echo.Context) error {
	var user models.UserData
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: true,
	}
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		fmt.Print(err)
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Print(err)
	}
	err = userH.userUseCase.Logout(user)
	if err != nil {
		return err
	}
	err = userH.userUseCase.Deletecoockieinredis(user)
	if err != nil {
		return err
	}
	ctx.SetCookie(&cookie)
	//err = ctx.Redirect(0, "/")
	//if err != nil {
	//	return err
	//}
	return ctx.NoContent(http.StatusOK)

}

func (userH UserHandlers) InitHandlers(server *echo.Echo) {

	server.PUT("/register", userH.Register)
	server.GET("/login", userH.Login)
	server.DELETE("/logout", userH.Logout)
}
