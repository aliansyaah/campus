package usecase

import (
	"context"
	"time"
	"fmt"
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
func (k *kelasUsecase) Fetch(c context.Context, cursor string, num int64) (res domain.Response, err error) {
	// Jika param num = 0, tampilkan 10 data
	if num == 0 {
		num = 10
	}

	// context.WithTimeout utk proses cancellation jika proses query membutuhkan waktu yg lama
	ctx, cancel := context.WithTimeout(c, k.contextTimeout)
	defer cancel()

	// Panggil fungsi Fetch di repository dosen
	result, nextCursor, err := k.kelasRepo.Fetch(ctx, cursor, num)
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

func (k *kelasUsecase) CheckIfExists(c context.Context, dk *domain.Kelas) (res domain.Response, err error) {
	ctx, cancel := context.WithTimeout(c, k.contextTimeout)
	defer cancel()

	result, err := k.kelasRepo.CheckIfExists(ctx, dk)
	fmt.Println("Usecase CheckIfExists result: ", result)
	fmt.Println("Usecase CheckIfExists err: ", err)

	if err != nil {
		res.Message = err.Error()
		return 
	}

	res.Status = true
	// res.Message = "Data found"
	res.Message = domain.SuccDataFound
	res.Data = map[string]interface{}{
		"data": result,
	}

	return
}

func (k *kelasUsecase) Store(c context.Context, dk *domain.Kelas) (res domain.Response, err error) {
	ctx, cancel := context.WithTimeout(c, k.contextTimeout)
	defer cancel()
	
	// Cek jika ada data yg sama
	existedKelas, _ := k.CheckIfExists(ctx, dk)
	fmt.Println("Usecase existedKelas result: ", existedKelas)

	if existedKelas.Status == true {
		err = domain.ErrConflict
		res.Message = "This kelas data already exists"
		return 
	}

	err = k.kelasRepo.Store(ctx, dk)
	if err != nil {
		res.Message = err.Error()
		return 
	}

	// fmt.Println(dk)
	// fmt.Println(&dk)
	// fmt.Println(*dk)

	res.Status = true
	// res.Message = "Data successfully created"
	res.Message = domain.SuccCreateData
	res.Data = dk
	// res.Data = map[string]interface{}{
	// 	"data": dk,
	// }

	return
}
