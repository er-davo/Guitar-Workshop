package models

import (
	audiopb "tabgen/internal/audioproto"
	"tabgen/internal/logger"

	"github.com/er-davo/guitar"

	"go.uber.org/zap"
)

func GenerateTab(audio *audiopb.AudioResponse) (string, error) {
	events := newNotesEvent(audio)
	tun, err := guitar.ParseTuning(guitar.StandardTuning)
	fb, err := guitar.NewFingerBoard(tun, 24)
	if err != nil {
		logger.Error("failed to create FingerBoard", zap.Error(err))
		return "", err
	}

	builder, err := guitar.NewTabBuilder(tun.NoteNames())
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
		builder.WriteNotes(note)
	}

	return builder.Tab(), nil
}
