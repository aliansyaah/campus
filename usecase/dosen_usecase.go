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
	}

	res.Status = true
	res.Message = "Data was successfully obtained"
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
	res.Message = "Data found"
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
	res.Message = "Data found"
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
		res.Message = "NIP sudah terpakai"
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
	res.Message = "Data successfully created"
	res.Data = dd
	// res.Data = map[string]interface{}{
	// 	"data": dd,
	// }

	return
}

// func (m *dosenUsecase) Update(c context.Context, dm *domain.Mahasiswa) (err error) {
// 	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
// 	defer cancel()

// 	// Cek jika ada NIM yg sama
// 	// existedMahasiswa, _ := m.GetByNIM(ctx, dm.Nim)
// 	// if existedMahasiswa != (domain.Mahasiswa{}) {
// 	// 	return domain.ErrConflict
// 	// }

// 	// dm.UpdatedAt = time.Now()
// 	return m.mahasiswaRepo.Update(ctx, dm)
// }

// func (m *dosenUsecase) Delete(c context.Context, id int64) (err error) {
// 	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
// 	defer cancel()

// 	existedMahasiswa, err := m.mahasiswaRepo.GetByID(ctx, id)
// 	if err != nil {
// 		return 
// 	}
// 	// fmt.Println(existedMahasiswa)
// 	// fmt.Println(domain.Mahasiswa{})

// 	if existedMahasiswa == (domain.Mahasiswa{}) {
// 		return domain.ErrNotFound
// 	}

// 	return m.mahasiswaRepo.Delete(ctx, id)
// }