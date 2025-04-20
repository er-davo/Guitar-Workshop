package guitar

type Tuning []Note

func (t *Tuning) String() *[]string {
	notes := []string{}
	for _, note := range *t {
		notes = append(notes, note.Note)
	}
	return &notes
}

var StandartTuning = Tuning{
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
}

type FingerBoard struct {
	tuning Tuning
	frets  int
}

func NewFingerBoard(tuning Tuning, frets int) *FingerBoard {
	return &FingerBoard{
		tuning: tuning,
		frets:  frets,
	}
}

func (fb *FingerBoard) GetNotes(note string, octave int) Notes {
	notes := Notes{}
	tuningStrings := fb.tuning

	for i := range tuningStrings {
		for range fb.frets + 1 {
			if tuningStrings[i].Note == note && tuningStrings[i].Octave == octave {
				notes = append(notes, tuningStrings[i])
			}
			tuningStrings[i].AddFret()
		}
	}

	return notes
}
