package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"zharchiver/models"
)

func downloadImage(imgURL string, answerID string, index int) (string, error) {
	saveDir := filepath.Join("storage", "images", answerID)
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return "", err
	}

	ext := ".jpg"
	if strings.Contains(imgURL, ".png") {
		ext = ".png"
	} else if strings.Contains(imgURL, ".webp") {
		ext = ".webp"
	} else if strings.Contains(imgURL, ".gif") {
		ext = ".gif"
	} else if strings.Contains(imgURL, ".mp4") {
		ext = ".mp4"
	}

	fileName := fmt.Sprintf("img_%d%s", index, ext)
	localPath := filepath.Join(saveDir, fileName)

	req, err := http.NewRequest("GET", imgURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	// 视频文件可能比较大，增加超时时间到 60 秒
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("图片下载失败，状态码: %d", resp.StatusCode)
	}

	out, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return "/" + filepath.ToSlash(localPath), nil
}

func ProcessImages(data *models.AnswerData) error {
	for i, rawURL := range data.ImageURLs {
		localRelPath, err := downloadImage(rawURL, data.AnswerID, i+1)
		if err != nil {
			log.Printf("图片下载失败 [%s]: %v", rawURL, err)
			continue
		}

		target1 := fmt.Sprintf(`data-actualsrc="%s"`, rawURL)
		target2 := fmt.Sprintf(`data-original="%s"`, rawURL)
		replaceStr := fmt.Sprintf(`src="%s"`, localRelPath)

		data.ContentHTML = strings.ReplaceAll(data.ContentHTML, target1, replaceStr)
		data.ContentHTML = strings.ReplaceAll(data.ContentHTML, target2, replaceStr)
	}
	return nil
}
