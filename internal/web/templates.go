package web

// HTML templates for the web UI

var indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Git-Migrator</title>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <header>
        <h1>Git-Migrator</h1>
        <nav>
            <a href="/">Dashboard</a>
            <a href="/new">New Migration</a>
            <a href="/config">Configuration</a>
        </nav>
    </header>
    <main>
        <section id="dashboard">
            <h2>Recent Migrations</h2>
            <div id="migrations-list">
                <p>Loading migrations...</p>
            </div>
            <a href="/new" class="button">Start New Migration</a>
        </section>
    </main>
    <script src="/static/app.js"></script>
</body>
</html>
`

var newMigrationHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Migration - Git-Migrator</title>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <header>
        <h1>Git-Migrator</h1>
        <nav>
            <a href="/">Dashboard</a>
            <a href="/new">New Migration</a>
            <a href="/config">Configuration</a>
        </nav>
    </header>
    <main>
        <section id="new-migration">
            <h2>New Migration</h2>
            <form id="migration-form">
                <div class="form-group">
                    <label for="sourceType">Source Type</label>
                    <select id="sourceType" name="sourceType" required>
                        <option value="cvs">CVS</option>
                        <option value="svn">SVN (Coming Soon)</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="sourcePath">Source Path</label>
                    <input type="text" id="sourcePath" name="sourcePath" required>
                </div>
                <div class="form-group">
                    <label for="targetPath">Target Path</label>
                    <input type="text" id="targetPath" name="targetPath" required>
                </div>
                <div class="form-group">
                    <label>
                        <input type="checkbox" id="dryRun" name="dryRun">
                        Dry Run (preview only)
                    </label>
                </div>
                <button type="submit">Start Migration</button>
            </form>
        </section>
    </main>
    <script src="/static/app.js"></script>
</body>
</html>
`

var configHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Configuration - Git-Migrator</title>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <header>
        <h1>Git-Migrator</h1>
        <nav>
            <a href="/">Dashboard</a>
            <a href="/new">New Migration</a>
            <a href="/config">Configuration</a>
        </nav>
    </header>
    <main>
        <section id="configuration">
            <h2>Configuration</h2>
            <form id="config-form">
                <div class="form-group">
                    <label for="chunkSize">Chunk Size</label>
                    <input type="number" id="chunkSize" name="chunkSize" value="100">
                </div>
                <div class="form-group">
                    <label>
                        <input type="checkbox" id="verbose" name="verbose">
                        Verbose Logging
                    </label>
                </div>
                <button type="submit">Save Configuration</button>
            </form>
        </section>
    </main>
    <script src="/static/app.js"></script>
</body>
</html>
`

var migrationHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Migration - Git-Migrator</title>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <header>
        <h1>Git-Migrator</h1>
        <nav>
            <a href="/">Dashboard</a>
            <a href="/new">New Migration</a>
            <a href="/config">Configuration</a>
        </nav>
    </header>
    <main>
        <section id="migration-status">
            <h2>Migration Progress</h2>
            <div class="progress-container">
                <div class="progress-bar">
                    <div id="progress-fill" class="progress-fill" style="width: 0%"></div>
                </div>
                <span id="progress-text">0%</span>
            </div>
            <div id="migration-info">
                <p><strong>Status:</strong> <span id="status">Loading...</span></p>
                <p><strong>Current Step:</strong> <span id="currentStep">-</span></p>
                <p><strong>Commits:</strong> <span id="commits">0 / 0</span></p>
            </div>
            <div id="errors" class="hidden">
                <h3>Errors</h3>
                <ul id="error-list"></ul>
            </div>
            <div class="actions">
                <button id="stop-btn" class="danger">Stop Migration</button>
                <a href="/" class="button">Back to Dashboard</a>
            </div>
        </section>
    </main>
    <script src="/static/app.js"></script>
</body>
</html>
`
