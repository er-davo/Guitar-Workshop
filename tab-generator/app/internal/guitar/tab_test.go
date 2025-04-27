package guitar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteSingleNote(t *testing.T) {
	note := Note{
		Note:   "E",
		Octave: 2,
		Fret:   12,
		String: 5,
		Time:   4,
	}
	tun, _ := GetTuning(StandartTuning)
	tb, _ := NewTabBuilder(GuitarType, tun.NoteNames())
	err := tb.WriteSingleNote(note)

	assert.NoError(t, err, "got error")

	assert.Equal(t, "e|----------------------\nB|----------------------\nG|----------------------\nD|----------------------\nA|----------------------\nE|--------------------12\n",
		tb.Tab(), "something wrong",
	)
}
