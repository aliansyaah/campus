package repository

import (
	"context"
	"database/sql"
	// "time"
	"fmt"
	"campus/domain"
	"github.com/sirupsen/logrus"
)

type dosenRepository struct {
	Conn *sql.DB
}

func NewDosenRepository(Conn *sql.DB) domain.DosenRepository {
	return &dosenRepository{Conn}
}

func (q *dosenRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Dosen, err error) {
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
	result = make([]domain.Dosen, 0)	// tipe = slice domain.Dosen

	// fmt.Println("Looping:")

	// Looping rows hasil query
	for rows.Next() {
		toBeAdded := domain.Dosen{}		// struct dosen pada layer domain/model
		// fmt.Println("toBeAdded: ", toBeAdded)
		// mhsID := int64(0)

		err = rows.Scan(
			&toBeAdded.ID,
			&toBeAdded.Nip,
			&toBeAdded.Name,
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

func (m *dosenRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Dosen, nextCursor string, err error) {
	// Query fetch Dosen
	query := `SELECT id_dosen, nip, name, created_at, updated_at
				FROM dosen
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

func (d *dosenRepository) GetByID(ctx context.Context, id int64) (res domain.Dosen, err error) {
	query := `SELECT id_dosen, nip, name, created_at, updated_at
				FROM dosen 
				WHERE id_dosen = ?`

	list, err := d.fetch(ctx, query, id)
	if err != nil {
		return domain.Dosen{}, err
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

func (d *dosenRepository) GetByNIP(ctx context.Context, nip int32) (res domain.Dosen, err error) {
	query := `SELECT id_dosen, nip, name, created_at, updated_at
				FROM dosen 
				WHERE nip = ?`

	list, err := d.fetch(ctx, query, nip)
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

func (d *dosenRepository) Store(ctx context.Context, dd *domain.Dosen) (err error) {
	// Insert datetime bisa pakai "time.Now()" atau langsung di query pakai "now()"
	// query := `INSERT dosen SET nip=?, name=?, created_at=?`
	query := `INSERT dosen SET nip=?, name=?, created_at=now()`
	stmt, err := d.Conn.PrepareContext(ctx, query)
	if err != nil {
		return 
	}

	// res, err := stmt.ExecContext(ctx, dd.Nip, dd.Name, time.Now())
	res, err := stmt.ExecContext(ctx, dd.Nip, dd.Name)
	if err != nil {
		return 
	}

	lastID, err := res.LastInsertId()	// ambil id terakhir
	if err != nil {
		return 
	}

	dd.ID = lastID 		// property "ID" pada struct "dosen" akan berisi ID terakhir
	return
}

// func (m *mahasiswaRepository) Update(ctx context.Context, dm *domain.Mahasiswa) (err error) {
// 	query := `UPDATE mahasiswa SET nim=?, name=?, semester=?, updated_at=now() WHERE ID = ?`
// 	stmt, err := m.Conn.PrepareContext(ctx, query)
// 	if err != nil {
// 		return 
// 	}

// 	// res, err := stmt.ExecContext(ctx, dm.Nim, dm.Name, dm.Semester, dm.UpdatedAt, dm.ID)
// 	res, err := stmt.ExecContext(ctx, dm.Nim, dm.Name, dm.Semester, dm.ID)
// 	if err != nil {
// 		return 
// 	}

// 	affect, err := res.RowsAffected()
// 	if err != nil {
// 		return 
// 	}
// 	fmt.Println("Row affected: ", affect)

// 	if affect != 1 {
// 		err = fmt.Errorf("Weird behavior. Total affected: %d", affect)
// 		return
// 	}

// 	return
// }

// func (m *mahasiswaRepository) Delete(ctx context.Context, id int64) (err error) {
// 	query := `DELETE FROM mahasiswa WHERE id = ?`
// 	stmt, err := m.Conn.PrepareContext(ctx, query)
// 	if err != nil {
// 		return 
// 	}

// 	res, err := stmt.ExecContext(ctx, id)
// 	if err != nil {
// 		return 
// 	}

// 	affect, err := res.RowsAffected()
// 	if err != nil {
// 		return 
// 	}
// 	fmt.Println("Row affected: ", affect)

// 	if affect != 1 {
// 		err = fmt.Errorf("Weird behavior. Total affected: %d", affect)
// 		return
// 	}

// 	return
// }