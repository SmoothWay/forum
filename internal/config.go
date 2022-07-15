package internal

import (
	"html/template"
	"log"

	"github.com/SmoothWay/forum/pkg/models/sqlite"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Session  *sqlite.Session
	User     *sqlite.UserModel
	Posts    *sqlite.PostModel

	TemplateCache map[string]*template.Template
}
