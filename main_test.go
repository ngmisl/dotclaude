package main

import (
	"encoding/json"
	"os"
	"testing"
)

func setup(m *testing.M) func() {
	// Create a dummy config file
	configContent := `
ignored_directories:
  - "node_modules"
supported_extensions:
  - ".txt"
`
	os.WriteFile("config.yaml", []byte(configContent), 0644)

	// Create a dummy test file
	content := `
This is a test file.
Here is a zero-width space: ​
Here is a homoglyph: а (Cyrillic 'a')
Here is some base64: aW1wb3J0IG9zCg==
`
	os.MkdirAll("testdata", 0755)
	os.WriteFile("testdata/testfile.txt", []byte(content), 0644)
	os.WriteFile("testdata/clean.txt", []byte("This is a clean file."), 0644)

	return func() {
		// Clean up
		os.Remove("config.yaml")
		os.RemoveAll("testdata")
	}
}

func TestMain(m *testing.M) {
	teardown := setup(m)
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

func TestScanDirectory(t *testing.T) {
	// Reset issues slice
	issues = []Issue{}

	config, err := loadConfig("config.yaml")
	if err != nil {
		t.Fatalf("Error loading config.yaml: %v", err)
	}

	for _, ext := range config.SupportedExtensions {
		supportedExts[ext] = true
	}

	scanDirectory("testdata", config.IgnoredDirectories)

	if len(issues) != 3 {
		t.Errorf("Expected 3 issues, but got %d", len(issues))
	}

	foundZWC := false
	foundHomoglyph := false
	foundBase64 := false

	for _, issue := range issues {
		if issue.Type == "Zero-Width Char" {
			foundZWC = true
		}
		if issue.Type == "Homoglyph" {
			foundHomoglyph = true
		}
		if issue.Type == "Base64" {
			foundBase64 = true
		}
	}

	if !foundZWC {
		t.Errorf("Expected to find 'Zero-Width Char' issue, but didn't")
	}
	if !foundHomoglyph {
		t.Errorf("Expected to find 'Homoglyph' issue, but didn't")
	}
	if !foundBase64 {
		t.Errorf("Expected to find 'Base64' issue, but didn't")
	}
}

func TestJSONOutput(t *testing.T) {
	// Reset issues slice
	issues = []Issue{}

	config, err := loadConfig("config.yaml")
	if err != nil {
		t.Fatalf("Error loading config.yaml: %v", err)
	}

	for _, ext := range config.SupportedExtensions {
		supportedExts[ext] = true
	}

	scanDirectory("testdata", config.IgnoredDirectories)

	jsonIssues, err := json.Marshal(issues)
	if err != nil {
		t.Fatalf("Error marshalling issues to JSON: %v", err)
	}

	var decodedIssues []Issue
	err = json.Unmarshal(jsonIssues, &decodedIssues)
	if err != nil {
		t.Fatalf("Error unmarshalling issues from JSON: %v", err)
	}

	if len(decodedIssues) != 3 {
		t.Errorf("Expected 3 issues in JSON output, but got %d", len(decodedIssues))
	}
}
