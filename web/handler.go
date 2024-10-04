package web

import (
	"net/http"
	"text/template"

	"github.com/alaalser/goreddit"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func NewHandler(store goreddit.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		store: store,
	}

	h.Route("/threads", func(r chi.Router) {
		r.Get("/", h.threadList())
		r.Get("/new", h.ThreadsCreate())
		r.Post("/", h.ThreadsStore())
		r.Post("/{id}/delete", h.ThreadsDelete()) // Fixed route pattern
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
	<dt><strong>{{.Title}}</strong></dt> <!-- Fixed closing tag -->
	<dt>{{.Description}}</dt> <!-- Fixed closing tag -->
	<dd>
	<form action="/threads/{{.ID}}/delete" method="POST"> <!-- Fixed action attribute -->
		<button type="submit">Delete</button>
		</form>
	</dd>
{{end}}
</dl>
<a href="threads/new">Create thread</a>
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

const threadsCreateHTML = `
<h1>Thread</h1>
<form action="/threads" method="POST">
   <table>
	 	<tr>
		<td>Title</td>
		<td><input type="text" name="title" /></td>
		</tr>
		<tr>
		<td>Description</td>
		<td><input type="text" name="description" /></td>
		</tr>
	 </table>
	 <button type="submit">Create thread</button>
</form>
`

func (h *Handler) ThreadsCreate() http.HandlerFunc {

	temp := template.Must(template.New("").Parse(threadsCreateHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		temp.Execute(w, nil)
	}
}

func (h *Handler) ThreadsStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		description := r.FormValue("description")

		if err := h.store.CreateThread(&goreddit.Thread{
			ID:          uuid.New(),
			Title:       title,
			Description: description,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}

func (h *Handler) ThreadsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.store.DeleteThread(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}
