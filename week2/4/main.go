package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// LogAnalyzerInterface defines the methods required for log analysis.
type LogAnalyzerInterface interface {
	Worker()
	ConcurrentWorker(filePath string)
}

// LogAnalyzer is the struct that implements LogAnalyzerInterface.
type LogAnalyzer struct {
	severity   string
	outputPath string
	logsDir    string
	skipPath   []string
	buggyLogs  map[string]int
	mu         sync.Mutex
}

// NewLogAnalyzer initializes a LogAnalyzer instance.
func NewLogAnalyzer(severity string, outputPath string, logsDir string) LogAnalyzerInterface {
	return &LogAnalyzer{
		severity:   severity,
		outputPath: outputPath,
		logsDir:    logsDir,
		buggyLogs:  make(map[string]int),
	}
}

// Worker processes log files sequentially.
func (la *LogAnalyzer) Worker() {
	var allValidLogs []string
	// Get list of log files from the logsDir
	logFiles := getLogFiles(la.logsDir)

	// Process each file in sequential order
outerLoop:
	for _, file := range logFiles {
		for _, skippath := range la.skipPath {
			fullPath := filepath.Join(la.logsDir, skippath)
			if file == fullPath {
				continue outerLoop

			}

		}
		validLogs, bugCount := parseLogFile(file, la)
		allValidLogs = append(allValidLogs, validLogs...)
		la.buggyLogs[file] = bugCount
	}

	// Write valid logs to output
	writeToFile(la.outputPath, allValidLogs, la.buggyLogs)
}

var fileCompletion sync.Map // Track file processing completion
var cond = sync.NewCond(&sync.Mutex{})

// ConcurrentWorker processes a specific log file concurrently.
func (la *LogAnalyzer) ConcurrentWorker(filePath string) {

	la.processFileWithDependencies(filePath)
	// fmt.Println("BUGGY MAP: ", la.buggyLogs)
	cond.L.Lock()
	fileCompletion.Store(filePath, true)
	cond.Broadcast() // Notify waiting goroutines

	cond.L.Unlock()
}

// parseLogFile parses a log file and counts valid and buggy logs.
func parseLogFile(logFile string, analyzer *LogAnalyzer) ([]string, int) {
	bugCount := 0
	validLogs := []string{}

	file, err := os.Open(logFile)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		return validLogs, bugCount
	}
	defer file.Close()

	// Read the log file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "...") {
			trimmedText := strings.TrimPrefix(line, "...")
			analyzer.skipPath = append(analyzer.skipPath, trimmedText)

			fPath := filepath.Join(analyzer.logsDir, trimmedText)
			validLogs2, bugCount := parseLogFile(fPath, analyzer)
			validLogs = append(validLogs, validLogs2...)
			analyzer.buggyLogs[trimmedText] = bugCount

		} else {

			// Check if log matches the format
			if isValidLog(line) {
				// Check if the log matches the required severity
				if strings.HasPrefix(line, analyzer.severity) {
					// fmt.Println(line)
					validLogs = append(validLogs, line)
				}
			} else {
				bugCount++
			}
		}

	}

	return validLogs, bugCount
}

// parseLogFile parses a log file and counts valid and buggy logs. with dependency
func (la *LogAnalyzer) processFileWithDependencies(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	var validLogs []string
	bugCount := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "...") {
			trimmedText := strings.TrimPrefix(line, "...")
			la.skipPath = append(la.skipPath, trimmedText)
			// The line indicates a dependency on another log file
			dependentFile := filepath.Join(la.logsDir, strings.TrimPrefix(line, "..."))

			// Wait for the dependent file to complete
			cond.L.Lock()
			for {
				if done, _ := fileCompletion.Load(dependentFile); done != nil {
					break
				}
				cond.Wait() // Wait until the dependency is done
			}
			cond.L.Unlock()

		} else {
			// Process the current line for logs and bugs
			if isValidLog(line) {
				if strings.HasPrefix(line, la.severity) {
					validLogs = append(validLogs, line)
				}
			} else {
				bugCount++
			}
		}
	}

	// Append results to the global output in a thread-safe manner
	if _, exists := la.buggyLogs[filePath]; !exists {
		la.buggyLogs[filePath] = bugCount
	}
	la.mu.Lock()
	la.appendToOutput(validLogs, filePath, bugCount) // Make sure this function appends and doesn't overwrite
	la.mu.Unlock()

}

// Ensure logs are appended to the output
func (la *LogAnalyzer) appendToOutput(logs []string, filePath string, bugCount int) {
	f, err := os.OpenFile(la.outputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening output file:", err)
		return
	}
	defer f.Close()

	for _, log := range logs {
		if _, err := f.WriteString(log + "\n"); err != nil {
			fmt.Println("Error writing to output file:", err)
		}
	}
	f.WriteString(fmt.Sprintf("%s: %d bugs\n", filepath.Base(filePath), bugCount))
	// Write buggy logs summary

	// Use filepath.Base for consistency

	// for fileName, bugs := range la.buggyLogs {
	// 	f.WriteString(fmt.Sprintf("%s: %d bugs\n", filepath.Base(fileName), bugs)) // Use filepath.Base for consistency
	// 	delete(la.buggyLogs, filepath.Base(fileName))
	// }
}

// getLogFiles retrieves all log files in the specified directory.
func getLogFiles(logsDir string) []string {
	var files []string
	filepath.Walk(logsDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// isValidLog checks if a log line is valid based on a regex pattern.
func isValidLog(line string) bool {
	// Define a regex pattern to match the log line structure with spaces around the pipes
	pattern := `^(INFO|WARNING|ERROR)\s*\|\s*([0-9-]+\s[0-9:]+)\s*\|\s*(.+)$`
	re := regexp.MustCompile(pattern)

	// Check if the log line matches the pattern
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		// fmt.Println("Invalid format")
		return false
	}

	// Extract the second part (date string)
	dateStr := matches[2]

	// Define the date layout you expect ("2006-01-02 15:04:05" for the format YYYY-MM-DD HH:MM:SS)
	layout := "2006-01-02 15:04:05"

	// Try to parse the date string
	_, err := time.Parse(layout, dateStr)
	if err != nil {
		// fmt.Println("Invalid date format:", err)
		return false
	}

	// fmt.Println("Valid log line")
	return true
}

// writeToFile writes the results to the output file.
func writeToFile(outputPath string, validLogs []string, buggyLogs map[string]int) {
	// Open the output file for writing
	outFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %s\n", err)
		return
	}
	defer outFile.Close()

	// Write valid logs first
	for _, log := range validLogs {
		outFile.WriteString(log + "\n")
	}

	// Write buggy logs summary
	for fileName, bugs := range buggyLogs {
		outFile.WriteString(fmt.Sprintf("%s: %d bugs\n", filepath.Base(fileName), bugs)) // Use filepath.Base for consistency
	}
}
