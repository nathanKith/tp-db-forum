package app

import (
	"tp-db-forum/internal/app/models"
)

type Repository interface {
	InsertUser(user models.User) error
	SelectUserByNickname(nickname string) (models.User, error)
	SelectUserByEmail(email string) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	SelectUsersByNickAndEmail(nickname, email string) ([]models.User, error)

	InsertForum(forum models.Forum) (models.Forum, error)
	SelectForumBySlug(slug string) (models.Forum, error)
	InsertThread(thread models.Thread) (models.Thread, error)
	SelectThreadBySlug(slug string) (models.Thread, error)
	SelectThreadById(id int) (models.Thread, error)
	InsertPosts(posts []models.Post) ([]models.Post, error)
	UpdateThread(thread models.Thread) (models.Thread, error)
	InsertVote(vote models.Vote) (models.Vote, error)
	UpdateVote(vote models.Vote) (models.Vote, error)
	GetServiceStatus() (map[string]int, error)
	ClearDatabase() error
}

type UseCase interface {
	CreateUser(user models.User) (models.User, error)
	CheckUserByEmail(email string) (models.User, error)
	CheckUserByNickname(nickname string) (models.User, error)
	HasUser(user models.User) ([]models.User, error)
	EditUser(oldUser, newUser models.User) (models.User, error)

	CreateForum(forum models.Forum) (models.Forum, error)
	CheckForumBySlug(slug string) (models.Forum, error)
	CreateForumThread(thread models.Thread) (models.Thread, error)
	CheckThreadBySlug(slug string) (models.Thread, error)
	CheckThreadById(id int) (models.Thread, error)
	CreatePosts(posts []models.Post, thread models.Thread) ([]models.Post, error)
	EditThread(thread models.Thread) (models.Thread, error)
	AddVote(vote models.Vote) (models.Vote, error)
	UpdateVote(vote models.Vote) (models.Vote, error)
	GetServiceStatus() (map[string]int, error)
	ClearDatabase() error
}
