# Migration Guide

This comprehensive guide covers everything you need to know about migrating repositories to Git using Git-Migrator. Whether you're migrating a single small project or a large enterprise repository with years of history, this guide will help you succeed.

## 📋 Table of Contents

1. [Migration Planning](#migration-planning)
2. [Pre-Migration Checklist](#pre-migration-checklist)
3. [Migration Strategies](#migration-strategies)
4. [Step-by-Step Migration](#step-by-step-migration)
5. [Large Repository Migration](#large-repository-migration)
6. [Multiple Repositories](#multiple-repositories)
7. [Post-Migration Tasks](#post-migration-tasks)
8. [Verification](#verification)
9. [Rollback Procedures](#rollback-procedures)
10. [Troubleshooting](#troubleshooting)

## Migration Planning

### Assess Your Repository

Before starting, gather information about your source repository:

```bash
# Analyze CVS repository
git-migrator analyze --source-type cvs --source /path/to/cvs/repo

# This provides:
# - Total commit count
# - Branch count and names
# - Tag count and names
# - Unique authors count
# - Repository size
# - File count
# - Date range of commits
```

**Key Metrics to Consider:**

| Metric | Small | Medium | Large | Enterprise |
|--------|-------|--------|-------|------------|
| Commits | <1,000 | 1,000-10,000 | 10,000-100,000 | >100,000 |
| Branches | <10 | 10-50 | 50-200 | >200 |
| Authors | <10 | 10-50 | 50-200 | >200 |
| History | <1 year | 1-5 years | 5-10 years | >10 years |
| Size | <100 MB | 100 MB - 1 GB | 1 GB - 10 GB | >10 GB |

### Timeline Estimation

Estimate migration time based on repository size:

- **Small (<1,000 commits)**: 5-15 minutes
- **Medium (1,000-10,000 commits)**: 15-60 minutes
- **Large (10,000-100,000 commits)**: 1-6 hours
- **Enterprise (>100,000 commits)**: 6-24 hours

**Factors affecting time:**
- Disk I/O speed
- CPU performance
- Network speed (if accessing remote CVS)
- Complexity of file changes
- Number of branches/tags

### Resource Requirements

Ensure adequate resources:

| Component | Minimum | Recommended | Enterprise |
|-----------|---------|-------------|------------|
| RAM | 2 GB | 4 GB | 8+ GB |
| Disk Space | 3x repo size | 5x repo size | 10x repo size |
| CPU | 2 cores | 4 cores | 8+ cores |

## Pre-Migration Checklist

### ✅ Source Repository

- [ ] Verify CVS repository is accessible
- [ ] Confirm read permissions on all files
- [ ] Check for locked files or active checkouts
- [ ] Document any custom CVS configuration
- [ ] Identify all modules to migrate
- [ ] Note any symbolic links or special files

### ✅ Target Environment

- [ ] Git 2.0 or later installed
- [ ] Sufficient disk space (5x source size)
- [ ] Write permissions to target directory
- [ ] Network access (if pushing to remote)
- [ ] SSH keys configured (for remote push)

### ✅ Author Information

- [ ] Extract complete author list
- [ ] Create comprehensive author mapping
- [ ] Verify email addresses
- [ ] Handle generic accounts (buildbot, root, etc.)

### ✅ Branch and Tag Strategy

- [ ] Document all branch names
- [ ] Document all tag names
- [ ] Create branch mapping plan
- [ ] Create tag mapping plan
- [ ] Identify branches to skip (if any)

### ✅ Configuration

- [ ] Create migration configuration file
- [ ] Validate configuration syntax
- [ ] Test with dry run
- [ ] Save configuration to version control

## Migration Strategies

### Strategy 1: Big Bang Migration

Migrate everything at once, switch over immediately.

**Pros:**
- Clean cutover
- Simple process
- Quick completion

**Cons:**
- Higher risk
- Requires coordination
- Limited testing time

**Best for:**
- Small teams
- Active development can pause
- Simple repository structure

**Process:**
```bash
# 1. Prepare configuration
git-migrator validate --config config.yaml

# 2. Perform dry run
git-migrator migrate --config config.yaml --dry-run

# 3. Announce freeze
# "Development frozen for migration starting [time]"

# 4. Perform migration
git-migrator migrate --config config.yaml

# 5. Verify
git-migrator verify --config config.yaml

# 6. Announce completion
# "Migration complete, Git repository ready at [URL]"
```

### Strategy 2: Parallel Running

Run CVS and Git in parallel during transition period.

**Pros:**
- Lower risk
- Extended testing period
- Gradual transition

**Cons:**
- More complex
- Requires synchronization
- Potential for confusion

**Best for:**
- Large teams
- Critical systems
- Complex workflows

**Process:**
```bash
# 1. Migrate to Git
git-migrator migrate --config config.yaml

# 2. Push to shared repository
cd /path/to/git/repo
git remote add origin git@github.com:org/repo.git
git push -u origin --all
git push origin --tags

# 3. Team evaluates Git repository
# Developers test cloning, branching, committing

# 4. Parallel development period
# Some use CVS, some use Git

# 5. Final cutover
# Announce CVS read-only, Git primary
```

### Strategy 3: Incremental Migration

Migrate modules or subsystems one at a time.

**Pros:**
- Lowest risk
- Learn and adapt
- Spread workload

**Cons:**
- Longest timeline
- More coordination needed
- Complex dependencies

**Best for:**
- Very large repositories
- Modular architecture
- Multiple teams

**Process:**
```bash
# 1. Identify modules
MODULES="core api web tools docs"

# 2. Migrate each module
for module in $MODULES; do
  echo "Migrating $module..."
  
  # Create module-specific config
  cat > "config-${module}.yaml" << EOF
source:
  type: cvs
  path: /cvs/main
  module: $module
target:
  type: git
  path: /git/${module}.git
# ... rest of config
EOF
  
  # Migrate
  git-migrator migrate --config "config-${module}.yaml"
  
  # Verify
  if [ $? -eq 0 ]; then
    echo "✓ $module migrated"
  else
    echo "✗ $module failed"
    exit 1
  fi
done

# 3. Create monorepo or separate repos
# Depends on your organization's needs
```

## Step-by-Step Migration

### Phase 1: Preparation (Day 1-2)

#### Day 1: Analysis and Planning

```bash
# Morning: Analyze repository
git-migrator analyze --source-type cvs --source /path/to/cvs > analysis.txt

# Review analysis
cat analysis.txt

# Afternoon: Extract authors
git-migrator authors extract --source /path/to/cvs > authors.txt

# Create author mapping
cat > author-mapping.yaml << 'EOF'
# Author mapping
cvsuser1: "Full Name <email@example.com>"
cvsuser2: "Full Name <email@example.com>"
# ... add all authors
EOF
```

#### Day 2: Configuration

```bash
# Morning: Create configuration
cat > migration-config.yaml << 'EOF'
source:
  type: cvs
  path: /path/to/cvs
  module: mymodule

target:
  type: git
  path: /path/to/output/myapp-git

mapping:
  authors_file: author-mapping.yaml
  
  branches:
    "MAIN": "main"
    "DEV": "develop"
    "RELEASE_1_0": "release/1.0"
  
  tags:
    "V1_0": "v1.0.0"
    "V2_0": "v2.0.0"

options:
  dryRun: true  # Preview mode
  verbose: true
  chunkSize: 100
EOF

# Afternoon: Validate and test
git-migrator validate --config migration-config.yaml
git-migrator migrate --config migration-config.yaml --dry-run
```

### Phase 2: Migration (Day 3)

#### Day 3: Execute Migration

```bash
# Morning: Final checks
git-migrator validate --config migration-config.yaml

# Announce freeze
echo "Repository frozen for migration" | mail -s "CVS Migration" team@example.com

# Start migration
nohup git-migrator migrate \
  --config migration-config.yaml \
  --verbose \
  > migration.log 2>&1 &

# Monitor progress
tail -f migration.log

# Or via web UI
git-migrator web --port 8080
# Open http://localhost:8080

# Afternoon: Verify migration
# (See Verification section below)
```

### Phase 3: Validation (Day 4-5)

#### Day 4: Testing

```bash
# Developer testing
# 1. Clone migrated repository
git clone /path/to/output/myapp-git test-clone
cd test-clone

# 2. Verify history
git log --oneline --graph --all

# 3. Test branching
git checkout -b test-branch
# Make changes
git commit -m "Test commit"
git checkout main
git branch -D test-branch

# 4. Test common workflows
# - Pulling updates
# - Creating branches
# - Merging
# - Tagging
```

#### Day 5: Finalize

```bash
# Morning: Address any issues found during testing

# Afternoon: Push to remote
cd /path/to/output/myapp-git
git remote add origin git@github.com:myorg/myapp.git
git push -u origin --all
git push origin --tags

# Announce completion
echo "Migration complete! New Git repo: git@github.com:myorg/myapp.git" \
  | mail -s "Migration Complete" team@example.com
```

## Large Repository Migration

### Optimizing Performance

For repositories with 10,000+ commits:

```yaml
# config.yaml
options:
  # Increase chunk size to reduce state saves
  chunkSize: 500
  
  # Disable verbose output for speed
  verbose: false
  
  # Use RCS mode if available (faster than CVS binary)
  cvsMode: rcs
```

### Parallel Processing

For multiple modules:

```bash
# Create GNU Parallel script
cat > migrate-parallel.sh << 'EOF'
#!/bin/bash
modules=(module1 module2 module3 module4)

echo "modules" | parallel -j 4 \
  'git-migrator migrate --config config-{}.yaml'
EOF

chmod +x migrate-parallel.sh
./migrate-parallel.sh
```

### Incremental Migration with Chunking

```bash
# Migrate in chunks of 10,000 commits
git-migrator migrate \
  --config config.yaml \
  --chunk-size 10000 \
  --start-commit 0 \
  --end-commit 10000

# Resume with next chunk
git-migrator migrate \
  --config config.yaml \
  --resume \
  --start-commit 10001 \
  --end-commit 20000
```

### Handling Binary Files

Large binary files can slow migration:

```yaml
# config.yaml
options:
  # Skip large binaries
  excludePatterns:
    - "*.zip"
    - "*.tar.gz"
    - "*.jar"
    - "*.war"
  
  # Or use Git LFS for binaries
  lfsPatterns:
    - "*.psd"
    - "*.bin"
    - "*.exe"
```

## Multiple Repositories

### Batch Migration Script

```bash
#!/bin/bash
# migrate-all.sh - Migrate multiple repositories

REPOS=(
  "myapp:MyApp:module"
  "lib1:Lib1:lib/module1"
  "lib2:Lib2:lib/module2"
  "tools:Tools:tools"
)

for repo in "${REPOS[@]}"; do
  IFS=':' read -r name module desc <<< "$repo"
  
  echo "==================================="
  echo "Migrating: $name ($module)"
  echo "==================================="
  
  # Create config
  cat > "config-${name}.yaml" << EOF
source:
  type: cvs
  path: /cvs/main
  module: $module
target:
  type: git
  path: /git/repos/${name}.git
mapping:
  authors_file: /config/author-mapping.yaml
options:
  dryRun: false
  verbose: true
EOF
  
  # Migrate
  if git-migrator migrate --config "config-${name}.yaml"; then
    echo "✓ $name migrated successfully"
    echo "$name" >> migration-success.log
  else
    echo "✗ $name migration failed"
    echo "$name" >> migration-failed.log
    # Continue with other repos
  fi
done

echo "==================================="
echo "Migration Summary"
echo "==================================="
echo "Successful: $(wc -l < migration-success.log)"
echo "Failed: $(wc -l < migration-failed.log)"
```

### Monorepo Migration

Combine multiple CVS modules into single Git monorepo:

```bash
# 1. Migrate each module to separate repos
./migrate-all.sh

# 2. Create monorepo
mkdir /git/monorepo
cd /git/monorepo
git init

# 3. Import each module as subdirectory
for repo in myapp lib1 lib2 tools; do
  git subtree add --prefix=$repo /git/repos/${repo}.git main
done

# 4. Result:
# /git/monorepo/
#   ├── myapp/
#   ├── lib1/
#   ├── lib2/
#   └── tools/
```

## Post-Migration Tasks

### Repository Setup

```bash
cd /path/to/git/repo

# Add standard Git files
cat > .gitignore << 'EOF'
# Build artifacts
*.o
*.so
*.exe

# IDE files
.idea/
.vscode/
*.swp

# OS files
.DS_Store
Thumbs.db

# Environment
.env
.env.local
EOF

# Add README if missing
cat > README.md << 'EOF'
# Project Name

Brief description of the project.

## Setup

Instructions for getting started.

## Development

Development workflow and guidelines.

## History

This repository was migrated from CVS on [date].
Original CVS module: [module-name]
EOF

# Add LICENSE file
cp /path/to/license/LICENSE .

# Create development branch
git checkout -b develop

# Push everything
git add .gitignore README.md LICENSE
git commit -m "Add standard Git files"
git push origin main develop
```

### Configure Branch Protection

```bash
# Using GitHub CLI
gh api repos/:owner/:repo/branches/main/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"contexts":[]}' \
  --field enforce_admins=true \
  --field required_pull_request_reviews='{"dismiss_stale_reviews":true}' \
  --field restrictions=null
```

### Set Up CI/CD

```bash
# Create GitHub Actions workflow
mkdir -p .github/workflows

cat > .github/workflows/ci.yml << 'EOF'
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Build
      run: go build -v ./...
    
    - name: Test
      run: go test -v ./...
EOF

git add .github/workflows/ci.yml
git commit -m "Add CI workflow"
git push origin main
```

## Verification

### Automated Verification

```bash
#!/bin/bash
# verify-migration.sh

CVS_PATH="/path/to/cvs"
GIT_PATH="/path/to/git/repo"

echo "=== Verification Report ==="

# 1. Commit count
CVS_COMMITS=$(git-migrator analyze --source-type cvs --source $CVS_PATH | grep "Commits:" | awk '{print $2}')
GIT_COMMITS=$(cd $GIT_PATH && git rev-list --all --count)

echo "Commit count:"
echo "  CVS:  $CVS_COMMITS"
echo "  Git:  $GIT_COMMITS"

if [ "$CVS_COMMITS" -eq "$GIT_COMMITS" ]; then
  echo "  ✓ Match"
else
  echo "  ✗ Mismatch!"
fi

# 2. Branches
echo ""
echo "Branches:"
cd $GIT_PATH
git branch -a | sed 's/^[* ]*//' | while read branch; do
  echo "  - $branch"
done

# 3. Tags
echo ""
echo "Tags:"
git tag -l | while read tag; do
  echo "  - $tag"
done

# 4. Authors
echo ""
echo "Authors:"
git log --format='%an <%ae>' | sort -u | while read author; do
  echo "  - $author"
done

# 5. Date range
echo ""
echo "Date range:"
FIRST=$(git log --reverse --format='%ai' | head -1)
LAST=$(git log --format='%ai' | head -1)
echo "  First commit: $FIRST"
echo "  Last commit:  $LAST"

# 6. File integrity
echo ""
echo "File count:"
FILES=$(find . -type f -not -path './.git/*' | wc -l)
echo "  $FILES files"
```

### Manual Verification Checklist

- [ ] **History Integrity**
  - [ ] All commits present
  - [ ] Commit order preserved
  - [ ] Commit messages intact
  - [ ] Commit timestamps correct
  
- [ ] **Author Information**
  - [ ] All authors mapped correctly
  - [ ] Email addresses valid
  - [ ] No unknown/generic authors
  
- [ ] **Branches**
  - [ ] All expected branches exist
  - [ ] Branch names correct
  - [ ] Branch points accurate
  
- [ ] **Tags**
  - [ ] All expected tags exist
  - [ ] Tag names correct
  - [ ] Tagged commits correct
  
- [ ] **File Content**
  - [ ] All files present
  - [ ] File content matches
  - [ ] Binary files intact
  - [ ] File permissions appropriate
  
- [ ] **Functionality**
  - [ ] Clone works
  - [ ] Branch operations work
  - [ ] Merge operations work
  - [ ] Push/pull works

### Comparing Specific Commits

```bash
# Find corresponding commits
# CVS commit: 1.45 on 2023-05-15 10:30:00
# Git commit: Find by date/message

cd /path/to/git/repo
GIT_COMMIT=$(git log --all --format='%H %ai %s' \
  | grep "2023-05-15" \
  | grep "commit message pattern" \
  | awk '{print $1}')

# Compare file content
cvs update -r 1.45 -p file.txt > /tmp/cvs-file.txt
git show $GIT_COMMIT:file.txt > /tmp/git-file.txt
diff /tmp/cvs-file.txt /tmp/git-file.txt
```

## Rollback Procedures

### Scenario: Migration Failed

```bash
# 1. Check migration state
git-migrator status --config config.yaml

# 2. If recoverable, resume
git-migrator migrate --config config.yaml --resume

# 3. If unrecoverable, clean up
rm -rf /path/to/git/repo
rm -f /path/to/git/repo/.migration-state.db

# 4. Fix issue and retry
# (address whatever caused the failure)
git-migrator migrate --config config.yaml
```

### Scenario: Issues Found After Completion

```bash
# Option 1: Create fixes in Git
cd /path/to/git/repo
git checkout -b fix-author-mapping

# Apply corrections
git filter-branch --env-filter '
if [ "$GIT_AUTHOR_EMAIL" = "wrong@email.com" ]; then
    export GIT_AUTHOR_EMAIL="correct@email.com"
fi
'

# Force push (with caution!)
git push origin --force --all

# Option 2: Re-migrate from scratch
# Only if issues are severe
rm -rf /path/to/git/repo
git-migrator migrate --config corrected-config.yaml
```

## Troubleshooting

### Problem: "Out of memory"

**Symptoms:**
- Migration crashes
- System becomes unresponsive
- "Cannot allocate memory" errors

**Solutions:**
```bash
# Increase chunk size (fewer state saves)
# config.yaml
options:
  chunkSize: 1000

# Exclude large files
options:
  excludePatterns:
    - "*.zip"
    - "*.tar.gz"
    - "*.bin"

# Use 64-bit system with more RAM
# Or migrate on machine with more resources
```

### Problem: "Commit ordering wrong"

**Symptoms:**
- Commits appear out of order
- Branches created at wrong times
- History looks incorrect

**Solutions:**
```bash
# Check CVS timestamps
cvs log file.txt | head -20

# Verify timezone handling
# config.yaml
options:
  timezone: "UTC"  # or your timezone

# Manual ordering verification
git log --graph --oneline --date-order --all
```

### Problem: "Binary files corrupted"

**Symptoms:**
- Images won't open
- Executables fail
- Archives won't extract

**Solutions:**
```bash
# Ensure CVS mode handles binaries correctly
# config.yaml
source:
  type: cvs
  cvsMode: auto  # or 'binary' for binary-heavy repos

# Verify with checksums
cvs update -p -r 1.5 image.png | md5sum
git show commit-hash:image.png | md5sum
```

### Problem: "Missing commits"

**Symptoms:**
- Git has fewer commits than CVS
- Some history appears truncated
- Git log shows gaps

**Solutions:**
```bash
# Check for empty commits
# config.yaml
options:
  preserveEmptyCommits: true

# Verify branch coverage
cvs rlog -h -l /path/to/cvs/module | grep "head:"
git log --all --oneline | wc -l

# Check for vendor branches
cvs log -h file.txt | grep "vendor"
```

### Problem: "Author mapping incomplete"

**Symptoms:**
- Unknown authors in Git log
- Generic email addresses
- Missing author information

**Solutions:**
```bash
# Extract all authors again
git-migrator authors extract --source /path/to/cvs --force > all-authors.txt

# Find unmapped authors
cd /path/to/git/repo
git log --format='%an <%ae>' | sort -u | while read author; do
  if ! grep -q "$author" author-mapping.yaml; then
    echo "Unmapped: $author"
  fi
done

# Add missing mappings
# Then re-migrate or use git filter-branch
```

## Migration Patterns

### Pattern: Clean History

Create cleaner history by squashing trivial commits:

```bash
# After migration
cd /path/to/git/repo

# Interactive rebase to squash commits
git rebase -i --root

# Or use filtering
git filter-branch --tree-filter '
  # Combine consecutive whitespace-only changes
  if git diff --cached --quiet; then
    git reset --soft HEAD~1
  fi
' HEAD
```

### Pattern: Modular History

Split monolithic CVS module into logical Git repositories:

```bash
# 1. Migrate entire CVS module
git-migrator migrate --config full-config.yaml

# 2. Split into separate repos
cd /git/full-repo

# Create repo for subdirectory
git subtree split --prefix=module1 --branch module1-branch
mkdir ../module1-repo
cd ../module1-repo
git init
git pull ../full-repo module1-branch

# Repeat for each module
```

### Pattern: Historical Import

Import CVS history into existing Git repository:

```bash
# 1. Migrate CVS to temporary repo
git-migrator migrate --config config.yaml
cd /tmp/cvs-migrated

# 2. Create graft point
git checkout --orphan cvs-history
git commit -m "Historical CVS import starting point"

# 3. Rebase onto existing repo
cd /path/to/existing/git/repo
git remote add cvs-history /tmp/cvs-migrated
git fetch cvs-history
git rebase --onto main --root cvs-history/main
```

## Best Practices

### DO ✓

- **Plan thoroughly**: Document entire migration process
- **Test extensively**: Use dry-run multiple times
- **Verify everything**: Check commits, branches, tags, authors
- **Communicate clearly**: Keep team informed throughout
- **Back up source**: Keep CVS repository accessible
- **Document decisions**: Record why certain choices were made
- **Time appropriately**: Schedule during low-activity periods
- **Have rollback plan**: Know how to undo if needed
- **Train team**: Ensure developers know Git workflows
- **Monitor post-migration**: Watch for issues after cutover

### DON'T ✗

- **Rush the process**: Take time to verify
- **Skip dry run**: Always test first
- **Ignore warnings**: Address all warnings and errors
- **Migrate blindly**: Understand your repository first
- **Forget branches/tags**: Include all historical references
- **Neglect authors**: Map all contributors properly
- **Delete CVS early**: Keep source until fully validated
- **Migrate alone**: Involve team in planning and testing
- **Skip documentation**: Record the migration process
- **Set unrealistic deadlines**: Allow buffer time

## Timeline Template

### Week 1: Planning
- Day 1-2: Repository analysis
- Day 3-4: Configuration creation
- Day 5: Dry run testing

### Week 2: Migration
- Day 1: Announce and freeze
- Day 2: Execute migration
- Day 3-4: Verification and testing
- Day 5: Go-live and support

### Week 3: Stabilization
- Day 1-3: Monitor and fix issues
- Day 4-5: Decommission CVS (read-only archive)

## Support Resources

- **Documentation**: [docs/](./)
- **Issues**: [GitHub Issues](https://github.com/adamf123git/git-migrator/issues)
- **Discussions**: [GitHub Discussions](https://github.com/adamf123git/git-migrator/discussions)
- **Email**: support@example.com (if applicable)

## Conclusion

Migration to Git is a significant change that requires careful planning and execution. By following this guide and adapting it to your specific needs, you can ensure a successful migration that preserves your project's valuable history while enabling modern development workflows.

Remember: **measure twice, cut once**. Thorough preparation prevents problems during migration. Take the time to plan, test, and verify, and your migration will be successful.

Good luck with your migration! 🚀