package internal

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/SmoothWay/forum/pkg/forms"
	"github.com/SmoothWay/forum/pkg/models"
)

type TemplateData struct {
	AuthenticatedUser int
	CurrentYear       int
	Flash             string
	Form              *forms.Form
	Post              *models.Post
	Posts             []*models.Post
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func NewTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*-page.html"))
	if err != nil {
		return nil, err
	}

	// loop through the pages one by one.

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*-layout.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*-partial.html"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts

	}
	return cache, nil
}
