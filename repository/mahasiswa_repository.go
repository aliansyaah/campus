package repository

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
	// "campus/repository"
	"campus/domain"
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

	result = make([]domain.Mahasiswa, 0)
	for rows.Next() {
		toBeAdded := domain.Mahasiswa{}
		// mhsID := int64(0)
		err = rows.Scan(
			&toBeAdded.ID,
			&toBeAdded.Nim,
			// &mhsID,
			&toBeAdded.Name,
			&toBeAdded.Semester,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		

		result = append(result, toBeAdded)
	}

	return result, nil
}

func (m *mahasiswaRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Mahasiswa, nextCursor string, err error) {
	query := `SELECT id, nim, name, semester 
				FROM mahasiswa
				WHERE created_at > ?
				ORDER BY created_at LIMIT ?`

	// decodedCursor, err := repository.DecodeCursor(cursor)
	decodedCursor, err := DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		// nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
		nextCursor = EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}

func (m *mahasiswaRepository) GetByID(ctx context.Context, id int64) (res domain.Mahasiswa, err error) {
	query := `SELECT id, nim, name, semester 
				FROM mahasiswa 
				WHERE id = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Mahasiswa{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}