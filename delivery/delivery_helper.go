package delivery

import (
	"net/http"
	"campus/domain"
	"github.com/sirupsen/logrus"
)

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

// Custom error validation notification
func validationNotif(namespace string, tag string) string {
	var fieldMessage string
	var tagMessage string

	switch {
		case (namespace == "Mahasiswa.Name"):
			fieldMessage = "Kolom nama mahasiswa "
		case (namespace == "Dosen.Name"):
			fieldMessage = "Kolom nama dosen "
		case (namespace == "Dosen.Nip"):
			fieldMessage = "Kolom NIP "
		default:
			fieldMessage = "(Kolom error tidak terdefinisi) "
	}

	switch {
		case (tag == "required"):
			tagMessage = "harus diisi"
		case (tag == "gte"):
			tagMessage = "terlalu pendek"
		case (tag == "lte"):
			tagMessage = "terlalu panjang"
		default:
			tagMessage = "(Tag error tidak terdefinisi)"
	}

	return fieldMessage+tagMessage
}