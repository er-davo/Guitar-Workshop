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

func GenerateTab(audio *audiopb.AudioResponse) (string, error) {
	events := newNotesEvent(audio)
	fb := guitar.NewFingerBoard(guitar.StandartTuning, 24)
	builder, err := guitar.NewTabBuilder(guitar.GuitarType, *fb.GetTuningNotes())
	if err != nil {
		// TODO
	}

	for _, event := range events.notes {
		notes := fb.GetNotes(event.mainNote, event.octave)
		note := notes.ClosestTo(guitar.Note{
			Note:   event.mainNote,
			Octave: event.octave,
			Time:   event.time,
		})
		builder.WriteSingleNote(note)
	}

	return builder.Tab(), nil
}
