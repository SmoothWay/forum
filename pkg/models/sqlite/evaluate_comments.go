package sqlite

import (
	"database/sql"
	"log"

	"github.com/SmoothWay/forum/pkg/models"
)

func (m *PostModel) InsertCommentEvaluate(userid, commentID, vote int) error {
	stmt := `INSERT INTO comments_evaluate(vote, userid, commentid)
			VALUES(?, ?, ?)`
	_, err := m.DB.Exec(stmt, vote, userid, commentID)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostModel) IsVotedComment(userid, commentID int) (int, error) {
	var vote int
	stmt := `SELECT vote FROM comments_evaluate WHERE userid = ? and commentid = ?`
	row := m.DB.QueryRow(stmt, userid, commentID)
	err := row.Scan(&vote)
	if err == sql.ErrNoRows {
		return 0, models.ErrNoRecord
	} else if err != nil {
		return 0, err
	}
	return vote, nil
}

func (m *PostModel) DelCommentVote(userid, commentID int) error {
	stmt := `DELETE FROM comments_evaluate WHERE userid = ? AND commentid = ?`
	_, err := m.DB.Exec(stmt, userid, commentID)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostModel) GetCommentEvaluate(commentID int) (int, int, error) {
	var like int
	var dislike int
	likestmt := `SELECT COUNT(vote) FROM comments_evaluate WHERE vote = 1 AND commentid = ?`
	dislikestmt := `SELECT COUNT(vote) FROM comments_evaluate WHERE vote = -1 AND commentid = ?`

	rows, err := m.DB.Query(likestmt, commentID)
	if err != nil {
		return 0, 0, err
	}

	for rows.Next() {
		err = rows.Scan(&like)
		if err != nil {
			return 0, 0, err
		}
	}
	rows, err = m.DB.Query(dislikestmt, commentID)
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

func (u *PostModel) GetVotesCountByCommentID(commentID int64) (*models.Comment, error) {
	rows, err := u.DB.Query(`SELECT "vote", count("vote") 
			FROM "comments_evaluate"
			WHERE commentid = ? 
			GROUP BY "vote"
			ORDER BY "vote" desc`, commentID)
	if err != nil {
		return nil, err
	}
	votes := &models.Comment{
		Like:    0,
		Dislike: 0,
	}
	for rows.Next() {
		var voteType int64
		var cnt int64
		err := rows.Scan(&voteType, &cnt)
		if err != nil {
			return nil, err
		}
		switch voteType {
		case 1:
			votes.Like = uint64(cnt)
		case -1:
			votes.Dislike = uint64(cnt)
		default:
			log.Println("Get Votes count bug:", voteType, cnt)
		}
	}
	return votes, nil
}

func (m *PostModel) CountByCommentID(commentID int) (*models.Comment, error) {
	rows, err := m.DB.Query(`SELECT "vote", count("vote") 
			FROM "comments_evaluate"
			WHERE commentid = ? 
			GROUP BY "vote"
			ORDER BY "vote" DESC`, commentID)

	if err != nil {
		return nil, err
	}
	votes := &models.Comment{
		Like:    0,
		Dislike: 0,
	}
	for rows.Next() {
		var voteType int64
		var cnt int64
		err := rows.Scan(&voteType, &cnt)
		if err != nil {
			return nil, err
		}
		switch voteType {
		case 1:
			votes.Like = uint64(cnt)
		case -1:
			votes.Dislike = uint64(cnt)
		default:
			log.Println("Get Votes count bug:", voteType, cnt)
		}
	}
	return votes, nil
}
