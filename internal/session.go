package internal

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func (app *Application) AddSession(w http.ResponseWriter, r *http.Request, id int) {
	u := uuid.NewV4()

	app.Session.Insert(id, u.String())
	http.SetCookie(w,
		&http.Cookie{
			Name:   "forum",
			Value:  u.String(),
			MaxAge: 1800,
			Path:   "/",
		})
}

func (app *Application) GetSession(r *http.Request) int {
	s, err := r.Cookie("forum")
	if err != nil {
		return 0
	}
	u, err := app.Session.GetUserByUUID(s.Value)
	if err != nil {
		return 0
	}
	err = app.Session.Exists(u.ID)
	if err != nil {
		return 0
	}
	return u.ID
}

func (app *Application) RemoveSession(w http.ResponseWriter, r *http.Request) error {
	c, _ := r.Cookie("forum")
	u, err := app.Session.GetUserByUUID(c.Value)
	if err != nil {
		return err
	}
	app.Session.Delete(u.ID)
	http.SetCookie(w, &http.Cookie{
		Name:   "forum",
		Value:  "",
		MaxAge: -1,
	})
	return nil
}
