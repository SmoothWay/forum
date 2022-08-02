package sqlite

import (
	"database/sql"
	"strings"
	"time"

	"github.com/SmoothWay/forum/pkg/models"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title, content string, tags []string, userID int) (int, error) {
	var id int
	var tag string
	stmt := `INSERT INTO posts (title, content, created, userid)
			VALUES($1, $2, $3, $4) RETURNING id`
	stmt2 := `INSERT INTO category(id, name)
			VALUES($1, $2)`
	err := m.DB.QueryRow(stmt, title, content, time.Now().Format("01-02-2006 15:04:05"), userID).Scan(&id)
	if err != nil {
		return 0, err
	}
	for i := 0; i < len(tags); i++ {
		tag = strings.TrimSpace(tags[i])
		_, err = m.DB.Exec(stmt2, id, tag)
		if err != nil {
			return 0, err
		}
	}
	return int(id), nil
}

func (m *PostModel) Get(id int) (*models.Post, error) {
	Post := &models.Post{}
	var tag string
	stmt := `SELECT id, title, content, created FROM posts
			WHERE id = $1`
	stmt2 := `SELECT name
			FROM category
			WHERE id = $1`

	// Extracting post

	err := m.DB.QueryRow(stmt, id).Scan(&Post.ID, &Post.Title, &Post.Content, &Post.Created)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	// Extracting post tags
	rows, err := m.DB.Query(stmt2, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		Post.Categories = append(Post.Categories, tag)
	}
	// Extracting comments

	// Extracting post evaluations
	like, dislike, err := m.GetPostEvaluate(Post.ID)
	if err != nil {
		return nil, err
	}
	Post.Like = like
	Post.Dislike = dislike

	return Post, nil
}

func (m *PostModel) Latest() ([]*models.Post, error) {
	stmt := `SELECT id, title, content, created FROM posts
	ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		s := &models.Post{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *PostModel) InsertComment(content string, userid, postid int) error {
	stmt := `INSERT INTO comments(content, userid, postid)
			VALUES(?, ?, ?)`
	_, err := m.DB.Exec(stmt, content, userid, postid)
	if err != nil {
		return err
	}
	return nil
}
