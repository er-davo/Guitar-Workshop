package models

import (
	"slices"
)

type noteEvent struct {
	time     float32
	pitch    float32
	mainNote string
	octave   int
	notes    []string
}

func (n *noteEvent) equalsTo(other noteEvent, timeDiff float32) bool {
	if n.mainNote == other.mainNote &&
		slices.Equal(n.notes, other.notes) &&
		n.time-other.time <= timeDiff {
		return true
	}
	return false
}

type notesEvent struct {
	notes []noteEvent
}

func newNotesEvent(events *audiopb.AudioResponse) *notesEvent {
	notes := notesEvent{}

	for _, e := range events.NoteFeatures {
		// skip noise and silence
		if len(e.ChromaNotes) > 6 || e.MainNote == "" {
			continue
		}
		// if !slices.Contains(e.ChromaNotes, e.MainNote) {
		// 	continue
		// }

		notes.notes = append(notes.notes, *createNoteEvent(e))
	}

	notes.removeLongDurations()

	return &notes
}

func (n *notesEvent) removeLongDurations() {
	var timeDiff float32 = 0.8
	for i := 1; i < len(n.notes); i++ {
		if n.notes[i-1].equalsTo(n.notes[i], timeDiff) {
			start := i - 1
			end := -1

			for j := i; j < len(n.notes)-1; j++ {
				if !n.notes[start].equalsTo(n.notes[j+1], timeDiff) {
					end = j + 1
					break
				}
			}

			if end == -1 {
				end = len(n.notes) - 1
			}

			// logger.Info("deleting from - to",
			// 	zap.String(
			// 		"start_note_event",
			// 		fmt.Sprintf("%+v", n.notes[start]),
			// 	),
			// 	zap.String(
			// 		"end_note_event",
			// 		fmt.Sprintf("%+v", n.notes[end]),
			// 	),
			// )

			n.notes = slices.Delete(n.notes, start, end)

			i = start + 1
		}
	}
}
