package repository

import (
	"github.com/jackc/pgx"
	"log"
	repo "tp-db-forum/internal/app"
	"tp-db-forum/internal/app/models"
)

type postgresAppRepository struct {
	Conn *pgx.ConnPool
}

func NewPostgresAppRepository(conn *pgx.ConnPool) repo.Repository {
	return &postgresAppRepository{
		Conn: conn,
	}
}

func (p postgresAppRepository) CreateUser(user models.User) error {
	_, err := p.Conn.Exec(`INSERT INTO users VALUES ($1, $2, $3, $4)`, user.Nickname, user.FullName, user.About, user.Email)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}

func (p postgresAppRepository) GetUserByNickname(nickname string) (models.User, error) {
	panic("implement me")
}

func (p postgresAppRepository) EditUser(about, email, fullName string) (models.User, error) {
	panic("implement me")
}
