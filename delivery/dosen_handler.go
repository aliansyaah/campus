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
	// validator "gopkg.in/go-playground/validator.v9"
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
	// e.POST("/", handler.Store, middleware.IsAuthenticated)
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

// func isRequestMahasiswaValid(m *domain.Mahasiswa) (bool, error) {
// 	validate := validator.New()
// 	err := validate.Struct(m)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// func validateSemester(c echo.Context) (sql.NullInt32, error) {
// 	var mahasiswa domain.Mahasiswa
// 	var sem sql.NullInt32

// 	reqSem := c.FormValue("semester")
// 	// fmt.Println(reqSem)
// 	if reqSem != "" {
// 		fmt.Println("semester not empty")
// 		if err := sem.Scan(reqSem); err != nil {
// 			// panic(err)
// 			// return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
// 			return sem, err
// 		}
// 	}
// 	mahasiswa.Semester = sem

// 	// fmt.Println(&mahasiswa)
// 	// fmt.Println(mahasiswa.Semester)
// 	// fmt.Printf("var semester = %T\n", semester)
// 	// fmt.Printf("var mahasiswa.Semester = %T\n", mahasiswa.Semester)

// 	// return true, nil
// 	return mahasiswa.Semester, nil
// }

// func validateMahasiswa(c echo.Context, m domain.Mahasiswa) (domain.Mahasiswa, error) {
// 	// var mahasiswa domain.Mahasiswa
// 	var sem sql.NullInt32						// deklarasi var sem dgn tipe data NullInt32

// 	reqSem := c.FormValue("semester")			// ambil request dgn nama "semester"
	
// 	if reqSem != "" {							// jika request "semester" tidak kosong
// 		// fmt.Println("semester not empty")
// 		if err := sem.Scan(reqSem); err != nil {
// 			// panic(err)
// 			return m, err
// 		}
// 	}
// 	// mahasiswa.Semester = sem
// 	m.Semester = sem 	// variabel "sem" dimasukkan ke property "Semester" pada struct "mahasiswa"

// 	// fmt.Printf("var semester = %T\n", semester)
// 	// fmt.Printf("var mahasiswa.Semester = %T\n", mahasiswa.Semester)

// 	return m, nil
// }

// func (m *DosenHandler) Store(c echo.Context) (err error) {
// 	var mahasiswa domain.Mahasiswa

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
// 	// var sem sql.NullInt32
// 	// if sem, err = validateSemester(c); err != nil {
// 	// 	return c.JSON(http.StatusBadRequest, err.Error())
// 	// }
// 	// mahasiswa.Semester = sem
// 	// fmt.Println(&mahasiswa)
// 	// fmt.Println(mahasiswa.Semester)

// 	/* Handle request tipe data DB sql.NullInt32 agar tidak null ketika insert */
// 	if mahasiswa, err = validateMahasiswa(c, mahasiswa); err != nil {
// 		return c.JSON(http.StatusBadRequest, err.Error())
// 	}
// 	// fmt.Println(&mahasiswa)
// 	// fmt.Println(mahasiswa)

// 	ctx := c.Request().Context()
// 	err = m.DUsecase.Store(ctx, &mahasiswa)
// 	// fmt.Println(&mahasiswa)

// 	if err != nil {
// 		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
// 	}

// 	return c.JSON(http.StatusCreated, mahasiswa)
// }

// func (m *DosenHandler) Update(c echo.Context) (err error) {
// 	var mahasiswa domain.Mahasiswa

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