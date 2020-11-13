package domain

// import "context"
// import "time"
// import "github.com/go-sql-driver/mysql"
// import "database/sql"

type Users struct {
	IdUser int `json:"id_user"`
	Username string `json:"username"`
	Name string `json:"name"`
	Password string `json:"password"`
	DefaultPassword string `json:"password"`
}