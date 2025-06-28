package music

import "github.com/er-davo/guitar"

func GenerateTab(notes NoteSequence) (string, error) {
	tun, err := guitar.ParseTuning(guitar.StandardTuning)
	if err != nil {
		return "", err
	}
	fb, err := guitar.NewFingerBoard(tun, 24)
	if err != nil {
		return "", err
	}
	tab, err := guitar.NewTabWriter(tun.NoteNames())
	if err != nil {
		return "", err
	}
	opt := guitar.NewFingeringOptimizer(*fb)

	layers, err := opt.TimeLayers(notes.guitarSequence())
	if err != nil {
		return "", err
	}

	path, err := opt.OptimizePath(layers)
	if err != nil {
		return "", err
	}

	err = tab.Write(path...)
	if err != nil {
		return "", err
	}

	return tab.Tab(), nil
}
