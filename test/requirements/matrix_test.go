package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// RequirementsMatrix maps requirement IDs to their test functions
var RequirementsMatrix = map[string][]string{
	"REQ-007": {
		"TestCLICommands",
		"TestCLIVersion",
		"TestCLIHelp",
	},
	"REQ-009": {
		"TestTDDWorkflow",
		"TestCoverageEnforcement",
		"TestRegressionSuite",
	},
	"REQ-010": {
		"TestRequirementsMatrix",
		"TestNoOrphanTests",
		"TestRequirementsCoverage",
	},
}

// findRequirementDir finds the directory for a requirement ID (e.g., REQ-007 matches REQ-007-cli-interface)
func findRequirementDir(reqID string) string {
	entries, err := os.ReadDir(".")
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), reqID) {
			return entry.Name()
		}
	}
	return ""
}

// TestAllRequirementsCovered ensures every requirement has tests
func TestAllRequirementsCovered(t *testing.T) {
	for reqID, tests := range RequirementsMatrix {
		if len(tests) == 0 {
			t.Errorf("Requirement %s has no tests defined", reqID)
		}

		for _, test := range tests {
			if test == "" {
				t.Errorf("Empty test name for requirement %s", reqID)
			}
		}
	}
}

// TestNoOrphanTests ensures every test maps to a requirement
func TestNoOrphanTests(t *testing.T) {
	// This will be expanded as we add more tests
	// For now, we just verify the matrix is valid
	requirementCount := len(RequirementsMatrix)
	if requirementCount == 0 {
		t.Error("Requirements matrix is empty")
	}
}

// TestRequirementsCoverage checks that all requirement directories exist
func TestRequirementsCoverage(t *testing.T) {
	for reqID := range RequirementsMatrix {
		reqDir := findRequirementDir(reqID)
		if reqDir == "" {
			t.Errorf("Requirement directory for %s does not exist", reqID)
			continue
		}

		// Check for requirement.md
		reqFile := filepath.Join(reqDir, "requirement.md")
		if _, err := os.Stat(reqFile); os.IsNotExist(err) {
			t.Errorf("Requirement file %s does not exist", reqFile)
		}
	}
}

// Helper function to check if a test exists (to be expanded later)
func testExists(testName string) bool {
	// This will be implemented to actually check if tests exist
	return testName != ""
}

// Helper function to get all tests from the matrix
func getAllMappedTests() map[string]bool {
	mapped := make(map[string]bool)
	for _, tests := range RequirementsMatrix {
		for _, test := range tests {
			mapped[test] = true
		}
	}
	return mapped
}
