package api

import (
	"errors"
	"net/http"

	"github.com/fmiskovic/eth-auth/store"
)

type Api struct {
	http.Handler
}

func New(secret string) *Api {
	mux := http.NewServeMux()

	storer, err := store.New(1024)
	if err != nil {
		panic(err)
	}
	h := newHandler(storer, secret)

	// serve index.html at the root "/"
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// serve other static files (CSS, JS, etc.)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Add api routes here
	mux.HandleFunc("/nonce", httpHandlerFunc(h.handleNonce))
	mux.HandleFunc("/auth", httpHandlerFunc(h.handleAuth))

	return &Api{mux}
}

// HandlerFunc is a function that can be used as an HTTP handler.
type handlerFunc func(w http.ResponseWriter, r *http.Request) error

// HttpHandlerFunc wraps a HandlerFunc and returns a http.HandlerFunc.
func httpHandlerFunc(handle handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Accept", "application/json")
		w.Header().Set("Content-Type", "application/json")

		err := handle(w, r)

		// handle error if any
		if err != nil {
			code := http.StatusInternalServerError
			var apiErr *Error
			if errors.As(err, &apiErr) {
				code = apiErr.Code
			}

			http.Error(w, err.Error(), code)
		}
	}
}
