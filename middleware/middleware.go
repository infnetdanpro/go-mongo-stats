package middleware

import "net/http"

// TODO: auth mw

func SetJSONHeader(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		f(w, r)
	}
}

func SetHTMLHeader(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		f(w, r)
	}
}
