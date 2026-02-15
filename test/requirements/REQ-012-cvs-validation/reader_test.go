package requirements

import (
	"testing"

	"github.com/adamf123git/git-migrator/internal/vcs"
	"github.com/adamf123git/git-migrator/internal/vcs/cvs"
)

// TestCVSReaderValidation tests repository validation
func TestCVSReaderValidation(t *testing.T) {
	reader := cvs.NewReader("../../../test/fixtures/cvs/simple")
	err := reader.Validate()
	if err != nil {
		t.Errorf("Expected valid repository, got error: %v", err)
	}
}

// TestCVSReaderInvalidPath tests validation with invalid path
func TestCVSReaderInvalidPath(t *testing.T) {
	reader := cvs.NewReader("/nonexistent/path")
	err := reader.Validate()
	if err == nil {
		t.Error("Expected error for invalid path")
	}
}

// TestCVSReaderGetCommits tests commit extraction
func TestCVSReaderGetCommits(t *testing.T) {
	reader := cvs.NewReader("../../../test/fixtures/cvs/simple")
	iter, err := reader.GetCommits()
	if err != nil {
		t.Fatalf("GetCommits failed: %v", err)
	}

	var commits []*vcs.Commit
	for iter.Next() {
		commits = append(commits, iter.Commit())
	}

	if iter.Err() != nil {
		t.Errorf("Iterator error: %v", iter.Err())
	}

	if len(commits) == 0 {
		t.Error("Expected at least one commit")
	}

	// Check first commit
	if len(commits) > 0 {
		c := commits[0]
		if c.Author == "" {
			t.Error("Expected author to be set")
		}
		if c.Message == "" {
			t.Error("Expected message to be set")
		}
	}
}

// TestCVSReaderGetBranches tests branch extraction
func TestCVSReaderGetBranches(t *testing.T) {
	reader := cvs.NewReader("../../../test/fixtures/cvs/branches")
	branches, err := reader.GetBranches()
	if err != nil {
		t.Fatalf("GetBranches failed: %v", err)
	}

	if len(branches) == 0 {
		t.Error("Expected at least one branch")
	}

	// Check for FEATURE_X branch
	found := false
	for _, b := range branches {
		if b == "FEATURE_X" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected FEATURE_X branch")
	}
}

// TestCVSReaderGetTags tests tag extraction
func TestCVSReaderGetTags(t *testing.T) {
	reader := cvs.NewReader("../../../test/fixtures/cvs/tags")
	tags, err := reader.GetTags()
	if err != nil {
		t.Fatalf("GetTags failed: %v", err)
	}

	expectedTags := []string{"RELEASE_1_0", "RELEASE_2_0", "BETA_1"}
	for _, tag := range expectedTags {
		if _, ok := tags[tag]; !ok {
			t.Errorf("Expected tag %s", tag)
		}
	}
}
