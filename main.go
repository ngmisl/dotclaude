package main

import (
	"encoding/base64"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

type charDetector struct {
	chars    []rune
	name     string
	severity string
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

	base64Pattern = regexp.MustCompile(`[A-Za-z0-9+/]{20,}={0,2}`)
	
	supportedExts = map[string]bool{
		".xml": true, ".md": true, ".yaml": true, ".yml": true, ".txt": true,
	}

	totalFiles      = 0
	scannedFiles    = 0
	filesWithIssues = 0
	issueCount      = 0
)

func main() {
	rootDir := "."
	if len(os.Args) > 1 {
		rootDir = os.Args[1]
	}

	fmt.Printf("🔍 Scanning directory: %s\n\n", rootDir)
	scanDirectory(rootDir)
	printSummary()
}

func scanDirectory(root string) {
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			name := d.Name()
			skip := []string{"node_modules", "vendor", ".git", "dist", "build", ".next"}
			for _, s := range skip {
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

func scanFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	if !utf8.Valid(content) {
		printIssue(path, 0, "HIGH", "Invalid UTF-8", "File contains invalid UTF-8 sequences")
		return
	}

	scannedFiles++
	lines := strings.Split(string(content), "\n")
	fileIssues := 0

	for lineNum, line := range lines {
		lineNum++
		
		// Character-based detection
		for i, r := range line {
			// Control characters (skip normal whitespace)
			if r != '\t' && r != '\n' && r != '\r' && ((r >= 0x00 && r <= 0x1F) || (r >= 0x7F && r <= 0x9F)) {
				if fileIssues == 0 {
					fmt.Printf("\n📄 %s\n", path)
				}
				fileIssues++
				fmt.Printf("  🟡 L%d:C%d MEDIUM Control Char U+%04X\n", lineNum, i+1, r)
			}

			// Run through detectors
			for _, det := range detectors {
				for _, char := range det.chars {
					if r == char {
						if fileIssues == 0 {
							fmt.Printf("\n📄 %s\n", path)
						}
						fileIssues++
						symbol := severitySymbol(det.severity)
						fmt.Printf("  %s L%d:C%d %s %s U+%04X\n", symbol, lineNum, i+1, det.severity, det.name, r)
					}
				}
			}
		}

		// Base64 detection
		matches := base64Pattern.FindAllStringIndex(line, -1)
		for _, match := range matches {
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
				
				// Check for dangerous patterns
				dangerous := []string{"exec(", "eval(", "system(", "__import__", "<script", "javascript:", "data:text/html"}
				for _, keyword := range dangerous {
					if strings.Contains(strings.ToLower(decodedStr), keyword) {
						severity = "HIGH"
						break
					}
				}

				if fileIssues == 0 {
					fmt.Printf("\n📄 %s\n", path)
				}
				fileIssues++
				symbol := severitySymbol(severity)
				snippet := decodedStr
				if len(snippet) > 60 {
					snippet = snippet[:60] + "..."
				}
				fmt.Printf("  %s L%d:C%d %s Base64: %s\n", symbol, lineNum, match[0]+1, severity, snippet)
			}
		}
	}

	if fileIssues > 0 {
		filesWithIssues++
		issueCount += fileIssues
	}
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

func printIssue(path string, line int, severity, issueType, details string) {
	symbol := severitySymbol(severity)
	fmt.Printf("\n📄 %s\n", path)
	fmt.Printf("  %s L%d %s %s: %s\n", symbol, line, severity, issueType, details)
	filesWithIssues++
	issueCount++
}

func printSummary() {
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
