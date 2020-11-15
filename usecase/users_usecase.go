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

func GenerateToken(du *domain.Users) (t string, err error) {
	// Generate token
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["username"] = du.Username
	claims["level"] = "application"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err = token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println("Error generate token")
		return t, err
	}

	fmt.Println("Token: ", t)
	return t, err
}

func (u *usersUsecase) CheckLogin(c context.Context, du *domain.Users) (res domain.Response, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	fmt.Println("Usecase domain users: ", du)

	result, err := u.usersRepo.CheckLogin(ctx, du)
	fmt.Println("Usecase CheckLogin result: ", result)
	fmt.Println("Usecase CheckLogin err: ", err)
	if err != nil {
		res.Status = false
		res.Message = "Username not found"
		return 
	}

	match, err := repository.CheckPasswordHash(du.Password, result.Password)
	fmt.Println("Usecase CheckPasswordHash match: ", match)
	fmt.Println("Usecase CheckPasswordHash err: ", err)
	if !match {
		fmt.Println("Hash and password doesn't match")

		res.Status = false
		res.Message = "Hash and password doesn't match"
		// return res, err
		return
	}

	token, err := GenerateToken(du)
	if err != nil {
		// res.Status = false
		// res.Message = err
		// res.Data = map[string]string{
		// 	"token": token,
		// }
		return res, err
	}

	res.Status = true
	res.Message = "Generate token success"
	res.Data = map[string]string{
		"token": token,
	}

	return res, nil
}

