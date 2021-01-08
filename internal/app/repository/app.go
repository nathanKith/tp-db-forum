package repository

import (
	"errors"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	"log"
	"time"
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

func (p *postgresAppRepository) InsertUser(user models.User) error {
	_, err := p.Conn.Exec(`INSERT INTO users(nickname, fullname, about, email) VALUES ($1, $2, $3, $4)`, user.Nickname, user.FullName, user.About, user.Email)

	return err
}

func (p *postgresAppRepository) SelectUserByNickname(nickname string) (models.User, error) {
	row := p.Conn.QueryRow(`SELECT fullname, about, email FROM users WHERE nickname=$1 LIMIT 1;`, nickname)

	var fullName, about, email string
	err := row.Scan(&fullName, &about, &email)
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}

	return models.User{
		About:    about,
		Email:    email,
		FullName: fullName,
		Nickname: nickname,
	}, nil
}

func (p *postgresAppRepository) SelectUserByEmail(email string) (models.User, error) {
	row := p.Conn.QueryRow(`SELECT nickname, fullname, about FROM users WHERE email=$1 LIMIT 1;`, email)

	var nickname, fullName, about string
	err := row.Scan(&nickname, &fullName, &about)
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		About:    about,
		Email:    email,
		FullName: fullName,
		Nickname: nickname,
	}, nil
}

