package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type ffprobeOutput struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("couldn't run ffprobe: %w", err)
	}

	var output ffprobeOutput
	err = json.Unmarshal(stdout.Bytes(), &output)
	if err != nil {
		return "", fmt.Errorf("couldn't unmarshal ffprobe output: %w", err)
	}

	if len(output.Streams) == 0 {
		return "", fmt.Errorf("no streams found in video")
	}

	width := output.Streams[0].Width
	height := output.Streams[0].Height

	ratio := float64(width) / float64(height)
	const landscape = 16.0 / 9.0
	const portrait = 9.0 / 16.0
	const tolerance = 0.05

	if ratio > landscape-tolerance && ratio < landscape+tolerance {
		return "16:9", nil
	}
	if ratio > portrait-tolerance && ratio < portrait+tolerance {
		return "9:16", nil
	}
	return "other", nil
}
