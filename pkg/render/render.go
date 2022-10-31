package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/dee-d-dev/gobookings/pkg/config"
	"github.com/dee-d-dev/gobookings/pkg/models"
)

// var functions = template.funcMap{

// }
var functions = template.FuncMap{}
var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	//crete template cache

	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	//get template from cache
	t, ok := tc[tmpl]

	if !ok {
		log.Fatal("could not render template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	//render the template

	_, err := buf.WriteTo(w)

	if err != nil {
		log.Println(err)
	}

	parsedTemplate, err := template.ParseFiles("../../templates/"+tmpl, "../../templates/base.layout.html")

	// if err != nil {
	// 	fmt.Println(err)
	// }

	_ = parsedTemplate.Execute(w, nil)

	if err != nil {
		fmt.Println("error parsing template:", err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("../../templates/*.page.html")

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("../../templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("../../templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts

	}
	return myCache, err
}
