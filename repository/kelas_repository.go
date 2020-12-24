package repository

import (
	"context"
	"database/sql"
	// "time"
	"fmt"
	"campus/domain"
	"github.com/sirupsen/logrus"
)

type kelasRepository struct {
	Conn *sql.DB
}

func NewKelasRepository(Conn *sql.DB) domain.KelasRepository {
	return &kelasRepository{Conn}
}

func (q *kelasRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Kelas, err error) {
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
	result = make([]domain.Kelas, 0)	// tipe = slice domain.Kelas

	// fmt.Println("Looping:")

	// Looping rows hasil query
	for rows.Next() {
		toBeAdded := domain.Kelas{}		// struct kelas pada layer domain/model
		// fmt.Println("toBeAdded: ", toBeAdded)
		ruangID := int64(0)
		ruangName := string(0)
		mataKuliahID := int64(0)
		matakuliahName := string(0)
		dosenID := int64(0)
		dosenName := string(0)
		mahasiswaID := int64(0)
		mahasiswaName := string(0)

		err = rows.Scan(
			&toBeAdded.ID,
			&toBeAdded.Name,
			&ruangID,
			&mataKuliahID,
			&dosenID,
			&mahasiswaID,
			&toBeAdded.CreatedBy,
			&toBeAdded.CreatedAt,
			&toBeAdded.UpdatedBy,
			&toBeAdded.UpdatedAt,
			&ruangName,
			&matakuliahName,
			&dosenName,
			&mahasiswaName,
		)
		// fmt.Println("err: ", err)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		toBeAdded.RuangID = domain.Ruang{
			ID: ruangID,
			Name: ruangName,
		}
		toBeAdded.MataKuliahID = domain.MataKuliah{
			ID: mataKuliahID,
			Name: matakuliahName,
		}
		toBeAdded.DosenID = domain.Dosen{
			ID: dosenID,
			Name: dosenName,
		}
		toBeAdded.MahasiswaID = domain.Mahasiswa{
			ID: mahasiswaID,
			Name: mahasiswaName,
		}
		// fmt.Println("toBeAdded 2: ", toBeAdded)
		
		// Jika tdk error, masukkan hasil ke slice result
		result = append(result, toBeAdded)
		// fmt.Println("result: ", result)
		// fmt.Println()
	}

	return result, nil
}

func (m *kelasRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Kelas, nextCursor string, err error) {
	// Query fetch Kelas

	// query := `SELECT id_kelas, name, ruang_id, mata_kuliah_id, dosen_id, mahasiswa_id, 
	// 			created_by, created_at, updated_by, updated_at
	// 			FROM kelas
	// 			WHERE created_at > ?
	// 			ORDER BY created_at LIMIT ?`

	query := `SELECT id_kelas, kelas.name, ruang_id, mata_kuliah_id, dosen_id, mahasiswa_id, 
				created_by, kelas.created_at, updated_by, kelas.updated_at,
				ruang.name AS ruang_name, mata_kuliah.name AS mata_kuliah_name,
				dosen.name AS dosen_name, mahasiswa.name AS mahasiswa_name
				FROM kelas
				LEFT JOIN ruang ON ruang.id_ruang = kelas.ruang_id
				LEFT JOIN mata_kuliah ON mata_kuliah.id_mata_kuliah = kelas.mata_kuliah_id
				LEFT JOIN dosen ON dosen.id_dosen = kelas.dosen_id
				LEFT JOIN mahasiswa ON mahasiswa.id = kelas.mahasiswa_id
				WHERE kelas.created_at > ?
				ORDER BY kelas.created_at LIMIT ?`

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
