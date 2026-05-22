package service

import (
	"fmt"
	"notes-api/internal/model"
	"notes-api/internal/store"
	"time"
)

type NoteService struct {
	store *store.NoteStore
}

func NewNoteService(s *store.NoteStore) *NoteService {
	return &NoteService{store: s}
}

func (svc *NoteService) Create(req model.CreateNoteRequest) *model.Note {
	note := &model.Note{
		ID:     fmt.Sprintf("%d", time.Now().UnixNano()),
		Title:  req.Title,
		Tags:   req.Tags,
		Status: "active",
	}
	svc.store.Create(note)
	return note
}

func (svc *NoteService) GetActive() []*model.Note {
	return svc.store.FindActive()
}

func (svc *NoteService) MarkDone(id string) (*model.Note, error) {
	return svc.store.MarkDone(id)
}
