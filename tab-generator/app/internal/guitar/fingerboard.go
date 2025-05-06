package guitar

import (
	"errors"
	"strings"
)

type TuningType int

const (
	StandartTuning TuningType = iota
)

type Tuning []Note

func (t *Tuning) NoteNames() []string {
	names := make([]string, len(*t))
	for i := range *t {
		names[i] = (*t)[i].Name
	}
	names[0] = strings.ToLower(names[0])
	return names
}

func GetTuning(t TuningType) (Tuning, error) {
	switch t {
	case StandartTuning:
		return Tuning{
			{
				Name:   "E",
				Octave: 4,
				Fret:   0,
				String: 0,
			},
			{
				Name:   "B",
				Octave: 3,
				Fret:   0,
				String: 1,
			},
			{
				Name:   "G",
				Octave: 3,
				Fret:   0,
				String: 2,
			},
			{
				Name:   "D",
				Octave: 3,
				Fret:   0,
				String: 3,
			},
			{
				Name:   "A",
				Octave: 2,
				Fret:   0,
				String: 4,
			},
			{
				Name:   "E",
				Octave: 2,
				Fret:   0,
				String: 5,
			},
		}, nil
	default:
		return Tuning{}, errors.ErrUnsupported
	}
}

type FingerBoard struct {
	tuning Tuning
	frets  int
}

func NewFingerBoard(tuningType TuningType, frets int) (*FingerBoard, error) {
	if frets < 0 {
		return nil, errors.New("frets value can not be negative")
	}

	tun, err := GetTuning(tuningType)
	if err != nil {
		return nil, err
	}

	return &FingerBoard{
		tuning: tun,
		frets:  frets,
	}, nil
}

func (fb *FingerBoard) GetTuningNotes() []string {
	return fb.tuning.NoteNames()
}

// TODO
// fix error
// returns empty list for # notes
func (fb FingerBoard) GetNotes(targetNote string, targetOctave int) Notes {
	notes := Notes{}
	currentNote := Note{}

	for i := range fb.tuning {
		currentNote = fb.tuning[i]

		for fret := 0; fret < fb.frets; fret++ {
			if currentNote.Name == targetNote && currentNote.Octave == targetOctave {
				notes = append(notes, currentNote)
			}
			currentNote.AddFret()
		}
	}

	return notes
}

/*
e|----------0-------
B|--------0---0-----
G|------0-------0---
D|------------------
A|------------------
E|----0-----------0-

e|-----------------------------------------------------------
B|------------------------0--0-------------------------------
G|---------------------------------0-0----------------------0
D|-----------------------------2-1---------------------------
A|-----------------------------------------------------------
E|-----------------0--0----------------4-0---0-------0-0--0--
*/
