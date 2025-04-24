package guitar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClosestTo(t *testing.T) {
	notes := Notes{
		{
			Note:   "E",
			Octave: 2,
			Fret:   0,
			String: 5,
			Time:   4,
		},
		{
			Note:   "E",
			Octave: 2,
			Fret:   2,
			String: 2,
			Time:   4,
		},
		{
			Note:   "E",
			Octave: 2,
			Fret:   6,
			String: 4,
			Time:   4,
		},
		{
			Note:   "E",
			Octave: 2,
			Fret:   8,
			String: 0,
			Time:   4,
		},
	}

	note := Note{
		Note:   "E",
		Octave: 2,
		Fret:   1,
		String: 0,
		Time:   4,
	}

	closest, err := notes.ClosestTo(note)

	assert.NoError(t, err)

	assert.Equal(t, notes[1], closest)
}

func TestAddFret(t *testing.T) {
	note := Note{
		Note:   "E",
		Octave: 2,
		Fret:   0,
		String: 0,
		Time:   0,
	}

	err := note.AddFret()

	assert.NoError(t, err)

	assert.Equal(t, Note{
		Note:   "F",
		Octave: 2,
		Fret:   1,
		String: 0,
		Time:   0,
	}, note)

	noteToChangeOctave := Note{
		Note:   "B",
		Octave: 2,
		Fret:   0,
		String: 0,
		Time:   0,
	}

	err = noteToChangeOctave.AddFret()

	assert.NoError(t, err)

	assert.Equal(t, Note{
		Note:   "C",
		Octave: 3,
		Fret:   1,
		String: 0,
		Time:   0,
	}, noteToChangeOctave)
}
