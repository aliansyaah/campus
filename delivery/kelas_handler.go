package delivery

import (
	"net/http"
	"strconv"
	"fmt"
	// "encoding/json"
	// "database/sql"
	"campus/domain"
	// "campus/delivery/middleware"

	"github.com/labstack/echo"
	// "github.com/sirupsen/logrus"
	// validator "gopkg.in/go-playground/validator.v9"
)

// KelasHandler represent the httphandler for dosen
type KelasHandler struct {
	KelasUc domain.KelasUsecase
}

func NewKelasHandler(e *echo.Echo, du domain.KelasUsecase) {
	handler := &KelasHandler{
		KelasUc: du,
	}

	// Using auth
	// e.GET("/kelas", handler.FetchKelas, middleware.IsAuthenticated)		// http://localhost:8080/kelas
	e.GET("/kelas", handler.FetchKelas)		// http://localhost:8080/kelas
	// e.GET("/kelas/:id", handler.GetByID, middleware.IsAuthenticated)	// http://localhost:8080/kelas/2
	// e.POST("/kelas", handler.Store, middleware.IsAuthenticated)
	// e.PUT("/kelas", handler.Update, middleware.IsAuthenticated)
	// e.DELETE("/kelas/:id", handler.Delete, middleware.IsAuthenticated)	// http://localhost:8080/8
}

// FetchKelas will fetch the kelas based on given params
func (d *KelasHandler) FetchKelas(c echo.Context) error {
	numS := c.QueryParam("num")			// ambil param dgn key "num"
	num, _ := strconv.Atoi(numS)		// convert string to int
	cursor := c.QueryParam("cursor")	// ambil param dgn key "cursor"
	ctx := c.Request().Context()

	// Panggil fungsi Fetch di usecase kelas
	// listDosen, nextCursor, err := d.KelasUc.Fetch(ctx, cursor, int64(num))
	res, err := d.KelasUc.Fetch(ctx, cursor, int64(num))
	fmt.Println("Handler res: ", res)
	fmt.Println("Handler err: ", err)

	// Mengambil object "nextCursor" dari variable res.Data yg bertipe data "map[string]interface{}"
	nextCursorInterface := res.Data.(map[string]interface{})["nextCursor"]
	fmt.Println("Handler nextCursorInterface: ", nextCursorInterface)

	// Convert interface to string
	nextCursor := fmt.Sprintf("%v", nextCursorInterface)
	fmt.Println("Handler nextCursor: ", nextCursor)

	// bs, _ := json.Marshal(res.Data)
	// fmt.Println(string(bs))

	fmt.Println()

	if err != nil {
		// return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return c.JSON(getStatusCode(err), res)
	}

	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, res)
}
