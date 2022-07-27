package sqlite

import (
	"database/sql"

	"github.com/SmoothWay/forum/pkg/models"
)

func (m *PostModel) InsertPostEvaluate(userid, postid int, vote int) error {
	stmt := `INSERT INTO posts_evaluate(vote, userid, postid)
			VALUES(?, ?, ?)`
	_, err := m.DB.Exec(stmt, vote, userid, postid)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostModel) IsVotedPost(userid, postid int) (int, error) {
	var vote int
	stmt := `SELECT vote FROM posts_evaluate WHERE userid = ? and postid = ?`
	row := m.DB.QueryRow(stmt, userid, postid)
	err := row.Scan(&vote)
	if err == sql.ErrNoRows {
		return 0, models.ErrNoRecord
	} else if err != nil {
		return 0, err
	}
	return vote, nil
}

func (m *PostModel) DelPostVote(userid, postid int) error {
	stmt := `DELETE FROM posts_evaluate WHERE userid = ? AND postid = ?`
	_, err := m.DB.Exec(stmt, userid, postid)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostModel) GetPostEvaluate(postid int) (int, int, error) {
	var like int
	var dislike int
	likestmt := `SELECT COUNT(vote) FROM posts_evaluate WHERE vote = 1 AND postid = ?`
	dislikestmt := `SELECT COUNT(vote) FROM posts_evaluate WHERE vote = -1 AND postid = ?`

	rows, err := m.DB.Query(likestmt, postid)
	if err != nil {
		return 0, 0, err
	}

	for rows.Next() {
		err = rows.Scan(&like)
		if err != nil {
			return 0, 0, err
		}
	}
	rows, err = m.DB.Query(dislikestmt, postid)
	if err != nil {
		return 0, 0, err
	}
	for rows.Next() {
		err = rows.Scan(&dislike)
		if err != nil {
			return 0, 0, err
		}
	}

	return like, dislike, nil
}
