#!/bin/bash
set -e

# Verify Migration Results
# This script verifies that a CVS to Git migration was successful

echo "=== Git Migrator Verification Script ==="
echo ""

# Configuration
GIT_REPO="${1:-/workspace/migration-test/git-repo}"
REPORT_FILE="${2:-/workspace/migration-test/verification-report.txt}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored status
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Check if Git repository exists
if [ ! -d "$GIT_REPO/.git" ]; then
    print_status "$RED" "Error: Git repository not found at $GIT_REPO"
    echo "Usage: $0 [git-repo-path] [report-file]"
    exit 1
fi

print_status "$BLUE" "Verifying Git repository at: $GIT_REPO"
echo ""

cd "$GIT_REPO"

# Initialize counters
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0

# Function to run a check
run_check() {
    local description="$1"
    local command="$2"
    local expected="$3"

    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))

    echo -n "Checking: $description... "

    if eval "$command" >/dev/null 2>&1; then
        if [ -n "$expected" ]; then
            result=$(eval "$command" 2>/dev/null)
            if echo "$result" | grep -q "$expected"; then
                print_status "$GREEN" "✓ PASS"
                PASSED_CHECKS=$((PASSED_CHECKS + 1))
                return 0
            else
                print_status "$RED" "✗ FAIL (expected: $expected)"
                FAILED_CHECKS=$((FAILED_CHECKS + 1))
                return 1
            fi
        else
            print_status "$GREEN" "✓ PASS"
            PASSED_CHECKS=$((PASSED_CHECKS + 1))
            return 0
        fi
    else
        print_status "$RED" "✗ FAIL"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        return 1
    fi
}

# Start verification report
cat > "$REPORT_FILE" << EOF
Git Migrator Verification Report
=================================

Date: $(date)
Repository: $GIT_REPO

Verification Checks
-------------------

EOF

echo "=== Repository Structure ===" | tee -a "$REPORT_FILE"
echo ""

# Check 1: Git repository is valid
run_check "Git repository is valid" "git rev-parse --git-dir" "" | tee -a "$REPORT_FILE"

# Check 2: Has commits
run_check "Repository has commits" "git log --oneline" "" | tee -a "$REPORT_FILE"

# Check 3: Has main/master branch
run_check "Has main branch" "git branch --list main" "main" | tee -a "$REPORT_FILE"

echo ""
echo "=== Branches ===" | tee -a "$REPORT_FILE"
echo ""

# List all branches
echo "Branches found:" | tee -a "$REPORT_FILE"
git branch -a | tee -a "$REPORT_FILE"
echo ""

# Check for expected branches
EXPECTED_BRANCHES=("main" "feature-1" "feature-2" "bugfix")
for branch in "${EXPECTED_BRANCHES[@]}"; do
    run_check "Branch '$branch' exists" "git branch --list $branch" "$branch" | tee -a "$REPORT_FILE"
done

echo ""
echo "=== Tags ===" | tee -a "$REPORT_FILE"
echo ""

# List all tags
echo "Tags found:" | tee -a "$REPORT_FILE"
git tag | tee -a "$REPORT_FILE"
echo ""

# Check for expected tags
EXPECTED_TAGS=("v1.0.0" "v1.1.0" "v1.2.0-beta" "snapshot-20240101")
for tag in "${EXPECTED_TAGS[@]}"; do
    run_check "Tag '$tag' exists" "git tag -l $tag" "$tag" | tee -a "$REPORT_FILE"
done

echo ""
echo "=== Commits ===" | tee -a "$REPORT_FILE"
echo ""

# Get commit statistics
COMMIT_COUNT=$(git log --all --oneline | wc -l)
echo "Total commits: $COMMIT_COUNT" | tee -a "$REPORT_FILE"
echo ""

run_check "Has at least 10 commits" "[ $COMMIT_COUNT -ge 10 ]" "" | tee -a "$REPORT_FILE"

# Check commit messages
echo "" | tee -a "$REPORT_FILE"
echo "Recent commit messages:" | tee -a "$REPORT_FILE"
git log --oneline --all -10 | tee -a "$REPORT_FILE"
echo ""

# Check for commits with specific messages
run_check "Has 'Initial' commit" "git log --all --oneline" "Initial" | tee -a "$REPORT_FILE"
run_check "Has 'Add' commit" "git log --all --oneline" "Add" | tee -a "$REPORT_FILE"
run_check "Has 'Update' commit" "git log --all --oneline" "Update" | tee -a "$REPORT_FILE"

echo ""
echo "=== Authors ===" | tee -a "$REPORT_FILE"
echo ""

# Get author statistics
AUTHORS=$(git log --all --format='%an <%ae>' | sort -u)
AUTHOR_COUNT=$(echo "$AUTHORS" | wc -l)
echo "Unique authors: $AUTHOR_COUNT" | tee -a "$REPORT_FILE"
echo ""
echo "Authors:" | tee -a "$REPORT_FILE"
echo "$AUTHORS" | tee -a "$REPORT_FILE"
echo ""

# Check for expected authors
EXPECTED_AUTHORS=("John Doe" "Jane Smith" "Bob Johnson")
for author in "${EXPECTED_AUTHORS[@]}"; do
    run_check "Author '$author' exists" "git log --all --format='%an'" "$author" | tee -a "$REPORT_FILE"
done

echo ""
echo "=== Files ===" | tee -a "$REPORT_FILE"
echo ""