func (p *postgresAppRepository) SelectUsersByNickAndEmail(nickname, email string) ([]models.User, error) {
	rows, err := p.Conn.Query(`SELECT * FROM users WHERE email=$1 OR nickname=$2 LIMIT 2;`, email, nickname)
	if err != nil {
		return nil, err
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.Nickname, &user.FullName, &user.About, &user.Email)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (p *postgresAppRepository) UpdateUser(user models.User) (models.User, error) {
	var newUser models.User
	err := p.Conn.QueryRow(
		`UPDATE users SET email=$1, about=$2, fullname=$3 WHERE nickname=$4 RETURNING *`,
		user.Email,
		user.About,
		user.FullName,
		user.Nickname,
	).Scan(&newUser.Nickname, &newUser.FullName, &newUser.About, &newUser.Email)

	return newUser, err
}

func (p *postgresAppRepository) InsertForum(forum models.Forum) (models.Forum, error) {
	_, err := p.Conn.Exec(
		`INSERT INTO forum(slug, title, "user") VALUES ($1, $2, $3)`,
		forum.Slug,
		forum.Title,
		forum.User,
	)

	return forum, err
}

func (p *postgresAppRepository) SelectForumBySlug(slug string) (models.Forum, error) {
	var forum models.Forum
	err := p.Conn.QueryRow(
		`SELECT * FROM forum WHERE slug=$1 LIMIT 1;`,
		slug).Scan(
		&forum.Slug,
		&forum.Title,
		&forum.User,
		&forum.Posts,
		&forum.Threads,
	)

	return forum, err
}

func (p *postgresAppRepository) InsertThread(thread models.Thread) (models.Thread, error) {
	query := `INSERT INTO thread(slug, author, created, message, title, forum) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`

	var row *pgx.Row
	if thread.Created != "" {
		row = p.Conn.QueryRow(
			query,
			thread.Slug,
			thread.Author,
			thread.Created,
			thread.Message,
			thread.Title,
			thread.Forum,
		)
	} else {
		row = p.Conn.QueryRow(
			query,
			thread.Slug,
			thread.Author,
			time.Time{},
			thread.Message,
			thread.Title,
			thread.Forum,
		)
	}

	var thr models.Thread
	var created time.Time
	err := row.Scan(&thr.Id, &thr.Author, &created, &thr.Forum, &thr.Message, &thr.Slug, &thr.Title, &thr.Votes)

	thr.Created = strfmt.DateTime(created.UTC()).String()

	return thr, err
}

func (p *postgresAppRepository) SelectThreadBySlug(slug string) (models.Thread, error) {
	row := p.Conn.QueryRow(`SELECT * FROM thread WHERE slug=$1 LIMIT 1;`, slug)

	var thread models.Thread
	var created time.Time
	err := row.Scan(&thread.Id, &thread.Author, &created, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)

	thread.Created = strfmt.DateTime(created.UTC()).String()

	return thread, err
}

func (p *postgresAppRepository) SelectThreadById(id int) (models.Thread, error) {
	row := p.Conn.QueryRow(`SELECT * FROM thread WHERE id=$1 LIMIT 1;`, id)

	var thread models.Thread
	var created time.Time
	err := row.Scan(&thread.Id, &thread.Author, &created, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)

	thread.Created = strfmt.DateTime(created.UTC()).String()

	return thread, err
}

func (p *postgresAppRepository) InsertPosts(posts []models.Post) ([]models.Post, error) {
	tx, err := p.Conn.Begin()
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO post(author, created, forum, message, parent, thread) VALUES ($1, $2, $3, $4, NULLIF($5, 0), $6) RETURNING *;`
	currentTime := time.Now()

	var resultPosts []models.Post

	for _, post := range posts {
		var currentPost models.Post
		var created time.Time
		err := tx.QueryRow(
			query,
			post.Author,
			currentTime,
			post.Forum,
			post.Message,
			post.Parent,
			post.Thread,
		).Scan(
			&currentPost.Id,
			&currentPost.Author,
			&created,
			&currentPost.Forum,
			&currentPost.Message,
			&currentPost.IsEdited,
			&currentPost.Parent,
			&currentPost.Thread,
		)
		if err != nil {
			tx.Rollback()

			return nil, err
		}

		currentPost.Created = strfmt.DateTime(created.UTC()).String()
		if !currentPost.Parent.Valid {
			currentPost.Parent.Int64 = 0
			currentPost.Parent.Valid = true
		}
		resultPosts = append(resultPosts, currentPost)
	}

	tx.Commit()

	return resultPosts, nil
}


func (p *postgresAppRepository) UpdateThread(thread models.Thread) (models.Thread, error) {
	query := `UPDATE thread SET title=$1, message=$2 WHERE %s RETURNING *`

	var row *pgx.Row
	if thread.Slug == "" {
		query = fmt.Sprintf(query, `id=$3`)
		row = p.Conn.QueryRow(query, thread.Title, thread.Message, thread.Id)
	} else {
		query = fmt.Sprintf(query, `slug=$3`)
		row = p.Conn.QueryRow(query, thread.Title, thread.Message, thread.Slug)
	}

	var newThread models.Thread
	var created time.Time
	err := row.Scan(
		&newThread.Id,
		&newThread.Author,
		&created,
		&newThread.Forum,
		&newThread.Message,
		&newThread.Slug,
		&newThread.Title,
		&newThread.Votes,
	)

	if err != nil {
		return models.Thread{}, err
	}

	newThread.Created = strfmt.DateTime(created.UTC()).String()

	return newThread, nil
}

func (p *postgresAppRepository) InsertVote(vote models.Vote) (models.Vote, error) {
	_, err := p.Conn.Exec(
		`INSERT INTO votes(nickname, voice, id_thread) VALUES ($1, $2, $3)`,
		vote.Nickname,
		vote.Voice,
		vote.IdThread,
	)

	return vote, err
}

func (p *postgresAppRepository) UpdateVote(vote models.Vote) (models.Vote, error) {
	_, err := p.Conn.Exec(
		`UPDATE vote SET voice=$1 WHERE id_thread=$2 AND nickname=$3`,
		vote.Voice,
		vote.IdThread,
		vote.Nickname,
	)

	return vote, err
}

func (p *postgresAppRepository) GetServiceStatus() (map[string]int, error) {
	info, err := p.Conn.Query(
		`SELECT * FROM (SELECT COUNT(*) FROM forum) as forumCount,
		(SELECT COUNT(*) FROM post) as postCount,
		(SELECT COUNT(*) FROM thread) as threadCount, 
		(SELECT COUNT(*) FROM users) as usersCount;`,
	)
	if err != nil {
		return nil, err
	}

	if info.Next() {
		forumCount, postCount, threadCount, usersCount := 0, 0, 0, 0
		err := info.Scan(&forumCount, &postCount, &threadCount, &usersCount)
		if err != nil {
			return nil, err
		}

		return map[string]int{
			"forum":  forumCount,
			"post":   postCount,
			"thread": threadCount,
			"user":   usersCount,
		}, nil
	}

	return nil, errors.New("have not information")
}

func (p *postgresAppRepository) ClearDatabase() error {
	_, err := p.Conn.Exec(`TRUNCATE users, thread, forum, post, vote`)

	return err
}