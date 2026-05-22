// Package handler - HTTP обработчики запросов
package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"notes-api/internal/model"
	"notes-api/internal/service"
)

type NoteHandler struct {
	svc *service.NoteService
}

func NewNoteHandler(svc *service.NoteService) *NoteHandler {
	return &NoteHandler{svc: svc}
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid body"}`, http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, `{"error": "title is required"}`, http.StatusBadRequest)
		return
	}

	note := h.svc.Create(req)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func (h *NoteHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	notes := h.svc.GetActive()
	if notes == nil {
		notes = []*model.Note{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *NoteHandler) MarkDone(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	note, err := h.svc.MarkDone(id)
	if err != nil {
		http.Error(w, `{"error": "note not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func extractID(path string) string {
	// Пример пути: /api/v1/notes/123/done
	parts := strings.Split(strings.TrimSuffix(path, "/done"), "/")
	return parts[len(parts)-1]
}
