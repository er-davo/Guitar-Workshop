package service

import (
	"context"
	"fmt"

	"tabgen/internal/clients"
	"tabgen/internal/music"
	note_analyzer "tabgen/internal/proto/note-analyzer"
	"tabgen/internal/proto/tab"

	"go.uber.org/zap"
)

type TabService struct {
	analyzer clients.NoteAnalyzer
	log      *zap.Logger
}

func NewTabService(analyzer clients.NoteAnalyzer, log *zap.Logger) *TabService {
	return &TabService{
		analyzer: analyzer,
		log:      log,
	}
}

func (s *TabService) GenerateTab(ctx context.Context, req *tab.TabRequest) (*tab.TabResponse, error) {
	AudioReq := &note_analyzer.AudioRequest{
		AudioData: &note_analyzer.AudioFileData{
			FileName:   req.Audio.FileName,
			AudioBytes: req.Audio.AudioBytes,
		},
	}
	s.log.Debug("analyzing for audio", zap.Int("size", len(AudioReq.AudioData.AudioBytes)))
	notes, err := s.analyzer.Analyze(ctx, AudioReq)
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
			EndTime:   note.DurationSeconds,
			Velocity:  note.Velocity,
		}
	}

	processedSeq := seq.Processed()

	noti := ""

	for _, note := range processedSeq.Notes {
		noti += fmt.Sprintf("%+v\n", note)
	}

	tabs, err := music.GenerateTab(*processedSeq)
	if err != nil {
		return nil, err
	}

	return &tab.TabResponse{Tab: tabs}, nil
}
