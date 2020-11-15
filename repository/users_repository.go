package repository

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
	// "campus/repository"
	"campus/domain"
	// "fmt"
	// "time"
	"fmt"
)

type usersRepository struct {
	Conn *sql.DB
}

func NewUsersRepository(Conn *sql.DB) domain.UsersRepository {
	return &usersRepository{Conn}
}

func (q *usersRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Users, err error) {
	rows, err := q.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(err)
		}
	}()

	// Membuat slice baru dari awal
	// make(type slice, panjang, kapasitas)
	result = make([]domain.Users, 0)	// tipe = slice domain.Users

	// fmt.Println("Looping:")

	// Looping rows hasil query
	for rows.Next() {
		toBeAdded := domain.Users{}		// struct mahasiswa pada layer domain/model
		// fmt.Println("toBeAdded: ", toBeAdded)
		// mhsID := int64(0)

		err = rows.Scan(
			&toBeAdded.IdUser,
			&toBeAdded.Username,
			// &mhsID,
			&toBeAdded.Name,
			&toBeAdded.Password,
			&toBeAdded.DefaultPassword,
		)
		// fmt.Println("err: ", err)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		// fmt.Println("toBeAdded 2: ", toBeAdded)
		
		// Jika tdk error, masukkan hasil ke slice result
		result = append(result, toBeAdded)
		fmt.Println()
		fmt.Println("result: ", result)
	}

	return result, nil
}

func (u *usersRepository) CheckLogin(ctx context.Context, dm *domain.Users) (res domain.Users, err error) {
	query := `SELECT * FROM users WHERE username=?`
	list, err := u.fetch(ctx, query, dm.Username)
	if err != nil {
		return domain.Users{}, err
	}
	fmt.Println("List: ", list)

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}