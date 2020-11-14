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

	// e.GET("/", handler.FetchMahasiswa)	// http://localhost:8080/
	// // e.GET("/", handler.FetchMahasiswa, middleware.IsAuthenticated)	// http://localhost:8080/
	// e.GET("/:id", handler.GetByID)		// http://localhost:8080/2
	// e.POST("/", handler.Store)
	// e.PUT("/", handler.Update)
	// e.DELETE("/:id", handler.Delete)	// http://localhost:8080/8
}

func GenerateHashPassword(c echo.Context) error {
	password := c.Param("password")
	hash, _ := repository.HashPassword(password)

	return c.JSON(http.StatusOK, hash)
}

func isRequestValid2(m *domain.Users) (bool, error) {
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
	if ok, err = isRequestValid2(&users); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	res, err := u.UsersUC.CheckLogin(ctx, &users)
	fmt.Println("Handler res: ", res)
	fmt.Println()

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, res)
}

// func getStatusCode(err error) int {
// 	if err == nil {
// 		return http.StatusOK
// 	}

// 	logrus.Error(err)
// 	switch err {
// 	case domain.ErrInternalServerError:
// 		return http.StatusInternalServerError
// 	case domain.ErrNotFound:
// 		return http.StatusNotFound
// 	case domain.ErrConflict:
// 		return http.StatusConflict
// 	default:
// 		return http.StatusInternalServerError
// 	}
// }
