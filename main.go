package main

import (
	"bufio"
	"flag"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"log/slog"
)

type LinkCheckResult struct {
	Link       string
	IsValid    bool
	StatusCode int
}

func main() {
	filename := flag.String("file", "", "Markdown file to test")
	verbose := flag.Bool("verbose", false, "Enable verbose mode")
	failedOnly := flag.Bool("failed-only", false, "Return only failed links")
	flag.Parse()

	if *filename == "" {
		slog.Error("Please provide a markdown file.")
		os.Exit(1)
	}

	file, err := os.Open(*filename)
	if err != nil {
		slog.Error("Failed to open file: ", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)

	results := []LinkCheckResult{}

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			link := strings.TrimSpace(match[2])
			resp, err := http.Head(link)
			if err != nil {
				slog.Error("Error checking link: ", err)
			} else {
				result := LinkCheckResult{
					Link:       link,
					IsValid:    resp.StatusCode == http.StatusOK,
					StatusCode: resp.StatusCode,
				}
				results = append(results, result)
				if *verbose || (!*failedOnly && !result.IsValid) {
					statusColor := color.GreenString
					if !result.IsValid {
						statusColor = color.RedString
					}
					slog.Info(link + ": " + statusColor("%d", result.StatusCode))
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Info("Error scanning file: ", err)
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
	if *failedOnly {
		for _, result := range results {
			if !result.IsValid {
				slog.Info("Failed link: " + result.Link)
			}
		}
	} else {
		summaryColor := color.GreenString
		if invalidCount > 0 {
			summaryColor = color.RedString
		}
		slog.Info(summaryColor("Summary: %d valid links, %d invalid links", validCount, invalidCount))
	}

	if invalidCount > 0 {
		os.Exit(1)
	}
}
