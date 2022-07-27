package internal

import (
	"net/http"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	if r.Method != http.MethodGet {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	s, err := app.Posts.Latest()
	if err != nil {
		app.serverError(w, err)
	}
	app.render(w, r, "home-page.html", &TemplateData{Posts: s})
}
