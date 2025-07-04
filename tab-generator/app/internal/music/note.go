package music

import (
	"fmt"
	"slices"

	"github.com/er-davo/guitar"
)

var noteNames = [...]string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

func MidiToNote(pitch int) (string, int) {
	if pitch < 0 || pitch > 127 {
		return "Invalid pitch", -1
	}

	note := noteNames[pitch%12]
	octave := pitch / 12

	return note, octave
}

type NoteEvent struct {
	Name      string
	Octave    int
	MidiPitch int
	StartTime float32
	Velocity  float32
}

type NoteSequence struct {
	Notes []NoteEvent
}

func NewNoteSequence(size int) NoteSequence {
	return NoteSequence{Notes: make([]NoteEvent, size)}
}

func (n *NoteSequence) Append(notes ...NoteEvent) {
	n.Notes = append(n.Notes, notes...)
}

func (n *NoteSequence) InsertInto(index int, notes ...NoteEvent) error {
	if index < 0 || index > len(n.Notes) {
		return fmt.Errorf("invalid index: %d", index)
	}
	n.Notes = append(n.Notes[:index], append(notes, n.Notes[index:]...)...)
	return nil
}

func (n *NoteSequence) Merge(seq NoteSequence) {
	if len(seq.Notes) == 0 {
		return
	}
	if len(n.Notes) == 0 {
		n.Append(seq.Notes...)
		return
	}

	const timeEps float32 = 0.01

	overlap := 0
	for overlap < len(seq.Notes) &&
		n.Notes[len(n.Notes)-1].StartTime-seq.Notes[overlap].StartTime > timeEps {
		overlap++
	}

	n.Append(seq.Notes[overlap:]...)
}

func (n *NoteSequence) guitarSequence() [][]guitar.Playable {
	frames := [][]guitar.Playable{}

	slices.SortFunc(n.Notes, func(left, right NoteEvent) int {
		if left.StartTime < right.StartTime {
			return -1
		}
		if left.StartTime > right.StartTime {
			return 1
		}
		return 0
	})

	var timeIndex int
	const eps = 0.01
	for timeIndex < len(n.Notes) {
		curTime := n.Notes[timeIndex].StartTime
		frame := []guitar.Playable{}

		for timeIndex < len(n.Notes) && n.Notes[timeIndex].StartTime-curTime < eps {
			note := n.Notes[timeIndex]
			frame = append(frame, guitar.Note{
				Name:   note.Name,
				Octave: note.Octave,
				Time:   note.StartTime,
			})
			timeIndex++
		}

		frames = append(frames, frame)
	}

	return frames
}
