package data_test

import (
	"errors"
	"testing"

	"github.com/addetz/secure-code-go/demo4/data"
	"github.com/addetz/secure-code-go/demo4/db"
	"github.com/addetz/secure-code-go/demo4/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddNote(t *testing.T) {
	mockDB := new(mocks.DatabaseServiceMock)
	notes := data.NewSecretNoteService(mockDB)
	user := "user1"
	note := data.SecretNote{
		Text: "My Secret Note",
	}
	mockDB.On("AddNote", mock.AnythingOfType("string"), user, note.Text).Return(nil)
	err := notes.Add(user, note)
	assert.Nil(t, err)
}

func TestGetAllNotes(t *testing.T) {
	t.Run("no notes found", func(t *testing.T) {
		user := "user1"
		mockDB := new(mocks.DatabaseServiceMock)
		noteService := data.NewSecretNoteService(mockDB)
		mockDB.On("GetUserNotes", user).Return(nil, errors.New("no notes found"))
		notes, err := noteService.GetAll(user)
		assert.Nil(t, notes)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no notes found")
	})
	t.Run("notes found", func(t *testing.T) {
		user := "user1"
		dbNotes := []db.Note{
			{
				ID: uuid.New().String(),
				Username: user,
				Text: "My first note",
			},
			{
				ID: uuid.New().String(),
				Username: user,
				Text: "My second note",
			},
		}
		mockDB := new(mocks.DatabaseServiceMock)
		noteService := data.NewSecretNoteService(mockDB)
		mockDB.On("GetUserNotes", user).Return(dbNotes, nil)
		notes, err := noteService.GetAll(user)
		assert.Nil(t, err)
		assert.NotNil(t, notes)
		assert.Equal(t, len(dbNotes), len(notes))
		for i, n := range dbNotes {
			assert.Equal(t, n.ID, notes[i].ID)
			assert.Equal(t, n.Username, notes[i].Username)
			assert.Equal(t, n.Text, notes[i].Text)

		}
	})
}
