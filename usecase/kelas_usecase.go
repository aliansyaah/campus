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

func NewKelasUsecase(k domain.KelasRepository, timeout time.Duration) domain.KelasUsecase {
	return &kelasUsecase{
		kelasRepo: k,
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
	// fmt.Println("Usecase CheckIfExists domain.Kelas: ", domain.Kelas{})
	fmt.Println("Usecase CheckIfExists err: ", err)

	if err != nil {
		res.Message = err.Error()
		return 
	}

	// result harus sama dgn domain.Kelas (result kosong & domain.Kelas juga kosong)
	if result != (domain.Kelas{}) {
		res.Message = "This kelas data already exists"
		return
	}

	// Return result
	res.Status = true
	res.Message = "Your item is not exists"
	res.Data = map[string]interface{}{
		"data": result,
	}

	return
}

func (k *kelasUsecase) Store(c context.Context, dk *domain.Kelas) (res domain.Response, err error) {
	ctx, cancel := context.WithTimeout(c, k.contextTimeout)
	defer cancel()
	
	// Cek jika ada data yg sama
	existedKelas, err := k.CheckIfExists(ctx, dk)
	fmt.Println("Usecase existedKelas result: ", existedKelas)

	if err != nil {
		res.Message = err.Error()
		return 
	}

	if existedKelas.Status == false {
		err = domain.ErrConflict
		res.Message = existedKelas.Message
		return 
	}

	// Insert data
	err = k.kelasRepo.Store(ctx, dk)
	if err != nil {
		res.Message = err.Error()
		return 
	}

	// Return result
	res.Status = true
	res.Message = domain.SuccCreateData
	res.Data = dk

	return
}
