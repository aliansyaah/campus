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

// Coba bikin custom notifikasi error
func notifError(field string, tag string) string {
	var pesan string

	switch {
		case (field == "Nama" || field == "Name") && (tag == "required"):
			pesan = "Kolom nama harus diisi"
		case (field == "Alamat") && (tag == "required"):
			pesan = "Kolom alamat harus diisi"
		case (field == "Telepon") && (tag == "gte"):
			pesan = "Nomor telepon terlalu pendek"
		case (field == "Telepon") && (tag == "lte"):
			pesan = "Nomor telepon terlalu panjang"
		case (field == "Nip") && (tag == "required"):
			pesan = "Kolom NIP harus diisi"
		default:
			pesan = "Pesan error tidak terdefinisi"
	}

	return pesan
}