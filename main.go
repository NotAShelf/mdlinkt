package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type LinkCheckResult struct {
	Link       string
	IsValid    bool
	StatusCode int
}

var wg sync.WaitGroup
var verboseMode bool

func logWithColor(level string, msg string, args ...interface{}) {
	if verboseMode {
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
}

func worker(jobs <-chan string, results chan<- LinkCheckResult, verboseFlag *bool, failedOnly *bool) {
	for link := range jobs {
		resp, err := http.Head(link)
		if err != nil {
			if *verboseFlag {
				logWithColor("ERROR", "Invalid link: %s", link)
			}
			results <- LinkCheckResult{
				Link:       link,
				IsValid:    false,
				StatusCode: http.StatusBadRequest,
			}
		} else {
			isValid := resp.StatusCode == http.StatusOK
			result := LinkCheckResult{
				Link:       link,
				IsValid:    isValid,
				StatusCode: resp.StatusCode,
			}
			results <- result
			if *verboseFlag || (!*failedOnly && !isValid) {
				statusColor := color.GreenString
				if !isValid {
					statusColor = color.RedString
				}
				logWithColor("INFO", "%s: %s", link, statusColor("%d", resp.StatusCode))
			}
		}
		wg.Done()
	}
}

func main() {
	filename := flag.String("file", "", "Markdown file to test")
	verboseFlag := flag.Bool("verbose", false, "Enable verbose mode")
	failedOnly := flag.Bool("failed-only", false, "Return only failed links")
	flag.Parse()

	verboseMode = *verboseFlag

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

	jobs := make(chan string, 10000)
	results := make(chan LinkCheckResult, 10000)

	// Start workers
	for i := 1; i <= 10; i++ {
		go worker(jobs, results, verboseFlag, failedOnly)
	}

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			wg.Add(1)
			jobs <- strings.TrimSpace(match[2])
		}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var invalidCount int

	for res := range results {
		if *failedOnly && res.IsValid {
			continue
		}
		if *verboseFlag || (!*failedOnly && !res.IsValid) {
			if res.IsValid {
				logWithColor("INFO", "Link %s is valid with status code %d", res.Link, res.StatusCode)
			} else {
				logWithColor("ERROR", "Link %s is invalid with status code %d", res.Link, res.StatusCode)
				invalidCount++
			}
		}
		wg.Done()
	}

	if invalidCount > 0 {
		os.Exit(1)
	}

	if err := scanner.Err(); err != nil {
		logWithColor("ERROR", "Error scanning file: %v", err)
		os.Exit(1)
	}
}
