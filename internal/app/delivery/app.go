package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"log"
	"net/http"
	"strings"
	"tp-db-forum/internal/app"
	"tp-db-forum/internal/app/models"
)

type AppHandler struct {
	appUseCase app.UseCase
}

func NewAppHandler(router *mux.Router, appUseCase app.UseCase) {
	handler := &AppHandler{
		appUseCase: appUseCase,
	}

	router.HandleFunc("/api/user/{nickname}/create", handler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/api/user/{nickname}/profile", handler.UserProfile).Methods(http.MethodGet, http.MethodPost)

	router.HandleFunc("/api/forum/create", handler.CreateForum).Methods(http.MethodPost)
	router.HandleFunc("/api/forum/{slug}/details", handler.ForumDetails).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/create", handler.CreateThread).Methods(http.MethodPost)
	//router.HandleFunc("/api/forum/{slug}/threads", handler.ForumThreads).Methods(http.MethodGet)
	//router.HandleFunc("/api/forum/{slug}/users", handler.ForumUsers).Methods(http.MethodGet)
	//
	//router.HandleFunc("/api/thread/{slug_or_id}/create", handler.CreatePosts).Methods(http.MethodPost)
	//router.HandleFunc("/api/thread/{slug_or_id}/vote", handler.VoteThread).Methods(http.MethodPost)
	//router.HandleFunc("/api/thread/{slug_or_id}/details", handler.ThreadDetails).Methods(http.MethodGet, http.MethodPost)
	//router.HandleFunc("/api/thread/{slug_or_id}/details", handler.ThreadDetails).Methods(http.MethodGet, http.MethodPost)
	//router.HandleFunc("/api/thread/{slug_or_id}/posts", handler.ThreadPosts).Methods(http.MethodGet)
	//
	//router.HandleFunc("/api/post/{id}/details", handler.PostDetails).Methods(http.MethodGet, http.MethodPost)
	//router.HandleFunc("/api/service/status", handler.StatusHandler).Methods(http.MethodGet)
	//router.HandleFunc("/api/service/clear", handler.ClearHandler).Methods(http.MethodPost)
}

func errorMarshal(message string) ([]byte, error) {
	e := models.Error{Message: message}
	body, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (h AppHandler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	nickname := strings.TrimSuffix(strings.TrimPrefix(request.URL.Path, "/api/user/"), "/create")

	var user models.User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		return
	}
	user.Nickname = nickname

	_, err = h.appUseCase.CreateUser(user)
	if err != nil {
		if users, err := h.appUseCase.HasUser(user); err == nil {
			body, err := json.Marshal(users)
			if err != nil {
				log.Println(err)
				return
			}

			writer.WriteHeader(http.StatusConflict)
			writer.Write(body)

			return
		}

		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Write(body)
}

func (h AppHandler) UserProfile(writer http.ResponseWriter, request *http.Request) {
	nickname := strings.TrimSuffix(strings.TrimPrefix(request.URL.Path, "/api/user/"), "/profile")

	if request.Method == "GET" {
		user, err := h.appUseCase.CheckUserByNickname(nickname)
		if err != nil {
			body, err := errorMarshal("Can't find user\n")
			if err != nil {
				log.Println(err)
				return
			}

			writer.WriteHeader(http.StatusNotFound)
			writer.Write(body)

			return
		}
		body, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
			return
		}

		writer.WriteHeader(http.StatusOK)
		writer.Write(body)

		return
	}

	var user models.User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		return
	}
	user.Nickname = nickname

	oldUser, err := h.appUseCase.CheckUserByNickname(user.Nickname)
	if err != nil {
		body, err := errorMarshal("Can't find user\n")
		if err != nil {
			log.Println(err)
			return
		}

		writer.WriteHeader(http.StatusNotFound)
		writer.Write(body)

		return
	}

	result, err := h.appUseCase.EditUser(oldUser, user)
	// check conflict with constraint unique on email
	if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
		body, err := errorMarshal("Conflict email\n")
		if err != nil {
			log.Println(err)
			return
		}

		writer.WriteHeader(http.StatusConflict)
		writer.Write(body)

		return

	}

	body, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(body)
}

func (h AppHandler) CreateForum(writer http.ResponseWriter, request *http.Request) {
	var forum models.Forum
	err := json.NewDecoder(request.Body).Decode(&forum)
	if err != nil {
		log.Println(err)
		return
	}

	f, err := h.appUseCase.CreateForum(forum)
	if pgErr, ok := err.(pgx.PgError); ok {
		switch pgErr.Code {
		case "23505":
			forumSlug, err := h.appUseCase.CheckForumBySlug(forum.Slug)
			if err != nil {
				log.Println(err)
				return
			}

			body, err := json.Marshal(forumSlug)
			if err != nil {
				log.Println(err)
				return
			}

			writer.WriteHeader(http.StatusConflict)
			writer.Write(body)
		case "23503":
			body, err := errorMarshal("Can't find user")
			if err != nil {
				log.Println(err)
				return
			}

			writer.WriteHeader(http.StatusNotFound)
			writer.Write(body)
		}

		return
	}

	body, err := json.Marshal(f)
	if err != nil {
		log.Println(err)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Write(body)
}

func (h AppHandler) ForumDetails(writer http.ResponseWriter, request *http.Request) {
	slug := strings.TrimSuffix(strings.TrimPrefix(request.URL.Path, "/api/forum/"), "/details")

	forum, err := h.appUseCase.CheckForumBySlug(slug)
	if err != nil {
		body, err := errorMarshal("Can't find forum")
		if err != nil {
			log.Println(err)
			return
		}

		writer.WriteHeader(http.StatusNotFound)
		writer.Write(body)

		return
	}

	body, err := json.Marshal(forum)
	if err != nil {
		log.Println(err)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(body)
}

func (h AppHandler) CreateThread(writer http.ResponseWriter, request *http.Request) {
	slug := strings.TrimSuffix(strings.TrimPrefix(request.URL.Path, "/api/forum/"), "/create")

	thread := models.Thread{Forum: slug}
	err := json.NewDecoder(request.Body).Decode(&thread)
	if err != nil {
		log.Println(err)
		return
	}

	newThread, err := h.appUseCase.CreateForumThread(thread)
	if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
		oldThread, err := h.appUseCase.CheckThreadBySlug(slug)
		if err != nil {
			log.Println(err)
			return
		}

		body, err := json.Marshal(oldThread)
		if err != nil {
			log.Println(err)
			return
		}

		writer.WriteHeader(http.StatusConflict)
		writer.Write(body)

		return
	}

	if err != nil {
		body, err := errorMarshal("Can't find forum or user\n")
		if err != nil {
			log.Println(err)
			return
		}

		writer.WriteHeader(http.StatusNotFound)
		writer.Write(body)

		return
	}

	body, err := json.Marshal(newThread)
	if err != nil {
		log.Println(err)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Write(body)
}
