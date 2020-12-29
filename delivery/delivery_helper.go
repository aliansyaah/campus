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

		case (namespace == "Dosen.Nip"):
			fieldMessage = "Kolom NIP "
		case (namespace == "Dosen.Name"):
			fieldMessage = "Kolom nama dosen "

		case (namespace == "Ruang.Name"):
			fieldMessage = "Kolom nama ruang "

		case (namespace == "MataKuliah.Name"):
			fieldMessage = "Kolom nama mata kuliah "

		case (namespace == "Kelas.Name"):
			fieldMessage = "Kolom nama kelas "
		case (namespace == "Kelas.Ruang.Name"):
			fieldMessage = "Kolom nama ruang "
		case (namespace == "Kelas.MataKuliah.Name"):
			fieldMessage = "Kolom nama mata kuliah "
		case (namespace == "Kelas.Dosen.Nip"):
			fieldMessage = "Kolom NIP "
		case (namespace == "Kelas.Dosen.Name"):
			fieldMessage = "Kolom nama dosen "
		case (namespace == "Kelas.Mahasiswa.Name"):
			fieldMessage = "Kolom nama mahasiswa "
		case (namespace == "Kelas.CreatedBy"):
			fieldMessage = "Kolom created by "
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