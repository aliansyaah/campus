package domain

import "context"
// import "time"
// import "github.com/go-sql-driver/mysql"
// import "database/sql"

type Users struct {
	IdUser int `json:"id_user"`
	Username string `json:"username" validate:"required"`
	Name string `json:"name"`
	Password string `json:"password" validate:"required"`
	DefaultPassword string `json:"default_password"`
}

type UsersUsecase interface {
	CheckLogin(ctx context.Context, u *Users) (Response, error)

	// Fetch(ctx context.Context, cursor string, num int64) ([]Mahasiswa, string, error)
	// GetByID(ctx context.Context, id int64) (Mahasiswa, error)
	// GetByNIM(ctx context.Context, nim int32) (Mahasiswa, error)
	// Store(ctx context.Context, m *Mahasiswa) error
	// Update(ctx context.Context, m *Mahasiswa) error
	// Delete(ctx context.Context, id int64) error
}

type UsersRepository interface {
	GetByUsername(ctx context.Context, u *Users) (Users, error)
	
	// Fetch(ctx context.Context, cursor string, num int64) (res []Mahasiswa, nextCursor string, err error)
	// GetByID(ctx context.Context, id int64) (Mahasiswa, error)
	// GetByNIM(ctx context.Context, nim int32) (Mahasiswa, error)
	// Store(ctx context.Context, m *Mahasiswa) error
	// Update(ctx context.Context, m *Mahasiswa) error
	// Delete(ctx context.Context, id int64) error
}