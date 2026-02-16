# Git-Migrator System Test

This directory contains a comprehensive end-to-end system test for the git-migrator tool. The test uses Docker to create a realistic CVS repository and migrates it to Git, verifying that all features work correctly.

## Overview

The system test:

1. **Sets up a realistic CVS repository** with:
   - Multiple files (source code, documentation, tests, configuration)
   - Multiple commits by different authors
   - Multiple branches (FEATURE_BRANCH_1, FEATURE_BRANCH_2, BUGFIX_BRANCH)
   - Multiple tags (RELEASE_1_0, RELEASE_1_1, BETA_1_2, SNAPSHOT_20240101)
   - Multiple authors (johndoe, janesmith, bobjohnson)

2. **Runs the migration** from CVS to Git using the git-migrator tool

3. **Verifies the migration** by checking:
   - All branches are migrated correctly
   - All tags are migrated correctly
   - All commits are preserved with correct history
   - Author mapping works correctly
   - All files are migrated
   - Branch content is accurate

## Prerequisites

- Docker installed and running
- Make (optional, for using Makefile targets)
- Sufficient disk space for Docker images and test data

## Quick Start

### Using Make (Recommended)

```bash
# Build and run the complete system test
make system-test-all

# Or run step by step:
make system-test-build  # Build the Docker image
make system-test        # Run the system test
```

### Using Docker Directly

```bash
# Build the Docker image
docker build -f test/system/Dockerfile -t git-migrator-system-test:latest .

# Run the system test
docker run --rm -v $(PWD)/test-results:/workspace/migration-test git-migrator-system-test:latest
```

## Available Make Targets

| Target | Description |
|--------|-------------|
| `make system-test-build` | Build the system test Docker image |
| `make system-test` | Run the full system test in Docker |
| `make system-test-all` | Build and run complete system test (convenience target) |
| `make system-test-interactive` | Run system test interactively for debugging |
| `make system-test-verify` | Run verification script on existing migration |
| `make system-test-clean` | Clean up system test artifacts |

## Test Results

After running the system test, results are saved to the `test-results/` directory:

```
test-results/
├── git-repo/                    # Migrated Git repository
├── migration-report.txt         # Detailed test report
├── verification-report.txt      # Verification results
├── migration-config.yaml        # Configuration used for migration
├── author-map.yaml             # Author mapping file
└── extracted-authors.yaml      # Authors extracted from CVS
```

## What Gets Tested

### Repository Structure
- Git repository is created and valid
- All files are migrated correctly
- Directory structure is preserved

### Branches
- Main branch exists
- Feature branches (feature-1, feature-2)
- Bugfix branch
- Branch-specific files are correct

### Tags
- Release tags (v1.0.0, v1.1.0)
- Beta tag (v1.2.0-beta)
- Snapshot tag (snapshot-20240101)
- Tags point to correct commits

### Commits
- All commits are migrated
- Commit messages are preserved
- Commit history is accurate
- Commit dates are correct

### Authors
- Author mapping works correctly
- All authors are mapped to Git format
- Email addresses are correct

### File Content
- README.md contains expected content
- Source files are accurate
- Documentation files are preserved
- Configuration files are migrated

## Interactive Testing

For debugging or exploring the test environment:

```bash
# Run interactively
make system-test-interactive

# Inside the container:
# - Explore the CVS repository: ls -la /cvs-repo/
# - Check the workspace: ls -la /workspace/
# - Run manual migrations: /app/bin/git-migrator --help
# - View logs: cat /workspace/migration-test/migration-report.txt
```

## Manual Testing

You can also run the test steps manually inside the Docker container:

```bash
# Start interactive session
make system-test-interactive

# Inside the container:

# 1. Setup CVS repository (if not already done)
/usr/local/bin/setup-cvs-repo.sh

# 2. Analyze the CVS repository
/app/bin/git-migrator analyze --source-type cvs --source /cvs-repo/testproject

# 3. Extract authors
/app/bin/git-migrator authors extract --source /cvs-repo/testproject --format yaml

# 4. Create migration config
cat > /workspace/test-config.yaml << 'EOF'
source:
  type: cvs
  path: /cvs-repo/testproject
target:
  path: /workspace/my-git-repo
mapping:
  authors:
    johndoe: John Doe <john.doe@example.com>
    janesmith: Jane Smith <jane.smith@example.com>
    bobjohnson: Bob Johnson <bob.johnson@example.com>
options:
  dryRun: false
  verbose: true
EOF

# 5. Run migration
/app/bin/git-migrator migrate --config /workspace/test-config.yaml --verbose

# 6. Verify results
cd /workspace/my-git-repo
git log --oneline --graph --all
git branch -a
git tag
```

## Verification

The verification script checks multiple aspects of the migration:

```bash
# Run verification on an existing migration
make system-test-verify

# Or manually:
docker run --rm -v $(PWD)/test-results:/workspace/migration-test \
  git-migrator-system-test:latest \
  /usr/local/bin/verify-migration.sh
```

