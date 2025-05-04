package guitar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClosestTo(t *testing.T) {
	notes := Notes{
		{Note: "E", Octave: 2, Fret: 0, String: 5},
		{Note: "E", Octave: 2, Fret: 2, String: 2},
		{Note: "E", Octave: 2, Fret: 6, String: 4},
		{Note: "E", Octave: 2, Fret: 8, String: 0},
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
			name:     "closest fret distance",
			target:   Note{Fret: 1, String: 5},
			expected: notes[0],
		},
		// {
		// 	name:     "prefer open string",
		// 	target:   Note{Fret: 0, String: 3},
		// 	expected: notes[0], // Should prefer open E string
		// },
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

func TestAddFret(t *testing.T) {
	testCases := []struct {
		name        string
		input       Note
		expected    Note
		expectError bool
	}{
		{
			name:     "regular increment",
			input:    Note{Note: "F", Octave: 2, Fret: 0},
			expected: Note{Note: "F#", Octave: 2, Fret: 1},
		},
		{
			name:     "octave change",
			input:    Note{Note: "B", Octave: 2, Fret: 0},
			expected: Note{Note: "C", Octave: 3, Fret: 1},
		},
		{
			name:     "G# to A",
			input:    Note{Note: "G#", Octave: 3, Fret: 11},
			expected: Note{Note: "A", Octave: 3, Fret: 12},
		},
		{
			name:        "invalid note",
			input:       Note{Note: "X", Octave: 1},
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
