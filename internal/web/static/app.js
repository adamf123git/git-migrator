// Git-Migrator Web UI JavaScript

// API helper
async function api(endpoint, options = {}) {
    const response = await fetch(endpoint, {
        headers: {
            'Content-Type': 'application/json',
            ...options.headers,
        },
        ...options,
    });

    const data = await response.json();
    if (!response.ok || !data.success) {
        throw new Error(data.error?.message || 'Request failed');
    }
    return data.data;
}

// Load migrations on dashboard
async function loadMigrations() {
    const list = document.getElementById('migrations-list');
    if (!list) return;

    try {
        const migrations = await api('/api/migrations');
        if (migrations.length === 0) {
            list.innerHTML = '<p>No migrations yet. <a href="/new">Start one</a></p>';
            return;
        }

        list.innerHTML = migrations.map(m => `
            <div class="migration-item">
                <div>
                    <strong>${m.id.substring(0, 8)}</strong>
                    <span class="migration-status ${m.status}">${m.status}</span>
                </div>
                <a href="/migration/${m.id}" class="button">View</a>
            </div>
        `).join('');
    } catch (err) {
        list.innerHTML = `<p class="error">Error loading migrations: ${err.message}</p>`;
    }
}

// Handle migration form
function setupMigrationForm() {
    const form = document.getElementById('migration-form');
    if (!form) return;

    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const formData = new FormData(form);
        const data = {
            sourceType: formData.get('sourceType'),
            sourcePath: formData.get('sourcePath'),
            targetPath: formData.get('targetPath'),
            options: {
                dryRun: formData.has('dryRun'),
            },
        };

        try {
            const result = await api('/api/migrations', {
                method: 'POST',
                body: JSON.stringify(data),
            });
            window.location.href = `/migration/${result.id}`;
        } catch (err) {
            alert(`Failed to start migration: ${err.message}`);
        }
    });
}

// Handle config form
function setupConfigForm() {
    const form = document.getElementById('config-form');
    if (!form) return;

    // Load current config
    api('/api/config').then(config => {
        document.getElementById('chunkSize').value = config.chunkSize || 100;
        document.getElementById('verbose').checked = config.verbose || false;
    }).catch(err => {
        console.error('Failed to load config:', err);
    });

    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const formData = new FormData(form);
        const data = {
            chunkSize: parseInt(formData.get('chunkSize'), 10),
            verbose: formData.has('verbose'),
        };

        try {
            await api('/api/config', {
                method: 'POST',
                body: JSON.stringify(data),
            });
            alert('Configuration saved!');
        } catch (err) {
            alert(`Failed to save config: ${err.message}`);
        }
    });
}

// Migration progress page
function setupMigrationProgress() {
    const section = document.getElementById('migration-status');
    if (!section) return;

    const migrationId = window.location.pathname.split('/').pop();
    let ws = null;

    function connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        ws = new WebSocket(`${protocol}//${window.location.host}/ws/progress/${migrationId}`);

        ws.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            handleProgressUpdate(msg);
        };

        ws.onerror = (err) => {
            console.error('WebSocket error:', err);
        };

        ws.onclose = () => {
            // Reconnect after 3 seconds
            setTimeout(connectWebSocket, 3000);
        };
    }

    function handleProgressUpdate(msg) {
        const data = msg.data || {};

        // Update progress bar
        const fill = document.getElementById('progress-fill');
        const text = document.getElementById('progress-text');
        if (fill && text) {
            fill.style.width = `${data.percentage || 0}%`;
            text.textContent = `${data.percentage || 0}%`;
        }

        // Update status
        const status = document.getElementById('status');
        if (status) {
            status.textContent = data.status || 'Unknown';
            status.className = `migration-status ${data.status}`;
        }

        // Update current step
        const currentStep = document.getElementById('currentStep');
        if (currentStep) {
            currentStep.textContent = data.currentStep || '-';
        }

        // Update commits
        const commits = document.getElementById('commits');
        if (commits) {
            commits.textContent = `${data.processedCommits || 0} / ${data.totalCommits || 0}`;
        }

        // Update errors
        const errorsSection = document.getElementById('errors');
        const errorList = document.getElementById('error-list');
        if (errorsSection && errorList && data.errors && data.errors.length > 0) {
            errorsSection.classList.remove('hidden');
            errorList.innerHTML = data.errors.map(e => `<li>${e}</li>`).join('');
        }

        // Handle completion
        if (data.status === 'completed' || data.status === 'failed' || data.status === 'stopped') {
            if (ws) ws.close();
        }
    }

    // Stop button
    const stopBtn = document.getElementById('stop-btn');
    if (stopBtn) {
        stopBtn.addEventListener('click', async () => {
            if (!confirm('Are you sure you want to stop this migration?')) return;

            try {
                await api(`/api/migrations/${migrationId}/stop`, { method: 'POST' });
                stopBtn.disabled = true;
            } catch (err) {
                alert(`Failed to stop migration: ${err.message}`);
            }
        });
    }

    connectWebSocket();
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    loadMigrations();
    setupMigrationForm();
    setupConfigForm();
    setupMigrationProgress();
});
