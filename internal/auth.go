package internal

import (
	"net/http"

	"github.com/SmoothWay/forum/pkg/forms"
	"github.com/SmoothWay/forum/pkg/models"
)

func (app *Application) signupUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.render(w, r, "signup-page.html", &TemplateData{
			Form: forms.New(nil),
		})
	} else if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		form := forms.New(r.PostForm)
		form.Required("name", "email", "password")
		form.MatchesPattern("email", forms.EmailRegexp)
		form.MinLength("password", 10)

		if !form.Valid() {
			app.render(w, r, "signup-page.html", &TemplateData{Form: form})
			return
		}
		err = app.User.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
		if err == models.ErrDuplicateEmail {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup-page.html", &TemplateData{Form: form})
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	} else {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *Application) loginUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/login" {
		app.notFound(w)
		return
	}
	if isSession(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	if r.Method == http.MethodGet {
		app.render(w, r, "login-page.html", &TemplateData{
			Form: forms.New(nil),
		})
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		form := forms.New(r.PostForm)
		id, err := app.User.Authenticate(form.Get("email"), form.Get("password"))
		if err == models.ErrInvalidCredentials {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login-page.html", &TemplateData{
				Form: form,
			})
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}
		app.Session.Delete(id)
		AddCookie(w, r, id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (app *Application) logoutUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	DeleteCookie(w, r)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
