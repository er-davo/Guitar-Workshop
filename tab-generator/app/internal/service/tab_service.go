package service

import (
	"context"

	"tabgen/internal/clients"
	"tabgen/internal/logger"
	"tabgen/internal/music"
	note_analyzer "tabgen/internal/proto/note-analyzer"
	"tabgen/internal/proto/tab"

	"go.uber.org/zap"
)

type TabService struct {
	tab.UnimplementedTabGenerateServer
}

func (s *TabService) GenerateTab(ctx context.Context, req *tab.TabRequest) (*tab.TabResponse, error) {
	AudioReq := note_analyzer.AudioRequest{
		AudioData: &note_analyzer.AudioFileData{
			FileName:   req.Audio.FileName,
			AudioBytes: req.Audio.AudioBytes,
		},
	}
	logger.Debug("analyzing for audio", zap.Int("size", len(AudioReq.AudioData.AudioBytes)))
	notes, err := clients.NoteAnalyzerClient.Analyze(context.Background(), &AudioReq)
	if err != nil {
		return nil, err
	}

	seq := music.NewNoteSequence(len(notes.Notes))

	for i, note := range notes.Notes {
		name, octave := music.MidiToNote(int(note.MidiPitch))
		seq.Notes[i] = music.NoteEvent{
			Name:      name,
			Octave:    octave,
			MidiPitch: int(note.MidiPitch),
			StartTime: note.StartSeconds,
			Velocity:  note.Velocity,
		}
	}

	tabs, err := music.GenerateTab(seq)
	if err != nil {
		return nil, err
	}

	return &tab.TabResponse{Tab: tabs}, nil
}
