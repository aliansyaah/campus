package domain

// import "context"
import "time"
import "github.com/go-sql-driver/mysql"
// import "database/sql"

type MataKuliah struct {
	ID int64 `json:"id_mata_kuliah"`
	Name string `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt mysql.NullTime `json:"updated_at"`	// mysql.NullTime handle null possible values
}
