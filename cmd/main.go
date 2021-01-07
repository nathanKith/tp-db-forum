package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"log"
	"net/http"
	"tp-db-forum/configs"
	_handler "tp-db-forum/internal/app/delivery"
	_repo "tp-db-forum/internal/app/repository"
	_useCase "tp-db-forum/internal/app/usecase"
)

func applicationJSONMiddleware(_ *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}

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

	repo := _repo.NewPostgresAppRepository(pool)
	useCase := _useCase.NewAppUseCase(repo)
	_handler.NewAppHandler(router, useCase)

	router.Use(applicationJSONMiddleware(router))

	log.Fatal(http.ListenAndServe(":5000", router))
}
