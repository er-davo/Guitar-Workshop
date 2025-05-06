package guitar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddFret(t *testing.T) {
	testCases := []struct {
		name        string
		input       Note
		expected    Note
		expectError bool
	}{
		{
			name:     "regular increment",
			input:    Note{Name: "F", Octave: 2, Fret: 0},
			expected: Note{Name: "F#", Octave: 2, Fret: 1},
		},
		{
			name:     "octave change",
			input:    Note{Name: "B", Octave: 2, Fret: 0},
			expected: Note{Name: "C", Octave: 3, Fret: 1},
		},
		{
			name:     "G# to A",
			input:    Note{Name: "G#", Octave: 3, Fret: 11},
			expected: Note{Name: "A", Octave: 3, Fret: 12},
		},
		{
			name:        "invalid note",
			input:       Note{Name: "X", Octave: 1},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n := tc.input
			err := n.AddFret()

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, n)
		})
	}
}

func TestCalculateScore(t *testing.T) {
	testCases := []struct {
		name          string
		note          Note
		target        Note
		expectedScore float64
	}{
		{
			name:          "exact match",
			note:          Note{Fret: 8, String: 0},
			target:        Note{Fret: 8, String: 0},
			expectedScore: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.note.calculateScore(tc.target)
			assert.Equal(t, tc.expectedScore, actual, "expected %f, found %f", tc.expectedScore, actual)
		})
	}
}

func TestClosestTo(t *testing.T) {
	notes := Notes{
		{Name: "E", Octave: 2, Fret: 0, String: 5},
		{Name: "E", Octave: 2, Fret: 2, String: 2},
		{Name: "E", Octave: 2, Fret: 6, String: 4},
		{Name: "E", Octave: 2, Fret: 8, String: 0},
	}

	testCases := []struct {
		name     string
		target   Note
		expected Note
	}{
		{
			name:     "exact match",
			target:   Note{Fret: 8, String: 0},
			expected: notes[3],
		},
		{
			name:     "prefer open string",
			target:   Note{Fret: 0, String: 3},
			expected: notes[0], // Should prefer open E string
		},
		{
			name:     "closest fret distance",
			target:   Note{Fret: 1, String: 5},
			expected: notes[0],
		},
		{
			name:     "tie breaker with string distance",
			target:   Note{Fret: 2, String: 3},
			expected: notes[1], // Closer string distance
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			closest, err := notes.ClosestTo(tc.target)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, closest)
		})
	}
}
