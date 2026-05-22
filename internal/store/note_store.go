package store

import (
	"errors"
	"notes-api/internal/model"
	"sync"
)

type NoteStore struct {
	mu    sync.RWMutex
	notes map[string]*model.Note
}

func NewNoteStore() *NoteStore {
	return &NoteStore{
		notes: make(map[string]*model.Note),
	}
}

func (s *NoteStore) Create(note *model.Note) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notes[note.ID] = note
}

func (s *NoteStore) FindActive() []*model.Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*model.Note
	for _, n := range s.notes {
		if n.Status == "active" {
			result = append(result, n)
		}
	}
	return result
}

func (s *NoteStore) MarkDone(id string) (*model.Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	note, ok := s.notes[id]
	if !ok {
		return nil, errors.New("note not found")
	}
	note.Status = "done"
	return note, nil
}
