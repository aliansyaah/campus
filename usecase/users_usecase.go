package usecase

import (
	"context"
	"time"
	"campus/domain"
)

type usersUsecase struct {
	usersRepo domain.UsersRepository
	contextTimeout time.Duration
}

func NewUsersUsecase(u domain.UsersRepository, timeout time.Duration) domain.UsersUsecase {
	return &usersUsecase{
		usersRepo: u,
		contextTimeout: timeout,
	}
}

func (u *usersUsecase) CheckLogin(c context.Context, du *domain.Users) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	err = u.usersRepo.CheckLogin(ctx, du)
	return
}