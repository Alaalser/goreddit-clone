package web

import (
	"net/http"
	"text/template"

	"github.com/alaalser/goreddit"
	"github.com/go-chi/chi"
)

func NewHandler(store goreddit.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		store: store,
	}

	h.Route("/threads", func(r chi.Router) {
		r.Get("/", h.threadList())
	})

	return h
}

type Handler struct {
	*chi.Mux

	store goreddit.Store
}

const threadsListHTML = `
<h1>Threads</h1>
<dl>
{{range .Threads}}
<dit><strong>{{.Title}}</strong><dit>
<dit>{{.Description}}<dit>
{{end}}
</dl>
`

func (h *Handler) threadList() http.HandlerFunc {
	type data struct {
		Threads []goreddit.Thread
	}
	temp := template.Must(template.New("").Parse(threadsListHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		tt, err := h.store.Threads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		temp.Execute(w, data{Threads: tt})
	}
}
