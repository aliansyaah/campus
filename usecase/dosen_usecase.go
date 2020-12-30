package usecase

import (
	"context"
	"time"
	"fmt"
	"campus/domain"
)

type dosenUsecase struct {
	dosenRepo domain.DosenRepository
	contextTimeout time.Duration
}

func NewDosenUsecase(d domain.DosenRepository, timeout time.Duration) domain.DosenUsecase {
	return &dosenUsecase{
		dosenRepo: d,
		contextTimeout: timeout,
	}
}

// func (d *dosenUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Dosen, nextCursor string, err error) {
func (d *dosenUsecase) Fetch(c context.Context, cursor string, num int64) (res domain.Response, err error) {
	// Jika param num = 0, tampilkan 10 data
	if num == 0 {
		num = 10
	}

	// context.WithTimeout utk proses cancellation jika proses query membutuhkan waktu yg lama
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	// Panggil fungsi Fetch di repository dosen
	result, nextCursor, err := d.dosenRepo.Fetch(ctx, cursor, num)
	// fmt.Println(res)

	if err != nil {
		nextCursor = ""
		res.Data = map[string]interface{}{
			"data": result,
			"nextCursor": nextCursor,
		}
		return
	}

	res.Status = true
	// res.Message = "Data was successfully obtained"
	res.Message = domain.SuccDataFound
	res.Data = map[string]interface{}{
		"data": result,
		"nextCursor": nextCursor,
	}

	// res.Data = map[string][]domain.Dosen{
	// 	"data": result,
	// }
	return res, nil
}

// func (d *dosenUsecase) GetByID(c context.Context, id int64) (res domain.Dosen, err error) {
func (d *dosenUsecase) GetByID(c context.Context, id int64) (res domain.Response, err error) {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	result, err := d.dosenRepo.GetByID(ctx, id)
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

// func (d *dosenUsecase) GetByNIP(c context.Context, nip int32) (res domain.Dosen, err error) {
func (d *dosenUsecase) GetByNIP(c context.Context, nip int32) (res domain.Response, err error) {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	result, err := d.dosenRepo.GetByNIP(ctx, nip)
	fmt.Println("Usecase GetByNIP result: ", result)
	fmt.Println("Usecase GetByNIP err: ", err)

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

func (d *dosenUsecase) Store(c context.Context, dd *domain.Dosen) (res domain.Response, err error) {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()
	
	// Cek jika ada NIP yg sama
	existedDosen, _ := d.GetByNIP(ctx, dd.Nip)
	fmt.Println("Usecase existedDosen result: ", existedDosen)

	if existedDosen.Status == true {
		err = domain.ErrConflict
		res.Message = "NIP already in use"
		return 
	}

	err = d.dosenRepo.Store(ctx, dd)
	if err != nil {
		res.Message = err.Error()
		return 
	}

	// fmt.Println(dd)
	// fmt.Println(&dd)
	// fmt.Println(*dd)

	res.Status = true
	// res.Message = "Data successfully created"
	res.Message = domain.SuccCreateData
	res.Data = dd
	// res.Data = map[string]interface{}{
	// 	"data": dd,
	// }

	return
}

func (d *dosenUsecase) Update(c context.Context, dd *domain.Dosen) (res domain.Response, err error) {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	err = d.dosenRepo.Update(ctx, dd)
	if err != nil {
		res.Message = err.Error()
		return 
	}

	res.Status = true
	res.Message = domain.SuccUpdateData
	res.Data = dd

	return 
}

func (d *dosenUsecase) Delete(c context.Context, id int64) (res domain.Response, err error) {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	existedDosen, err := d.dosenRepo.GetByID(ctx, id)
	if err != nil {
		res.Message = "No data can be deleted"
		return 
	}
	fmt.Println(existedDosen)
	// fmt.Println(domain.Dosen{})

	// if existedDosen == (domain.Dosen{}) {
	// 	err = domain.ErrNotFound
	// 	res.Message = "Entity tidak sama"
	// 	return
	// }

	err = d.dosenRepo.Delete(ctx, id)
	if err != nil {
		res.Message = err.Error()
		return 
	}

	res.Status = true
	// res.Message = "Data successfully deleted"
	res.Message = domain.SuccDeleteData
	res.Data = existedDosen

	return 
}