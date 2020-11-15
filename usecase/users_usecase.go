package usecase

import (
	"context"
	"time"
	"campus/domain"
	"campus/repository"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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

	fmt.Println("Usecase domain users: ", du)

	res, err = u.usersRepo.CheckLogin(ctx, du)
	fmt.Println("Usecase CheckLogin res: ", res)
	fmt.Println("Usecase CheckLogin err: ", err)
	if err != nil {
		return 
	}

	match, err := repository.CheckPasswordHash(du.Password, res.Password)
	fmt.Println("Usecase CheckPasswordHash match: ", match)
	fmt.Println("Usecase CheckPasswordHash err: ", err)
	if !match {
		fmt.Println("Hash and password doesn't match")
		return res, err
	}

	// Generate token
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["username"] = du.Username
	claims["level"] = "application"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return res, err
	}

	fmt.Println("Token: ", t)

	return
}