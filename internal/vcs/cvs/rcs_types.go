package cvs

import (
	"strings"
	"time"
)

// RCSFile represents a parsed RCS file
type RCSFile struct {
	Head        string
	Branch      string
	Access      []string
	Symbols     map[string]string
	Locks       map[string]string
	StrictLocks bool
	Comment     string
	Description string
	Deltas      map[string]*Delta
	DeltaOrder  []string // Order of deltas as they appear
}

// Delta represents a single revision in an RCS file
type Delta struct {
	Revision string
	Date     time.Time
	Author   string
	State    string
	Branches []string
	Next     string
	Log      string
	Text     string
}

// Commit represents a commit extracted from RCS deltas
type Commit struct {
	Revision string
	Author   string
	Date     time.Time
	Message  string
	Branch   string // Empty for trunk
}

// GetCommits returns commits in reverse chronological order
func (r *RCSFile) GetCommits() []*Commit {
	var commits []*Commit
	seen := make(map[string]bool)

	// Helper to add commits recursively
	var addCommit func(rev string, branch string)
	addCommit = func(rev string, branch string) {
		if rev == "" || seen[rev] {
			return
		}
		seen[rev] = true

		delta := r.Deltas[rev]
		if delta == nil {
			return
		}

		// Add this commit
		commits = append(commits, &Commit{
			Revision: rev,
			Author:   delta.Author,
			Date:     delta.Date,
			Message:  delta.Log,
			Branch:   branch,
		})

		// Add branches from this commit
		for _, branchRev := range delta.Branches {
			// Find branch symbol
			branchName := ""
			for sym, symRev := range r.Symbols {
				if symRev == branchRev || isBranchPrefix(symRev, branchRev) {
					branchName = sym
					break
				}
			}
			addCommit(branchRev, branchName)
		}

		// Add next (previous revision)
		addCommit(delta.Next, branch)
	}

	// Start from head
	addCommit(r.Head, "")

	return commits
}

// GetBranches returns the list of branch names
func (r *RCSFile) GetBranches() []string {
	var branches []string
	for sym, rev := range r.Symbols {
		// Branch numbers have odd number of components (e.g., 1.2.0.2)
		if isBranchNumber(rev) {
			branches = append(branches, sym)
		}
	}
	return branches
}

// GetTags returns the list of tag names (symbols pointing to trunk revisions)
func (r *RCSFile) GetTags() map[string]string {
	tags := make(map[string]string)
	for sym, rev := range r.Symbols {
		if !isBranchNumber(rev) {
			tags[sym] = rev
		}
	}
	return tags
}

func isBranchNumber(rev string) bool {
	// Magic branch numbers have ".0." in them (e.g., 1.2.0.2)
	// Regular branch commits have 4+ components without .0. (e.g., 1.2.2.1)
	// Trunk revisions have 2 components (e.g., 1.3)
	if strings.Contains(rev, ".0.") {
		return true // Magic branch number
	}
	
	// Count dots to determine component count
	dots := 0
	for _, c := range rev {
		if c == '.' {
			dots++
		}
	}
	
	// 4+ components (3+ dots) means it's on a branch
	// 2 components (1 dot) means it's on trunk
	return dots >= 3
}

func isBranchPrefix(branchNum, rev string) bool {
	// Check if rev starts with branchNum prefix
	if len(rev) < len(branchNum) {
		return false
	}
	return rev[:len(branchNum)] == branchNum
}
