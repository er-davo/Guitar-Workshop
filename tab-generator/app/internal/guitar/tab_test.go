package guitar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteSingleNote(t *testing.T) {
	tun, _ := GetTuning(StandartTuning)
	tuningNotes := tun.NoteNames()

	testCases := []struct {
		name        string
		notes       []Note
		expectedTab string
		expectError bool
	}{
		{
			name: "single note on E",
			notes: []Note{
				{Name: "E", Octave: 2, Fret: 12, String: 5, Time: 0},
			},
			expectedTab: "e|---\nB|---\nG|---\nD|---\nA|---\nE|12-\n",
		},
		{
			name: "multiple notes with timing",
			notes: []Note{
				{Name: "E", Fret: 0, String: 5, Time: 0},
				{Name: "B", Fret: 1, String: 1, Time: 0.2},
				{Name: "G", Fret: 3, String: 2, Time: 0.4},
			},
			expectedTab: "e|------\nB|--1---\nG|----3-\nD|------\nA|------\nE|0-----\n",
		},
		{
			name: "invalid note time",
			notes: []Note{
				{Name: "E", Time: 0.5},
				{Name: "B", Time: 0.3},
			},
			expectError: true,
		},
		{
			name: "two-digit fret formatting",
			notes: []Note{
				{Name: "E", Fret: 10, String: 5, Time: 0},
			},
			expectedTab: "e|---\nB|---\nG|---\nD|---\nA|---\nE|10-\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tb, _ := NewTabBuilder(GuitarType, tuningNotes)

			var err error
			for _, n := range tc.notes {
				err = tb.WriteSingleNote(n)
				if err != nil {
					break
				}
			}

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedTab, tb.Tab())
		})
	}
}
