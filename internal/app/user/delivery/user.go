package delivery

import (
	"github.com/gorilla/mux"
	"net/http"
	"tp-db-forum/internal/app/user"
)

type UserHandler struct {
	userRepository user.Repository
}

func NewUserHandler(router *mux.Router, userRepository user.Repository) {
	handler := &UserHandler{
		userRepository: userRepository,
	}

	router.HandleFunc("/user/{nickname}/create", handler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{nickname}/profile", handler.GetUserByNickname).Methods(http.MethodGet)
	router.HandleFunc("/user/{nickname}/profile", handler.EditUser).Methods(http.MethodPost)

}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (uh *UserHandler) GetUserByNickname(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (uh *UserHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
