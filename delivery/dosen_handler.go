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
	// e.PUT("/", handler.Update, middleware.IsAuthenticated)
	// e.DELETE("/:id", handler.Delete, middleware.IsAuthenticated)	// http://localhost:8080/8
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

func isDosenRequestValid(m *domain.Dosen) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *DosenHandler) Store(c echo.Context) (err error) {
	var dosen domain.Dosen

	err = c.Bind(&dosen)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	fmt.Println(&dosen)

	var ok bool
	if ok, err = isDosenRequestValid(&dosen); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	res, err := m.DosenUc.Store(ctx, &dosen)
	fmt.Println("Handler res: ", res)
	fmt.Println("Handler err: ", err)

	if err != nil {
		// return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return c.JSON(getStatusCode(err), res)
	}

	return c.JSON(http.StatusCreated, res)
}

// func (m *DosenHandler) Update(c echo.Context) (err error) {
// 	var mahasiswa domain.Dosen

// 	err = c.Bind(&mahasiswa)
// 	if err != nil {
// 		return c.JSON(http.StatusUnprocessableEntity, err.Error())
// 	}
// 	// fmt.Println(&mahasiswa)
// 	// fmt.Println(mahasiswa.Semester)

// 	var ok bool
// 	if ok, err = isRequestMahasiswaValid(&mahasiswa); !ok {
// 		return c.JSON(http.StatusBadRequest, err.Error())
// 	}

// 	/* Handle request tipe data DB sql.NullInt32 agar tidak null ketika insert */
// 	if mahasiswa, err = validateMahasiswa(c, mahasiswa); err != nil {
// 		return c.JSON(http.StatusBadRequest, err.Error())
// 	}
// 	// fmt.Println(&mahasiswa)
// 	// fmt.Println(mahasiswa)
// 	// os.Exit(1)

// 	ctx := c.Request().Context()
// 	err = m.DUsecase.Update(ctx, &mahasiswa)
// 	// fmt.Println(&mahasiswa)

// 	if err != nil {
// 		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
// 	}

// 	return c.JSON(http.StatusCreated, mahasiswa)
// }

// func (m *DosenHandler) Delete(c echo.Context) error {
// 	idP, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
// 	}

// 	id := int64(idP)
// 	ctx := c.Request().Context()

// 	err = m.DUsecase.Delete(ctx, id)
// 	if err != nil {
// 		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
// 	}
// 	fmt.Println(&m)

// 	return c.NoContent(http.StatusNoContent)
// }