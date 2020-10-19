package domain

import "context"
import "time"
import "github.com/go-sql-driver/mysql"
import "database/sql"

type Mahasiswa struct {
	ID int64 `json:"id"`
	Nim int32
	Name string
	// Semester int32
	Semester sql.NullInt32	`json:"semester"`		// sql.NullInt32 handle null possible values
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at"`
	UpdatedAt mysql.NullTime `json:"updated_at"`	// mysql.NullTime handle null possible values
}

// MahasiswaUsecase
type MahasiswaUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Mahasiswa, string, error)
	GetByID(ctx context.Context, id int64) (Mahasiswa, error)
	GetByNIM(ctx context.Context, nim int32) (Mahasiswa, error)
	// Store(context.Context, *Mahasiswa)
	Store(ctx context.Context, m *Mahasiswa) error
	Update(ctx context.Context, m *Mahasiswa) error
	Delete(ctx context.Context, id int64) error
}

// MahasiswaRepository
type MahasiswaRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Mahasiswa, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (Mahasiswa, error)
	GetByNIM(ctx context.Context, nim int32) (Mahasiswa, error)
	Store(ctx context.Context, m *Mahasiswa) error
	Update(ctx context.Context, m *Mahasiswa) error
	Delete(ctx context.Context, id int64) error
}