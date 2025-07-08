package music

import (
	"fmt"
	"math"
	"slices"

	"github.com/er-davo/guitar"
)

var noteNames = [...]string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

func MidiToNote(pitch int) (string, int) {
	if pitch < 0 || pitch > 127 {
		return "Invalid pitch", -1
	}

	note := noteNames[pitch%12]
	octave := pitch/12 - 1

	return note, octave
}

type NoteEvent struct {
	Name      string
	Octave    int
	MidiPitch int
	StartTime float64
	EndTime   float64
	Velocity  float64
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

	const timeEps = 0.01

	overlap := 0
	for overlap < len(seq.Notes) &&
		n.Notes[len(n.Notes)-1].StartTime-seq.Notes[overlap].StartTime > timeEps {
		overlap++
	}

	n.Append(seq.Notes[overlap:]...)
}

func (n *NoteSequence) Sort() {
	slices.SortFunc(n.Notes, func(left, right NoteEvent) int {
		if left.StartTime < right.StartTime {
			return -1
		}
		if left.StartTime > right.StartTime {
			return 1
		}
		return 0
	})
}

func (n *NoteSequence) MergeRepeatedNotes() *NoteSequence {
	if len(n.Notes) < 2 {
		return n
	}

	merged := NoteSequence{
		Notes: make([]NoteEvent, 0, len(n.Notes)),
	}
	merged.Append(n.Notes[0])

	velocityDiffThreshold := 0.2
	timeDiffThreshold := 0.05

	for i := 1; i < len(n.Notes); i++ {
		last := &merged.Notes[len(merged.Notes)-1]

		midiIsEqual := last.MidiPitch == n.Notes[i].MidiPitch
		timeIsAlmostEqual := last.EndTime-
			n.Notes[i].StartTime < timeDiffThreshold
		velocityIsAlmostEqual := math.Abs(last.Velocity-
			n.Notes[i].Velocity) < velocityDiffThreshold

		if midiIsEqual && timeIsAlmostEqual && velocityIsAlmostEqual {
			last.EndTime = n.Notes[i].EndTime
			last.Velocity = (last.Velocity + n.Notes[i].Velocity) / 2
			continue
		}

		merged.Append(n.Notes[i])
	}

	return &merged
}

// octave diff: 1, duration: 0.15, velocity: 0.43
// octave diff: 3, duration: 0.31, velocity: 0.34
// octave diff: 3, duration: 0.15, velocity: 0.41
// octave diff: 3, duration: 0.35, velocity: 0.34
//

func abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}

func (n *NoteSequence) RemoveNoisyNotes() *NoteSequence {
	if len(n.Notes) < 2 {
		return n
	}

	octaveDiffThreshold := 2
	velocityThreshold := 0.5
	durationThreshold := 0.35

	seq := NoteSequence{
		Notes: make([]NoteEvent, 0, len(n.Notes)),
	}
	seq.Append(n.Notes[0])

	for i := 1; i < len(n.Notes); i++ {
		last := seq.Notes[len(seq.Notes)-1]
		note := &n.Notes[i]
		noiseValue := 0

		if abs(last.Octave-note.Octave) > octaveDiffThreshold {
			noiseValue++
		}
		if note.EndTime-note.StartTime < durationThreshold {
			noiseValue++
		}
		if note.Velocity < velocityThreshold {
			noiseValue++
		}

		if noiseValue < 2 {
			seq.Append(*note)
		}
	}

	return &seq
}

func (n *NoteSequence) Processed() *NoteSequence {
	n.Sort()
	return n.RemoveNoisyNotes().MergeRepeatedNotes()
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
