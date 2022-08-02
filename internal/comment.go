package internal

import (
	"database/sql"
	"fmt"

	"github.com/SmoothWay/forum/pkg/models"
)

func (app *Application) GetCommentsByPostID(id int) ([]*models.Comment, error) {
	comments, err := app.Posts.GetCommentsByPostID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRecord
		}
		return nil, fmt.Errorf("GetCommentsByPostID: %w", err)
	}

	for _, comment := range comments {
		votes, err := app.Posts.CountByCommentID(comment.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, models.ErrNoRecord
			}
			return nil, fmt.Errorf("GetVotesCountByCommentID: %w", err)
		}
		comment.Like = votes.Like
		comment.Dislike = votes.Dislike
	}
	return comments, nil
}
