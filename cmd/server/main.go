package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"notes-api/internal/handler"
	"notes-api/internal/service"
	"notes-api/internal/store"
)

func main() {
	noteStore := store.NewNoteStore()
	noteService := service.NewNoteService(noteStore)
	noteHandler := handler.NewNoteHandler(noteService)

	mux := http.NewServeMux()

	// Раздача статических файлов (фронтенд)
	mux.Handle("/", http.FileServer(http.Dir("static")))

	// Swagger UI
	mux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile("api/swagger.html")
		if err != nil {
			http.Error(w, "Swagger not found", 404)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(content)
	})

	// Отдача openapi.yaml
	mux.HandleFunc("/api/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, "api/openapi.yaml")
	})

	// API
	mux.HandleFunc("/api/v1/notes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			return
		}

		switch r.Method {
		case http.MethodPost:
			noteHandler.Create(w, r)
		case http.MethodGet:
			noteHandler.GetActive(w, r)
		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/notes/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			return
		}

		if r.Method == http.MethodPut && strings.HasSuffix(r.URL.Path, "/done") {
			noteHandler.MarkDone(w, r)
		} else {
			http.Error(w, `{"error": "not found"}`, http.StatusNotFound)
		}
	})

	log.Println("===========================================")
	log.Println("Фронтенд:  http://localhost:8080")
	log.Println("Swagger:   http://localhost:8080/swagger")
	log.Println("===========================================")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
