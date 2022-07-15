package internal

import (
	"net/http"
)

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	// mux.Handle("/post/create", app.requireAuthenticatedUser((app.createSnippetForm)))
	mux.Handle("/post/create", app.requireAuthenticatedUser((app.createPost)))
	mux.HandleFunc("/post", (app.showPost))
	mux.HandleFunc("/user/signup", (app.signupUser))
	mux.HandleFunc("/user/login", (app.loginUser))
	mux.Handle("/user/logout", app.requireAuthenticatedUser((app.logoutUser)))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return app.recoverPanic(app.logRequest(app.secureHeaders((mux))))
}
