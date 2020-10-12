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
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	res, nextCursor, err = m.mahasiswaRepo.Fetch(ctx, cursor, num)
	if err != nil {
		nextCursor = ""
	}
	return
}
