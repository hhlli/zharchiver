package utils

import (
	"os"

	"github.com/abema/go-mp4"
)

// GetVideoDimensions 解析 MP4 文件并返回宽高
func GetVideoDimensions(filePath string) (width int, height int, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	info, err := mp4.Probe(file)
	if err != nil {
		return 0, 0, err
	}

	for _, track := range info.Tracks {
		if track.AVC != nil {
			return int(track.AVC.Width), int(track.AVC.Height), nil
		}
	}

	return 0, 0, nil
}
