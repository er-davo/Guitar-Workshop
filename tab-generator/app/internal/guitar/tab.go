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
	time     float32
	timeStep float32

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

func (tb *tabBuilder) Tab() string {
	tab := strings.Builder{}

	for i := range len(tb.tabStrings) {
		tab.WriteString(tb.tabStrings[i].String() + "\n")
	}

	return tab.String()
}

func (tb *tabBuilder) WriteSingleNote(n Note) error {
	if n.Time < tb.time {
		// TODO
		return errors.New("")
	}
	silence := int((n.Time - tb.time) / tb.timeStep)
	tb.addSilence(silence)
	tb.time += tb.timeStep * (float32(silence + 1))

	for i := range tb.tabStrings {
		skipToOtherStrings := "-"

		if i == n.String {
			if n.Fret/10 > 0 {
				skipToOtherStrings += "-"
			}
			tb.tabStrings[i].WriteString(fmt.Sprintf("%d", n.Fret))
			continue
		}

		tb.tabStrings[i].WriteString(skipToOtherStrings)
	}

	return nil
}

// TODO
func (tb *tabBuilder) WriteChord() {

}

func (tb *tabBuilder) addNotes(notes []string) {
	for i := range notes {
		tb.tabStrings[i].WriteString(notes[i] + "|")
	}
}

func (tb *tabBuilder) addSilence(n int) {
	for range n {
		for i := range len(tb.tabStrings) {
			tb.tabStrings[i].WriteRune('-')
		}
	}
}
