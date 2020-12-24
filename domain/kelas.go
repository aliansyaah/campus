package domain

import "context"
import "time"
import "github.com/go-sql-driver/mysql"
import "database/sql"

type Kelas struct {
	ID int64 `json:"id_kelas"`
	Name string `json:"name" validate:"required"`	// tag validate utk payload validation
	RuangID Ruang `json:"ruang_id" validate:"required"`
	MataKuliahID MataKuliah `json:"mata_kuliah_id" validate:"required"`
	DosenID Dosen `json:"dosen_id" validate:"required"`
	MahasiswaID Mahasiswa `json:"mahasiswa_id" validate:"required"`
	CreatedBy string `json:"created_by" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedBy sql.NullString `json:"updated_by"`
	UpdatedAt mysql.NullTime `json:"updated_at"`	// mysql.NullTime handle null possible values
	// RuangName Ruang `json:"ruang_name"`
}

// KelasUsecase
type KelasUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) (Response, error)
	// GetByID(ctx context.Context, id int64) (Response, error)
	// GetByNIP(ctx context.Context, nim int32) (Response, error)
	// Store(ctx context.Context, d *Dosen) (Response, error)
	// Update(ctx context.Context, d *Dosen) (Response, error)
	// Delete(ctx context.Context, id int64) (Response, error)
}

// KelasRepository
type KelasRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Kelas, nextCursor string, err error)
	// GetByID(ctx context.Context, id int64) (Dosen, error)
	// GetByNIP(ctx context.Context, nim int32) (Dosen, error)
	// Store(ctx context.Context, d *Dosen) error
	// Update(ctx context.Context, d *Dosen) error
	// Delete(ctx context.Context, id int64) error
}