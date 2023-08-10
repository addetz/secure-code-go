package data

import (
	"errors"
	"fmt"
)

type SecretNote struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// SecretNoteService maintains the user notes.
type SecretNoteService struct {
	notes map[string][]SecretNote
}

// NewSecretNoteService creates a SecretNoteService that is ready to use.
func NewSecretNoteService() *SecretNoteService {
	return &SecretNoteService{
		notes: make(map[string][]SecretNote),
	}
}

// Add adds a new SecretNote for the given user by using the SecretNoteService.
func (ns *SecretNoteService) Add(user string, n SecretNote) {
	existing := ns.notes[user]
	n.ID = fmt.Sprintf("note-%d", len(existing)+1)
	existing = append(existing, n)
	ns.notes[user] = existing
}

// Get returns all the SecretNotes of a given user by using the SecretNoteService.
func (ns *SecretNoteService) GetAll(user string) ([]SecretNote, error) {
	existing, ok := ns.notes[user]
	if !ok {
		return nil, errors.New("no notes found for given user")
	}

	return existing, nil
}
