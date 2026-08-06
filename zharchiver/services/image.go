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
	"zharchiver/utils"
)

type ProgressWriter struct {
	Total    int64
	Written  int64
	FileName string
	LastTime time.Time
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.Written += int64(n)
	
	if pw.Total > 0 && time.Since(pw.LastTime) > 500*time.Millisecond {
		percent := int((pw.Written * 100) / pw.Total)
		if percent > 100 {
			percent = 100
		}
		// 广播进度 PROGRESS|percent|filename
		utils.BroadcastLog("PROGRESS", fmt.Sprintf("%d|%s", percent, pw.FileName))
		pw.LastTime = time.Now()
	}
	return n, nil
}

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

	pw := &ProgressWriter{
		Total:    resp.ContentLength,
		FileName: fileName,
		LastTime: time.Now(),
	}

	// 开始时发送一次0%
	if pw.Total > 0 {
		utils.BroadcastLog("PROGRESS", fmt.Sprintf("0|%s", fileName))
	}

	_, err = io.Copy(out, io.TeeReader(resp.Body, pw))
	
	// 完成时发送100%
	if pw.Total > 0 && err == nil {
		utils.BroadcastLog("PROGRESS", fmt.Sprintf("100|%s", fileName))
	}
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
