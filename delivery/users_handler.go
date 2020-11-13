package delivery

import (
	"net/http"
	// "strconv"
	// "campus/domain"
	"campus/repository"

	"github.com/labstack/echo"
	// "github.com/sirupsen/logrus"
	// validator "gopkg.in/go-playground/validator.v9"

	// "path"
	// "strings"
	// "text/template"
	// "github.com/julienschmidt/httprouter"
	// "fmt"
	// "database/sql"
	// "os"
)

func GenerateHashPassword(c echo.Context) error {
	password := c.Param("password")
	hash, _ := repository.HashPassword(password)

	return c.JSON(http.StatusOK, hash)
}

func NewUsersHandler(e *echo.Echo) {
	// handler := &MahasiswaHandler{
	// 	MUsecase: us,
	// }

	e.GET("/generate-hash/:password", GenerateHashPassword)	// http://localhost:9000/generate-hash/

	// e.GET("/", handler.FetchMahasiswa)	// http://localhost:8080/
	// // e.GET("/", handler.FetchMahasiswa, middleware.IsAuthenticated)	// http://localhost:8080/
	// e.GET("/:id", handler.GetByID)		// http://localhost:8080/2
	// e.POST("/", handler.Store)
	// e.PUT("/", handler.Update)
	// e.DELETE("/:id", handler.Delete)	// http://localhost:8080/8
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
