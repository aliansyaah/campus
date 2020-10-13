package domain

import "context"
import "time"
import "github.com/go-sql-driver/mysql"
import "database/sql"

type Mahasiswa struct {
	ID int64
	Nim int32
	Name string
	Semester sql.NullInt32							// sql.NullInt32 handle null possible values
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at"`
	UpdatedAt mysql.NullTime `json:"updated_at"`	// mysql.NullTime handle null possible values
}

// MahasiswaUsecase
type MahasiswaUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Mahasiswa, string, error)
}

// MahasiswaRepository
type MahasiswaRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Mahasiswa, nextCursor string, err error)
}