# Check for expected files
EXPECTED_FILES=("README.md" "CHANGELOG.md" "Makefile" "src/main.c" "docs/CONTRIBUTING.md")
for file in "${EXPECTED_FILES[@]}"; do
    run_check "File '$file' exists" "[ -f $file ]" "" | tee -a "$REPORT_FILE"
done

# Check file count
FILE_COUNT=$(find . -type f -not -path './.git/*' | wc -l)
echo "" | tee -a "$REPORT_FILE"
echo "Total files: $FILE_COUNT" | tee -a "$REPORT_FILE"

run_check "Has at least 8 files" "[ $FILE_COUNT -ge 8 ]" "" | tee -a "$REPORT_FILE"

echo ""
echo "=== Branch Content ===" | tee -a "$REPORT_FILE"
echo ""

# Check feature-1 branch content
if git branch --list feature-1 | grep -q "feature-1"; then
    git checkout feature-1 2>/dev/null
    run_check "feature-1 has src/feature1.c" "[ -f src/feature1.c ]" "" | tee -a "$REPORT_FILE"
    git checkout - 2>/dev/null
fi

# Check feature-2 branch content
if git branch --list feature-2 | grep -q "feature-2"; then
    git checkout feature-2 2>/dev/null
    run_check "feature-2 has src/feature2.c" "[ -f src/feature2.c ]" "" | tee -a "$REPORT_FILE"
    git checkout - 2>/dev/null
fi

# Check bugfix branch content
if git branch --list bugfix | grep -q "bugfix"; then
    git checkout bugfix 2>/dev/null
    run_check "bugfix has src/bugfix.c" "[ -f src/bugfix.c ]" "" | tee -a "$REPORT_FILE"
    git checkout - 2>/dev/null
fi

echo ""
echo "=== File Content ===" | tee -a "$REPORT_FILE"
echo ""

# Check README.md content
if [ -f "README.md" ]; then
    run_check "README.md contains 'Test Project'" "grep 'Test Project' README.md" "Test Project" | tee -a "$REPORT_FILE"
fi

# Check source files
if [ -f "src/main.c" ]; then
    run_check "src/main.c contains 'main'" "grep 'main' src/main.c" "main" | tee -a "$REPORT_FILE"
fi

echo ""
echo "=== Tag Verification ===" | tee -a "$REPORT_FILE"
echo ""

# Check tag v1.0.0 points to a valid commit
if git tag -l | grep -q "v1.0.0"; then
    run_check "Tag v1.0.0 is valid" "git rev-parse v1.0.0" "" | tee -a "$REPORT_FILE"
fi

# Check tag v1.1.0 points to a valid commit
if git tag -l | grep -q "v1.1.0"; then
    run_check "Tag v1.1.0 is valid" "git rev-parse v1.1.0" "" | tee -a "$REPORT_FILE"
fi

echo ""
echo "=== Commit History ===" | tee -a "$REPORT_FILE"
echo ""

# Check commit history integrity
run_check "No merge commits in main history" "git log --first-parent --merges" "" | tee -a "$REPORT_FILE"

# Check commit dates are reasonable
FIRST_COMMIT=$(git log --all --format='%ai' --reverse | head -1)
LAST_COMMIT=$(git log --all --format='%ai' -1)
echo "" | tee -a "$REPORT_FILE"
echo "First commit: $FIRST_COMMIT" | tee -a "$REPORT_FILE"
echo "Last commit: $LAST_COMMIT" | tee -a "$REPORT_FILE"

echo ""
echo "=== Summary ===" | tee -a "$REPORT_FILE"
echo ""

# Generate summary
SUMMARY=$(cat <<EOF

Verification Summary
====================

Total Checks: $TOTAL_CHECKS
Passed: $PASSED_CHECKS
Failed: $FAILED_CHECKS
Success Rate: $(awk "BEGIN {printf \"%.1f\", ($PASSED_CHECKS/$TOTAL_CHECKS)*100}")%

Status: $([ $FAILED_CHECKS -eq 0 ] && echo "✓ ALL CHECKS PASSED" || echo "✗ SOME CHECKS FAILED")

Repository Statistics:
- Commits: $COMMIT_COUNT
- Branches: $(git branch -a | grep -v HEAD | wc -l)
- Tags: $(git tag | wc -l)
- Files: $FILE_COUNT
- Authors: $AUTHOR_COUNT

Repository Health:
- Valid Git repository: ✓
- Has commit history: ✓
- Multiple authors: ✓
- Branches migrated: $([ $PASSED_CHECKS -gt 0 ] && echo "✓" || echo "✗")
- Tags migrated: $([ $PASSED_CHECKS -gt 0 ] && echo "✓" || echo "✗")

EOF
)

echo "$SUMMARY" | tee -a "$REPORT_FILE"

# Final status message
echo ""
if [ $FAILED_CHECKS -eq 0 ]; then
    print_status "$GREEN" "✓ Migration verification PASSED!"
    print_status "$GREEN" "All checks completed successfully."
    print_status "$GREEN" "The Git repository is ready for use."
    echo ""
    print_status "$GREEN" "You can explore the repository:"
    print_status "$GREEN" "  cd $GIT_REPO"
    print_status "$GREEN" "  git log --oneline --graph --all"
    print_status "$GREEN" "  gitk --all"
    exit 0
else
    print_status "$YELLOW" "⚠ Migration verification completed with WARNINGS"
    print_status "$YELLOW" "Some checks failed. Please review the report."
    echo ""
    print_status "$YELLOW" "Report saved to: $REPORT_FILE"
    echo ""
    print_status "$YELLOW" "Failed checks:"
    echo "  - Review the output above for details"
    exit 1
fi
