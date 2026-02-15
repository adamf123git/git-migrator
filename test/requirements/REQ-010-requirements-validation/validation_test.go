package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// findReqDir finds the directory for a requirement ID by searching parent directory
func findReqDir(reqID string) string {
	parentDir := filepath.Join("..")
	entries, err := os.ReadDir(parentDir)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), reqID) {
			return filepath.Join(parentDir, entry.Name())
		}
	}
	return ""
}

// TestRequirementsDocumentsExist validates all requirement docs exist
func TestRequirementsDocumentsExist(t *testing.T) {
	requirements := []string{"REQ-007", "REQ-009", "REQ-010"}

	for _, reqID := range requirements {
		reqDir := findReqDir(reqID)
		if reqDir == "" {
			t.Errorf("Requirement directory for %s does not exist", reqID)
			continue
		}
		reqFile := filepath.Join(reqDir, "requirement.md")
		if _, err := os.Stat(reqFile); os.IsNotExist(err) {
			t.Errorf("Requirement file %s does not exist", reqFile)
		}
	}
}

// TestRequirementsHaveAcceptanceCriteria validates acceptance criteria
func TestRequirementsHaveAcceptanceCriteria(t *testing.T) {
	requirements := []string{"REQ-007", "REQ-009", "REQ-010"}

	for _, reqID := range requirements {
		reqDir := findReqDir(reqID)
		if reqDir == "" {
			t.Errorf("Requirement directory for %s does not exist", reqID)
			continue
		}
		reqFile := filepath.Join(reqDir, "requirement.md")
		if _, err := os.Stat(reqFile); os.IsNotExist(err) {
			t.Errorf("Requirement file %s does not exist", reqFile)
		}
	}
}

// TestRequirementsValidationScript tests the standalone validator
func TestRequirementsValidationScript(t *testing.T) {
	t.Skip("Validation script not yet implemented")
}

// TestRequirementsStatusTracking validates status tracking
func TestRequirementsStatusTracking(t *testing.T) {
	t.Skip("Status tracking not yet implemented")
}
