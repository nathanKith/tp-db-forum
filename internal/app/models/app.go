package models

import (
	"github.com/jackc/pgx/pgtype"
	"encoding/json"
	"database/sql"
)

type Error struct {
	Message string `json:"message"`
}

type User struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	FullName string `json:"fullname"`
	Nickname string `json:"nickname"`
}

type Forum struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int    `json:"posts"`
	Threads int    `json:"threads"`
}

type Thread struct {
	Id      int        `json:"id"`
	Author  string     `json:"author"`
	Created string     `json:"created"`
	Forum   string     `json:"forum"`
	Title   string 	   `json:"title"`
	Message string 	   `json:"message"`
	Slug    string     `json:"slug"`
	Votes   int        `json:"votes"`
}
