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
		names[i] = (*t)[i].Note
	}
	names[0] = strings.ToLower(names[0])
	return names
}

func GetTuning(t TuningType) (Tuning, error) {
	switch t {
	case StandartTuning:
		return Tuning{
			{
				Note:   "E",
				Octave: 4,
				Fret:   0,
				String: 0,
			},
			{
				Note:   "B",
				Octave: 3,
				Fret:   0,
				String: 1,
			},
			{
				Note:   "G",
				Octave: 3,
				Fret:   0,
				String: 2,
			},
			{
				Note:   "D",
				Octave: 3,
				Fret:   0,
				String: 3,
			},
			{
				Note:   "A",
				Octave: 2,
				Fret:   0,
				String: 4,
			},
			{
				Note:   "E",
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
		return &FingerBoard{}, err
	}

	return &FingerBoard{
		tuning: tun,
		frets:  frets,
	}, nil
}

func (fb *FingerBoard) GetTuningNotes() []string {
	return fb.tuning.NoteNames()
}

func (fb *FingerBoard) GetNotes(targetNote string, targetOctave int) Notes {
	notes := Notes{}
	tuningStrings := make(Tuning, len(fb.tuning))

	for i := range fb.tuning {
		copy(tuningStrings, fb.tuning)
		for fret := 0; fret < fb.frets+1; fret++ {
			if tuningStrings[i].Note == targetNote && tuningStrings[i].Octave == targetOctave {
				notes = append(notes, tuningStrings[i])
			}
			tuningStrings[i].AddFret()
		}
	}

	return notes
}
