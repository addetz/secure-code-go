package data

import (
	"github.com/addetz/secure-code-go/demo4/db"
	"github.com/google/uuid"
)

type SecretNote struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

// SecretNoteService maintains the user notes.
type SecretNoteService struct {
	dbService db.DatabaseService
}

// NewSecretNoteService creates a SecretNoteService that is ready to use.
func NewSecretNoteService(dbService db.DatabaseService) *SecretNoteService {
	return &SecretNoteService{
		dbService: dbService,
	}
}

// Add adds a new SecretNote for the given user by using the SecretNoteService.
func (ns *SecretNoteService) Add(user string, n SecretNote) error {
	id := uuid.New().String()
	return ns.dbService.AddNote(id, user, n.Text)
}

// Get returns all the SecretNotes of a given user by using the SecretNoteService.
func (ns *SecretNoteService) GetAll(user string) ([]SecretNote, error) {
	dbNotes, err := ns.dbService.GetUserNotes(user)
	if err != nil {
		return nil, err
	}
	var notes []SecretNote
	for _, n := range dbNotes {
		notes = append(notes, SecretNote{
			ID:       n.ID,
			Username: n.Username,
			Text:     n.Text,
		})
	}

	return notes, nil
}
