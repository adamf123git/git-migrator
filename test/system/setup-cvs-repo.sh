#!/bin/bash
set -e

# Setup CVS Repository for System Testing
# This script creates a comprehensive CVS repository with branches, tags, and multiple commits

echo "=== Setting up CVS Repository for System Testing ==="

# Configuration
MODULE_NAME="testproject"
WORK_DIR="/workspace/cvs-work"
EXPORT_DIR="/workspace/cvs-export"

# Ensure CVSROOT is set
if [ -z "$CVSROOT" ]; then
    echo "Error: CVSROOT not set"
    exit 1
fi

echo "CVSROOT: $CVSROOT"
echo "Module: $MODULE_NAME"

# Create working directory
mkdir -p "$WORK_DIR"
cd "$WORK_DIR"

# Import initial project
echo "=== Creating initial project structure ==="
mkdir -p "$MODULE_NAME"
cd "$MODULE_NAME"

# Create initial files
cat > README.md << 'EOF'
# Test Project

This is a test project for CVS to Git migration.

## Features
- Multiple files
- Multiple authors
- Branches and tags
- Various commit types
EOF

cat > src/main.c << 'EOF'
#include <stdio.h>

int main() {
    printf("Hello, World!\n");
    return 0;
}
EOF
mkdir -p src

cat > src/utils.h << 'EOF'
#ifndef UTILS_H
#define UTILS_H

int add(int a, int b);
int subtract(int a, int b);

#endif
EOF

cat > src/utils.c << 'EOF'
#include "utils.h"

int add(int a, int b) {
    return a + b;
}

int subtract(int a, int b) {
    return a - b;
}
EOF

cat > Makefile << 'EOF'
CC=gcc
CFLAGS=-Wall -g

all: main

main: src/main.c src/utils.c
	$(CC) $(CFLAGS) -o main src/main.c src/utils.c

clean:
	rm -f main *.o

test: main
	./main
EOF

# Import to CVS
echo "=== Importing project to CVS ==="
cd ..
cvs import -m "Initial import" "$MODULE_NAME" testuser start

# Checkout the module to work with it
echo "=== Checking out module ==="
cd "$WORK_DIR"
cvs checkout "$MODULE_NAME"
cd "$MODULE_NAME"

# Function to create a commit by a specific author
create_commit() {
    local author="$1"
    local message="$2"
    local date="$3"

    echo "Creating commit by $author: $message"

    # Set CVS user
    export CVSUSER="$author"

    # Modify a file to create a change
    echo "// Modified by $author at $date" >> README.md

    # Commit with specified date
    cvs commit -m "$message" README.md
}

# Create commits by different authors
echo "=== Creating commits by different authors ==="

# Commit 1: John Doe
export CVSUSER="johndoe"
echo "## Documentation Update" >> README.md
echo "Added by John Doe for better documentation." >> README.md
cvs add -m "johndoe" README.md 2>/dev/null || true
cvs commit -m "Improve documentation" README.md

# Add new file
cat > docs/CONTRIBUTING.md << 'EOF'
# Contributing

Please read the guidelines before contributing.

1. Fork the repository
2. Create a feature branch
3. Submit a pull request
EOF
mkdir -p docs
cvs add docs
cvs add docs/CONTRIBUTING.md
cvs commit -m "Add contributing guidelines"

# Commit 2: Jane Smith
export CVSUSER="janesmith"
cat > src/math.c << 'EOF'
#include <stdio.h>

int multiply(int a, int b) {
    return a * b;
}

int divide(int a, int b) {
    if (b == 0) return 0;
    return a / b;
}
EOF
cvs add src/math.c
cvs commit -m "Add math operations"

# Update main.c to use math operations
cat > src/main.c << 'EOF'
#include <stdio.h>
#include "utils.h"

extern int multiply(int a, int b);
extern int divide(int a, int b);

int main() {
    printf("Hello, World!\n");
    printf("5 + 3 = %d\n", add(5, 3));
    printf("10 - 4 = %d\n", subtract(10, 4));
    printf("6 * 7 = %d\n", multiply(6, 7));
    printf("20 / 4 = %d\n", divide(20, 4));
    return 0;
}
EOF
cvs commit -m "Update main to use all operations"

# Commit 3: Bob Johnson
export CVSUSER="bobjohnson"
cat > config.h << 'EOF'
#ifndef CONFIG_H
#define CONFIG_H

#define VERSION "1.0.0"
#define DEBUG 1

#endif
EOF
cvs add config.h
cvs commit -m "Add configuration header"

# Update Makefile
cat > Makefile << 'EOF'
CC=gcc
CFLAGS=-Wall -g -DVERSION=\"1.0.0\"

all: main

main: src/main.c src/utils.c src/math.c
	$(CC) $(CFLAGS) -o main src/main.c src/utils.c src/math.c

clean:
	rm -f main *.o

test: main
	./main

install: main
	install -m 755 main /usr/local/bin/testproject
EOF
cvs commit -m "Update Makefile with version info"

# Create branches
echo "=== Creating branches ==="

# Create a feature branch
cvs tag -b FEATURE_BRANCH_1
echo "Created branch: FEATURE_BRANCH_1"

# Create another feature branch
cvs tag -b FEATURE_BRANCH_2
echo "Created branch: FEATURE_BRANCH_2"

# Work on FEATURE_BRANCH_1
cvs update -r FEATURE_BRANCH_1
cat > src/feature1.c << 'EOF'
#include <stdio.h>

void feature1() {
    printf("Feature 1 implementation\n");
}
EOF
cvs add src/feature1.c
cvs commit -m "Add feature 1 on branch"

