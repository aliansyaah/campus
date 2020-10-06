package repository

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
)

type mahasiswaRepository struct {
	Conn *sql.DB
}

func NewMahasiswaRepository(Conn *sql.DB) domain.MahasiswaRepository {
	return &questionRepository{Conn}
}

func (q *questionRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Question, err error) {
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

	result = make([]domain.Question, 0)
	for rows.Next() {
		toBeAdded := domain.Question{}
		err = rows.Scan(
			&toBeAdded.ID,
			&toBeAdded.Nim,
			&toBeAdded.Name,
			&toBeAdded.Semester,
		)

		if err != nil {
			logrus.Error(err)
		}
		result = append(result, toBeAdded)
	}

	return result, nil
}
