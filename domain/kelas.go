package domain

import "context"
import "time"
import "github.com/go-sql-driver/mysql"
import "database/sql"

type Kelas struct {
	ID int64 `json:"id_kelas"`
	Name string `json:"name" validate:"required"`	// tag validate utk payload validation
	Ruang Ruang `json:"ruang_id" validate:"required"`
	MataKuliah MataKuliah `json:"mata_kuliah_id" validate:"required"`
	Dosen Dosen `json:"dosen_id" validate:"required"`
	Mahasiswa Mahasiswa `json:"mahasiswa_id" validate:"required"`
	CreatedBy string `json:"created_by" validate:"required"`
	// CreatedBy string `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedBy sql.NullString `json:"updated_by"`
	UpdatedAt mysql.NullTime `json:"updated_at"`	// mysql.NullTime handle null possible values
	// RuangName Ruang `json:"ruang_name"`
}

// KelasUsecase
type KelasUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) (Response, error)
	// GetByID(ctx context.Context, id int64) (Response, error)
	CheckIfExists(ctx context.Context, k *Kelas) (Response, error)
	Store(ctx context.Context, d *Kelas) (Response, error)
	// Update(ctx context.Context, d *Dosen) (Response, error)
	// Delete(ctx context.Context, id int64) (Response, error)
}

// KelasRepository
type KelasRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Kelas, nextCursor string, err error)
	// GetByID(ctx context.Context, id int64) (Dosen, error)
	CheckIfExists(ctx context.Context, k *Kelas) (Kelas, error)
	Store(ctx context.Context, d *Kelas) error
	// Update(ctx context.Context, d *Dosen) error
	// Delete(ctx context.Context, id int64) error
}