package delivery

import (
	"net/http"
	"strconv"
	"campus/domain"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	// "path"
	// "strings"
	// "text/template"
	// "github.com/julienschmidt/httprouter"
	"fmt"
	"database/sql"
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

	e.GET("/", handler.FetchMahasiswa)	// http://localhost:8080/
	e.GET("/:id", handler.GetByID)		// http://localhost:8080/2
	e.POST("/", handler.Store)
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

func isRequestValid(m *domain.Mahasiswa) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func validateInput(c echo.Context) (sql.NullInt32, error) {
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

func validateInputNull(c echo.Context, m domain.Mahasiswa) (domain.Mahasiswa, error) {
	// var mahasiswa domain.Mahasiswa
	var sem sql.NullInt32

	reqSem := c.FormValue("semester")
	// fmt.Println(reqSem)
	if reqSem != "" {
		fmt.Println("semester not empty")
		if err := sem.Scan(reqSem); err != nil {
			// panic(err)
			// return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
			return m, err
		}
	}
	// mahasiswa.Semester = sem
	m.Semester = sem

	// fmt.Println(&mahasiswa)
	// fmt.Println(mahasiswa.Semester)
	// fmt.Printf("var semester = %T\n", semester)
	// fmt.Printf("var mahasiswa.Semester = %T\n", mahasiswa.Semester)

	// return true, nil
	return m, nil
}

func (m *MahasiswaHandler) Store(c echo.Context) (err error) {
	var mahasiswa domain.Mahasiswa

	err = c.Bind(&mahasiswa)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	fmt.Println(&mahasiswa)
	fmt.Println(mahasiswa.Semester)

	var ok bool
	if ok, err = isRequestValid(&mahasiswa); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	/* Handle request tipe data DB sql.NullInt32 agar tidak null ketika insert */
	// var sem sql.NullInt32
	// if sem, err = validateInput(c); err != nil {
	// 	return c.JSON(http.StatusBadRequest, err.Error())
	// }
	// mahasiswa.Semester = sem
	// fmt.Println(&mahasiswa)
	// fmt.Println(mahasiswa.Semester)

	if mahasiswa, err = validateInputNull(c, mahasiswa); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println(&mahasiswa)
	fmt.Println(mahasiswa)

	ctx := c.Request().Context()
	err = m.MUsecase.Store(ctx, &mahasiswa)
	fmt.Println(&mahasiswa)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, mahasiswa)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
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