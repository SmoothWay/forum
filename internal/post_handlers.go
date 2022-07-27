package internal

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/SmoothWay/forum/pkg/forms"
	"github.com/SmoothWay/forum/pkg/models"
)

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
		form.Required("title", "content", "categories")
		tags := strings.Split(form.Get("categories"), ",")
		form.MaxLength("title", 100)
		form.TagField(tags, len(tags))
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
		id, err := app.Posts.Insert(form.Get("title"), form.Get("content"), tags, u.ID)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)
	} else {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *Application) showPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		postid, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || postid < 1 {
			app.notFound(w)
			return
		}
		s, err := app.Posts.Get(postid)
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
	} else if r.Method == http.MethodPost {
		postid, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || postid < 1 {
			app.notFound(w)
			return
		}
		c, err := r.Cookie("forum")
		if err != nil {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		u, err := app.Session.GetUserByUUID(c.Value)
		if err == models.ErrNoRecord {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		err = r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		form := forms.New(r.PostForm)
		form.Required("comment")
		form.MaxLength("comment", 2000)
		if !form.Valid() {
			app.render(w, r, "/show-page.html", &TemplateData{Form: form})
			return
		}
		err = app.Posts.InsertComment(form.Get("comment"), u.ID, postid)
		if err != nil {

			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postid), http.StatusSeeOther)
	} else {
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (app *Application) voteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || postID < 1 {
			app.notFound(w)
			return
		}
		userID, err := GetUserIDByCookie(r)
		if err == models.ErrNoRecord {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		vote, err := strconv.Atoi(r.URL.Query().Get("vote"))
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return

		}
		commentID, err := strconv.Atoi(r.URL.Query().Get("comm"))
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		err = app.commentVote(userID, commentID, vote)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postID), http.StatusSeeOther)
	}
}

func (app *Application) votePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || postID < 1 {
			app.notFound(w)
			return
		}
		userID, err := GetUserIDByCookie(r)
		if err == models.ErrNoRecord {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		vote, err := strconv.Atoi(r.URL.Query().Get("vote"))
		if err != nil {
			app.clientError(w, http.StatusBadRequest)

		}
		err = app.postVote(userID, postID, vote)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postID), http.StatusSeeOther)
	}
}
