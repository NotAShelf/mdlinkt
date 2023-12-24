package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
)

type LinkCheckResult struct {
	Link       string
	IsValid    bool
	StatusCode int
}

func logWithColor(level string, msg string, args ...interface{}) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	colorFunc := color.New(color.FgWhite).SprintFunc()
	switch level {
	case "ERROR":
		colorFunc = color.New(color.FgRed).SprintFunc()
	case "WARN":
		colorFunc = color.New(color.FgYellow).SprintFunc()
	case "INFO":
		colorFunc = color.New(color.FgCyan).SprintFunc()
	}
	fmt.Printf("%s %s %s\n", timestamp, colorFunc(level), fmt.Sprintf(msg, args...))
}

func main() {
	filename := flag.String("file", "", "Markdown file to test")
	verbose := flag.Bool("verbose", false, "Enable verbose mode")
	failedOnly := flag.Bool("failed-only", false, "Return only failed links")
	flag.Parse()

	if *filename == "" {
		logWithColor("INFO", "Please provide a markdown file.")
		os.Exit(1)
	}

	file, err := os.Open(*filename)
	if err != nil {
		logWithColor("ERROR", "Failed to open file: %v", err)
		os.Exit(1)
	}
	if file == nil {
		logWithColor("ERROR", "File is nil")
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner == nil {
		logWithColor("ERROR", "Scanner is nil")
		os.Exit(1)
	}
	re := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)

	results := []LinkCheckResult{}

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			link := strings.TrimSpace(match[2])
			resp, err := http.Head(link)
			if err != nil || resp == nil {
				logWithColor("ERROR", "Invalid link: %s", link)
				results = append(results, LinkCheckResult{
					Link:       link,
					IsValid:    false,
					StatusCode: http.StatusBadRequest,
				})
			} else {
				isValid := resp.StatusCode == http.StatusOK
				result := LinkCheckResult{
					Link:       link,
					IsValid:    isValid,
					StatusCode: resp.StatusCode,
				}
				results = append(results, result)
				if *verbose || (!*failedOnly && !isValid) {
					statusColor := color.GreenString
					if !isValid {
						statusColor = color.RedString
					}
					logWithColor("INFO", "%s: %s", link, statusColor("%d", resp.StatusCode))
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logWithColor("ERROR", "Error scanning file: %v", err)
		os.Exit(1)
	}

	// Print summary
	validCount := 0
	invalidCount := 0
	for _, result := range results {
		if result.IsValid {
			validCount++
		} else {
			invalidCount++
		}
	}
	summaryColor := color.GreenString
	if invalidCount > 0 {
		summaryColor = color.RedString
	}
	logWithColor("INFO", summaryColor("Summary: %d valid links, %d invalid links"), validCount, invalidCount)

	if *failedOnly {
		for _, result := range results {
			if !result.IsValid {
				logWithColor("ERROR", "Failed link: %s", result.Link)
			}
		}
	}

	if invalidCount > 0 {
		os.Exit(1)
	}
}
