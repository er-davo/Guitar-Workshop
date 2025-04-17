package models

import (
	"tabgen/internal/audioproto"
)

type TabRequest struct {
	AudioURL string `json:"audio_url"`
}

type TabResponse struct {
	Tab    string `json:"tab"`
	Status string `json:"status"`
}

func GenerateTab(audio *audioproto.AudioResponse) (string, error) {
	chromo, err := decodeChromagram(audio.Chromagram)
	if err != nil {
		return "", err
	}

	rows, cols := chromo.Dims()
	for frmIdx, frame := range chromo {

	}

	return "", nil
}