# Work on FEATURE_BRANCH_2
cvs update -r FEATURE_BRANCH_2
cat > src/feature2.c << 'EOF'
#include <stdio.h>

void feature2() {
    printf("Feature 2 implementation\n");
}
EOF
cvs add src/feature2.c
cvs commit -m "Add feature 2 on branch"

# Go back to main branch
cvs update -A

# More commits on main branch
export CVSUSER="johndoe"
cat > CHANGELOG.md << 'EOF'
# Changelog

## 1.0.0 - Initial Release
- Basic functionality
- Multiple operations
- Documentation
EOF
cvs add CHANGELOG.md
cvs commit -m "Add changelog"

# Create tags
echo "=== Creating tags ==="

# Create release tags
cvs tag RELEASE_1_0
echo "Created tag: RELEASE_1_0"

# Make some more changes
export CVSUSER="janesmith"
cat >> README.md << 'EOF'

## Installation

Run `make install` to install the program.

## Usage

Run `./main` or `testproject` after installation.
EOF
cvs commit -m "Update README with installation instructions"

# Create another tag
cvs tag RELEASE_1_1
echo "Created tag: RELEASE_1_1"

# Create a beta tag
cvs tag BETA_1_2
echo "Created tag: BETA_1_2"

# Add more complex file changes
export CVSUSER="bobjohnson"
cat > tests/test_utils.c << 'EOF'
#include <stdio.h>
#include "../src/utils.h"

int main() {
    printf("Testing utils...\n");

    if (add(2, 3) != 5) {
        printf("FAIL: add(2, 3)\n");
        return 1;
    }

    if (subtract(5, 3) != 2) {
        printf("FAIL: subtract(5, 3)\n");
        return 1;
    }

    printf("All tests passed!\n");
    return 0;
}
EOF
mkdir -p tests
cvs add tests tests/test_utils.c
cvs commit -m "Add unit tests"

# Update Makefile to include tests
cat > Makefile << 'EOF'
CC=gcc
CFLAGS=-Wall -g -DVERSION=\"1.0.0\"

all: main test

main: src/main.c src/utils.c src/math.c
	$(CC) $(CFLAGS) -o main src/main.c src/utils.c src/math.c

test: tests/test_utils.c src/utils.c
	$(CC) $(CFLAGS) -o test_utils tests/test_utils.c src/utils.c

clean:
	rm -f main test_utils *.o

run-test: test
	./test_utils

install: main
	install -m 755 main /usr/local/bin/testproject

.PHONY: all clean test run-test install
EOF
cvs commit -m "Add test target to Makefile"

# Create a bugfix branch
cvs tag -b BUGFIX_BRANCH
echo "Created branch: BUGFIX_BRANCH"

# Work on bugfix branch
cvs update -r BUGFIX_BRANCH
cat > src/bugfix.c << 'EOF'
#include <stdio.h>

void apply_bugfix() {
    printf("Bugfix applied\n");
}
EOF
cvs add src/bugfix.c
cvs commit -m "Add bugfix on branch"

# Back to main
cvs update -A

# Final commits
export CVSUSER="johndoe"
echo "" >> CHANGELOG.md
echo "## 1.1.0 - Current Development" >> CHANGELOG.md
echo "- Added unit tests" >> CHANGELOG.md
echo "- Improved documentation" >> CHANGELOG.md
cvs commit -m "Update changelog for 1.1.0"

# Create final tag
cvs tag SNAPSHOT_20240101
echo "Created tag: SNAPSHOT_20240101"

# Export the repository for migration testing
echo "=== Exporting repository for migration ==="
mkdir -p "$EXPORT_DIR"
cd "$WORK_DIR"
cvs export -r HEAD -d "$EXPORT_DIR/cvs-repo" "$MODULE_NAME"

echo "=== CVS Repository Setup Complete ==="
echo ""
echo "Repository Statistics:"
echo "  - Module: $MODULE_NAME"
echo "  - Branches: FEATURE_BRANCH_1, FEATURE_BRANCH_2, BUGFIX_BRANCH"
echo "  - Tags: RELEASE_1_0, RELEASE_1_1, BETA_1_2, SNAPSHOT_20240101"
echo "  - Authors: johndoe, janesmith, bobjohnson"
echo "  - Files: Multiple source, doc, and config files"
echo "  - Location: $CVSROOT"
echo "  - Export location: $EXPORT_DIR/cvs-repo"
echo ""
echo "The repository is ready for migration testing!"

# Create a summary file
cat > "$EXPORT_DIR/repository-info.txt" << EOF
CVS Repository Test Data
========================

Module: $MODULE_NAME
Location: $CVSROOT

Branches:
  - FEATURE_BRANCH_1
  - FEATURE_BRANCH_2
  - BUGFIX_BRANCH

Tags:
  - RELEASE_1_0
  - RELEASE_1_1
  - BETA_1_2
  - SNAPSHOT_20240101

Authors:
  - johndoe (John Doe)
  - janesmith (Jane Smith)
  - bobjohnson (Bob Johnson)

Files:
  - README.md
  - CHANGELOG.md
  - Makefile
  - config.h
  - src/main.c
  - src/utils.h
  - src/utils.c
  - src/math.c
  - src/feature1.c (on branch)
  - src/feature2.c (on branch)
  - src/bugfix.c (on branch)
  - docs/CONTRIBUTING.md
  - tests/test_utils.c

Expected Migration Results:
  - All branches should be migrated
  - All tags should be migrated
  - All commit history should be preserved
  - Author mapping should work correctly
  - File changes should be accurate
EOF

echo "Repository information saved to: $EXPORT_DIR/repository-info.txt"
