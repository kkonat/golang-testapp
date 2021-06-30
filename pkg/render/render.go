package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/kkonat/hotel-webapp/pkg/config"
	"github.com/kkonat/hotel-webapp/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	var templCache map[string]*template.Template

	if app.UseCache { // production don't re read, since .tmpl files don't change
		templCache = app.TemplateCache
	} else {
		templCache, _ = CreateTemplateCache() // re-read each time casuse .tmpl files could be changed externally
		app.TemplateCache = templCache
		app.UseCache = true
	}

	t, ok := templCache[tmpl] // get cached data from the map
	if !ok {
		log.Fatal("Counld not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)
	_ = t.Execute(buf, td) // This combines data with template

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Fatal("Error writing template to browser", err)
	}
}

//  CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	//get all page template files matching the *.page.tmpl pattern
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		//create template basing on each page
		nTempl, err := template.New(name).Funcs(functions).ParseFiles(page)
		log.Println("Building Template Cache", page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		// if there are any templates, then parse them
		if len(matches) > 0 {
			nTempl, err = nTempl.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}

			myCache[name] = nTempl // store each new template in the cache
		}
	}
	return myCache, nil
}
