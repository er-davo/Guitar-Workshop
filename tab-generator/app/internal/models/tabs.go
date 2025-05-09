package models

import (
	audiopb "tabgen/internal/audioproto"
	"tabgen/internal/logger"

	"github.com/DavidCage31/guitar"

	"go.uber.org/zap"
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
	fb, err := guitar.NewFingerBoard(guitar.StandartTuning, 24)
	if err != nil {
		logger.Error("failed to create FingerBoard", zap.Error(err))
		return "", err
	}

	builder, err := guitar.NewTabBuilder(guitar.GuitarType, fb.GetTuningNotes())
	if err != nil {
		logger.Error("failed to create tabBuilder", zap.Error(err))
		return "", err
	}

	for _, event := range events.notes {
		notes := fb.GetNotes(event.mainNote, event.octave)
		if len(notes) == 0 {
			logger.Error("got empty notes list", zap.String("note", event.mainNote), zap.Int("octave", event.octave))
			continue
		}

		note, err := notes.ClosestTo(guitar.Note{
			Name:   event.mainNote,
			Octave: event.octave,
			Time:   event.time,
		})
		if err != nil {
			logger.Error("failed to get closest note", zap.Error(err))
			return "", err
		}
		builder.WriteSingleNote(note)
	}

	return builder.Tab(), nil
}
