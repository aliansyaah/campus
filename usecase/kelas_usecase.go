package usecase

import (
	"context"
	"time"
	// "fmt"
	"campus/domain"
)

type kelasUsecase struct {
	kelasRepo domain.KelasRepository
	contextTimeout time.Duration
}

func NewKelasUsecase(d domain.KelasRepository, timeout time.Duration) domain.KelasUsecase {
	return &kelasUsecase{
		kelasRepo: d,
		contextTimeout: timeout,
	}
}

// func (d *kelasUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Kelas, nextCursor string, err error) {
func (d *kelasUsecase) Fetch(c context.Context, cursor string, num int64) (res domain.Response, err error) {
	// Jika param num = 0, tampilkan 10 data
	if num == 0 {
		num = 10
	}

	// context.WithTimeout utk proses cancellation jika proses query membutuhkan waktu yg lama
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	// Panggil fungsi Fetch di repository dosen
	result, nextCursor, err := d.kelasRepo.Fetch(ctx, cursor, num)
	// fmt.Println(res)

	if err != nil {
		nextCursor = ""
	}

	res.Status = true
	// res.Message = "Data was successfully obtained"
	res.Message = domain.SuccDataFound
	res.Data = map[string]interface{}{
		"data": result,
		"nextCursor": nextCursor,
	}
	
	return res, nil
}
