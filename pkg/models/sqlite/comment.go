package sqlite

import (
	"github.com/SmoothWay/forum/pkg/models"
)

func (m *PostModel) GetCommentsByPostID(id int) ([]*models.Comment, error) {

	stmt := `SELECT comments.ID, comments.content, comments.postid, users.nickname
			FROM comments JOIN users 
			ON comments.userid = users.id 
			WHERE comments.postid = ?
			ORDER BY comments.id`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}
	for rows.Next() {
		comment := models.Comment{}
		err = rows.Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.Nickname)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}
