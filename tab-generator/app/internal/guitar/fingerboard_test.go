package guitar

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFingerBoard(t *testing.T) {
	_, err := NewFingerBoard(StandartTuning, 24)

	assert.NoError(t, err)
}

func TestGetNotes(t *testing.T) {
	fb, _ := NewFingerBoard(StandartTuning, 24)

	notes := fb.GetNotes("F#", 4)

	assert.True(t, slices.Contains(notes, Note{
		Note:   "E",
		Octave: 2,
		Fret:   0,
		String: 5,
		Time:   0,
	}))

	assert.Equal(t, Notes{}, notes)
}
