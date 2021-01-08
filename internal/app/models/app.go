package models

import (
	"database/sql"
	"encoding/json"
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

type Post struct {
	Id       int    `json:"id"`
	Author   string `json:"author"`
	Created  string `json:"created"`
	Forum    string `json:"forum"`
	Message  string `json:"message"`
	IsEdited bool   `json:"isEdited"`
	Parent   JsonNullInt    `json:"parent"`
	Thread   int    `json:"thread"`
}

type JsonNullInt struct {
	sql.NullInt64
}

func (v JsonNullInt) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullInt) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

var PostParentError = `insert or update on table "post" violates foreign key constraint "post_parent_fkey"`

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int32  `json:"voice"`
	IdThread int64  `json:"-"`
}
