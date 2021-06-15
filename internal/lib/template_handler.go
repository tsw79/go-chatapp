package lib

import (
	handler "chatapp/internal/app/handlers"
	"chatapp/internal/lib/config"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/objx"
)

/* Represents a single template */
type TemplateHandler struct {
	once  sync.Once
	Files []string
	Data  handler.Page
	templ *template.Template
}

/* ServeHTTP handles the HTTP request.
 * 	This will ensure we compile the template once inside the `ServeHTTP` function.
 *	`sync.Once` ensures the function passed as an argument will be executed only once.
 */
func (this *TemplateHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	this.once.Do(func() {
		len := len(this.Files)
		for i := 0; i < len; i++ {
			// Append templates directory to all file paths
			this.Files[i] = filepath.Join(config.GetInstance().Dir.Templates, this.Files[i])
		}
		this.templ = template.Must(template.ParseFiles(this.Files...))
	})

	// data["Host"]
	data := map[string]interface{}{
		"Host": req.Host,
	}
	// data["UserData"]
	if authCookie, err := req.Cookie(config.GetInstance().Auth.Cookie); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	// Add Page data
	data["Page"] = *&this.Data
	// rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := this.templ.Execute(rw, data); err != nil {
		log.Fatalln(err)
	}
}
