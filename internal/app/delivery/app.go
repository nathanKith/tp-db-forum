package delivery

import (
	"github.com/gorilla/mux"
	"net/http"
	"tp-db-forum/internal/app"
)

type AppHandler struct {
	appRepository app.Repository
}

func NewAppHandler(router *mux.Router, appRepository app.Repository) {
	handler := &AppHandler{
		appRepository: appRepository,
	}

	router.HandleFunc("/api/user/{nickname}/create", handler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/api/user/{nickname}/profile", handler.UserProfile).Methods(http.MethodGet)
	router.HandleFunc("/api/user/{nickname}/profile", handler.EditUser).Methods(http.MethodPost)

	router.HandleFunc("/api/forum/create", handler.CreateForum).Methods(http.MethodPost)
	router.HandleFunc("/api/forum/{slug}/details", handler.ForumDetails).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/create", handler.CreateForumSlug).Methods(http.MethodPost)
	router.HandleFunc("/api/forum/{slug}/threads", handler.ForumThreads).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/users", handler.ForumUsers).Methods(http.MethodGet)

	router.HandleFunc("/api/thread/{slug_or_id}/create", handler.CreatePosts).Methods(http.MethodPost)
	router.HandleFunc("/api/thread/{slug_or_id}/vote", handler.VoteThread).Methods(http.MethodPost)
	router.HandleFunc("/api/thread/{slug_or_id}/details", handler.ThreadDetails).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/api/thread/{slug_or_id}/details", handler.ThreadDetails).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/api/thread/{slug_or_id}/posts", handler.ThreadPosts).Methods(http.MethodGet)

	router.HandleFunc("/api/post/{id}/details", handler.PostDetails).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/api/service/status", handler.StatusHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/service/clear", handler.ClearHandler).Methods(http.MethodPost)
}

func (h AppHandler) CreateUser(writer http.ResponseWriter, request *http.Request) {

}

func (h AppHandler) UserProfile(writer http.ResponseWriter, request *http.Request) {

}

func (h AppHandler) EditUser(writer http.ResponseWriter, request *http.Request) {

}