package usecase

import (
	"github.com/google/uuid"
	"log"
	"tp-db-forum/internal/app"
	"tp-db-forum/internal/app/models"
)

type appUseCase struct {
	appRepository app.Repository
}

func NewAppUseCase(ar app.Repository) app.UseCase {
	return &appUseCase{
		appRepository: ar,
	}
}

func (a appUseCase) CreateUser(user models.User) (models.User, error) {
	err := a.appRepository.InsertUser(user)

	return user, err
}

func (a appUseCase) CheckUserByEmail(email string) (models.User, error) {
	user, err := a.appRepository.SelectUserByEmail(email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (a appUseCase) CheckUserByNickname(nickname string) (models.User, error) {
	user, err := a.appRepository.SelectUserByNickname(nickname)

	return user, err
}

func (a appUseCase) HasUser(user models.User) ([]models.User, error) {
	users, err := a.appRepository.SelectUsersByNickAndEmail(user.Nickname, user.Email)

	return users, err
}

func (a appUseCase) EditUser(oldUser, newUser models.User) (models.User, error) {
	if newUser.Email == "" {
		newUser.Email = oldUser.Email
	}

	if newUser.About == "" {
		newUser.About = oldUser.About
	}

	if newUser.FullName == "" {
		newUser.FullName = oldUser.FullName
	}

	u, err := a.appRepository.UpdateUser(newUser)

	return u, err
}

func (a appUseCase) CreateForum(forum models.Forum) (models.Forum, error) {
	f, err := a.appRepository.InsertForum(forum)
	if err != nil {
		log.Println(err)
	}

	return f, err
}

func (a appUseCase) CheckForumBySlug(slug string) (models.Forum, error) {
	forum, err := a.appRepository.SelectForumBySlug(slug)
	if err != nil {
		log.Println(err)
	}

	return forum, err
}

func (a appUseCase) CreateForumThread(thread models.Thread) (models.Thread, error) {
	if thread.Slug == "" {
		u, err := uuid.NewRandom()
		if err != nil {
			log.Println(err)
			panic("AAAAAAAAAAAAAAAAA")
		}
		thread.Slug = u.String()
	}
	thr, err := a.appRepository.InsertThread(thread)

	return thr, err
}

func (a appUseCase) CheckThreadBySlug(slug string) (models.Thread, error) {
	thread, err := a.appRepository.SelectThreadBySlug(slug)

	return thread, err
}
