set -e

# Run Migration Test
# This script runs the full system test for git-migrator

echo "=== Git Migrator System Test ==="
echo ""

# Configuration
CVS_REPO="/cvs-repo"
CVS_MODULE="testproject"
WORK_DIR="/workspace"
TEST_DIR="$WORK_DIR/migration-test"
GIT_REPO="$TEST_DIR/git-repo"
AUTHOR_MAP_FILE="$TEST_DIR/author-map.yaml"
CONFIG_FILE="$TEST_DIR/migration-config.yaml"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored status
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Cleanup previous test runs
echo "=== Cleaning up previous test runs ==="
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"

# Step 1: Setup CVS Repository
print_status "$YELLOW" "Step 1: Setting up CVS repository..."
if [ ! -d "$CVS_REPO" ]; then
    print_status "$YELLOW" "CVS repository not found, creating it..."
    /usr/local/bin/setup-cvs-repo.sh
else
    print_status "$GREEN" "CVS repository already exists at $CVS_REPO"
fi

# Step 2: Create Author Mapping
print_status "$YELLOW" "Step 2: Creating author mapping..."
cat > "$AUTHOR_MAP_FILE" << 'EOF'
johndoe: John Doe <john.doe@example.com>
janesmith: Jane Smith <jane.smith@example.com>
bobjohnson: Bob Johnson <bob.johnson@example.com>
testuser: Test User <testuser@example.com>
EOF
print_status "$GREEN" "Author mapping created at $AUTHOR_MAP_FILE"

# Step 3: Create Migration Configuration
print_status "$YELLOW" "Step 3: Creating migration configuration..."
cat > "$CONFIG_FILE" << EOF
source:
  type: cvs
  path: $CVS_REPO/$CVS_MODULE

target:
  path: $GIT_REPO

mapping:
  authors:
    johndoe: John Doe <john.doe@example.com>
    janesmith: Jane Smith <jane.smith@example.com>
    bobjohnson: Bob Johnson <bob.johnson@example.com>
    testuser: Test User <testuser@example.com>
  branches:
    FEATURE_BRANCH_1: feature-1
    FEATURE_BRANCH_2: feature-2
    BUGFIX_BRANCH: bugfix
  tags:
    RELEASE_1_0: v1.0.0
    RELEASE_1_1: v1.1.0
    BETA_1_2: v1.2.0-beta
    SNAPSHOT_20240101: snapshot-20240101

options:
  dryRun: false
  verbose: true
  chunkSize: 10
EOF
print_status "$GREEN" "Migration configuration created at $CONFIG_FILE"

# Step 4: Analyze the CVS repository
print_status "$YELLOW" "Step 4: Analyzing CVS repository..."
echo ""
/app/bin/git-migrator analyze --source-type cvs --source "$CVS_REPO/$CVS_MODULE"
echo ""

# Step 5: Extract authors
print_status "$YELLOW" "Step 5: Extracting authors from repository..."
echo ""
/app/bin/git-migrator authors extract --source "$CVS_REPO/$CVS_MODULE" --format yaml > "$TEST_DIR/extracted-authors.yaml"
echo ""
print_status "$GREEN" "Authors extracted to $TEST_DIR/extracted-authors.yaml"

# Step 6: Run the migration
print_status "$YELLOW" "Step 6: Running migration..."
echo ""
/app/bin/git-migrator migrate --config "$CONFIG_FILE" --verbose
echo ""

# Step 7: Verify the migration
print_status "$YELLOW" "Step 7: Verifying migration results..."
echo ""

# Check if Git repository was created
if [ ! -d "$GIT_REPO/.git" ]; then
    print_status "$RED" "FAILED: Git repository was not created at $GIT_REPO"
    exit 1
fi
print_status "$GREEN" "✓ Git repository created successfully"

# Check branches
echo ""
echo "Checking branches..."
cd "$GIT_REPO"
BRANCHES=$(git branch -a | grep -v HEAD | wc -l)
echo "Found $BRANCHES branches"

EXPECTED_BRANCHES=("main" "feature-1" "feature-2" "bugfix")
BRANCHES_OK=true

for branch in "${EXPECTED_BRANCHES[@]}"; do
    if git branch -a | grep -q "$branch"; then
        print_status "$GREEN" "✓ Branch '$branch' found"
    else
        print_status "$RED" "✗ Branch '$branch' NOT found"
        BRANCHES_OK=false
    fi
done

# Check tags
echo ""
echo "Checking tags..."
TAGS=$(git tag | wc -l)
echo "Found $TAGS tags"

EXPECTED_TAGS=("v1.0.0" "v1.1.0" "v1.2.0-beta" "snapshot-20240101")
TAGS_OK=true

for tag in "${EXPECTED_TAGS[@]}"; do
    if git tag | grep -q "$tag"; then
        print_status "$GREEN" "✓ Tag '$tag' found"
    else
        print_status "$RED" "✗ Tag '$tag' NOT found"
        TAGS_OK=false
    fi
done

# Check commits
echo ""
echo "Checking commits..."
COMMIT_COUNT=$(git log --all --oneline | wc -l)
echo "Found $COMMIT_COUNT commits"

if [ "$COMMIT_COUNT" -lt 10 ]; then
    print_status "$RED" "✗ Too few commits found (expected at least 10, got $COMMIT_COUNT)"
else
    print_status "$GREEN" "✓ Sufficient commits found ($COMMIT_COUNT)"
fi

