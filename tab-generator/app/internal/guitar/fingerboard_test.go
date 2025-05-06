package guitar

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFingerBoard(t *testing.T) {
	t.Run("valid frets", func(t *testing.T) {
		fb, err := NewFingerBoard(StandartTuning, 24)
		assert.NoError(t, err)
		assert.Equal(t, 24, fb.frets)
	})

	t.Run("negative frets", func(t *testing.T) {
		_, err := NewFingerBoard(StandartTuning, -5)
		assert.ErrorContains(t, err, "frets value can not be negative")
	})

	t.Run("invalid tuning", func(t *testing.T) {
		_, err := NewFingerBoard(TuningType(99), 12)
		assert.ErrorIs(t, err, errors.ErrUnsupported)
	})
}

func TestGetNotes(t *testing.T) {
	fb, _ := NewFingerBoard(StandartTuning, 24)

	testCases := []struct {
		name         string
		targetNote   string
		targetOctave int
		expected     []Note
	}{
		{
			name:         "A2",
			targetNote:   "A",
			targetOctave: 2,
			expected: []Note{
				{Name: "A", Octave: 2, Fret: 0, String: 4},
				{Name: "A", Octave: 2, Fret: 5, String: 5},
			},
		},
		{
			name:         "C#3",
			targetNote:   "C#",
			targetOctave: 3,
			expected: []Note{
				{Name: "C#", Octave: 3, Fret: 4, String: 4},
				{Name: "C#", Octave: 3, Fret: 9, String: 5},
			},
		},
		{
			name:         "F4",
			targetNote:   "F",
			targetOctave: 4,
			expected: []Note{
				{Name: "F", Octave: 4, Fret: 1, String: 0},
				{Name: "F", Octave: 4, Fret: 6, String: 1},
				{Name: "F", Octave: 4, Fret: 10, String: 2},
				{Name: "F", Octave: 4, Fret: 15, String: 3},
				{Name: "F", Octave: 4, Fret: 20, String: 4},
			},
		},
		{
			name:         "non-existent note",
			targetNote:   "H",
			targetOctave: 2,
			expected:     []Note{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			notes := fb.GetNotes(tc.targetNote, tc.targetOctave)

			assert.Equal(t, len(tc.expected), len(notes), "unexpected number of notes")
			for _, expectedNote := range tc.expected {
				assert.True(t, containsNote(notes, expectedNote),
					"expected note %+v not found", expectedNote)
			}
		})
	}
}

func containsNote(notes Notes, target Note) bool {
	for _, n := range notes {
		if n.Name == target.Name &&
			n.Octave == target.Octave &&
			n.Fret == target.Fret &&
			n.String == target.String {
			return true
		}
	}
	return false
}
