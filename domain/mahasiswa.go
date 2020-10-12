package domain

import "context"
import "time"

type Mahasiswa struct {
	ID int64
	Nim int32
	Name string
	Semester int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MahasiswaUsecase
type MahasiswaUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Mahasiswa, string, error)
}

// MahasiswaRepository
type MahasiswaRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Mahasiswa, nextCursor string, err error)
}