package repository

import (
	"github.com/jackc/pgx"
	"tp-db-forum/internal/app/user"
	"tp-db-forum/internal/app/user/models"
)

type postgresUserRepository struct {
	Conn *pgx.ConnPool
}

func NewPostgresUserRepository(conn *pgx.ConnPool) user.Repository {
	return &postgresUserRepository{
		Conn: conn,
	}
}

func (p postgresUserRepository) CreateUser(user models.User) error {
	panic("implement me")
}

func (p postgresUserRepository) GetUserByNickname(nickname string) (models.User, error) {
	panic("implement me")
}

func (p postgresUserRepository) EditUser(about, email, fullName string) (models.User, error) {
	panic("implement me")
}
