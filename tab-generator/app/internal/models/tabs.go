package models

import (
	audiopb "tabgen/internal/audioproto"
	"tabgen/internal/guitar"
)

type TabRequest struct {
	AudioURL string `json:"audio_url"`
}

type TabResponse struct {
	Tab    string `json:"tab"`
	Status string `json:"status"`
}

func newNote(n *audiopb.AudioEvent) *guitar.Note {
	return &guitar.Note{
		Note:   n.MainNote,
		Octave: int(n.Octave),
		Time:   n.Time,
	}
}

func GenerateTab(audio *audiopb.AudioResponse) (string, error) {
	events := newNotesEvent(audio)
	fb := guitar.NewFingerBoard(guitar.StandartTuning, 24)
	builder, err := guitar.NewTabBuilder(guitar.GuitarType)

	for _, event := range events.notes {
		notes := fb.GetNotes(event.mainNote, event.octave)
		note := notes.ClosestTo()

	}

	return events.notes[0].mainNote, nil
}
