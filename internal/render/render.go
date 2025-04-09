package render

import (
	"bytes"
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/sabrodigan/bookings-app/internal/config"
	"github.com/sabrodigan/bookings-app/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders a template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
func CheckTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	cwd, _ := os.Getwd()
	log.Println("Current working directory:", cwd)

	templatesPath := filepath.Join(cwd, "templates")
	pages, err := filepath.Glob(filepath.Join(templatesPath, "*.page.tmpl"))
	if err != nil {
		log.Println("Error finding page templates:", err)
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			log.Printf("Error parsing page template %s: %v", name, err)
			return myCache, err
		}

		matches, err := filepath.Glob(filepath.Join(templatesPath, "*.layout.tmpl"))
		if err != nil {
			log.Println("Error finding layout templates:", err)
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(filepath.Join(templatesPath, "*.layout.tmpl"))
			if err != nil {
				log.Printf("Error parsing layout templates for %s: %v", name, err)
				return myCache, err
			}
		}

		myCache[name] = ts
		log.Printf("Template cached: %s", name)
	}

	return myCache, nil
}
