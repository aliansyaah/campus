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

	h.GET("/", handler.Index)
}

func (m *MahasiswaHandler) FetchMahasiswa(c echo.Context) error {
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()

	listMhs, nextCursor, err := a.MUsecase.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listMhs)
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