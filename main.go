package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		uuid_, err := uuid.NewUUID()
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			fmt.Println("uuid err: ", err)
			return
		}
		randomId := uuid_.String()

		renderEditorPage(w, map[string]string{
			"id":   randomId,
			"data": "",
		})
	})

	r.Post("/save/{id}", func(w http.ResponseWriter, r *http.Request) {
		bx, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "read error", http.StatusInternalServerError)
			fmt.Println("ioutil.ReadAll(r.Body) err: ", err)
			return
		}
		fmt.Println("got new data: ", string(bx))
		docId := chi.URLParam(r, "id")
		fileName := fmt.Sprintf("./storage/%s.json", docId)
		err = os.WriteFile(fileName, bx, fs.FileMode(0664))
		if err != nil {
			http.Error(w, "write error", http.StatusInternalServerError)
			fmt.Println("os.WriteFile err: ", err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		docId := chi.URLParam(r, "id")
		fileName := fmt.Sprintf("./storage/%s.json", docId)

		bx, err := os.ReadFile(fileName)
		if err != nil {
			http.Error(w, "doc not found", http.StatusNotFound)
			fmt.Println("doc not found err: ", err)
			return
		}
		renderEditorPage(w, map[string]string{
			"id":   docId,
			"data": string(bx),
		})
	})

	http.ListenAndServe(":8000", r)
}

func renderEditorPage(w http.ResponseWriter, data map[string]string) {
	tmpl, err := template.ParseFiles("./templates/editor.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("index route template err: ", err)
		return
	}
	err = tmpl.ExecuteTemplate(w, "editor.html", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("index route template err: ", err)
		return
	}
}
