package repository

import (
	"github.com/jackc/pgx"
	"log"
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
	ct, err := p.Conn.Exec(`INSERT INTO users VALUES ($1, $2, $3, $4)`, user.Nickname, user.FullName, user.About, user.Email)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	log.Print(ct.RowsAffected())

	return nil
}

func (p postgresUserRepository) GetUserByNickname(nickname string) (models.User, error) {
	panic("implement me")
}

func (p postgresUserRepository) EditUser(about, email, fullName string) (models.User, error) {
	panic("implement me")
}
