package render

import (
	"Ebook/internal/config"
	"Ebook/internal/models"
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// Gán token CSRF
	td.CSRFToken = nosurf.Token(r)
	// Kiểm tra và khởi tạo map td.Data nếu cần
	if td.Data == nil {
		td.Data = make(map[string]interface{})
	}

	// Kiểm tra context để lấy user
	user, ok := r.Context().Value("user").(models.User)
	if ok {
		td.Data["User"] = user
	}

	return td
}

// Template renders a template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		var err error
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Fatal("Cannot create template cache:", err)
		}
	}

	// Kiểm tra xem template có tồn tại không
	t, ok := tc[tmpl]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	if td == nil {
		td = &models.TemplateData{Data: map[string]interface{}{}}
	}
	td = AddDefaultData(td, r)

	// Thực hiện template
	err := t.Execute(buf, td)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
		return
	}

	// Ghi buffer vào response
	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, "Error writing template to browser", http.StatusInternalServerError)
		log.Println("Error writing buffer to response:", err)
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

		// Kiểm tra layout templates
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