# Check files
echo ""
echo "Checking files..."
FILES_OK=true

if [ -f "README.md" ]; then
    print_status "$GREEN" "✓ README.md exists"
else
    print_status "$RED" "✗ README.md NOT found"
    FILES_OK=false
fi

if [ -f "CHANGELOG.md" ]; then
    print_status "$GREEN" "✓ CHANGELOG.md exists"
else
    print_status "$RED" "✗ CHANGELOG.md NOT found"
    FILES_OK=false
fi

if [ -f "src/main.c" ]; then
    print_status "$GREEN" "✓ src/main.c exists"
else
    print_status "$RED" "✗ src/main.c NOT found"
    FILES_OK=false
fi

if [ -f "docs/CONTRIBUTING.md" ]; then
    print_status "$GREEN" "✓ docs/CONTRIBUTING.md exists"
else
    print_status "$RED" "✗ docs/CONTRIBUTING.md NOT found"
    FILES_OK=false
fi

# Check authors
echo ""
echo "Checking author mapping..."
AUTHORS=$(git log --all --format='%an' | sort -u)
echo "Authors in Git repository:"
echo "$AUTHORS"

EXPECTED_AUTHORS=("John Doe" "Jane Smith" "Bob Johnson")
AUTHORS_OK=true

for author in "${EXPECTED_AUTHORS[@]}"; do
    if echo "$AUTHORS" | grep -q "$author"; then
        print_status "$GREEN" "✓ Author '$author' found"
    else
        print_status "$RED" "✗ Author '$author' NOT found"
        AUTHORS_OK=false
    fi
done

# Check commit messages
echo ""
echo "Sample commit messages:"
git log --oneline --all -10

echo ""

# Step 8: Test checkout of branches
print_status "$YELLOW" "Step 8: Testing branch checkouts..."

git checkout feature-1 2>/dev/null
if [ -f "src/feature1.c" ]; then
    print_status "$GREEN" "✓ feature-1 branch has feature1.c"
else
    print_status "$RED" "✗ feature-1 branch missing feature1.c"
fi

git checkout feature-2 2>/dev/null
if [ -f "src/feature2.c" ]; then
    print_status "$GREEN" "✓ feature-2 branch has feature2.c"
else
    print_status "$RED" "✗ feature-2 branch missing feature2.c"
fi

git checkout bugfix 2>/dev/null
if [ -f "src/bugfix.c" ]; then
    print_status "$GREEN" "✓ bugfix branch has bugfix.c"
else
    print_status "$RED" "✗ bugfix branch missing bugfix.c"
fi

git checkout main 2>/dev/null

# Step 9: Generate test report
print_status "$YELLOW" "Step 9: Generating test report..."
REPORT_FILE="$TEST_DIR/migration-report.txt"

cat > "$REPORT_FILE" << EOF
Git Migrator System Test Report
================================

Date: $(date)
CVS Repository: $CVS_REPO/$CVS_MODULE
Git Repository: $GIT_REPO

Results:
--------
Branches: $([ "$BRANCHES_OK" = true ] && echo "PASS" || echo "FAIL")
Tags: $([ "$TAGS_OK" = true ] && echo "PASS" || echo "FAIL")
Files: $([ "$FILES_OK" = true ] && echo "PASS" || echo "FAIL")
Authors: $([ "$AUTHORS_OK" = true ] && echo "PASS" || echo "FAIL")
Commit Count: $COMMIT_COUNT (minimum expected: 10)

Branches Found:
$(git branch -a)

Tags Found:
$(git tag)

Authors Found:
$AUTHORS

Recent Commits:
$(git log --oneline --all -10)

Overall Status: $([ "$BRANCHES_OK" = true ] && [ "$TAGS_OK" = true ] && [ "$FILES_OK" = true ] && [ "$AUTHORS_OK" = true ] && echo "SUCCESS" || echo "PARTIAL SUCCESS")
EOF

print_status "$GREEN" "Test report saved to $REPORT_FILE"

# Final status
echo ""
echo "=== Migration Test Summary ==="
echo ""

OVERALL_OK=true

if [ "$BRANCHES_OK" = true ] && [ "$TAGS_OK" = true ] && [ "$FILES_OK" = true ] && [ "$AUTHORS_OK" = true ]; then
    print_status "$GREEN" "✓ ALL TESTS PASSED!"
    echo ""
    print_status "$GREEN" "Migration completed successfully!"
    print_status "$GREEN" "Git repository is available at: $GIT_REPO"
    echo ""
    print_status "$GREEN" "To explore the migrated repository:"
    print_status "$GREEN" "  cd $GIT_REPO"
    print_status "$GREEN" "  git log --oneline --graph --all"
    print_status "$GREEN" "  git branch -a"
    print_status "$GREEN" "  git tag"
    exit 0
else
    print_status "$YELLOW" "⚠ SOME TESTS FAILED OR INCOMPLETE"
    echo ""
    print_status "$YELLOW" "Migration completed with some issues."
    print_status "$YELLOW" "Please check the report at: $REPORT_FILE"
    echo ""

    [ "$BRANCHES_OK" = false ] && print_status "$RED" "✗ Branch migration incomplete"
    [ "$TAGS_OK" = false ] && print_status "$RED" "✗ Tag migration incomplete"
    [ "$FILES_OK" = false ] && print_status "$RED" "✗ File migration incomplete"
    [ "$AUTHORS_OK" = false ] && print_status "$RED" "✗ Author mapping incomplete"

    exit 1
fi
