package user

import (
	"tp-db-forum/internal/app/user/models"
)

type Repository interface {
	CreateUser(user models.User) error
	GetUserByNickname(nickname string) (models.User, error)
	EditUser(about, email, fullName string) (models.User, error)
}
