package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

type charDetector struct {
	chars    []rune
	name     string
	severity string
}

type Config struct {
	IgnoredDirectories []string `yaml:"ignored_directories"`
	SupportedExtensions []string `yaml:"supported_extensions"`
}

type Issue struct {
	FilePath    string `json:"file_path"`
	Line        int    `json:"line"`
	Column      int    `json:"column"`
	Severity    string `json:"severity"`
	Type        string `json:"type"`
	Details     string `json:"details"`
}

var (
	detectors = []charDetector{
		{
			chars:    []rune{'\u200B', '\u200C', '\u200D', '\uFEFF', '\u180E'},
			name:     "Zero-Width Char",
			severity: "HIGH",
		},
		{
			chars:    []rune{'\u202A', '\u202B', '\u202C', '\u202D', '\u202E', '\u2066', '\u2067', '\u2068', '\u2069'},
			name:     "Direction Override",
			severity: "HIGH",
		},
		{
			chars:    []rune{'\u0430', '\u0435', '\u043E', '\u0440', '\u0441', '\u0445', '\u0443', '\u0391', '\u0392', '\u0395', '\u0397'},
			name:     "Homoglyph",
			severity: "MEDIUM",
		},
	}

	base64Pattern = regexp.MustCompile(`[A-Za-z0-9+/]{10,}={0,2}`)
	
	supportedExts = make(map[string]bool)
	issues        = []Issue{}

	totalFiles      = 0
	scannedFiles    = 0
	filesWithIssues = 0
	issueCount      = 0
	hasHighSeverityIssue = false
)

func loadConfig(path string) (*Config, error) {
	config := &Config{}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func main() {
	outputFormat := flag.String("output", "text", "Output format (text or json)")
	flag.Parse()

	config, err := loadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Error loading config.yaml: %v\n", err)
		os.Exit(1)
	}

	for _, ext := range config.SupportedExtensions {
		supportedExts[ext] = true
	}

	rootDir := "."
	if flag.NArg() > 0 {
		rootDir = flag.Arg(0)
	}

	if *outputFormat == "text" {
		fmt.Printf("🔍 Scanning directory: %s\n\n", rootDir)
	}
	
	scanDirectory(rootDir, config.IgnoredDirectories)
	printSummary(*outputFormat)

	if hasHighSeverityIssue {
		os.Exit(1)
	}
}

func scanDirectory(root string, ignoredDirs []string) {
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			name := d.Name()
			for _, s := range ignoredDirs {
				if name == s {
					return fs.SkipDir
				}
			}
			return nil
		}

		if supportedExts[strings.ToLower(filepath.Ext(path))] {
			totalFiles++
			scanFile(path)
		}
		return nil
	})
}

func addIssue(path string, line, col int, severity, issueType, details string) {
	issues = append(issues, Issue{
		FilePath: path,
		Line:     line,
		Column:   col,
		Severity: severity,
		Type:     issueType,
		Details:  details,
	})
	if severity == "HIGH" {
		hasHighSeverityIssue = true
	}
}

func scanFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	if !utf8.Valid(content) {
		addIssue(path, 0, 0, "HIGH", "Invalid UTF-8", "File contains invalid UTF-8 sequences")
		return
	}

	scannedFiles++
	lines := strings.Split(string(content), "\n")
	
	for lineNum, line := range lines {
		lineNum++
		
		// Character-based detection
		for i, r := range line {
			col := i + 1
			// Control characters (skip normal whitespace)
			if r != '\t' && r != '\n' && r != '\r' && ((r >= 0x00 && r <= 0x1F) || (r >= 0x7F && r <= 0x9F)) {
				addIssue(path, lineNum, col, "MEDIUM", "Control Char", fmt.Sprintf("U+%04X", r))
			}

			// Run through detectors
			for _, det := range detectors {
				for _, char := range det.chars {
					if r == char {
						addIssue(path, lineNum, col, det.severity, det.name, fmt.Sprintf("U+%04X", r))
					}
				}
			}
		}

		// Base64 detection
		matches := base64Pattern.FindAllStringIndex(line, -1)
		for _, match := range matches {
			col := match[0] + 1
			candidate := line[match[0]:match[1]]
			decoded, err := base64.StdEncoding.DecodeString(candidate)
			if err != nil {
				decoded, err = base64.URLEncoding.DecodeString(candidate)
				if err != nil {
					continue
				}
			}

			if len(decoded) > 0 {
				decodedStr := string(decoded)
				severity := "LOW"
				
				dangerous := []string{"exec(", "eval(", "system(", "__import__", "<script", "javascript:", "data:text/html"}
				for _, keyword := range dangerous {
					if strings.Contains(strings.ToLower(decodedStr), keyword) {
						severity = "HIGH"
						break
					}
				}
				
				snippet := decodedStr
				if len(snippet) > 60 {
					snippet = snippet[:60] + "..."
				}
				addIssue(path, lineNum, col, severity, "Base64", snippet)
			}
		}
	}
}

func printSummary(format string) {
	if format == "json" {
		jsonIssues, err := json.MarshalIndent(issues, "", "  ")
		if err != nil {
			fmt.Printf("Error marshalling issues to JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonIssues))
		return
	}

	// Human-readable summary
	issueCount = len(issues)
	filesWithIssuesMap := make(map[string]bool)
	for _, issue := range issues {
		filesWithIssuesMap[issue.FilePath] = true
	}
	filesWithIssues = len(filesWithIssuesMap)

	for _, issue := range issues {
		symbol := severitySymbol(issue.Severity)
		fmt.Printf("\n📄 %s\n", issue.FilePath)
		fmt.Printf("  %s L%d:%d %s %s: %s\n", symbol, issue.Line, issue.Column, issue.Severity, issue.Type, issue.Details)
	}

	fmt.Println("\n═══════════════════════════════════════════════════════════════")
	fmt.Println("                    SCAN SUMMARY")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("Files found:          %d\n", totalFiles)
	fmt.Printf("Files scanned:        %d\n", scannedFiles)
	fmt.Printf("Clean files:          %d\n", scannedFiles-filesWithIssues)
	fmt.Printf("Files with issues:    %d\n", filesWithIssues)
	fmt.Printf("Total issues:         %d\n", issueCount)

	if issueCount == 0 {
		fmt.Println("\n✅ No suspicious content detected!")
	} else {
		fmt.Printf("\n⚠️  Found %d issue(s) in %d file(s)\n", issueCount, filesWithIssues)
		fmt.Println("   Review the files listed above")
	}
	fmt.Println()
}

func severitySymbol(severity string) string {
	switch severity {
	case "HIGH":
		return "🔴"
	case "MEDIUM":
		return "🟡"
	case "LOW":
		return "🔵"
	default:
		return "⚪"
	}
}
