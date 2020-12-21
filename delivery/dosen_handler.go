package delivery

import (
	"net/http"
	"strconv"
	"fmt"
	// "encoding/json"
	// "database/sql"
	"campus/domain"
	"campus/delivery/middleware"

	"github.com/labstack/echo"
	// "github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

// DosenHandler represent the httphandler for dosen
type DosenHandler struct {
	DosenUc domain.DosenUsecase
}

func NewDosenHandler(e *echo.Echo, du domain.DosenUsecase) {
	handler := &DosenHandler{
		DosenUc: du,
	}

	// Using auth
	e.GET("/dosen", handler.FetchDosen, middleware.IsAuthenticated)		// http://localhost:8080/dosen
	e.GET("/dosen/:id", handler.GetByID, middleware.IsAuthenticated)	// http://localhost:8080/dosen/2
	e.POST("/dosen", handler.Store, middleware.IsAuthenticated)
	e.PUT("/dosen", handler.Update, middleware.IsAuthenticated)
	e.DELETE("/dosen/:id", handler.Delete, middleware.IsAuthenticated)	// http://localhost:8080/8
}

// FetchDosen will fetch the dosen based on given params
func (d *DosenHandler) FetchDosen(c echo.Context) error {
	numS := c.QueryParam("num")			// ambil param dgn key "num"
	num, _ := strconv.Atoi(numS)		// convert string to int
	cursor := c.QueryParam("cursor")	// ambil param dgn key "cursor"
	ctx := c.Request().Context()

	// Panggil fungsi Fetch di usecase dosen
	// listDosen, nextCursor, err := d.DosenUc.Fetch(ctx, cursor, int64(num))
	res, err := d.DosenUc.Fetch(ctx, cursor, int64(num))
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

func (d *DosenHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	res, err := d.DosenUc.GetByID(ctx, id)
	fmt.Println("Handler res: ", res)
	fmt.Println("Handler err: ", err)

	if err != nil {
		// return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return c.JSON(getStatusCode(err), res)
	}

	return c.JSON(http.StatusOK, res)
}

func isRequestDosenValid(d *domain.Dosen) (bool, error) {
	validate := validator.New()
	err := validate.Struct(d)	// validasi struct
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *DosenHandler) Store(c echo.Context) (err error) {
	var dosen domain.Dosen

	// Method .Bind() milik echo.Context digunakan dgn disisipi param pointer object hasil cetakan struct Dosen
	err = c.Bind(&dosen)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	fmt.Println(&dosen)

	var ok bool
	if ok, err = isRequestDosenValid(&dosen); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	res, err := d.DosenUc.Store(ctx, &dosen)
	fmt.Println("Handler res: ", res)
	fmt.Println("Handler err: ", err)

	if err != nil {
		// return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return c.JSON(getStatusCode(err), res)
	}

	return c.JSON(http.StatusCreated, res)
}

func (d *DosenHandler) Update(c echo.Context) (err error) {
	var dosen domain.Dosen

	err = c.Bind(&dosen)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	// fmt.Println(&dosen)

	var ok bool
	if ok, err = isRequestDosenValid(&dosen); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// os.Exit(1)

	ctx := c.Request().Context()
	res, err := d.DosenUc.Update(ctx, &dosen)
	fmt.Println("Handler res: ", res)
	fmt.Println("Handler err: ", err)

	if err != nil {
		// return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return c.JSON(getStatusCode(err), res)
	}

	return c.JSON(http.StatusCreated, res)
}

func (d *DosenHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	res, err := d.DosenUc.Delete(ctx, id)
	fmt.Println("Handler res: ", res)
	fmt.Println("Handler err: ", err)

	if err != nil {
		// return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return c.JSON(getStatusCode(err), res)
	}

	// return c.NoContent(http.StatusNoContent)
	return c.JSON(http.StatusOK, res)
}