package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

// TODO: Fix structs
// TODO: decide to use time.Time or string

type User struct {
	ID             int
	Nickname       string
	Email          string
	HashedPassword []byte
}

type Post struct {
	ID         int
	Like       int
	Dislike    int
	Title      string
	Content    string
	Created    string
	Categories []string
	User       *User
	Comments   []Comment
}

type Category struct {
	ID   int
	Name string
}

type Comment struct {
	Content  string
	Nickname string
}

type Session struct {
	ID      int
	UUID    string
	Expires time.Time
	UserID  int
}
