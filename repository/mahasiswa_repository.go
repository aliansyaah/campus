package repository

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
	// "campus/repository"
	"campus/domain"
	// "fmt"
	"time"
	"fmt"
)

type mahasiswaRepository struct {
	Conn *sql.DB
}

func NewMahasiswaRepository(Conn *sql.DB) domain.MahasiswaRepository {
	return &mahasiswaRepository{Conn}
}

func (q *mahasiswaRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Mahasiswa, err error) {
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
	result = make([]domain.Mahasiswa, 0)	// tipe = slice domain.Mahasiswa

	// fmt.Println("Looping:")

	// Looping rows hasil query
	for rows.Next() {
		toBeAdded := domain.Mahasiswa{}		// struct mahasiswa pada layer domain/model
		// fmt.Println("toBeAdded: ", toBeAdded)
		// mhsID := int64(0)

		err = rows.Scan(
			&toBeAdded.ID,
			&toBeAdded.Nim,
			// &mhsID,
			&toBeAdded.Name,
			&toBeAdded.Semester,
			&toBeAdded.CreatedAt,
			&toBeAdded.UpdatedAt,
		)
		// fmt.Println("err: ", err)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		// fmt.Println("toBeAdded 2: ", toBeAdded)
		
		// Jika tdk error, masukkan hasil ke slice result
		result = append(result, toBeAdded)
		// fmt.Println("result: ", result)
		// fmt.Println()
	}

	return result, nil
}

func (m *mahasiswaRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Mahasiswa, nextCursor string, err error) {
	// Query fetch mahasiswa
	query := `SELECT id, nim, name, semester, created_at, updated_at
				FROM mahasiswa
				WHERE created_at > ?
				ORDER BY created_at LIMIT ?`

	// Decoding cursor
	decodedCursor, err := DecodeCursor(cursor)
	// decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	fmt.Println("num:", num)
	fmt.Println("cursor:", cursor)
	fmt.Println("decodedCursor:", decodedCursor)

	// Panggil fungsi fetch
	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}
	// fmt.Println(res)
	// fmt.Println(num)

	// Jika jumlah row result = query params num
	if len(res) == int(num) {
		nextCursor = EncodeCursor(res[len(res)-1].CreatedAt)
		// nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}
	fmt.Println("nextCursor:", nextCursor)
	fmt.Println()

	return
}

func (m *mahasiswaRepository) GetByID(ctx context.Context, id int64) (res domain.Mahasiswa, err error) {
	query := `SELECT id, nim, name, semester, created_at, updated_at
				FROM mahasiswa 
				WHERE id = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Mahasiswa{}, err
	}
	fmt.Println("GetByID: ", list)

	if len(list) > 0 {
		res = list[0]
		// return res, nil
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *mahasiswaRepository) GetByNIM(ctx context.Context, nim int32) (res domain.Mahasiswa, err error) {
	query := `SELECT id, nim, name, semester, created_at, updated_at
				FROM mahasiswa 
				WHERE nim = ?`

	list, err := m.fetch(ctx, query, nim)
	if err != nil {
		return 
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *mahasiswaRepository) Store(ctx context.Context, dm *domain.Mahasiswa) (err error) {
	// Insert datetime bisa pakai "time.Now()" atau langsung di query pakai "now()"
	// query := `INSERT mahasiswa SET nim=?, name=?, semester=?, created_at=?, updated_at=now()`
	query := `INSERT mahasiswa SET nim=?, name=?, semester=?, created_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return 
	}

	res, err := stmt.ExecContext(ctx, dm.Nim, dm.Name, dm.Semester, time.Now())
	if err != nil {
		return 
	}

	lastID, err := res.LastInsertId()	// ambil id terakhir
	if err != nil {
		return 
	}

	dm.ID = lastID 		// property "ID" pada struct "mahasiswa" akan berisi ID terakhir
	return
}

func (m *mahasiswaRepository) Update(ctx context.Context, dm *domain.Mahasiswa) (err error) {
	query := `UPDATE mahasiswa SET nim=?, name=?, semester=?, updated_at=now() WHERE ID = ?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return 
	}

	// res, err := stmt.ExecContext(ctx, dm.Nim, dm.Name, dm.Semester, dm.UpdatedAt, dm.ID)
	res, err := stmt.ExecContext(ctx, dm.Nim, dm.Name, dm.Semester, dm.ID)
	if err != nil {
		return 
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return 
	}
	fmt.Println("Row affected: ", affect)

	if affect != 1 {
		err = fmt.Errorf("Weird behavior. Total affected: %d", affect)
		return
	}

	return
}

func (m *mahasiswaRepository) Delete(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM mahasiswa WHERE id = ?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return 
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return 
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return 
	}
	fmt.Println("Row affected: ", affect)

	if affect != 1 {
		err = fmt.Errorf("Weird behavior. Total affected: %d", affect)
		return
	}

	return
}