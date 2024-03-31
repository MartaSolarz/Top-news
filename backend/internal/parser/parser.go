package parser

import (
	"log"
	"os/exec"
	"time"
)

func ParseDate(date string) string {
	t, err := time.Parse(time.RFC1123, date)
	if err != nil {
		log.Println(err)
	}
	return t.Format("2006-01-02 15:04:05")
}

func ExtractContent(url string) string {
	cmd := exec.Command("venv/bin/python", "backend/python/extract.py", url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to extract content: %v", err)
		return ""
	}
	return string(output)
}
