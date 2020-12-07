package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"log"
	"net/http"
	_userHandler "tp-db-forum/internal/app/user/delivery"
	_userRepo "tp-db-forum/internal/app/user/repository"

	"tp-db-forum/configs"

	//userRepository "tp-db-forum/internal/app/user/repository"
)

func main() {
	router := mux.NewRouter()

	connString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable port=%s",
		configs.PostgresConfig.User,
		configs.PostgresConfig.Password,
		configs.PostgresConfig.DB,
		configs.PostgresConfig.Port,
	)

	pgxConnConfig, err := pgx.ParseConnectionString(connString)
	if err != nil {
		log.Fatal(err.Error())
	}
	pgxConnConfig.PreferSimpleProtocol = true

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     pgxConnConfig,
		MaxConnections: 10,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	pool, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}

	userRepo := _userRepo.NewPostgresUserRepository(pool)
	_userHandler.NewUserHandler(router, userRepo)

	log.Fatal(http.ListenAndServe(":5000", router))
}