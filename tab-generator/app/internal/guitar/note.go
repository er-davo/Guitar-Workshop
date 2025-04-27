package guitar

import (
	"errors"
	"fmt"
	"math"
	"slices"
)

const (
	fretDistanceThreshold   = 4.0
	stringDistanceThreshold = 3.0
)

var notesChromo = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

type Note struct {
	Note   string
	Octave int

	Fret   int
	String int

	Time float32
}

func (n *Note) AddFret() error {
	found := -1

	for i := range len(notesChromo) {
		if n.Note == notesChromo[i] {
			found = i
			break
		}
	}

	if found == -1 {
		return fmt.Errorf("invalid note: %s", n.Note)
	}

	n.Note = notesChromo[(found+1)%len(notesChromo)]

	if found == len(notesChromo)-1 {
		n.Octave++
	}

	n.Fret++
	return nil
}

type Notes []Note

func (n *Notes) ClosestTo(target Note) (Note, error) {
	if len(*n) == 0 {
		return Note{}, errors.New("empty notes list")
	}

	minStringDistance := math.MaxFloat32
	minFretDistance := math.MaxFloat32
	closest := Note{}

	for _, candidate := range *n {
		curStringDistance := math.Abs(float64(candidate.String - target.String))
		curFretDistance := math.Abs(float64(candidate.Fret - target.Fret))

		stringAndFret := curStringDistance < float64(minStringDistance) &&
			curFretDistance < float64(minFretDistance)
		stringOnly := curStringDistance < float64(minStringDistance) &&
			curFretDistance-float64(minFretDistance) <= fretDistanceThreshold
		fretOnly := curStringDistance-float64(minStringDistance) <= stringDistanceThreshold &&
			curFretDistance < float64(minFretDistance)

		if stringAndFret || stringOnly || fretOnly {
			closest = candidate
			minFretDistance = curFretDistance
			minStringDistance = curStringDistance
		}
	}

	closest.Time = target.Time

	return closest, nil
}

func noteIsValid(n Note) bool {
	return slices.Contains(notesChromo, n.Note)
}
