package repository

import (
	"context"
	// "database/sql"
	"github.com/go-sql-driver/mysql"
	"time"
	"campus/domain"

	"testing"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockDosen := []domain.Dosen{
		domain.Dosen{
			ID: 1, Nip: 22222, Name: "Nama Dosen 1",
			CreatedAt: time.Now(), UpdatedAt: mysql.NullTime{Time: time.Now(), Valid: true},
		},
		domain.Dosen{
			ID: 2, Nip: 66666, Name: "Nama Dosen 2",
			CreatedAt: time.Now(), UpdatedAt: mysql.NullTime{Time: time.Now(), Valid: true},
		},
	}

	rows := sqlmock.NewRows([]string{"id_dosen", "nip", "name", "created_at", "updated_at"}).
		AddRow(mockDosen[0].ID, mockDosen[0].Nip, mockDosen[0].Name,
			mockDosen[0].CreatedAt, time.Now()).
		AddRow(mockDosen[1].ID, mockDosen[1].Nip, mockDosen[1].Name,
			mockDosen[1].CreatedAt, time.Now())

	query := "SELECT id_dosen, nip, name, created_at, updated_at FROM dosen WHERE created_at > \\? ORDER BY created_at LIMIT \\?"

	// t.Log("mockDosen : ", mockDosen)
	// t.Log("rows : ", rows)
	// t.Log("query : ", query)

	mock.ExpectQuery(query).WillReturnRows(rows)
	d := NewDosenRepository(db)
	cursor := EncodeCursor(mockDosen[1].CreatedAt)
	num := int64(2)
	list, nextCursor, err := d.Fetch(context.TODO(), cursor, num)
	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}