package usecase

import (
	"context"
	"time"

	// "github.com/sirupsen/logrus"
	// "golang.org/x/sync/errgroup"
	"campus/domain"
)

type mahasiswaUsecase struct {
	mahasiswaRepo domain.MahasiswaRepository
	contextTimeout time.Duration
}

func NewMahasiswaUsecase(m domain.MahasiswaRepository, timeout time.Duration) domain.MahasiswaUsecase {
	return &mahasiswaUsecase{
		mahasiswaRepo: m,
		contextTimeout: timeout,
	}
}

func (m *mahasiswaUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Mahasiswa, nextCursor string, err error) {
	// Jika param num = 0, tampilkan 10 data
	if num == 0 {
		num = 10
	}

	// context.WithTimeout utk proses cancellation jika proses query membutuhkan waktu yg lama
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	// Panggil fungsi Fetch di repository mahasiswa
	res, nextCursor, err = m.mahasiswaRepo.Fetch(ctx, cursor, num)
	if err != nil {
		nextCursor = ""
	}
	return
}

func (m *mahasiswaUsecase) GetByID(c context.Context, id int64) (res domain.Mahasiswa, err error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	res, err = m.mahasiswaRepo.GetByID(ctx, id)
	if err != nil {
		return 
	}

	return
}

func (m *mahasiswaUsecase) GetByNIM(c context.Context, nim int32) (res domain.Mahasiswa, err error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	res, err = m.mahasiswaRepo.GetByNIM(ctx, nim)
	if err != nil {
		return 
	}

	return
}

func (m *mahasiswaUsecase) Store(c context.Context, dm *domain.Mahasiswa) (err error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	
	// Cek jika ada NIM yg sama
	existedMahasiswa, _ := m.GetByNIM(ctx, dm.Nim)
	if existedMahasiswa != (domain.Mahasiswa{}) {
		return domain.ErrConflict
	}

	err = m.mahasiswaRepo.Store(ctx, dm)
	return
}

func (m *mahasiswaUsecase) Update(c context.Context, dm *domain.Mahasiswa) (err error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	// Cek jika ada NIM yg sama
	// existedMahasiswa, _ := m.GetByNIM(ctx, dm.Nim)
	// if existedMahasiswa != (domain.Mahasiswa{}) {
	// 	return domain.ErrConflict
	// }

	// dm.UpdatedAt = time.Now()
	return m.mahasiswaRepo.Update(ctx, dm)
}
