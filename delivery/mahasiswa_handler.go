package delivery

import (
	"net/http"
	"strconv"
	"campus/domain"

	"github.com/labstack/echo"
	// "github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	// "path"
	// "strings"
	// "text/template"
	// "github.com/julienschmidt/httprouter"
	"fmt"
	"database/sql"
	// "os"
	"campus/delivery/middleware"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message`
}

// MahasiswaHandler represent the httphandler for mahasiswa
type MahasiswaHandler struct {
	MUsecase domain.MahasiswaUsecase
}

func NewMahasiswaHandler(e *echo.Echo, us domain.MahasiswaUsecase) {
	handler := &MahasiswaHandler{
		MUsecase: us,
	}

	// Not using auth
	// e.GET("/", handler.FetchMahasiswa)	// http://localhost:8080/
	// e.GET("/:id", handler.GetByID)		// http://localhost:8080/2
	// e.POST("/", handler.Store)
	// e.PUT("/", handler.Update)
	// e.DELETE("/:id", handler.Delete)		// http://localhost:8080/8

	// Using auth
	e.GET("/mahasiswa", handler.FetchMahasiswa, middleware.IsAuthenticated)		// http://localhost:8080/mahasiswa
	e.GET("/mahasiswa/:id", handler.GetByID, middleware.IsAuthenticated)		// http://localhost:8080/mahasiswa/2
	e.POST("/mahasiswa", handler.Store, middleware.IsAuthenticated)
	e.PUT("/mahasiswa", handler.Update, middleware.IsAuthenticated)
	e.DELETE("/mahasiswa/:id", handler.Delete, middleware.IsAuthenticated)		// http://localhost:8080/mahasiswa/8
}

// FetchMahasiswa will fetch the mahasiswa based on given params
func (m *MahasiswaHandler) FetchMahasiswa(c echo.Context) error {
	numS := c.QueryParam("num")			// ambil param dgn key "num"
	num, _ := strconv.Atoi(numS)		// convert string to int
	cursor := c.QueryParam("cursor")	// ambil param dgn key "cursor"
	ctx := c.Request().Context()

	// Panggil fungsi Fetch di usecase mahasiswa
	listMhs, nextCursor, err := m.MUsecase.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listMhs)
}

func (m *MahasiswaHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	art, err := m.MUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, art)
}

func isRequestMahasiswaValid(m *domain.Mahasiswa) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func validateSemester(c echo.Context) (sql.NullInt32, error) {
	var mahasiswa domain.Mahasiswa
	var sem sql.NullInt32

	reqSem := c.FormValue("semester")
	// fmt.Println(reqSem)
	if reqSem != "" {
		fmt.Println("semester not empty")
		if err := sem.Scan(reqSem); err != nil {
			// panic(err)
			// return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
			return sem, err
		}
	}
	mahasiswa.Semester = sem

	// fmt.Println(&mahasiswa)
	// fmt.Println(mahasiswa.Semester)
	// fmt.Printf("var semester = %T\n", semester)
	// fmt.Printf("var mahasiswa.Semester = %T\n", mahasiswa.Semester)

	// return true, nil
	return mahasiswa.Semester, nil
}

func validateMahasiswa(c echo.Context, m domain.Mahasiswa) (domain.Mahasiswa, error) {
	// var mahasiswa domain.Mahasiswa
	var sem sql.NullInt32						// deklarasi var sem dgn tipe data NullInt32

	reqSem := c.FormValue("semester")			// ambil request dgn nama "semester"
	
	if reqSem != "" {							// jika request "semester" tidak kosong
		// fmt.Println("semester not empty")
		if err := sem.Scan(reqSem); err != nil {
			// panic(err)
			return m, err
		}
	}
	// mahasiswa.Semester = sem
	m.Semester = sem 	// variabel "sem" dimasukkan ke property "Semester" pada struct "mahasiswa"

	// fmt.Printf("var semester = %T\n", semester)
	// fmt.Printf("var mahasiswa.Semester = %T\n", mahasiswa.Semester)

	return m, nil
}

func (m *MahasiswaHandler) Store(c echo.Context) (err error) {
	var mahasiswa domain.Mahasiswa

	err = c.Bind(&mahasiswa)	// pointer variable domain mahasiswa
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	// fmt.Println(&mahasiswa)
	// fmt.Println(mahasiswa.Semester)

	var ok bool
	if ok, err = isRequestMahasiswaValid(&mahasiswa); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	/* Handle request tipe data DB sql.NullInt32 agar tidak null ketika insert */
	// var sem sql.NullInt32
	// if sem, err = validateSemester(c); err != nil {
	// 	return c.JSON(http.StatusBadRequest, err.Error())
	// }
	// mahasiswa.Semester = sem
	// fmt.Println(&mahasiswa)
	// fmt.Println(mahasiswa.Semester)

	/* Handle request tipe data DB sql.NullInt32 agar tidak null ketika insert */
	if mahasiswa, err = validateMahasiswa(c, mahasiswa); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// fmt.Println(&mahasiswa)
	// fmt.Println(mahasiswa)

	ctx := c.Request().Context()
	err = m.MUsecase.Store(ctx, &mahasiswa)		// pointer variable domain mahasiswa
	// fmt.Println(&mahasiswa)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, mahasiswa)
}

func (m *MahasiswaHandler) Update(c echo.Context) (err error) {
	var mahasiswa domain.Mahasiswa

	err = c.Bind(&mahasiswa)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	// fmt.Println(&mahasiswa)
	// fmt.Println(mahasiswa.Semester)

	var ok bool
	if ok, err = isRequestMahasiswaValid(&mahasiswa); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	/* Handle request tipe data DB sql.NullInt32 agar tidak null ketika insert */
	if mahasiswa, err = validateMahasiswa(c, mahasiswa); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// fmt.Println(&mahasiswa)
	// fmt.Println(mahasiswa)
	// os.Exit(1)

	ctx := c.Request().Context()
	err = m.MUsecase.Update(ctx, &mahasiswa)
	// fmt.Println(&mahasiswa)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, mahasiswa)
}

func (m *MahasiswaHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	err = m.MUsecase.Delete(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	fmt.Println(&m)

	return c.NoContent(http.StatusNoContent)
}

/*func (h *MahasiswaHandler) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	list, err := h.mahasiswausecase.Fetch(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var html string

	if strings.Contains(r.RequestURI, "question") {
		html = "admin/question.html"
	} else {
		html = ".html"
	}
}*/