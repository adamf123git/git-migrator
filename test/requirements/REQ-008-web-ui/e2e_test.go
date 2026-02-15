package requirements_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamf123git/git-migrator/internal/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWebUIServesIndex tests that root path serves the web UI
func TestWebUIServesIndex(t *testing.T) {
	server := web.NewServer(web.ServerConfig{
		Port: 8080,
	})
	router := server.Router()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should serve HTML
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")
	assert.NotEmpty(t, rec.Body.String())
}

// TestWebUIServesStaticFiles tests static file serving
func TestWebUIServesStaticFiles(t *testing.T) {
	server := web.NewServer(web.ServerConfig{
		Port: 8080,
	})
	router := server.Router()

	// Test CSS
	req := httptest.NewRequest(http.MethodGet, "/static/style.css", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Test JS
	req = httptest.NewRequest(http.MethodGet, "/static/app.js", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestWebUIContainsRequiredElements tests that UI has required elements
func TestWebUIContainsRequiredElements(t *testing.T) {
	server := web.NewServer(web.ServerConfig{
		Port: 8080,
	})
	router := server.Router()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	body := rec.Body.String()

	// Check for essential UI elements
	assert.Contains(t, body, "Git-Migrator", "Should have title")
	assert.Contains(t, body, "migration", "Should have migration section")
}

// TestWebUINewMigrationPage tests the new migration page
func TestWebUINewMigrationPage(t *testing.T) {
	server := web.NewServer(web.ServerConfig{
		Port: 8080,
	})
	router := server.Router()

	req := httptest.NewRequest(http.MethodGet, "/new", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should serve HTML
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")
}

// TestWebUIConfigPage tests the configuration page
func TestWebUIConfigPage(t *testing.T) {
	server := web.NewServer(web.ServerConfig{
		Port: 8080,
	})
	router := server.Router()

	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should serve HTML
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")
}

// TestWebUIMigrationPage tests viewing a specific migration
func TestWebUIMigrationPage(t *testing.T) {
	server := web.NewServer(web.ServerConfig{
		Port: 8080,
	})
	router := server.Router()

	req := httptest.NewRequest(http.MethodGet, "/migration/test-id", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should serve HTML (even for non-existent migration)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")
}
