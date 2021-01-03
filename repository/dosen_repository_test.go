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
	assert.NotEmpty(t, nextCursor)	// nextCursor tidak boleh kosong
	assert.NoError(t, err)			// tidak boleh ada error
	assert.Len(t, list, 2)			// isi list harus ada 2, krn data yg di-mock di atas ada 2
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id_dosen", "nip", "name", "created_at", "updated_at"}).
		AddRow(1, 22222, "Nama Dosen 1", time.Now(), time.Now())

	query := "SELECT id_dosen, nip, name, created_at, updated_at FROM dosen WHERE id_dosen = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	d := NewDosenRepository(db)

	num := int64(5)
	result, err := d.GetByID(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, result)	// result tidak boleh nil
}

func TestGetByNIP(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id_dosen", "nip", "name", "created_at", "updated_at"}).
		AddRow(1, 22222, "Nama Dosen 1", time.Now(), time.Now())

	query := "SELECT id_dosen, nip, name, created_at, updated_at FROM dosen WHERE nip = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	d := NewDosenRepository(db)

	num := int32(22222)
	result, err := d.GetByNIP(context.TODO(), num)
	// t.Log("err : ", err)
	// t.Log("result : ", result)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestStore(t *testing.T) {
	ar := &domain.Dosen{
		Nip: 22222,
		Name: "Nama Dosen 2",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "INSERT dosen SET nip=\\?, name=\\?, created_at=now()"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Nip, ar.Name).WillReturnResult(sqlmock.NewResult(12, 1))

	d := NewDosenRepository(db)

	err = d.Store(context.TODO(), ar)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), ar.ID)
}

func TestUpdate(t *testing.T) {
	now := mysql.NullTime{Time: time.Now(), Valid: true}
	ar := &domain.Dosen{
		ID: 2,
		Nip: 33333,
		Name: "Nama Dosen 3",
		UpdatedAt: now, 
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "UPDATE dosen SET nip=\\?, name=\\?, updated_at=now() WHERE id_dosen = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Nip, ar.Name, ar.ID).WillReturnResult(sqlmock.NewResult(12, 1))

	d := NewDosenRepository(db)

	err = d.Update(context.TODO(), ar)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "DELETE FROM dosen WHERE id_dosen = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(12).WillReturnResult(sqlmock.NewResult(12, 1))

	d := NewDosenRepository(db)

	num := int64(12)
	err = d.Delete(context.TODO(), num)
	assert.NoError(t, err)
}
