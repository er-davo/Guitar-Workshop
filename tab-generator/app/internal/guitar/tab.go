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
	defaultTimeStep := float32(0.2)
	switch instrumentType {
	case GuitarType:
		tb := tabBuilder{
			time:           0.0,
			timeStep:       defaultTimeStep,
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
	if !noteIsValid(n) {
		return fmt.Errorf("invalid note: %s", n.Name)
	}
	if n.Time < tb.time {
		return fmt.Errorf("note time %v precedes current time %v", n.Time, tb.time)
	}

	silence := int((n.Time - tb.time) / tb.timeStep)
	tb.addSilence(silence)
	tb.time += tb.timeStep * (float32(silence + 1))
	skipToOtherStrings := "-"
	if n.Fret/10 > 0 {
		skipToOtherStrings = "--"
	}

	for i := range tb.tabStrings {
		if i == n.String {
			tb.tabStrings[i].WriteString(fmt.Sprintf("%d", n.Fret))
			continue
		}

		tb.tabStrings[i].WriteString(skipToOtherStrings)
	}

	// to escape situations like:
	// E|-3--123-----
	tb.addSilence(1)

	return nil
}

// TODO
func (tb *tabBuilder) WriteChord() {

}

func (tb *tabBuilder) addNotes(notes []string) error {
	if len(notes) != len(tb.tabStrings) {
		return fmt.Errorf("invalid tuning notes count")
	}

	for i := range notes {
		tb.tabStrings[i].WriteString(notes[i] + "|")
	}
	return nil
}

func (tb *tabBuilder) addSilence(n int) {
	for range n {
		for i := range len(tb.tabStrings) {
			tb.tabStrings[i].WriteString("-")
		}
	}
}
