package sqlite

import (
	"database/sql"
	"time"

	"github.com/SmoothWay/forum/pkg/models"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title, content string, userID int) (int, error) {
	var id int
	stmt := `INSERT INTO posts (title, content, created, userid)
			VALUES($1, $2, $3, $4) RETURNING id`
	err := m.DB.QueryRow(stmt, title, content, time.Now().Format("01-02-2006 15:04:05"), userID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *PostModel) Get(id int) (*models.Post, error) {
	s := &models.Post{}
	stmt := `SELECT id, title, content, created FROM posts
			WHERE id = $1`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
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
