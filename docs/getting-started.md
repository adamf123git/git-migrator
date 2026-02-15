# Getting Started with Git-Migrator

Welcome to Git-Migrator! This guide will walk you through migrating your first CVS repository to Git. By the end of this tutorial, you'll have successfully migrated a repository while preserving all commits, branches, tags, and author information.

## 📋 Prerequisites

Before you begin, ensure you have:

- **CVS repository access**: Path to your CVS repository
- **Git installed**: Git 2.0 or later
- **Storage space**: Enough disk space for both repositories
- **Optional**: Docker (for containerized migration)

## 🚀 Installation

Choose your preferred installation method:

### Option 1: Docker (Recommended)

```bash
# Pull the official image
docker pull adamf123docker/git-migrator:latest

# Verify installation
docker run --rm adamf123docker/git-migrator version
```

### Option 2: Binary Download

Download the latest release for your platform:

```bash
# Linux/macOS
curl -sL https://github.com/adamf123git/git-migrator/releases/latest/download/git-migrator-$(uname -s)-$(uname -m) -o git-migrator
chmod +x git-migrator
./git-migrator version

# macOS with Homebrew
brew tap adamf123git/git-migrator
brew install git-migrator
git-migrator version
```

### Option 3: Build from Source

```bash
git clone https://github.com/adamf123git/git-migrator.git
cd git-migrator
go build -o git-migrator ./cmd/git-migrator
./git-migrator version
```

## 📚 Your First Migration

Let's migrate a sample CVS repository to Git. We'll use a fictional project called "myapp".

### Step 1: Analyze Your CVS Repository

First, understand what you're migrating:

```bash
# If using binary
git-migrator analyze --source-type cvs --source /path/to/cvs/myapp

# If using Docker
docker run --rm \
  -v /path/to/cvs/myapp:/source \
  adamf123docker/git-migrator analyze \
  --source-type cvs \
  --source /source
```

This will show you:
- Number of commits to migrate
- Branches and tags found
- Unique authors (for mapping)
- Repository size and structure

### Step 2: Extract Author Information

CVS uses short usernames (like `jsmith`). Git requires full names and email addresses. Extract the list of authors:

```bash
# Extract unique authors
git-migrator authors extract --source /path/to/cvs/myapp > authors.txt

# View the results
cat authors.txt
```

Example output:
```
jsmith
mjones
twilliams
```

### Step 3: Create Author Mapping

Create a mapping file to convert CVS usernames to Git authors:

```bash
# Create authors.yaml
cat > authors.yaml << 'EOF'
jsmith: "John Smith <john.smith@example.com>"
mjones: "Mary Jones <mary.jones@example.com>"
twilliams: "Tom Williams <tom.williams@example.com>"
EOF
```

### Step 4: Create Configuration File

Create a `migration-config.yaml` file:

```yaml
source:
  type: cvs
  path: /path/to/cvs/myapp
  module: myapp  # Your CVS module name

target:
  type: git
  path: /path/to/output/myapp-git
  # Optional: Push to remote after migration
  # remote: git@github.com:myorg/myapp.git

mapping:
  authors:
    jsmith: "John Smith <john.smith@example.com>"
    mjones: "Mary Jones <mary.jones@example.com>"
    twilliams: "Tom Williams <tom.williams@example.com>"

  # Map branch names (optional)
  branches:
    "MAIN": "main"
    "DEV_BRANCH": "develop"
  
  # Map tag names (optional)
  tags:
    "RELEASE_1_0": "v1.0.0"
    "RELEASE_2_0": "v2.0.0"

options:
  dryRun: false      # Set to true to preview without changes
  verbose: true      # Show detailed progress
  chunkSize: 100     # Save state every 100 commits
  resume: false      # Set to true to resume interrupted migration
```

### Step 5: Preview Migration (Dry Run)

Before making actual changes, preview the migration:

```bash
# Using binary
git-migrator migrate --config migration-config.yaml --dry-run --verbose

# Using Docker
docker run --rm \
  -v /path/to/cvs/myapp:/source \
  -v /path/to/output:/target \
  -v $(pwd)/migration-config.yaml:/config.yaml \
  adamf123docker/git-migrator migrate \
  --config /config.yaml \
  --dry-run \
  --verbose
```

Review the output carefully. This shows what will happen without making changes.

### Step 6: Run Migration

Once you're satisfied with the preview, run the actual migration:

```bash
# Using binary
git-migrator migrate --config migration-config.yaml --verbose

# Using Docker
docker run --rm \
  -v /path/to/cvs/myapp:/source \
  -v /path/to/output:/target \
  -v $(pwd)/migration-config.yaml:/config.yaml \
  adamf123docker/git-migrator migrate \
  --config /config.yaml \
  --verbose
```

You'll see progress output like:
```
Starting migration...
[1/247] Processing commit abc123...
[2/247] Processing commit def456...
...
Migration complete!
```

### Step 7: Verify Migration

Check that everything migrated correctly:

```bash
cd /path/to/output/myapp-git

# Check commit count
git rev-list --all --count
# Should match your CVS commit count

# Verify branches
git branch -a
# Should include all your mapped branches

# Verify tags
git tag -l
# Should include all your mapped tags

# Check commit history
git log --oneline --graph --all
# Review commit messages and authors
```

