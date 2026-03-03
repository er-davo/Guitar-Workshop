package processor

import (
	"tabgen/internal/music"
)

type TabProcessor struct {
}

func NewTabProcessor() *TabProcessor {
	return &TabProcessor{}
}

func (p *TabProcessor) GenerateTab(notes *music.NoteSequence) (string, error) {
	processedNotes := notes.Processed()
	return music.GenerateTab(*processedNotes)
}
