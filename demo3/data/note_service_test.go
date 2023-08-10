package data_test

import (
	"testing"

	"github.com/addetz/secure-code-go/demo3/data"
	"github.com/stretchr/testify/assert"
)

func TestAddNote(t *testing.T) {
	notes := data.NewSecretNoteService()
	user := "user1"
	note := data.SecretNote{
		ID:   "note-1",
		Text: "My Secret Note",
	}
	otherNote := data.SecretNote{
		ID:   "note-2",
		Text: "My Other Secret Note",
	}
	notes.Add(user, note)
	userNotes, err := notes.GetAll(user)
	assert.Nil(t, err)
	assert.NotNil(t, userNotes)
	assert.Len(t, userNotes, 1)
	assert.Equal(t, note, userNotes[0])
	notes.Add(user, otherNote)
	userNotes, err = notes.GetAll(user)
	assert.Nil(t, err)
	assert.NotNil(t, userNotes)
	assert.Len(t, userNotes, 2)
	assert.Equal(t, note, userNotes[0])
	assert.Equal(t, otherNote, userNotes[1])
}

func TestGetAllNotes(t *testing.T) {
	t.Run("no notes found", func(t *testing.T) {
		noteService := data.NewSecretNoteService()
		notes, err := noteService.GetAll("user1")
		assert.Nil(t, notes)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no notes found")
	})
}
