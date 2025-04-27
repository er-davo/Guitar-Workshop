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

	notes := fb.GetNotes("A", 2)

	assert.True(t, slices.Contains(notes, Note{
		Note:   "A",
		Octave: 2,
		Fret:   0,
		String: 4,
		Time:   0,
	}))

	notes = fb.GetNotes("C#", 3)

	assert.True(t, slices.Contains(notes, Note{
		Note:   "C#",
		Octave: 3,
		Fret:   4,
		String: 4,
		Time:   0,
	}))

	assert.True(t, len(notes) != 0)
}