### Step 8: Push to Remote (Optional)

If you specified a remote repository:

```bash
cd /path/to/output/myapp-git

# Add remote (if not in config)
git remote add origin git@github.com:myorg/myapp.git

# Push all branches
git push -u origin --all

# Push all tags
git push origin --tags
```

## 🌐 Using the Web UI

Git-Migrator includes a web interface for easier migration management.

### Starting the Web UI

```bash
# Using binary
git-migrator web --port 8080

# Using Docker
docker run -d \
  -p 8080:8080 \
  -v /path/to/repos:/repos \
  --name git-migrator-web \
  adamf123docker/git-migrator web
```

Open http://localhost:8080 in your browser.

### Web UI Features

1. **Dashboard**: Overview of all migrations
2. **New Migration**: Wizard to configure and start migrations
3. **Progress Monitoring**: Real-time progress with WebSocket updates
4. **Configuration Editor**: Edit YAML configurations in the browser
5. **Log Viewer**: View migration logs in real-time
6. **Migration History**: Track all completed migrations

### Running Migration via Web UI

1. Click "New Migration"
2. Fill in source details (CVS path, module)
3. Configure target (Git path, optional remote)
4. Add author mappings
5. Configure branch/tag mappings
6. Click "Preview" to dry run
7. Click "Start Migration" to begin
8. Monitor progress on the dashboard

## 🔧 Common Scenarios

### Resuming Interrupted Migration

If a migration is interrupted (power loss, Ctrl+C, crash):

```bash
# Simply run with --resume flag
git-migrator migrate --config migration-config.yaml --resume
```

The tool will continue from the last checkpoint.

### Large Repository Migration

For repositories with thousands of commits:

```yaml
# In migration-config.yaml
options:
  chunkSize: 50      # Save state more frequently
  verbose: true      # Monitor progress
```

```bash
# Run in background with nohup
nohup git-migrator migrate --config migration-config.yaml > migration.log 2>&1 &

# Monitor progress
tail -f migration.log
```

### Multiple Repositories

Create a script to migrate multiple modules:

```bash
#!/bin/bash
MODULES="app1 app2 lib1 lib2"

for module in $MODULES; do
  echo "Migrating $module..."
  
  git-migrator migrate --config "configs/${module}-config.yaml"
  
  if [ $? -eq 0 ]; then
    echo "✓ $module migrated successfully"
  else
    echo "✗ $module migration failed"
    exit 1
  fi
done

echo "All migrations complete!"
```

## 🐛 Troubleshooting

### Issue: "CVS binary not found"

**Solution**: Install CVS or use RCS mode:

```yaml
source:
  type: cvs
  path: /path/to/cvs
  cvsMode: rcs  # Use RCS files directly
```

### Issue: "Permission denied"

**Solution**: Check file permissions:

```bash
# Ensure read access to CVS repository
ls -la /path/to/cvs

# Ensure write access to target
mkdir -p /path/to/output
touch /path/to/output/test && rm /path/to/output/test
```

### Issue: "Author mapping incomplete"

**Solution**: Extract all authors and create complete mapping:

```bash
# Re-extract authors
git-migrator authors extract --source /path/to/cvs > all-authors.txt

# Add missing mappings to your config
```

### Issue: "Migration appears stuck"

**Solution**: Check progress and enable verbose mode:

```bash
# Kill current migration (Ctrl+C)
# Resume with verbose output
git-migrator migrate --config config.yaml --resume --verbose
```

## 📊 Best Practices

1. **Always dry run first**: Preview changes before committing
2. **Save configuration**: Keep your YAML files in version control
3. **Verify thoroughly**: Check commits, branches, tags after migration
4. **Test with subset**: Try a small module first before large migrations
5. **Monitor progress**: Use web UI for long-running migrations
6. **Backup original**: Keep CVS repository accessible until verification complete
7. **Document mappings**: Save author/branch/tag mappings for future reference
8. **Use chunking**: Set appropriate `chunkSize` for large repositories

## 📖 Next Steps

Now that you've completed your first migration:

- **[Configuration Guide](./configuration.md)**: Learn all configuration options
- **[Migration Guide](./migration.md)**: Advanced migration scenarios
- **[Architecture](./software-architecture.md)**: Understand how Git-Migrator works
- **[Contributing](../CONTRIBUTING.md)**: Help improve Git-Migrator

## 🆘 Getting Help

- **Documentation**: Browse the [docs](./) directory
- **Issues**: Report bugs on [GitHub Issues](https://github.com/adamf123git/git-migrator/issues)
- **Discussions**: Ask questions in [GitHub Discussions](https://github.com/adamf123git/git-migrator/discussions)

## 🎉 Congratulations!

You've successfully migrated your first CVS repository to Git using Git-Migrator. Your repository now has:
- ✅ Complete commit history preserved
- ✅ All branches migrated
- ✅ All tags created
- ✅ Author information properly mapped
- ✅ Clean Git history ready for collaborative development

Welcome to the Git world! 🚀