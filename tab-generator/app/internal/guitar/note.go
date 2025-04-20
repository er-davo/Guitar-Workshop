package guitar

import "math"

var notesChromo = [12]string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

type Note struct {
	Note   string
	Octave int
	Fret   int
	String int
	Time   float32
}

func (n *Note) AddFret() *Note {
	for i := range 12 {
		if n.Note == notesChromo[i] {
			if i == 11 {
				n.Note = notesChromo[0]
				n.Octave++
			} else {
				n.Note = notesChromo[i+1]
			}
		}
	}

	n.Fret++
	return n
}

type Notes []Note

func (n *Notes) ClosestTo(note Note) Note {
	minStringDistanse := 100
	minFretDistanse := 100
	notePos := Note{}

	for _, possibleNote := range *n {
		curStringDistanse := math.Abs(float64(possibleNote.String - note.Octave))
		curFretDistanse := math.Abs(float64(possibleNote.Fret - note.Fret))

		stringAndFret := curStringDistanse < float64(minStringDistanse) &&
			curFretDistanse < float64(minFretDistanse)
		stringOnly := curStringDistanse < float64(minStringDistanse) &&
			curFretDistanse-float64(minFretDistanse) <= 4
		fretOnly := curStringDistanse-float64(minStringDistanse) <= 3 &&
			curFretDistanse < float64(minFretDistanse)

		if stringAndFret || stringOnly || fretOnly {
			notePos = possibleNote
			minFretDistanse = int(curFretDistanse)
			minStringDistanse = int(curStringDistanse)
		}
	}

	return notePos
}
