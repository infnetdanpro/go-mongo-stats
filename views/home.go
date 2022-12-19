package views

import (
	"html/template"
	"net/http"
	"os"
)

func Home(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles(os.Getenv("TEMPLATED_DIR") + "index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, "test")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
