# Configuration Reference

This document provides complete reference for Git-Migrator configuration options. All settings can be specified in a YAML configuration file or via command-line flags.

## 📋 Table of Contents

1. [Configuration File](#configuration-file)
2. [Source Configuration](#source-configuration)
3. [Target Configuration](#target-configuration)
4. [Mapping Configuration](#mapping-configuration)
5. [Options Configuration](#options-configuration)
6. [Complete Examples](#complete-examples)
7. [Environment Variables](#environment-variables)
8. [Validation](#validation)

## Configuration File

### File Format

Git-Migrator uses YAML for configuration files. Create a file named `migration-config.yaml` (or any name you prefer).

```yaml
# migration-config.yaml
source:
  type: cvs
  path: /path/to/cvs/repository
  
target:
  type: git
  path: /path/to/output/git/repository
  
mapping:
  authors:
    cvsuser1: "Full Name <email@example.com>"
    
options:
  dryRun: false
  verbose: true
```

### File Locations

Git-Migrator searches for configuration files in this order:

1. Path specified with `--config` flag
2. `./migration-config.yaml` (current directory)
3. `~/.git-migrator/config.yaml` (home directory)
4. `/etc/git-migrator/config.yaml` (system-wide)

### Configuration Precedence

Settings are applied in this order (later overrides earlier):

1. Default values
2. Configuration file
3. Environment variables
4. Command-line flags

## Source Configuration

Configure the source version control system.

### CVS Source

```yaml
source:
  # Required
  type: cvs                          # Source type
  
  # Repository location
  path: /path/to/cvs/repository      # Local path to CVS repository
  module: mymodule                   # CVS module name (optional if at root)
  
  # CVS-specific options
  cvsRoot: :local:/path/to/cvs       # CVSROOT override (optional)
  cvsMode: auto                      # CVS access mode: auto, rcs, binary
  
  # Authentication (for remote CVS)
  cvsServer: cvs.example.com         # CVS server hostname
  cvsUser: username                  # CVS username
  cvsPassword: password              # CVS password (use env var instead)
  
  # Filtering
  excludeModules:                    # Modules to skip
    - CVSROOT
    - test-data
  
  # Advanced
  encoding: UTF-8                    # Character encoding for filenames
  timezone: UTC                      # Timezone for commit dates
```

#### CVS Options Explained

**`path`** (required)
- Local filesystem path to CVS repository
- Can be absolute or relative
- Must contain CVSROOT directory

**`module`** (conditional)
- CVS module name to migrate
- Required if repository contains multiple modules
- Omit if migrating entire repository

**`cvsMode`**
- `auto` (default): Automatically detect best mode
- `rcs`: Parse RCS files directly (faster, no CVS binary needed)
- `binary`: Use CVS binary commands (slower, more compatible)

**`cvsRoot`**
- Override CVSROOT environment variable
- Format: `:method:user@host:path`
- Example: `:ext:user@cvs.server.com:/cvsroot`

**`encoding`**
- Character encoding for filenames and commit messages
- Common values: `UTF-8`, `ISO-8859-1`, `Windows-1252`
- Default: `UTF-8`

**`timezone`**
- Timezone for interpreting commit timestamps
- CVS doesn't store timezone information
- Values: `UTC`, `America/New_York`, `Europe/London`, etc.
- Default: `UTC`

### SVN Source (Future)

```yaml
source:
  type: svn                          # Source type (coming in v2.0)
  
  # Repository location
  url: http://svn.example.com/repo   # SVN repository URL
  path: /local/svn/checkout          # Or local checkout path
  
  # Authentication
  username: svnuser
  password: svnpass                  # Use environment variable
  
  # SVN-specific options
  trunk: trunk                       # Trunk directory name
  branches: branches                 # Branches directory name
  tags: tags                         # Tags directory name
  
  # Advanced
  includeExternals: true             # Include svn:externals
  followCopies: true                 # Track copy history
```

## Target Configuration

Configure the target Git repository.

### Git Target

```yaml
target:
  # Required
  type: git                          # Target type
  path: /path/to/output/git/repo     # Local path for Git repository
  
  # Remote configuration
  remote: git@github.com:org/repo.git    # Remote URL (optional)
  remoteName: origin                 # Remote name (default: origin)
  
  # Initial configuration
  config:
    user.name: "Migration Bot"       # Default Git author
    user.email: "migration@example.com"
    core.autocrlf: input             # Line ending handling
  
  # Repository settings
  initialBranch: main                # Initial branch name (default: main)
  bare: false                        # Create bare repository
  
  # Post-migration
  pushOnComplete: false              # Auto-push after migration
  pushTags: true                     # Include tags when pushing
  forcePush: false                   # Use force push (dangerous!)
  
  # Git LFS
  lfs:
    enabled: false                   # Enable Git LFS
    patterns:                        # File patterns for LFS
      - "*.psd"
      - "*.bin"
      - "*.zip"
```

#### Target Options Explained

**`path`** (required)
- Local filesystem path for Git repository
- Must not exist (will be created)
- Or must be empty directory

**`remote`**
- Git remote URL for pushing
- Supports SSH: `git@github.com:org/repo.git`
- HTTPS: `https://github.com/org/repo.git`
- Local: `/path/to/another/repo.git`

**`config`**
- Git configuration options to set in created repository
- Key-value pairs
- Applied before any commits

**`initialBranch`**
- Name of the initial branch
- Default: `main`
- Common alternatives: `master`, `trunk`, `develop`

**`bare`**
- Create bare repository (no working directory)
- Useful for server-side repositories
- Default: `false`

**`pushOnComplete`**
- Automatically push to remote after migration
- Requires `remote` to be set
- Requires proper authentication

**`lfs.enabled`**
- Enable Git Large File Storage
- Large files stored separately
- Requires Git LFS to be installed

**`lfs.patterns`**
- File patterns to store in LFS
- Git wildmatch patterns
- Example: `*.zip`, `path/to/large/**`

## Mapping Configuration

Configure mappings between source and target.

### Author Mapping

Map CVS/SVN usernames to Git authors.

```yaml
mapping:
  # Inline author mapping
  authors:
    jsmith: "John Smith <john.smith@example.com>"
    mjones: "Mary Jones <mary.jones@example.com>"
    twilliams: "Tom Williams <tom.williams@example.com>"
  
  # Or load from external file
  authors_file: /path/to/authors.yaml
  
  # Default for unmapped authors
  defaultAuthor:
    name: "Unknown Author"
    email: "unknown@example.com"
  
  # Author mapping rules
  authorRules:
    # Map generic accounts
    - pattern: "^(buildbot|jenkins|ci)$"
      replace: "CI System <ci@example.com>"
    
    # Map by email domain
    - pattern: "^(.+)$"
      emailTemplate: "$1@company.com"
      nameTemplate: "$1"
```

#### Author Mapping Format

**Inline Mapping**
```yaml
mapping:
  authors:
    cvsusername: "Full Name <email@example.com>"
```

**External File**
```yaml
# authors.yaml
cvsuser1: "Full Name <email@example.com>"
cvsuser2: "Another User <user@example.com>"
```

```yaml
# config.yaml
mapping:
  authors_file: authors.yaml
```

**Format Requirements**
- Key: CVS/SVN username (case-sensitive)
- Value: `"Full Name <email@example.com>"`
- Must include both name and email
- Email must be in angle brackets

**Special Cases**

```yaml
mapping:
  authors:
    # Generic build account
    buildbot: "Build System <build@example.com>"
    
    # Shared team account
    teamlead: "Team Lead <lead@example.com>"
    
    # Historical account (person left)
    olddev: "Former Developer <former@example.com>"
    
    # Service account
    cvsroot: "CVS Administrator <admin@example.com>"
```

### Branch Mapping

Map source branch names to Git branch names.

```yaml
mapping:
  branches:
    # Direct mapping
    "MAIN": "main"
    "DEV": "develop"
    "RELEASE_1_0": "release/1.0"
    "FEATURE_X": "feature/x"
    
    # Pattern-based mapping
    branchPatterns:
      - pattern: "^RELEASE_(.+)$"
        replace: "release/$1"
      
      - pattern: "^FEATURE_(.+)$"
        replace: "feature/$1"
      
      - pattern: "^BUGFIX_(.+)$"
        replace: "bugfix/$1"
    
    # Branches to skip
    skipBranches:
      - "test-*"
      - "temp-*"
      - "BACKUP_*"
```

#### Branch Mapping Rules

**Direct Mapping**
```yaml
mapping:
  branches:
    "OLD_NAME": "new-name"
```

**Pattern Mapping**
```yaml
mapping:
  branches:
    # Convert RELEASE_1_0 -> release/1.0
    # Convert RELEASE_2_1 -> release/2.1
    "^RELEASE_(.+)_(.+)$": "release/$1.$2"
```

**Skip Branches**
```yaml
mapping:
  skipBranches:
    - "test-*"        # Skip all test branches
    - "TEMP_*"        # Skip temporary branches
    - "BACKUP_*"      # Skip backup branches
```

### Tag Mapping

Map source tag names to Git tag names.

```yaml
mapping:
  tags:
    # Direct mapping
    "V1_0": "v1.0.0"
    "V2_0": "v2.0.0"
    "RELEASE_1_0": "release-1.0.0"
    
    # Pattern-based mapping
    tagPatterns:
      - pattern: "^V(.+)_(.+)_(.+)$"
        replace: "v$1.$2.$3"
      
      - pattern: "^RELEASE_(.+)$"
        replace: "release-$1"
    
    # Tags to skip
    skipTags:
      - "test-*"
      - "temp-*"
      - "BUILD_*"
    
    # Tag type
    tagType: lightweight            # lightweight or annotated
    tagMessage: "Release {tag}"     # Message template for annotated tags
```

#### Tag Mapping Examples

```yaml
mapping:
  tags:
    # Version tags
    "V1_0_0": "v1.0.0"
    "V2_1_3": "v2.1.3"
    
    # Release tags
    "RELEASE_1_0": "release-1.0"
    
    # Special tags
    "START": "project-start"
    "MIGRATION": "cvs-migration-point"
```

### File Path Mapping

Transform file paths during migration.

```yaml
mapping:
  paths:
    # Rename directories
    rename:
      "oldname/": "newname/"
      "src/deprecated/": "src/legacy/"
    
    # Move files
    move:
      "README": "README.md"
      "LICENSE.TXT": "LICENSE"
    
    # Pattern-based transformation
    patterns:
      - pattern: "^old_project/(.+)$"
        replace: "new_project/$1"
```

## Options Configuration

Control migration behavior.

```yaml
options:
  # Execution mode
  dryRun: false                      # Preview without changes
  verbose: false                     # Detailed output
  quiet: false                       # Minimal output
  
  # Resume capability
  resume: false                      # Resume interrupted migration
  chunkSize: 100                     # Save state every N commits
  stateFile: .migration-state.db     # State file path
  
  # History handling
  preserveEmptyCommits: false        # Keep commits with no changes
  includeBinaryFiles: true           # Include binary files
  
  # Performance
  parallelJobs: 1                    # Parallel processing (experimental)
  bufferSize: 65536                  # I/O buffer size
  
  # Filtering
  startDate: null                    # Only commits after date
  endDate: null                      # Only commits before date
  excludePatterns: []                # Files to exclude
  includePatterns: []                # Files to include (whitelist)
  
  # Verification
  verifyAfterMigration: true         # Verify migrated repository
  strictMode: false                  # Fail on any warning
  
  # Advanced
  interruptAt: 0                     # Testing: interrupt after N commits
  skipBranches: []                   # Branches to skip
  skipTags: []                       # Tags to skip
```

### Options Explained

**`dryRun`**
- Preview migration without making changes
- Shows what would happen
- Always test with dry-run first
- Default: `false`

**`verbose`**
- Show detailed progress information
- Lists each commit as it's processed
- Useful for monitoring progress
- Default: `false`

**`quiet`**
- Suppress non-essential output
- Only show errors
- Overrides `verbose`
- Default: `false`

**`resume`**
- Continue from last checkpoint
- Uses state file to track progress
- Safe to run multiple times
- Default: `false`

**`chunkSize`**
- Save state every N commits
- Lower = more frequent saves (safer)
- Higher = faster migration
- Default: `100`
- Recommended: 50-500

**`preserveEmptyCommits`**
- Keep commits with no file changes
- CVS may have commits that only changed metadata
- Usually safe to skip
- Default: `false`

**`parallelJobs`**
- Number of parallel workers
- Experimental feature
- Can speed up large migrations
- Default: `1` (sequential)

**`startDate` / `endDate`**
- Filter commits by date
- Format: `YYYY-MM-DD` or `YYYY-MM-DD HH:MM:SS`
- Useful for partial migrations
- Example:
```yaml
options:
  startDate: "2020-01-01"
  endDate: "2023-12-31"
```

**`excludePatterns`**
- Glob patterns for files to exclude
- Applied before include patterns
- Example:
```yaml
options:
  excludePatterns:
    - "*.zip"
    - "*.tar.gz"
    - "test-data/**"
    - "**/*.bak"
```

**`includePatterns`**
- Glob patterns for files to include
- If specified, only matching files included
- Example:
```yaml
options:
  includePatterns:
    - "src/**"
    - "include/**"
    - "*.md"
```

**`verifyAfterMigration`**
- Run verification after migration
- Checks commit counts, branches, tags
- Recommended for production migrations
- Default: `true`

**`strictMode`**
- Fail on any warning
- Useful for ensuring completeness
- Default: `false`

## Complete Examples

### Basic CVS to Git

```yaml
source:
  type: cvs
  path: /cvs/myproject
  module: myproject

target:
  type: git
  path: /git/myproject.git
  remote: git@github.com:myorg/myproject.git

mapping:
  authors:
    jsmith: "John Smith <john@company.com>"
    mjones: "Mary Jones <mary@company.com>"

options:
  dryRun: false
  verbose: true
  chunkSize: 100
```

### Large Repository with Performance Tuning

```yaml
source:
  type: cvs
  path: /cvs/large-project
  module: project
  cvsMode: rcs                        # Faster than binary

target:
  type: git
  path: /git/large-project.git
  initialBranch: main
  lfs:
    enabled: true
    patterns:
      - "*.zip"
      - "*.tar.gz"
      - "*.jar"

mapping:
  authors_file: authors-large.yaml   # External file for many authors
  
  branchPatterns:
    - pattern: "^RELEASE_(.+)$"
      replace: "release/$1"
  
  tags:
    "V1_0": "v1.0.0"
    "V2_0": "v2.0.0"

options:
  dryRun: false
  verbose: false                      # Faster without verbose
  chunkSize: 500                      # Larger chunks for speed
  parallelJobs: 4                     # Parallel processing
  excludePatterns:
    - "*.zip"
    - "*.tar.gz"
    - "test-data/**"
    - "**/*.bak"
    - "**/*.orig"
  bufferSize: 131072                  # 128KB buffer
```

### Incremental Migration

```yaml
source:
  type: cvs
  path: /cvs/project
  module: module1                     # Migrate one module at a time

target:
  type: git
  path: /git/module1.git

mapping:
  authors_file: ../common/authors.yaml

options:
  dryRun: false
  resume: true                        # Enable resume
  chunkSize: 50
  startDate: "2022-01-01"            # Only recent history
  verifyAfterMigration: true
```

### Enterprise Migration with Full Configuration

```yaml
# Enterprise migration configuration
# Project: Large Enterprise Application
# Date: 2024-01-15

source:
  type: cvs
  path: /cvs/enterprise-app
  cvsMode: rcs
  encoding: ISO-8859-1                # Legacy encoding
  timezone: America/New_York          # Company timezone
  
  excludeModules:
    - CVSROOT
    - archived
    - test-data
    - temp-backups

target:
  type: git
  path: /git/enterprise-app.git
  remote: git@github.enterprise.com:company/enterprise-app.git
  
  config:
    user.name: "Migration Bot"
    user.email: "migration@company.com"
    core.autocrlf: input
    core.safecrlf: warn
  
  initialBranch: main
  pushOnComplete: true
  pushTags: true
  
  lfs:
    enabled: true
    patterns:
      - "*.psd"
      - "*.ai"
      - "*.zip"
      - "*.tar.gz"
      - "*.war"
      - "*.jar"
      - "binaries/**"

mapping:
  authors_file: enterprise-authors.yaml
  
  defaultAuthor:
    name: "Unknown Developer"
    email: "unknown@company.com"
  
  branchPatterns:
    - pattern: "^MAIN$"
      replace: "main"
    - pattern: "^RELEASE_(.+)$"
      replace: "release/$1"
    - pattern: "^FEATURE_(.+)$"
      replace: "feature/$1"
    - pattern: "^BUGFIX_(.+)$"
      replace: "bugfix/$1"
  
  skipBranches:
    - "test-*"
    - "TEMP_*"
    - "BACKUP_*"
    - "old-*"
  
  tags:
    "V1_0_0": "v1.0.0"
    "V2_0_0": "v2.0.0"
  
  tagPatterns:
    - pattern: "^RELEASE_(.+)$"
      replace: "release-$1"
  
  skipTags:
    - "test-*"
    - "TEMP_*"
    - "BUILD_*"
  
  paths:
    rename:
      "old_module/": "legacy_module/"
      "deprecated/": "archived/"

options:
  dryRun: false
  verbose: true
  quiet: false
  
  resume: true
  chunkSize: 100
  stateFile: .migration-enterprise.db
  
  preserveEmptyCommits: false
  includeBinaryFiles: true
  
  parallelJobs: 1
  bufferSize: 65536
  
  startDate: "2015-01-01"
  endDate: null
  
  excludePatterns:
    - "*.log"
    - "*.tmp"
    - "*.bak"
    - "build/**"
    - "dist/**"
    - "node_modules/**"
    - ".metadata/**"
  
  verifyAfterMigration: true
  strictMode: true
```

### Configuration for Web UI

When using the Web UI, you can still provide default configuration:

```yaml
# ~/.git-migrator/web-config.yaml
# Default configuration for Web UI migrations

defaults:
  source:
    type: cvs
    cvsMode: auto
    encoding: UTF-8
  
  target:
    type: git
    initialBranch: main
    lfs:
      enabled: false
  
  mapping:
    authors_file: ~/.git-migrator/authors.yaml
  
  options:
    dryRun: false
    verbose: true
    chunkSize: 100
    verifyAfterMigration: true

web:
  port: 8080
  host: 0.0.0.0
  workers: 4
  
  storage:
    type: sqlite
    path: ~/.git-migrator/migrations.db
  
  security:
    corsOrigins:
      - "http://localhost:3000"
      - "https://migrate.company.com"
    authEnabled: false
    authType: basic
```

## Environment Variables

Override configuration with environment variables:

```bash
# Source configuration
export GIT_MIGRATOR_SOURCE_TYPE=cvs
export GIT_MIGRATOR_SOURCE_PATH=/path/to/cvs
export GIT_MIGRATOR_SOURCE_MODULE=mymodule

# Target configuration
export GIT_MIGRATOR_TARGET_PATH=/path/to/git
export GIT_MIGRATOR_TARGET_REMOTE=git@github.com:org/repo.git

# Mapping
export GIT_MIGRATOR_AUTHORS_FILE=/path/to/authors.yaml

# Options
export GIT_MIGRATOR_DRY_RUN=false
export GIT_MIGRATOR_VERBOSE=true
export GIT_MIGRATOR_CHUNK_SIZE=100

# CVS authentication
export CVSROOT=:ext:user@cvs.server.com:/cvsroot
export CVS_RSH=ssh

# Git authentication
export GIT_SSH_COMMAND="ssh -i /path/to/key"

# Proxy settings
export HTTP_PROXY=http://proxy.company.com:8080
export HTTPS_PROXY=http://proxy.company.com:8080
export NO_PROXY=localhost,127.0.0.1
```

### Environment Variable Format

Environment variables follow this pattern:

```
GIT_MIGRATOR_<SECTION>_<KEY>
```

Examples:
- `GIT_MIGRATOR_SOURCE_PATH` → `source.path`
- `GIT_MIGRATOR_TARGET_REMOTE` → `target.remote`
- `GIT_MIGRATOR_OPTIONS_CHUNK_SIZE` → `options.chunkSize`
- `GIT_MIGRATOR_MAPPING_AUTHORS_FILE` → `mapping.authorsFile`

## Validation

### Validate Configuration

Before running migration, validate your configuration:

```bash
# Validate configuration file
git-migrator validate --config migration-config.yaml

# Output example:
# ✓ Configuration valid
# ✓ Source accessible
# ✓ Target path available
# ✓ Author mapping complete (15/15 authors)
# ✓ Branch mapping valid (5 branches)
# ✓ Tag mapping valid (10 tags)
```

### Common Validation Errors

**Error: Source path not found**
```
Error: Source path '/path/to/cvs' does not exist
```
Solution: Verify the source path is correct and accessible.

**Error: Author mapping incomplete**
```
Warning: 3 unmapped authors found: jsmith, mjones, twilliams
```
Solution: Add missing author mappings to configuration.

**Error: Target path exists**
```
Error: Target path '/path/to/git' already exists
```
Solution: Remove existing target directory or use different path.

**Error: Invalid YAML syntax**
```
Error: YAML parse error at line 15: mapping values are not allowed here
```
Solution: Check YAML syntax, ensure proper indentation.

### Configuration Testing

Test configuration with dry-run:

```bash
# Preview migration
git-migrator migrate --config config.yaml --dry-run --verbose

# Check what would happen:
# - Number of commits to migrate
# - Branches that would be created
# - Tags that would be created
# - Authors that would be mapped
# - Files that would be excluded
```

## Configuration Best Practices

### DO ✓

- **Use version control**: Keep configuration files in Git
- **Use external files**: Store large author mappings separately
- **Test with dry-run**: Always preview before actual migration
- **Document decisions**: Comment your configuration
- **Use environment variables**: For secrets (passwords, tokens)
- **Validate first**: Run validation before migration
- **Start simple**: Begin with minimal config, add complexity as needed
- **Back up configs**: Save working configurations for reference
- **Use templates**: Create reusable configuration templates

### DON'T ✗

- **Hardcode secrets**: Never put passwords in config files
- **Skip validation**: Always validate before migration
- **Use absolute paths**: Prefer relative or configurable paths
- **Ignore warnings**: Address all validation warnings
- **Over-complicate**: Start with simple configuration
- **Skip dry-run**: Always test with dry-run first
- **Forget encoding**: Specify character encoding for legacy repos
- **Neglect timezone**: Set timezone for accurate timestamps

## Configuration Templates

### Minimal Template

```yaml
source:
  type: cvs
  path: /path/to/cvs
  module: mymodule

target:
  type: git
  path: /path/to/git

mapping:
  authors:
    user1: "User One <user1@example.com>"
    user2: "User Two <user2@example.com>"

options:
  dryRun: true
  verbose: true
```

### Standard Template

```yaml
source:
  type: cvs
  path: /path/to/cvs
  module: mymodule
  cvsMode: auto

target:
  type: git
  path: /path/to/git
  remote: git@github.com:org/repo.git

mapping:
  authors_file: authors.yaml
  
  branches:
    "MAIN": "main"
    "DEV": "develop"

options:
  dryRun: false
  verbose: true
  chunkSize: 100
  resume: true
  verifyAfterMigration: true
```

### Advanced Template

```yaml
# Advanced configuration template
# Includes all common options

source:
  type: cvs
  path: ${CVS_PATH}
  module: ${CVS_MODULE}
  cvsMode: rcs
  encoding: ${ENCODING:-UTF-8}
  timezone: ${TIMEZONE:-UTC}

target:
  type: git
  path: ${GIT_PATH}
  remote: ${GIT_REMOTE}
  
  config:
    user.name: "Migration Bot"
    user.email: "migration@example.com"
  
  lfs:
    enabled: true
    patterns:
      - "*.zip"
      - "*.tar.gz"

mapping:
  authors_file: ${AUTHORS_FILE}
  
  branchPatterns:
    - pattern: "^RELEASE_(.+)$"
      replace: "release/$1"
    - pattern: "^FEATURE_(.+)$"
      replace: "feature/$1"
  
  skipBranches:
    - "test-*"
    - "TEMP_*"
  
  tags:
    "V1_0": "v1.0.0"
    "V2_0": "v2.0.0"

options:
  dryRun: ${DRY_RUN:-false}
  verbose: ${VERBOSE:-true}
  chunkSize: ${CHUNK_SIZE:-100}
  resume: true
  verifyAfterMigration: true
  strictMode: false
  
  excludePatterns:
    - "*.log"
    - "*.tmp"
    - "build/**"
```

## Troubleshooting Configuration

### Issue: "Configuration file not found"

**Symptoms:**
```
Error: Configuration file 'config.yaml' not found
```

**Solutions:**
```bash
# Specify full path
git-migrator migrate --config /full/path/to/config.yaml

# Or use environment variable
export GIT_MIGRATOR_CONFIG=/path/to/config.yaml
```

### Issue: "Invalid YAML"

**Symptoms:**
```
Error: YAML parse error: mapping values are not allowed here
```

**Solutions:**
```bash
# Validate YAML syntax
yamllint config.yaml

# Common issues:
# 1. Incorrect indentation
# 2. Missing quotes around special characters
# 3. Tabs instead of spaces

# Fix: Ensure proper YAML formatting
```

### Issue: "Author mapping incomplete"

**Symptoms:**
```
Warning: 5 unmapped authors found
```

**Solutions:**
```yaml
# Option 1: Add all mappings
mapping:
  authors:
    user1: "User One <user1@example.com>"
    user2: "User Two <user2@example.com>"
    # ... add all users

# Option 2: Use default author
mapping:
  defaultAuthor:
    name: "Unknown"
    email: "unknown@example.com"

# Option 3: Use pattern rules
mapping:
  authorRules:
    - pattern: "^(.+)$"
      emailTemplate: "$1@company.com"
```

### Issue: "Path too long" (Windows)

**Symptoms:**
```
Error: The filename or extension is too long
```

**Solutions:**
```yaml
# Enable long paths in Git
target:
  config:
    core.longpaths: true

# Or use shorter target path
target:
  path: C:\git\repo
```

## Configuration Reference

### Complete Option Reference

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `source.type` | string | required | cvs, svn |
| `source.path` | string | required | Source repository path |
| `source.module` | string | optional | CVS module name |
| `source.cvsMode` | string | auto | auto, rcs, binary |
| `source.encoding` | string | UTF-8 | Character encoding |
| `source.timezone` | string | UTC | Timezone for dates |
| `target.type` | string | required | git |
| `target.path` | string | required | Target repository path |
| `target.remote` | string | optional | Git remote URL |
| `target.initialBranch` | string | main | Initial branch name |
| `target.bare` | boolean | false | Create bare repository |
| `mapping.authors` | map | optional | Inline author mapping |
| `mapping.authors_file` | string | optional | External author file |
| `mapping.branches` | map | optional | Branch name mapping |
| `mapping.tags` | map | optional | Tag name mapping |
| `options.dryRun` | boolean | false | Preview mode |
| `options.verbose` | boolean | false | Detailed output |
| `options.quiet` | boolean | false | Minimal output |
| `options.resume` | boolean | false | Resume capability |
| `options.chunkSize` | integer | 100 | State save interval |
| `options.preserveEmptyCommits` | boolean | false | Keep empty commits |
| `options.verifyAfterMigration` | boolean | true | Verify repository |
| `options.strictMode` | boolean | false | Fail on warnings |

## Next Steps

- **[Getting Started](./getting-started.md)** - Step-by-step tutorial
- **[Migration Guide](./migration.md)** - Comprehensive migration guide
- **[Architecture](./software-architecture.md)** - System design

## Support

- **Documentation**: [docs/](./)
- **Issues**: [GitHub Issues](https://github.com/adamf123git/git-migrator/issues)
- **Discussions**: [GitHub Discussions](https://github.com/adamf123git/git-migrator/discussions)

---

**Configuration is the foundation of successful migration. Take time to understand and test your configuration before running the actual migration.**