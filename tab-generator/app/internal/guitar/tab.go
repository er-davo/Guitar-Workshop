package guitar

import (
	"errors"
	"fmt"
	"strings"
)

type InstrumentType int

const (
	GuitarType InstrumentType = iota
)

type tabBuilder struct {
	time           float32
	instrumentType InstrumentType
	tabStrings     []strings.Builder
}

func NewTabBuilder(instrumentType InstrumentType, tuningNotes []string) (*tabBuilder, error) {
	switch instrumentType {
	case GuitarType:
		tb := tabBuilder{
			time:           0.0,
			instrumentType: instrumentType,
			tabStrings:     make([]strings.Builder, 6),
		}
		tb.addNotes(tuningNotes)

		return &tb, nil
	default:
		return nil, errors.ErrUnsupported
	}
}

func (tb *tabBuilder) addNotes(notes []string) {
	for i := range notes {
		tb.tabStrings[i].WriteString(notes[i] + "|")
	}
}

func (tb *tabBuilder) WriteSingleNote(n Note) {
	for i := range tb.tabStrings {
		if i == n.String {
			tb.tabStrings[i].WriteString(fmt.Sprintf("%d", n.Fret))
			continue
		}

		tb.tabStrings[i].WriteRune('-')
	}
}

// TODO
func (tb *tabBuilder) WriteChord() {

}
