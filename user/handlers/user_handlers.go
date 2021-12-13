package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/box-of-nails/BackendMeowDisk/models"
	"github.com/box-of-nails/BackendMeowDisk/user/usecase"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
)

type UserHandlers struct{
	userUseCase  usecase.UserUseCase
}

func NewUserHandlers(db *sql.DB) UserHandlers {
	userUseCase := usecase.NewUserUseCase(db)
	return UserHandlers{userUseCase: userUseCase}
}


func(userH UserHandlers) Register(ctx echo.Context) error{
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
	return ctx.NoContent(http.StatusOK)
}

func (userH UserHandlers) Login(ctx echo.Context) error {

	var users models.UserData
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		fmt.Print(err)
	}
	err = json.Unmarshal(body, &users)
	if err != nil {
		fmt.Print(err)
	}
	err = userH.userUseCase.Login(users)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, err)
	}
	return ctx.NoContent(http.StatusOK)
	//cookie:=http.Cookie{
	//	Name: "session_id",
	//	Value: "MMRN9FDZx02MMgVo",
	//	Expires: expiration,
	//	HttpOnly: true,
	//}
	//http.SetCookie(w,&cookie)
	//http.Redirect(w,r, "/",http.StatusFound)
	//w.Write([]byte {'h','e','l'})
}

func (userH UserHandlers) InitHandlers(server *echo.Echo) {

	server.PUT("/register", userH.Register)
	server.GET("/login", userH.Login)
	//server.DELETE("/logout", albumD.GetAlbums)
}
