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

func (k *kelasRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Kelas, err error) {
	rows, err := k.Conn.QueryContext(ctx, query, args...)
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
		dosenNip := int32(0)
		dosenName := string(0)
		mahasiswaID := int64(0)
		mahasiswaNim := int32(0)
		mahasiswaName := string(0)
		var mahasiswaSemester sql.NullInt32

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
			&dosenNip,
			&dosenName,
			&mahasiswaNim,
			&mahasiswaName,
			&mahasiswaSemester,
		)
		// fmt.Println("err: ", err)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		toBeAdded.Ruang = domain.Ruang{
			ID: ruangID,
			Name: ruangName,
		}
		toBeAdded.MataKuliah = domain.MataKuliah{
			ID: mataKuliahID,
			Name: matakuliahName,
		}
		toBeAdded.Dosen = domain.Dosen{
			ID: dosenID,
			Nip: dosenNip,
			Name: dosenName,
		}
		toBeAdded.Mahasiswa = domain.Mahasiswa{
			ID: mahasiswaID,
			Nim: mahasiswaNim,
			Name: mahasiswaName,
			Semester: mahasiswaSemester,
		}
		// fmt.Println("toBeAdded 2: ", toBeAdded)
		
		// Jika tdk error, masukkan hasil ke slice result
		result = append(result, toBeAdded)
		// fmt.Println("result: ", result)
		// fmt.Println()
	}

	return result, nil
}

func (k *kelasRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Kelas, nextCursor string, err error) {
	// Query fetch Kelas

	// query := `SELECT id_kelas, name, ruang_id, mata_kuliah_id, dosen_id, mahasiswa_id, 
	// 			created_by, created_at, updated_by, updated_at
	// 			FROM kelas
	// 			WHERE created_at > ?
	// 			ORDER BY created_at LIMIT ?`

	query := `SELECT id_kelas, kelas.name, ruang_id, mata_kuliah_id, dosen_id, mahasiswa_id, 
				created_by, kelas.created_at, updated_by, kelas.updated_at,
				ruang.name AS ruang_name, mata_kuliah.name AS mata_kuliah_name,
				dosen.nip AS dosen_nip, dosen.name AS dosen_name, 
				mahasiswa.nim AS mahasiswa_nim, mahasiswa.name AS mahasiswa_name, 
				mahasiswa.semester AS mahasiswa_semester
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
	res, err = k.fetch(ctx, query, decodedCursor, num)
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

func (k *kelasRepository) CheckIfExists(ctx context.Context, dk *domain.Kelas) (res domain.Kelas, err error) {
	query := `SELECT id_kelas, name, ruang_id, mata_kuliah_id, dosen_id, mahasiswa_id
				FROM kelas 
				WHERE ruang_id = ?
				AND mata_kuliah_id = ?
				AND dosen_id = ?
				AND mahasiswa_id = ?`

	list, err := k.fetch(ctx, query, dk.Ruang.ID, dk.MataKuliah.ID, dk.Dosen.ID, dk.Mahasiswa.ID)
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

func (k *kelasRepository) Store(ctx context.Context, dk *domain.Kelas) (err error) {
	query := `INSERT kelas SET name=?, ruang_id=?, mata_kuliah_id=?, dosen_id=?, mahasiswa_id=?, 
		created_by=?, created_at=now()`
	stmt, err := k.Conn.PrepareContext(ctx, query)
	if err != nil {
		return 
	}

	res, err := stmt.ExecContext(ctx, dk.Name, dk.Ruang.ID, dk.MataKuliah.ID, dk.Dosen.ID, dk.Mahasiswa.ID, dk.CreatedBy)
	if err != nil {
		return 
	}

	lastID, err := res.LastInsertId()	// ambil id terakhir
	if err != nil {
		return 
	}

	dk.ID = lastID 		// property "ID" pada struct "kelas" akan berisi ID terakhir
	return
}
