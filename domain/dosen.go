package domain

import "context"
import "time"
import "github.com/go-sql-driver/mysql"
// import "database/sql"

type Dosen struct {
	ID int64 `json:"id_dosen"`
	Nip int32 `json:"nip" validate:"required"`
	Name string `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at"`
	UpdatedAt mysql.NullTime `json:"updated_at"`	// mysql.NullTime handle null possible values
}

// DosenUsecase
type DosenUsecase interface {
	// Fetch(ctx context.Context, cursor string, num int64) ([]Dosen, string, error)
	Fetch(ctx context.Context, cursor string, num int64) (Response, error)
	// GetByID(ctx context.Context, id int64) (Dosen, error)
	GetByID(ctx context.Context, id int64) (Response, error)
	// GetByNIM(ctx context.Context, nim int32) (Mahasiswa, error)
	GetByNIP(ctx context.Context, nim int32) (Response, error)
	// Store(ctx context.Context, d *Dosen) error
	Store(ctx context.Context, d *Dosen) (Response, error)
	// Update(ctx context.Context, m *Mahasiswa) error
	// Delete(ctx context.Context, id int64) error
}

// DosenRepository
type DosenRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Dosen, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (Dosen, error)
	GetByNIP(ctx context.Context, nim int32) (Dosen, error)
	Store(ctx context.Context, d *Dosen) error
	// Update(ctx context.Context, m *Mahasiswa) error
	// Delete(ctx context.Context, id int64) error
}