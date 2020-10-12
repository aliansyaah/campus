package main

import (
	"database/sql"
	"fmt"
	"log"
	// "net/http"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	// "github.com/julienschmidt/httprouter"
	"github.com/labstack/echo"
	"github.com/spf13/viper"

	mhsRepo "campus/repository"
	mhsUsecase "campus/usecase"
	mhsDeliv "campus/delivery"
	mhsDelivMiddleware "campus/delivery/middleware"
)

func init() {
	viper.SetConfigFile("../config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	// fmt.Println("hello world")
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")

	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	// dbConn, err := sql.Open("mysql", dsn)
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatal(err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := mhsDelivMiddleware.InitMiddleware()
	e.Use(middL.CORS)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	// router := httprouter.New()

	/* Layer Repository */
	mr := mhsRepo.NewMahasiswaRepository(dbConn)

	/* Layer Usecase */
	mu := mhsUsecase.NewMahasiswaUsecase(mr, timeoutContext)

	/* Layer Delivery */
	mhsDeliv.NewMahasiswaHandler(e, mu)
	// mhsDeliv.NewMahasiswaHandler(router, mu)

	log.Fatal(e.Start(viper.GetString("server.address")))

	// router.ServeFiles("/static/*filepath", http.Dir("assets"))
	// log.Fatal(http.ListenAndServe(viper.GetString("server.address"), router))
}
