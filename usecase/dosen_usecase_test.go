package usecase

import (
	"context"
	"time"
	// "fmt"
	"errors"
	"campus/domain"
	"campus/domain/mocks"

	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetch(t *testing.T) {
	mockDosenRepo := new(mocks.DosenRepository)
	mockDosen := domain.Dosen{
		Nip:	22222,
		Name: 	"Billy",
	}

	mockListDosen := make([]domain.Dosen, 0)
	mockListDosen = append(mockListDosen, mockDosen)

	t.Run("success", func(t *testing.T) {
		mockDosenRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(mockListDosen, "next-cursor", nil).Once()

		u := NewDosenUsecase(mockDosenRepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		
		res, err := u.Fetch(context.TODO(), cursor, num)
		cursorExpected := "next-cursor"
		nextCursor := res.Data.(map[string]interface{})["nextCursor"]
		list := res.Data.(map[string]interface{})["data"]

		// t.Log("mockListDosen : ", mockListDosen)
		// t.Log("res : ", res)
		// t.Log("err : ", err)

		assert.Equal(t, cursorExpected, nextCursor)		// cursorExpected & nextCursor harus sama
		assert.NotEmpty(t, nextCursor)					// nextCursor tdk boleh kosong
		assert.NoError(t, err)							// tdk boleh ada error
		assert.Len(t, list, len(mockListDosen))			// isi list harus sejumlah / panjang mockListDosen

		mockDosenRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockDosenRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(nil, "", errors.New("Unexpexted Error")).Once()

		u := NewDosenUsecase(mockDosenRepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		
		res, err := u.Fetch(context.TODO(), cursor, num)
		nextCursor := res.Data.(map[string]interface{})["nextCursor"]
		list := res.Data.(map[string]interface{})["data"]

		// t.Log("res : ", res)
		// t.Log("err : ", err)

		assert.Empty(t, nextCursor)
		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockDosenRepo.AssertExpectations(t)
	})

}

func TestGetByID(t *testing.T) {
	mockDosenRepo := new(mocks.DosenRepository)
	mockDosen := domain.Dosen{
		Nip:	22222,
		Name:	"Billy",
	}

	t.Run("success", func(t *testing.T) {
		mockDosenRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockDosen, nil).Once()
		
		u := NewDosenUsecase(mockDosenRepo, time.Second*2)
		res, err := u.GetByID(context.TODO(), mockDosen.ID)

		assert.NoError(t, err)
		assert.NotNil(t, res.Data)
		mockDosenRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockDosenRepo.On("GetByID", mock.Anything, 
			mock.AnythingOfType("int64")).Return(domain.Dosen{}, errors.New("Unexpexted")).Once()
		
		u := NewDosenUsecase(mockDosenRepo, time.Second*2)
		res, err := u.GetByID(context.TODO(), mockDosen.ID)

		// t.Log("res : ", res)
		// t.Log("err : ", err)

		assert.Error(t, err)
		assert.Equal(t, res.Status, false)
		assert.Nil(t, res.Data)
		mockDosenRepo.AssertExpectations(t)
	})
}

/*func TestStore(t *testing.T) {
	mockDosenRepo := new(mocks.DosenRepository)
	mockDosen := domain.Dosen{
		Nip:	22222,
		Name:	"Billy",
	}

	t.Run("success", func(t *testing.T) {
		tempMockDosen := mockDosen
		tempMockDosen.ID = 0
		mockDosenRepo.On("GetByNip", mock.Anything, 
			mock.AnythingOfType("int32")).Return(domain.Dosen{}, domain.ErrNotFound).Once()
		mockDosenRepo.On("Store", mock.Anything,
			mock.AnythingOfType("*domain.Dosen")).Return(nil).Once()

		u := NewDosenUsecase(mockDosenRepo, time.Second*2)
		res, err := u.Store(context.TODO(), &tempMockDosen)

		t.Log("res : ", res)
		t.Log("err : ", err)

		assert.NoError(t, err)
		assert.Equal(t, mockDosen.Nip, tempMockDosen.Nip)
		mockDosenRepo.AssertExpectations(t)
	})
}*/

func TestUpdate(t *testing.T) {
	mockDosenRepo := new(mocks.DosenRepository)
	mockDosen := domain.Dosen{
		Nip:	22222,
		Name:	"Billy",
	}

	t.Run("success", func(t *testing.T) {
		mockDosenRepo.On("Update", mock.Anything, &mockDosen).Once().Return(nil)

		u := NewDosenUsecase(mockDosenRepo, time.Second*2)
		res, err := u.Update(context.TODO(), &mockDosen)

		// t.Log("res : ", res)
		// t.Log("err : ", err)

		assert.NoError(t, err)
		assert.NotNil(t, res.Data)
		assert.Equal(t, res.Status, true)
		mockDosenRepo.AssertExpectations(t)
	})
}
