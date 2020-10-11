package usecase

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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

