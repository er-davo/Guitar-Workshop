package music

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMidiToNote(t *testing.T) {
	testCases := []struct {
		pitch      int
		wantNote   string
		wantOctave int
	}{
		{60, "C", 4},
		{61, "C#", 4},
		{0, "C", -1},
		{127, "G", 9},
		{-1, "Invalid pitch", -1},
		{128, "Invalid pitch", -1},
	}

	for _, tc := range testCases {
		note, oct := MidiToNote(tc.pitch)
		require.Equal(t, tc.wantNote, note)
		require.Equal(t, tc.wantOctave, oct)
	}
}

func TestNoteSequence_InsertInto(t *testing.T) {
	n := NewNoteSequence(0)

	n.Append(
		NoteEvent{Name: "C", Octave: 4, StartTime: 0},
		NoteEvent{Name: "D", Octave: 4, StartTime: 1},
	)

	toInsert := NoteEvent{Name: "E", Octave: 4, StartTime: 0.5}

	err := n.InsertInto(1, toInsert)
	require.NoError(t, err)
	require.Equal(t, 3, len(n.Notes))
	require.Equal(t, "C", n.Notes[0].Name)
	require.Equal(t, "E", n.Notes[1].Name)
	require.Equal(t, "D", n.Notes[2].Name)

	err = n.InsertInto(0, NoteEvent{Name: "B", Octave: 3})
	require.NoError(t, err)
	require.Equal(t, "B", n.Notes[0].Name)

	err = n.InsertInto(len(n.Notes), NoteEvent{Name: "F", Octave: 4})
	require.NoError(t, err)
	require.Equal(t, "F", n.Notes[len(n.Notes)-1].Name)

	err = n.InsertInto(-1, NoteEvent{Name: "X"})
	require.Error(t, err)

	err = n.InsertInto(len(n.Notes)+1, NoteEvent{Name: "Y"})
	require.Error(t, err)
}

func TestNoteSequence_Merge(t *testing.T) {
	a := NewNoteSequence(0)
	a.Append(NoteEvent{StartTime: 1.0})

	b := NewNoteSequence(0)
	b.Append(
		NoteEvent{StartTime: 0.5},
		NoteEvent{StartTime: 1.01},
	)

	a.Merge(b)

	require.Equal(t, 2, len(a.Notes))
	require.Equal(t, 1.0, a.Notes[0].StartTime)
	require.Equal(t, 1.01, a.Notes[1].StartTime)
}

func TestNoteSequence_Sort(t *testing.T) {
	n := NewNoteSequence(0)
	n.Append(
		NoteEvent{StartTime: 2.0},
		NoteEvent{StartTime: 0.5},
		NoteEvent{StartTime: 1.0},
	)

	n.Sort()

	require.Equal(t, 0.5, n.Notes[0].StartTime)
	require.Equal(t, 1.0, n.Notes[1].StartTime)
	require.Equal(t, 2.0, n.Notes[2].StartTime)
}

func TestNoteSequence_MergeRepeatedNotes(t *testing.T) {
	n := NewNoteSequence(0)
	n.Append(
		NoteEvent{MidiPitch: 60, StartTime: 0, EndTime: 0.5, Velocity: 0.8},
		NoteEvent{MidiPitch: 60, StartTime: 0.51, EndTime: 1.0, Velocity: 0.9}, // almost same
		NoteEvent{MidiPitch: 62, StartTime: 1.1, EndTime: 1.6, Velocity: 0.8},  // diff
	)

	merged := n.MergeRepeatedNotes()

	require.Equal(t, 2, len(merged.Notes))

	require.Equal(t, 60, merged.Notes[0].MidiPitch)
	require.InEpsilon(t, 0.85, merged.Notes[0].Velocity, 0.01)
	require.Equal(t, 0.0, merged.Notes[0].StartTime)
	require.Equal(t, 1.0, merged.Notes[0].EndTime)

	require.Equal(t, 62, merged.Notes[1].MidiPitch)
}

func TestNoteSequence_RemoveNoisyNotes(t *testing.T) {
	n := NewNoteSequence(0)
	n.Append(
		NoteEvent{Octave: 4, StartTime: 0, EndTime: 0.5, Velocity: 0.9},
		NoteEvent{Octave: 8, StartTime: 0.6, EndTime: 0.7, Velocity: 0.4},
		NoteEvent{Octave: 5, StartTime: 0.8, EndTime: 1.2, Velocity: 0.9},
	)

	cleaned := n.RemoveNoisyNotes()

	require.Equal(t, 2, len(cleaned.Notes))
	require.Equal(t, 4, cleaned.Notes[0].Octave)
	require.Equal(t, 5, cleaned.Notes[1].Octave)
}
