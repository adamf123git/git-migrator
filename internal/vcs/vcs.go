// Package vcs provides version control system abstractions for git-migrator.
package vcs

import (
	"time"
)

// Commit represents a single commit in a VCS
type Commit struct {
	Revision string    // VCS-specific revision identifier
	Author   string    // Commit author
	Email    string    // Author email (if available)
	Date     time.Time // Commit timestamp
	Message  string    // Commit message
	Branch   string    // Branch name (empty for trunk/main)
	Files    []FileChange
}

// FileChange represents a file change in a commit
type FileChange struct {
	Path    string // File path
	Action  Action // Add, Modify, Delete
	Content []byte // File content (for Add/Modify)
}

// Action represents the type of file change
type Action int

const (
	ActionModify Action = iota
	ActionAdd
	ActionDelete
)

// VCSReader defines the interface for reading from a VCS repository
type VCSReader interface {
	// Validate checks if the repository is valid and accessible
	Validate() error

	// GetCommits returns an iterator over all commits
	GetCommits() (CommitIterator, error)

	// GetBranches returns a list of branch names
	GetBranches() ([]string, error)

	// GetTags returns a map of tag names to revision identifiers
	GetTags() (map[string]string, error)

	// Close releases any resources
	Close() error
}

// CommitIterator provides iteration over commits
type CommitIterator interface {
	// Next advances to the next commit, returns false when done
	Next() bool

	// Commit returns the current commit
	Commit() *Commit

	// Err returns any error that occurred during iteration
	Err() error
}

// VCSWriter defines the interface for writing to a VCS repository
type VCSWriter interface {
	// Init creates a new repository at the given path
	Init(path string) error

	// ApplyCommit applies a commit to the repository
	ApplyCommit(commit *Commit) error

	// CreateBranch creates a new branch
	CreateBranch(name, revision string) error

	// CreateTag creates a new tag
	CreateTag(name, revision string) error

	// Close releases any resources
	Close() error
}

// RepositoryInfo contains metadata about a repository
type RepositoryInfo struct {
	Path     string
	VCS      string // "cvs", "svn", "git"
	Root     string // Repository root
	Branches []string
	Tags     map[string]string
}