The verification includes:
- Repository validity checks
- Branch existence and content
- Tag existence and validity
- Commit count and messages
- Author mapping
- File existence and content
- Branch-specific content
- Tag integrity

## Expected Output

When the system test runs successfully, you should see:

```
=== Git Migrator System Test ===

Step 1: Setting up CVS repository...
✓ CVS repository already exists at /cvs-repo

Step 2: Creating author mapping...
✓ Author mapping created

Step 3: Creating migration configuration...
✓ Migration configuration created

Step 4: Analyzing CVS repository...
[Analysis output...]

Step 5: Extracting authors from repository...
[Author extraction output...]

Step 6: Running migration...
[Migration output...]

Step 7: Verifying migration results...
✓ Git repository created successfully
✓ Branch 'main' found
✓ Branch 'feature-1' found
✓ Branch 'feature-2' found
✓ Branch 'bugfix' found
✓ Tag 'v1.0.0' found
✓ Tag 'v1.1.0' found
✓ Tag 'v1.2.0-beta' found
✓ Tag 'snapshot-20240101' found
✓ Sufficient commits found (15)
✓ All tests passed!

=== Migration Test Summary ===
✓ ALL TESTS PASSED!

Migration completed successfully!
Git repository is available at: test-results/git-repo
```

## Troubleshooting

### Docker Issues

**Problem**: Docker image fails to build
```
Solution: Ensure Docker is running and you have sufficient disk space
Check Docker logs: docker logs <container-id>
```

**Problem**: Permission denied errors
```
Solution: On Linux, you may need to use sudo or add your user to the docker group
sudo make system-test-all
```

### CVS Issues

**Problem**: CVS repository not created properly
```
Solution: The setup script handles this automatically, but you can manually run:
docker run --rm -it git-migrator-system-test:latest /usr/local/bin/setup-cvs-repo.sh
```

**Problem**: CVS commands fail
```
Solution: CVS requires proper environment variables (CVSROOT) and user setup
These are configured in the Dockerfile automatically
```

### Migration Issues

**Problem**: Migration fails or produces errors
```
Solution: Check the migration report:
cat test-results/migration-report.txt

Run interactively to debug:
make system-test-interactive
```

**Problem**: Missing branches or tags
```
Solution: Verify the CVS repository structure:
docker run --rm -it git-migrator-system-test:latest /bin/bash
ls -la /cvs-repo/testproject/
cvs -d /cvs-repo rlog testproject
```

### Verification Issues

**Problem**: Verification fails
```
Solution: Check the verification report:
cat test-results/verification-report.txt

Manually inspect the Git repository:
cd test-results/git-repo
git log --oneline --graph --all
git branch -a
git tag
```

## Cleaning Up

To remove all test artifacts and Docker images:

```bash
make system-test-clean
```

This removes:
- `test-results/` directory
- Docker images created for system testing

## Test Data Details

The system test creates a CVS repository with:

### Files
- `README.md` - Project documentation
- `CHANGELOG.md` - Version history
- `Makefile` - Build configuration
- `config.h` - Configuration header
- `src/main.c` - Main program
- `src/utils.h` - Utility functions header
- `src/utils.c` - Utility functions implementation
- `src/math.c` - Math operations
- `src/feature1.c` - Feature 1 (on branch)
- `src/feature2.c` - Feature 2 (on branch)
- `src/bugfix.c` - Bugfix (on branch)
- `docs/CONTRIBUTING.md` - Contribution guidelines
- `tests/test_utils.c` - Unit tests

### Branches
- `main` (default)
- `FEATURE_BRANCH_1` → `feature-1`
- `FEATURE_BRANCH_2` → `feature-2`
- `BUGFIX_BRANCH` → `bugfix`

### Tags
- `RELEASE_1_0` → `v1.0.0`
- `RELEASE_1_1` → `v1.1.0`
- `BETA_1_2` → `v1.2.0-beta`
- `SNAPSHOT_20240101` → `snapshot-20240101`

### Authors
- `johndoe` → `John Doe <john.doe@example.com>`
- `janesmith` → `Jane Smith <jane.smith@example.com>`
- `bobjohnson` → `Bob Johnson <bob.johnson@example.com>`
- `testuser` → `Test User <testuser@example.com>`

## CI/CD Integration

To integrate the system test into a CI/CD pipeline:

```yaml
# Example GitLab CI
system_test:
  stage: test
  script:
    - make system-test-all
  artifacts:
    paths:
      - test-results/
    expire_in: 1 week
  tags:
    - docker

# Example GitHub Actions
- name: Run System Test
  run: make system-test-all
  
- name: Upload Test Results
  uses: actions/upload-artifact@v2
  with:
    name: system-test-results
    path: test-results/
```

## Support

For issues or questions:

1. Check the troubleshooting section above
2. Review the test reports in `test-results/`
3. Run interactively: `make system-test-interactive`
4. Check the main project documentation

## License

This system test is part of the git-migrator project and is licensed under the same license as the main project.