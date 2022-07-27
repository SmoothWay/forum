package internal

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/SmoothWay/forum/pkg/models"
)

func (app *Application) authenticatedUser(r *http.Request) bool {
	return isSession(r)
}

func (app *Application) addDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	if td == nil {
		td = &TemplateData{}
	}
	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()
	// td.Flash = app.Session.PopString(r, flash)
	return td
}

func (app *Application) render(w http.ResponseWriter, r *http.Request, name string, td *TemplateData) {
	ts, ok := app.TemplateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exits", name))
		return
	}
	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)
}

func (app *Application) postVote(userid, postid int, formVote int) error {
	vote, err := app.Posts.IsVotedPost(userid, postid)
	if err == models.ErrNoRecord {
		err = app.Posts.InsertPostEvaluate(userid, postid, formVote)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	if vote != formVote {
		app.Posts.DelPostVote(userid, postid)
		err = app.Posts.InsertPostEvaluate(userid, postid, formVote)
		if err != nil {
			return err
		}
	} else {
		app.Posts.DelPostVote(userid, postid)
	}
	return nil
}
func (app *Application) commentVote(userid, commentID, formVote int) error {
	vote, err := app.Posts.IsVotedComment(userid, commentID)
	if err == models.ErrNoRecord {
		err = app.Posts.InsertCommentEvaluate(userid, commentID, formVote)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	if vote != formVote {
		app.Posts.DelCommentVote(userid, commentID)
		err = app.Posts.InsertCommentEvaluate(userid, commentID, formVote)
		if err != nil {
			return err
		}
	} else {
		app.Posts.DelPostVote(userid, commentID)
	}
	return nil
}
func (app *Application) execTemp(w http.ResponseWriter, status int) {
	templates, templErr := template.ParseFiles("./ui/templates/errors.html")
	if templErr != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		app.ErrorLog.Fatal(templErr)
	}
	templates.Execute(w, status)
}

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	w.WriteHeader(http.StatusInternalServerError)
	app.execTemp(w, http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	app.execTemp(w, status)
}

func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
