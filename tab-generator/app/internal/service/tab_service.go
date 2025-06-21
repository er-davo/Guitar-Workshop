package service

import (
	"context"

	"tabgen/internal/clients"
	"tabgen/internal/models"
	"tabgen/internal/music"
	onsets_frames "tabgen/internal/proto/onsets-frames"
	"tabgen/internal/proto/tab"
)

type TabService struct {
	tab.UnimplementedTabGenerateServer
}

func (s *TabService) GenerateTab(ctx context.Context, req *tab.TabRequest) (*tab.TabResponse, error) {
	rawNoteSeq := music.NoteSequence{}
	for _, chunk := range req.Chunks {
		OandFReq := onsets_frames.OAFRequest{
			AudioData: &onsets_frames.AudioFileData{
				FileName:   "temp",
				AudioBytes: chunk.AudioData,
			},
		}
		notes, err := clients.OnsetsAndFramesClient.Analyze(context.Background(), &OandFReq)
		if err != nil {
			return nil, err
		}

		seq := music.NewNoteSequence(len(notes.Notes))

		for i, note := range notes.Notes {
			name, octave := music.MidiToNote(int(note.MidiPitch))
			seq.Notes[i] = music.NoteEvent{
				Name:      name,
				Octave:    octave,
				StartTime: note.OnsetSeconds + chunk.StartTime,
				Velocity:  note.Velocity,
			}
		}

		rawNoteSeq.Merge(seq)
	}

	tabs, err := models.GenerateTab(audioResp)
	if err != nil {
		return nil, err
	}

	return &tab.TabResponse{Tab: tabs}, nil
}
