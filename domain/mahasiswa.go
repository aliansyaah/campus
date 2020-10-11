package domain

import "context"

type Mahasiswa struct {
	ID int64
	Nim int32
	Name string
	Semester int32
	CreatedAt string
	UpdatedAt string
}

// MahasiswaUsecase
type MahasiswaUsecase interface {
	Fetch(ctx context.Context) ([]Question, error)
}

// MahasiswaRepository
type MahasiswaRepository interface {
	Fetch(ctx context.Context) ([]Question, error)
}