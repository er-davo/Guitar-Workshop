package convertor

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
)

func MP3ToWAVInMemory(file multipart.File) ([]byte, error) {
	cmd := exec.Command("ffmpeg",
		"-i", "pipe:0", // stdin = mp3 input
		"-f", "wav", // output format
		"-ar", "44100", // sample rate
		"-ac", "2", // stereo
		"pipe:1", // stdout = wav output
	)

	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// передаём mp3 в stdin ffmpeg
	go func() {
		defer stdin.Close()
		io.Copy(stdin, file)
	}()

	// читаем результат wav из stdout ffmpeg
	var buf bytes.Buffer
	_, err = io.Copy(&buf, stdout)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil // здесь — wav в памяти
}
