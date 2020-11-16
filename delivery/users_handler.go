package delivery

import (
	"net/http"
	// "strconv"
	"campus/domain"
	"campus/repository"

	"github.com/labstack/echo"
	// "github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	// "path"
	// "strings"
	// "text/template"
	// "github.com/julienschmidt/httprouter"
	"fmt"
	// "database/sql"
	// "os"
)

// type ResponseError struct {
// 	Message string `json:"message`
// }

type UsersHandler struct {
	UsersUC domain.UsersUsecase
}

func NewUsersHandler(e *echo.Echo, us domain.UsersUsecase) {
	handler := &UsersHandler{
		UsersUC: us,
	}

	e.GET("/generate-hash/:password", GenerateHashPassword)	// http://localhost:9000/generate-hash/
	e.POST("/login", handler.CheckLogin)
}

func GenerateHashPassword(c echo.Context) error {
	password := c.Param("password")
	hash, _ := repository.HashPassword(password)

	return c.JSON(http.StatusOK, hash)
}

func isRequestUserValid(m *domain.Users) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UsersHandler) CheckLogin(c echo.Context) (err error) {
	var users domain.Users

	// username := c.FormValue("username")
	// password := c.FormValue("password")
	
	err = c.Bind(&users)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestUserValid(&users); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	res, err := u.UsersUC.CheckLogin(ctx, &users)
	fmt.Println("Handler res: ", res)
	fmt.Println("Handler err: ", err)
	fmt.Println()

	if err != nil {
		// return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return c.JSON(getStatusCode(err), res)
	}

	return c.JSON(http.StatusOK, res)
}
