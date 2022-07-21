package sqlite

import (
	"database/sql"
	"time"

	"github.com/SmoothWay/forum/pkg/models"
)

type Session struct {
	DB *sql.DB
}

func (s *Session) GetUserByUUID(uuid string) (*models.User, error) {
	stmt := `SELECT userid FROM sessions WHERE uuid = ?`
	user := models.User{}
	row := s.DB.QueryRow(stmt, uuid)
	err := row.Scan(&user.ID)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	stmt2 := `SELECT nickname, email, hashed_password FROM users WHERE id = ?`

	user2 := &models.User{}

	row2 := s.DB.QueryRow(stmt2, user.ID)
	err = row2.Scan(&user2.Nickname, &user2.Email, &user2.HashedPassword)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	user2.ID = user.ID

	return user2, nil
}

func (s *Session) Insert(userID int, uuid string) error {
	stmt := `INSERT INTO sessions(uuid, expires, userid)
			VALUES(?, ?, ?)`
	now := time.Now().Add(time.Hour * 3).Format("01-02-2006 15:04:05")
	_, err := s.DB.Exec(stmt, uuid, now, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) Delete(userID int) {
	stmt := `DELETE FROM sessions WHERE userid = ?`
	s.DB.Exec(stmt, userID)
}

func (s *Session) Exists(userID int) error {
	stmt := `SELECT userid FROM sessions
			WHERE userid = $1`

	var id int
	row := s.DB.QueryRow(stmt, userID)
	err := row.Scan(&id)
	if err != nil {
		// if err == sql.ErrNoRows{
		// return models.ErrNoRecord
		// }
		return err
	}
	return nil
}
