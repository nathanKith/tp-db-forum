package repository

import (
	"errors"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	"log"
	"strings"
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

func (p *postgresAppRepository) SelectUsersByForum(slugForum string, parameters models.QueryParameters) ([]models.User, error) {
	query := `SELECT users.nickname, users.fullname, users.about, users.email FROM 
			  ((SELECT thread.author FROM thread WHERE thread.forum=$1) UNION 
			  (SELECT post.author FROM post WHERE post.author=$1)) AS union_users 
			  JOIN users ON union_users.author=users.nickname
			  WHERE LOWER(author) < LOWER($3) ORDER BY LOWER(author) %s LIMIT NULLIF($2)`

	if parameters.Desc {
		query = fmt.Sprintf(query, "DESC")
	} else {
		query = fmt.Sprintf(query, "ASC")
	}

	if parameters.Since == "" {
		query = strings.ReplaceAll(query, `WHERE LOWER(author) < LOWER($3)`, "")
	}

	rows, err := p.Conn.Query(query, slugForum, parameters.Limit, parameters.Since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (p *postgresAppRepository) SelectThreadsByForum(slugForum string, parameters models.QueryParameters) ([]models.Thread, error) {
	var rows *pgx.Rows
	var err error
	if parameters.Since != "" {
		if parameters.Desc {
			rows, err = p.Conn.Query(
				`SELECT * FROM thread WHERE LOWER(forum)=LOWER($1) AND created <= $2 
				ORDER BY created DESC LIMIT NULLIF($3, 0)`,
				slugForum, parameters.Since, parameters.Limit)
		} else {
			rows, err = p.Conn.Query(
				`SELECT * FROM thread WHERE LOWER(forum)=LOWER($1) AND created <= $2 
				ORDER BY created ASC LIMIT NULLIF($3, 0)`,
				slugForum, parameters.Since, parameters.Limit)
		}
	} else {
		if parameters.Desc {
			rows, err = p.Conn.Query(
				`SELECT * FROM thread WHERE LOWER(forum)=LOWER($1)
				ORDER BY created DESC LIMIT NULLIF($3, 0)`,
				slugForum, parameters.Since, parameters.Limit)
		} else {
			rows, err = p.Conn.Query(
				`SELECT * FROM thread WHERE LOWER(forum)=LOWER($1)
				ORDER BY created ASC LIMIT NULLIF($3, 0)`,
				slugForum, parameters.Since, parameters.Limit)
		}
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var thread models.Thread
		var created time.Time

		err := rows.Scan(
			&thread.Id,
			&thread.Author,
			&created,
			&thread.Forum,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes,
		)
		if err != nil {
			return nil, err
		}

		thread.Created = strfmt.DateTime(created.UTC()).String()

		threads = append(threads, thread)
	}

	return threads, nil
}

func (p *postgresAppRepository) SelectPostById(id int) (models.Post, error) {
	var post models.Post
	var created time.Time

	err := p.Conn.QueryRow(
		`SELECT * FROM post WHERE id=$1 LIMIT 1;`,
		id).Scan(
		&post.Id,
		&post.Author,
		&created,
		&post.Forum,
		&post.Message,
		&post.IsEdited,
		&post.Parent,
		&post.Thread,
	)
	if err != nil {
		return models.Post{}, err
	}

	post.Created = strfmt.DateTime(created.UTC()).String()

	return post, nil
}

func (p *postgresAppRepository) UpdatePost(id int, message string) (models.Post, error) {
	var post models.Post
	var created time.Time
	err := p.Conn.QueryRow(
		`UPDATE post SET message=$1, isEdited=true WHERE id=$2 RETURNING *`,
		message,
		id,
	).Scan(
		&post.Id,
		&post.Author,
		&created,
		&post.Forum,
		&post.Message,
		&post.IsEdited,
		&post.Parent,
		&post.Thread,
	)

	post.Created = strfmt.DateTime(created.UTC()).String()

	return post, err
}

func (p *postgresAppRepository) SelectPostsByThread(thread models.Thread, limit, since int, sort string, desc bool) ([]models.Post, error) {
	var threadId int
	if thread.Id == 0 {
		thr, err := p.SelectThreadBySlug(thread.Slug)
		if err != nil {
			return nil, err
		}

		threadId = thr.Id
	} else {
		threadId = thread.Id
	}

	switch sort {
	case "flat":
		posts, err := p.selectPostsByThreadFlat(threadId, limit, since, desc)

		return posts, err
	case "tree":
		posts, err := p.selectPostsByThreadTree(threadId, limit, since, desc)

		return posts, err
	case "parent_tree":
		posts, err := p.selectPostsByThreadParentTree(threadId, limit, since, desc)

		return posts, err
	default:
		return nil, errors.New("u gay")
	}
}

func (p *postgresAppRepository) selectPostsByThreadFlat(id, limit, since int, desc bool) ([]models.Post, error) {
	var rows *pgx.Rows
	var err error
	if since != 0 {
		if desc {
			rows, err = p.Conn.Query(`SELECT * FROM post WHERE thread=$1 ORDER BY id DESC LIMIT NULLIF($2, 0)`, id, limit)
		} else {
			rows, err = p.Conn.Query(`SELECT * FROM post WHERE thread=$1 ORDER BY id ASC LIMIT NULLIF($2, 0)`, id, limit)
		}
	} else {
		if desc {
			rows, err = p.Conn.Query(`SELECT * FROM post WHERE thread=$1 AND id < $2 ORDER BY id DESC LIMIT NULLIF($3, 0)`, id, since, limit)
		} else {
			rows, err = p.Conn.Query(`SELECT * FROM post WHERE thread=$1 AND id > $2 ORDER BY id DESC LIMIT NULLIF($3, 0)`, id, since, limit)
		}
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var created time.Time

		err = rows.Scan(
			&post.Id,
			&post.Author,
			&created,
			&post.Forum,
			&post.Message,
			&post.IsEdited,
			&post.Parent,
			&post.Thread,
		)
		if err != nil {
			return nil, err
		}

		post.Created = strfmt.DateTime(created.UTC()).String()

		posts = append(posts, post)
	}

	return posts, err
}

func (p *postgresAppRepository) selectPostsByThreadTree(id, limit, since int, desc bool) ([]models.Post, error) {
	var rows *pgx.Rows
	var err error

	if since == 0 {
		if desc {
			rows, err = p.Conn.Query(
				`SELECT * FROM post
				WHERE thread=$1 ORDER BY path DESC, id  DESC LIMIT $2;`,
				id, limit,
			)
		} else {
			rows, err = p.Conn.Query(
				`SELECT * FROM post
				WHERE thread=$1 ORDER BY path ASC, id  ASC LIMIT $2;`,
				id, limit,
			)
		}
	} else {
		if desc {
			rows, err = p.Conn.Query(
				`SELECT * FROM post
				WHERE thread=$1 AND PATH < (SELECT path FROM post WHERE id = $2)
				ORDER BY path DESC, id  DESC LIMIT $3;`,
				id, since, limit,
			)
		} else {
			rows, err = p.Conn.Query(
				`SELECT * FROM post
				WHERE thread=$1 AND PATH > (SELECT path FROM post WHERE id = $2)
				ORDER BY path ASC, id  ASC LIMIT $3;`,
				id, since, limit,
			)
		}
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var created time.Time

		err = rows.Scan(
			&post.Id,
			&post.Author,
			&created,
			&post.Forum,
			&post.Message,
			&post.IsEdited,
			&post.Parent,
			&post.Thread,
		)
		if err != nil {
			return nil, err
		}

		post.Created = strfmt.DateTime(created.UTC()).String()

		posts = append(posts, post)
	}

	return posts, nil
}

func (p *postgresAppRepository) selectPostsByThreadParentTree(id, limit, since int, desc bool) ([]models.Post, error) {
	var rows *pgx.Rows
	var err error

	if since == 0 {
		if desc {
			rows, err = p.Conn.Query(
				`SELECT * FROM post
				WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent IS NULL ORDER BY id DESC LIMIT $2)
				ORDER BY path[1] DESC, path, id;`,
				id, limit,
			)
		} else {
			rows, err = p.Conn.Query(
				`SELECT * FROM post
				WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent IS NULL ORDER BY id LIMIT $2)
				ORDER BY path, id;`,
				id, limit,
			)
		}
	} else {
		if desc {
			rows, err = p.Conn.Query(
				`SELECT * FROM post
				WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent IS NULL AND PATH[1] <
				(SELECT path[1] FROM post WHERE id = $2) ORDER BY id DESC LIMIT $3) ORDER BY path[1] DESC, path, id;`,
				id, since, limit,
			)
		} else {
			rows, err = p.Conn.Query(`SELECT * FROM post
				WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent IS NULL AND PATH[1] >
				(SELECT path[1] FROM post WHERE id = $2) ORDER BY id ASC LIMIT $3) ORDER BY path, id;`,
				id, since, limit,
			)
		}
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var created time.Time

		err = rows.Scan(
			&post.Id,
			&post.Author,
			&created,
			&post.Forum,
			&post.Message,
			&post.IsEdited,
			&post.Parent,
			&post.Thread,
		)
		if err != nil {
			return nil, err
		}

		post.Created = strfmt.DateTime(created.UTC()).String()

		posts = append(posts, post)
	}

	return posts, nil
}