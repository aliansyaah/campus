package usecase

import (
	"context"
	"time"
	"campus/domain"
	"campus/repository"
	"fmt"
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

func (u *usersUsecase) CheckLogin(c context.Context, du *domain.Users) (res domain.Users, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.usersRepo.CheckLogin(ctx, du)
	fmt.Println("Usecase domain users: ", du)
	fmt.Println("Usecase res: ", res)
	fmt.Println("Usecase err: ", err)
	if err != nil {
		return 
	}

	match, err := repository.CheckPasswordHash(du.Password, res.Password)
	if !match {
		fmt.Println("Hash and password doesn't match")
		return res, err
	}

	return
}