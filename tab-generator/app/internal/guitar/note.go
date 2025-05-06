package guitar

import (
	"errors"
	"fmt"
	"math"
	"slices"
)

const (
	stringlWeight   = 1.0  // Предпочтение вертикальным перемещениям (по струнам)
	fretWeight      = 1.0  // Вес для перемещений по ладам
	openStringBonus = -2.0 // Бонус за открытые струны
)

var notesChromo = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

type Note struct {
	Name   string
	Octave int

	Fret   int
	String int

	Time float32
}

func (n *Note) AddFret() error {
	found := -1

	for i := range len(notesChromo) {
		if n.Name == notesChromo[i] {
			found = i
			break
		}
	}

	if found == -1 {
		return fmt.Errorf("invalid note: %s", n.Name)
	}

	n.Name = notesChromo[(found+1)%len(notesChromo)]

	if found == len(notesChromo)-1 {
		n.Octave++
	}

	n.Fret++
	return nil
}

func (n *Note) calculateScore(target Note) float64 {
	// Расстояние по горизонтали (лады)
	fretDist := math.Abs(float64(n.Fret - target.Fret))

	// Расстояние по вертикали (строки)
	stringDist := math.Abs(float64(n.String - target.String))

	// Бонус за открытые струны
	openString := 0.0
	if n.Fret == 0 {
		openString = openStringBonus
	}

	score := (stringDist * stringlWeight) +
		(fretDist * fretWeight) +
		openString

	return score
}

type Notes []Note

func (n *Notes) ClosestTo(target Note) (Note, error) {
	if len(*n) == 0 {
		return Note{}, errors.New("empty notes list")
	}

	closest := Note{}
	minScore := math.MaxFloat64

	for _, candidate := range *n {
		currentScore := candidate.calculateScore(target)

		if currentScore < minScore {
			minScore = currentScore
			closest = candidate
		}
	}

	closest.Time = target.Time

	return closest, nil
}

func noteIsValid(n Note) bool {
	return slices.Contains(notesChromo, n.Name)
}
