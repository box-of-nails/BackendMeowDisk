package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/box-of-nails/BackendMeowDisk/models"
	"github.com/box-of-nails/BackendMeowDisk/user/usecase"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/labstack/echo"
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
		Value:    uuid.New().String(),
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
		Value:    uuid.New().String(),
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

func (userH UserHandlers) Upload(ctx echo.Context) error {
	fmt.Println("File Upload Endpoint Hit")

	ctx.Request().ParseMultipartForm(10 << 20)

	file, handler, err := ctx.Request().FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return nil
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	tempFile, err := ioutil.TempFile("/home/nikita/test", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(ctx.Response(), "Successfully Uploaded File\n")
	return ctx.NoContent(http.StatusFound)
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
	//err = ctx.Redirect(http.StatusPermanentRedirect, "/")
	//if err != nil {
	//	return ctx.JSON(http.StatusInternalServerError, err.Error())
	//}
	return ctx.NoContent(http.StatusOK)

}

func (userH UserHandlers) InitHandlers(server *echo.Echo) {

	server.PUT("/register", userH.Register)
	server.GET("/login", userH.Login)
	server.DELETE("/logout", userH.Logout)
	server.POST("/upload", userH.Upload)
}
