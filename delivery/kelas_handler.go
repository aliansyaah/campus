package delivery

import (
	"net/http"
	"strconv"
	"fmt"
	"encoding/json"
	// "database/sql"
	"campus/domain"
	"campus/delivery/middleware"

	"github.com/labstack/echo"
	// "github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

// KelasHandler represent the httphandler for dosen
type KelasHandler struct {
	KelasUc domain.KelasUsecase
}

func NewKelasHandler(e *echo.Echo, ku domain.KelasUsecase) {
	handler := &KelasHandler{
		KelasUc: ku,
	}

	// Using auth
	// e.GET("/kelas", handler.FetchKelas, middleware.IsAuthenticated)		// http://localhost:8080/kelas
	e.GET("/kelas", handler.FetchKelas)		// http://localhost:8080/kelas
	// e.GET("/kelas/:id", handler.GetByID, middleware.IsAuthenticated)	// http://localhost:8080/kelas/2
	e.POST("/kelas", handler.Store, middleware.IsAuthenticated)
	// e.PUT("/kelas", handler.Update, middleware.IsAuthenticated)
	// e.DELETE("/kelas/:id", handler.Delete, middleware.IsAuthenticated)	// http://localhost:8080/8
}

// FetchKelas will fetch the kelas based on given params
func (k *KelasHandler) FetchKelas(c echo.Context) error {
	numS := c.QueryParam("num")			// ambil param dgn key "num"
	num, _ := strconv.Atoi(numS)		// convert string to int
	cursor := c.QueryParam("cursor")	// ambil param dgn key "cursor"
	ctx := c.Request().Context()

	// Panggil fungsi Fetch di usecase kelas
	res, err := k.KelasUc.Fetch(ctx, cursor, int64(num))
	// fmt.Printf("Handler res: %+v", res)
	// fmt.Println()

	// Print pretty JSON
	/*resJSON, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		// log.Fatalf(err.Error())
		return c.JSON(getStatusCode(err), err.Error())
	}
	fmt.Printf("Handler res: Marshal function output %s\n", string(resJSON))*/
	fmt.Println("Handler err: ", err)

	// Mengambil object "nextCursor" dari variable res.Data yg bertipe data "map[string]interface{}"
	nextCursorInterface := res.Data.(map[string]interface{})["nextCursor"]
	fmt.Println("Handler nextCursorInterface: ", nextCursorInterface)

	// Convert interface to string
	nextCursor := fmt.Sprintf("%v", nextCursorInterface)
	fmt.Println("Handler nextCursor: ", nextCursor)
	fmt.Println()

	if err != nil {
		return c.JSON(getStatusCode(err), res)
	}

	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, res)
}

func isRequestKelasValid(k *domain.Kelas) (bool, error) {
	validate := validator.New()
	err := validate.Struct(k)	// validasi struct
	if err != nil {
		return false, err
	}
	return true, nil
}

func (k *KelasHandler) Store(c echo.Context) (err error) {
	var kelas domain.Kelas

	// Method .Bind() milik echo.Context digunakan dgn disisipi param pointer object hasil cetakan struct Kelas
	err = c.Bind(&kelas)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	// fmt.Println(&kelas)
	// fmt.Println()

	// Print pretty JSON
	resJSON, err := json.MarshalIndent(kelas, "", " ")
	if err != nil {
		// log.Fatalf(err.Error())
		return c.JSON(getStatusCode(err), err.Error())
	}
	fmt.Printf("Marshal function output %s\n", string(resJSON))

	// Validate input
	var ok bool
	if ok, err = isRequestKelasValid(&kelas); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Insert data
	ctx := c.Request().Context()
	res, err := k.KelasUc.Store(ctx, &kelas)
	fmt.Println("Handler res: ", res)
	fmt.Println("Handler err: ", err)

	if err != nil {
		return c.JSON(getStatusCode(err), res)
	}

	return c.JSON(http.StatusCreated, res)
}
