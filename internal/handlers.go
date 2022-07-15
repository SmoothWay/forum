package internal

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/SmoothWay/forum/pkg/forms"
	"github.com/SmoothWay/forum/pkg/models"
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
	log.Println(s)
	app.render(w, r, "home-page.html", &TemplateData{Posts: s})
}

func (app *Application) showPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.Posts.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "show-page.html", &TemplateData{
		Post: s,
	})
}

func (app *Application) createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.render(w, r, "create-page.html", &TemplateData{
			Form: forms.New(nil),
		})
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		form := forms.New(r.PostForm)
		form.Required("title", "content", "tags")
		form.MaxLength("title", 100)

		if !form.Valid() {
			app.render(w, r, "create-page.html", &TemplateData{Form: form})
			return
		}
		c, _ := r.Cookie("forum")

		u, err := app.Session.GetUserByUUID(c.Value)
		if err != nil {
			app.serverError(w, err)
			return
		}
		id, err := app.Posts.Insert(form.Get("title"), form.Get("content"), u.ID)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// app.Session.Put(r, "flash", "Post successfullt created!")
		http.Redirect(w, r, fmt.Sprintf("/post/?id=%d", id), http.StatusSeeOther)
	} else {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

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
	if app.GetSession(r) != 0 {
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
		app.AddSession(w, r, id)
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
	err := app.RemoveSession(w, r)
	if err != nil {
		app.serverError(w, err)
	}
	// app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
