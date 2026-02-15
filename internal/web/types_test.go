package web

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIResponseSuccess(t *testing.T) {
	response := APIResponse{
		Success: true,
		Data:    "test data",
	}

	if !response.Success {
		t.Error("Success should be true")
	}

	if response.Data != "test data" {
		t.Errorf("Data = %v, want 'test data'", response.Data)
	}

	if response.Error != nil {
		t.Error("Error should be nil for success response")
	}
}

func TestAPIResponseError(t *testing.T) {
	response := APIResponse{
		Success: false,
		Error: &APIError{
			Code:    "TEST_ERROR",
			Message: "Test error message",
		},
	}

	if response.Success {
		t.Error("Success should be false")
	}

	if response.Error == nil {
		t.Fatal("Error should not be nil")
	}

	if response.Error.Code != "TEST_ERROR" {
		t.Errorf("Error.Code = %q, want %q", response.Error.Code, "TEST_ERROR")
	}

	if response.Error.Message != "Test error message" {
		t.Errorf("Error.Message = %q, want %q", response.Error.Message, "Test error message")
	}
}

func TestAPIResponseJSON(t *testing.T) {
	tests := []struct {
		name     string
		response APIResponse
	}{
		{
			name: "success with string data",
			response: APIResponse{
				Success: true,
				Data:    "test",
			},
		},
		{
			name: "success with map data",
			response: APIResponse{
				Success: true,
				Data: map[string]interface{}{
					"id":   "123",
					"name": "test",
				},
			},
		},
		{
			name: "success with nil data",
			response: APIResponse{
				Success: true,
				Data:    nil,
			},
		},
		{
			name: "error response",
			response: APIResponse{
				Success: false,
				Error: &APIError{
					Code:    "ERR001",
					Message: "Something went wrong",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.response)
			require.NoError(t, err)

			var decoded APIResponse
			err = json.Unmarshal(data, &decoded)
			require.NoError(t, err)

			if decoded.Success != tt.response.Success {
				t.Errorf("Success = %v, want %v", decoded.Success, tt.response.Success)
			}
		})
	}
}

func TestAPIError(t *testing.T) {
	err := &APIError{
		Code:    "VALIDATION_ERROR",
		Message: "Invalid input",
	}

	if err.Code != "VALIDATION_ERROR" {
		t.Errorf("Code = %q, want %q", err.Code, "VALIDATION_ERROR")
	}

	if err.Message != "Invalid input" {
		t.Errorf("Message = %q, want %q", err.Message, "Invalid input")
	}
}

func TestAPIErrorJSON(t *testing.T) {
	apiErr := &APIError{
		Code:    "TEST_CODE",
		Message: "Test message",
	}
	data, err := json.Marshal(apiErr)
	require.NoError(t, err)

	var decoded APIError
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	if decoded.Code != "TEST_CODE" {
		t.Errorf("Code = %q, want %q", decoded.Code, "TEST_CODE")
	}
	if decoded.Message != "Test message" {
		t.Errorf("Message = %q, want %q", decoded.Message, "Test message")
	}
}

func TestStartMigrationRequest(t *testing.T) {
	req := StartMigrationRequest{
		SourceType: "cvs",
		SourcePath: "/tmp/source",
		TargetPath: "/tmp/target",
		Options: map[string]interface{}{
			"dryRun":  true,
			"verbose": true,
		},
	}

	if req.SourceType != "cvs" {
		t.Errorf("SourceType = %q, want %q", req.SourceType, "cvs")
	}

	if req.SourcePath != "/tmp/source" {
		t.Errorf("SourcePath = %q, want %q", req.SourcePath, "/tmp/source")
	}

	if req.TargetPath != "/tmp/target" {
		t.Errorf("TargetPath = %q, want %q", req.TargetPath, "/tmp/target")
	}

	if req.Options["dryRun"] != true {
		t.Error("dryRun option should be true")
	}
}

func TestStartMigrationRequestJSON(t *testing.T) {
	jsonData := `{
		"sourceType": "svn",
		"sourcePath": "/svn/repo",
		"targetPath": "/git/repo",
		"options": {
			"dryRun": false
		}
	}`

	var req StartMigrationRequest
	err := json.Unmarshal([]byte(jsonData), &req)
	require.NoError(t, err)

	if req.SourceType != "svn" {
		t.Errorf("SourceType = %q, want %q", req.SourceType, "svn")
	}

	if req.SourcePath != "/svn/repo" {
		t.Errorf("SourcePath = %q, want %q", req.SourcePath, "/svn/repo")
	}
}

func TestStartMigrationRequestEmpty(t *testing.T) {
	req := StartMigrationRequest{}

	if req.SourceType != "" {
		t.Errorf("SourceType should be empty, got %q", req.SourceType)
	}

	if req.Options != nil {
		t.Error("Options should be nil")
	}
}

func TestAnalyzeRequest(t *testing.T) {
	req := AnalyzeRequest{
		SourceType: "cvs",
		SourcePath: "/cvs/repo",
	}

	if req.SourceType != "cvs" {
		t.Errorf("SourceType = %q, want %q", req.SourceType, "cvs")
	}

	if req.SourcePath != "/cvs/repo" {
		t.Errorf("SourcePath = %q, want %q", req.SourcePath, "/cvs/repo")
	}
}

func TestAnalyzeRequestJSON(t *testing.T) {
	jsonData := `{
		"sourceType": "git",
"sourcePath": "/path/to/repo"
	}`

	var req AnalyzeRequest
	err := json.Unmarshal([]byte(jsonData), &req)
	require.NoError(t, err)

	if req.SourceType != "git" {
		t.Errorf("SourceType = %q, want %q", req.SourceType, "git")
	}
}

func TestMigrationStatus(t *testing.T) {
	now := time.Now()
	status := MigrationStatus{
		ID:               "migration-123",
		Status:           "running",
		Percentage:       50,
		CurrentStep:      "Processing commits",
		TotalCommits:     100,
		ProcessedCommits: 50,
		Errors:           []string{"error1", "error2"},
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if status.ID != "migration-123" {
		t.Errorf("ID = %q, want %q", status.ID, "migration-123")
	}

	if status.Status != "running" {
		t.Errorf("Status = %q, want %q", status.Status, "running")
	}

	if status.Percentage != 50 {
		t.Errorf("Percentage = %d, want 50", status.Percentage)
	}

	if status.CurrentStep != "Processing commits" {
		t.Errorf("CurrentStep = %q, want %q", status.CurrentStep, "Processing commits")
	}

	if status.TotalCommits != 100 {
		t.Errorf("TotalCommits = %d, want 100", status.TotalCommits)
	}

	if status.ProcessedCommits != 50 {
		t.Errorf("ProcessedCommits = %d, want 50", status.ProcessedCommits)
	}

	if len(status.Errors) != 2 {
		t.Errorf("Errors length = %d, want 2", len(status.Errors))
	}
}

func TestMigrationStatusJSON(t *testing.T) {
	status := MigrationStatus{
		ID:               "test-id",
		Status:           "completed",
		Percentage:       100,
		TotalCommits:     50,
		ProcessedCommits: 50,
		Errors:           []string{},
	}

	data, err := json.Marshal(status)
	require.NoError(t, err)

	var decoded MigrationStatus
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	if decoded.ID != "test-id" {
		t.Errorf("ID = %q, want %q", decoded.ID, "test-id")
	}

	if decoded.Status != "completed" {
		t.Errorf("Status = %q, want %q", decoded.Status, "completed")
	}
}

func TestMigrationStatusEmptyErrors(t *testing.T) {
	status := MigrationStatus{
		ID:     "test",
		Errors: []string{},
	}

	if status.Errors == nil {
		t.Error("Errors should be empty slice, not nil")
	}

	if len(status.Errors) != 0 {
		t.Errorf("Errors length = %d, want 0", len(status.Errors))
	}
}

func TestMigrationStatusNilErrors(t *testing.T) {
	status := MigrationStatus{
		ID:     "test",
		Errors: nil,
	}

	if status.Errors != nil {
		t.Error("Errors should be nil")
	}
}

func TestProgressEvent(t *testing.T) {
	event := ProgressEvent{
		Type: "progress",
		Data: ProgressData{
			MigrationID:      "migration-456",
			Status:           "running",
			Percentage:       75,
			CurrentStep:      "Step 3 of 4",
			TotalCommits:     200,
			ProcessedCommits: 150,
			Errors:           []string{},
		},
	}

	if event.Type != "progress" {
		t.Errorf("Type = %q, want %q", event.Type, "progress")
	}

	if event.Data.MigrationID != "migration-456" {
		t.Errorf("MigrationID = %q, want %q", event.Data.MigrationID, "migration-456")
	}

	if event.Data.Percentage != 75 {
		t.Errorf("Percentage = %d, want 75", event.Data.Percentage)
	}
}

func TestProgressEventJSON(t *testing.T) {
	event := ProgressEvent{
		Type: "status",
		Data: ProgressData{
			MigrationID: "test-id",
			Status:      "completed",
			Percentage:  100,
		},
	}

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var decoded ProgressEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	if decoded.Type != "status" {
		t.Errorf("Type = %q, want %q", decoded.Type, "status")
	}

	if decoded.Data.MigrationID != "test-id" {
		t.Errorf("MigrationID = %q, want %q", decoded.Data.MigrationID, "test-id")
	}
}

func TestProgressData(t *testing.T) {
	data := ProgressData{
		MigrationID:      "test-migration",
		Status:           "in_progress",
		Percentage:       33,
		CurrentStep:      "Processing",
		TotalCommits:     300,
		ProcessedCommits: 100,
		Errors:           []string{"warning1"},
	}

	if data.MigrationID != "test-migration" {
		t.Errorf("MigrationID = %q, want %q", data.MigrationID, "test-migration")
	}

	if data.Status != "in_progress" {
		t.Errorf("Status = %q, want %q", data.Status, "in_progress")
	}

	if data.Percentage != 33 {
		t.Errorf("Percentage = %d, want 33", data.Percentage)
	}
}

func TestProgressDataJSON(t *testing.T) {
	data := ProgressData{
		MigrationID:      "json-test",
		Status:           "pending",
		Percentage:       0,
		CurrentStep:      "Initializing",
		TotalCommits:     0,
		ProcessedCommits: 0,
		Errors:           nil,
	}

	jsonBytes, err := json.Marshal(data)
	require.NoError(t, err)

	var decoded ProgressData
	err = json.Unmarshal(jsonBytes, &decoded)
	require.NoError(t, err)

	if decoded.MigrationID != "json-test" {
		t.Errorf("MigrationID = %q, want %q", decoded.MigrationID, "json-test")
	}
}

func TestServerConfig(t *testing.T) {
	config := ServerConfig{
		Port:         8080,
		ConfigPath:   "/etc/migrator/config.yaml",
		DatabasePath: "/var/lib/migrator/state.db",
	}

	if config.Port != 8080 {
		t.Errorf("Port = %d, want 8080", config.Port)
	}

	if config.ConfigPath != "/etc/migrator/config.yaml" {
		t.Errorf("ConfigPath = %q, want %q", config.ConfigPath, "/etc/migrator/config.yaml")
	}

	if config.DatabasePath != "/var/lib/migrator/state.db" {
		t.Errorf("DatabasePath = %q, want %q", config.DatabasePath, "/var/lib/migrator/state.db")
	}
}

func TestServerConfigEmpty(t *testing.T) {
	config := ServerConfig{}

	if config.Port != 0 {
		t.Errorf("Port = %d, want 0", config.Port)
	}

	if config.ConfigPath != "" {
		t.Errorf("ConfigPath = %q, want empty", config.ConfigPath)
	}
}

func TestHealthStatus(t *testing.T) {
	status := HealthStatus{
		Status:  "ok",
		Version: "1.0.0",
	}

	if status.Status != "ok" {
		t.Errorf("Status = %q, want %q", status.Status, "ok")
	}

	if status.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", status.Version, "1.0.0")
	}
}

func TestHealthStatusJSON(t *testing.T) {
	status := HealthStatus{
		Status:  "healthy",
		Version: "2.0.0",
	}

	data, err := json.Marshal(status)
	require.NoError(t, err)

	var decoded HealthStatus
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	if decoded.Status != "healthy" {
		t.Errorf("Status = %q, want %q", decoded.Status, "healthy")
	}

	if decoded.Version != "2.0.0" {
		t.Errorf("Version = %q, want %q", decoded.Version, "2.0.0")
	}
}

func TestConfigData(t *testing.T) {
	config := ConfigData{
		ChunkSize: 100,
		Verbose:   true,
		DryRun:    false,
	}

	if config.ChunkSize != 100 {
		t.Errorf("ChunkSize = %d, want 100", config.ChunkSize)
	}

	if config.Verbose != true {
		t.Error("Verbose should be true")
	}

	if config.DryRun != false {
		t.Error("DryRun should be false")
	}
}

func TestConfigDataJSON(t *testing.T) {
	config := ConfigData{
		ChunkSize: 500,
		Verbose:   false,
		DryRun:    true,
	}

	data, err := json.Marshal(config)
	require.NoError(t, err)

	var decoded ConfigData
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	if decoded.ChunkSize != 500 {
		t.Errorf("ChunkSize = %d, want 500", decoded.ChunkSize)
	}

	if decoded.Verbose != false {
		t.Error("Verbose should be false")
	}

	if decoded.DryRun != true {
		t.Error("DryRun should be true")
	}
}

func TestErrorResponse(t *testing.T) {
	response := ErrorResponse("NOT_FOUND", "Resource not found")

	if response.Success != false {
		t.Error("Success should be false")
	}

	if response.Data != nil {
		t.Error("Data should be nil")
	}

	if response.Error == nil {
		t.Fatal("Error should not be nil")
	}

	if response.Error.Code != "NOT_FOUND" {
		t.Errorf("Error.Code = %q, want %q", response.Error.Code, "NOT_FOUND")
	}

	if response.Error.Message != "Resource not found" {
		t.Errorf("Error.Message = %q, want %q", response.Error.Message, "Resource not found")
	}
}

func TestErrorResponseEmptyCode(t *testing.T) {
	response := ErrorResponse("", "Error message")

	if response.Error.Code != "" {
		t.Errorf("Error.Code = %q, want empty", response.Error.Code)
	}

	if response.Error.Message != "Error message" {
		t.Errorf("Error.Message = %q, want %q", response.Error.Message, "Error message")
	}
}

func TestErrorResponseEmptyMessage(t *testing.T) {
	response := ErrorResponse("CODE", "")

	if response.Error.Code != "CODE" {
		t.Errorf("Error.Code = %q, want %q", response.Error.Code, "CODE")
	}

	if response.Error.Message != "" {
		t.Errorf("Error.Message = %q, want empty", response.Error.Message)
	}
}

func TestSuccessResponse(t *testing.T) {
	response := SuccessResponse("test data")

	if response.Success != true {
		t.Error("Success should be true")
	}

	if response.Data != "test data" {
		t.Errorf("Data = %v, want 'test data'", response.Data)
	}

	if response.Error != nil {
		t.Error("Error should be nil")
	}
}

func TestSuccessResponseWithMap(t *testing.T) {
	data := map[string]interface{}{
		"id":     "123",
		"status": "active",
		"count":  10,
	}

	response := SuccessResponse(data)

	if !response.Success {
		t.Error("Success should be true")
	}

	responseData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Data should be a map")
	}

	if responseData["id"] != "123" {
		t.Errorf("id = %v, want '123'", responseData["id"])
	}

	if responseData["status"] != "active" {
		t.Errorf("status = %v, want 'active'", responseData["status"])
	}
}

func TestSuccessResponseWithNil(t *testing.T) {
	response := SuccessResponse(nil)

	if !response.Success {
		t.Error("Success should be true")
	}

	if response.Data != nil {
		t.Error("Data should be nil")
	}
}

func TestSuccessResponseWithStruct(t *testing.T) {
	status := MigrationStatus{
		ID:     "struct-test",
		Status: "running",
	}

	response := SuccessResponse(status)

	if !response.Success {
		t.Error("Success should be true")
	}

	data, ok := response.Data.(MigrationStatus)
	if !ok {
		t.Fatal("Data should be MigrationStatus")
	}

	if data.ID != "struct-test" {
		t.Errorf("ID = %q, want %q", data.ID, "struct-test")
	}
}

func TestSuccessResponseJSONRoundTrip(t *testing.T) {
	original := SuccessResponse(map[string]string{
		"key": "value",
	})

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded APIResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	if !decoded.Success {
		t.Error("Success should be true after round trip")
	}
}

func TestErrorResponseJSONRoundTrip(t *testing.T) {
	original := ErrorResponse("TEST_ERROR", "Test error message")

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded APIResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	if decoded.Success {
		t.Error("Success should be false after round trip")
	}

	if decoded.Error == nil {
		t.Fatal("Error should not be nil after round trip")
	}

	if decoded.Error.Code != "TEST_ERROR" {
		t.Errorf("Error.Code = %q, want %q", decoded.Error.Code, "TEST_ERROR")
	}
}

func TestMigrationStatusZeroValues(t *testing.T) {
	status := MigrationStatus{}

	if status.ID != "" {
		t.Errorf("ID = %q, want empty", status.ID)
	}

	if status.Status != "" {
		t.Errorf("Status = %q, want empty", status.Status)
	}

	if status.Percentage != 0 {
		t.Errorf("Percentage = %d, want 0", status.Percentage)
	}

	if status.TotalCommits != 0 {
		t.Errorf("TotalCommits = %d, want 0", status.TotalCommits)
	}

	if status.ProcessedCommits != 0 {
		t.Errorf("ProcessedCommits = %d, want 0", status.ProcessedCommits)
	}
}

func TestProgressDataZeroValues(t *testing.T) {
	data := ProgressData{}

	if data.MigrationID != "" {
		t.Errorf("MigrationID = %q, want empty", data.MigrationID)
	}

	if data.Status != "" {
		t.Errorf("Status = %q, want empty", data.Status)
	}

	if data.Percentage != 0 {
		t.Errorf("Percentage = %d, want 0", data.Percentage)
	}
}

func TestMigrationStatusTimeFields(t *testing.T) {
	before := time.Now()
	time.Sleep(1 * time.Millisecond)
	status := MigrationStatus{
		ID:        "time-test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	time.Sleep(1 * time.Millisecond)
	after := time.Now()

	if status.CreatedAt.Before(before) {
		t.Error("CreatedAt should be after 'before'")
	}

	if status.CreatedAt.After(after) {
		t.Error("CreatedAt should be before 'after'")
	}

	if status.UpdatedAt.Before(before) {
		t.Error("UpdatedAt should be after 'before'")
	}

	if status.UpdatedAt.After(after) {
		t.Error("UpdatedAt should be before 'after'")
	}
}

func TestAPIResponseWithComplexData(t *testing.T) {
	complexData := map[string]interface{}{
		"string": "value",
		"int":    42,
		"float":  3.14,
		"bool":   true,
		"array":  []string{"a", "b", "c"},
		"nested": map[string]string{"key": "value"},
		"nil":    nil,
	}

	response := SuccessResponse(complexData)

	if !response.Success {
		t.Error("Success should be true")
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Data should be a map")
	}

	assert.Equal(t, "value", data["string"])
	assert.Equal(t,
		42, data["int"])
	assert.Equal(t, 3.14, data["float"])
	assert.Equal(t, true, data["bool"])
}

func TestStartMigrationRequestWithOptionsNil(t *testing.T) {
	req := StartMigrationRequest{
		SourceType: "cvs",
		SourcePath: "/source",
		TargetPath: "/target",
		Options:    nil,
	}

	if req.Options != nil {
		t.Error("Options should be nil")
	}
}

func TestStartMigrationRequestWithOptionsEmpty(t *testing.T) {
	req := StartMigrationRequest{
		SourceType: "cvs",
		SourcePath: "/source",
		TargetPath: "/target",
		Options:    map[string]interface{}{},
	}

	if req.Options == nil {
		t.Error("Options should not be nil")
	}

	if len(req.Options) != 0 {
		t.Errorf("Options length = %d, want 0", len(req.Options))
	}
}

func TestMigrationStatusWithMultipleErrors(t *testing.T) {
	errors := []string{
		"Failed to read file",
		"Invalid format",
		"Permission denied",
	}

	status := MigrationStatus{
		ID:     "error-test",
		Errors: errors,
	}

	if len(status.Errors) != 3 {
		t.Errorf("Errors length = %d, want 3", len(status.Errors))
	}

	if status.Errors[0] != "Failed to read file" {
		t.Errorf("First error = %q, want %q", status.Errors[0], "Failed to read file")
	}
}

func TestProgressEventTypes(t *testing.T) {
	types := []string{"progress", "status", "connected", "error", "completed"}

	for _, eventType := range types {
		t.Run(eventType, func(t *testing.T) {
			event := ProgressEvent{
				Type: eventType,
				Data: ProgressData{
					Status: eventType,
				},
			}

			if event.Type != eventType {
				t.Errorf("Type = %q, want %q", event.Type, eventType)
			}
		})
	}
}

func TestConfigDataZeroValues(t *testing.T) {
	config := ConfigData{}

	if config.ChunkSize != 0 {
		t.Errorf("ChunkSize = %d, want 0", config.ChunkSize)
	}

	if config.Verbose != false {
		t.Error("Verbose should be false")
	}

	if config.DryRun != false {
		t.Error("DryRun should be false")
	}
}

func TestAPIErrorEmpty(t *testing.T) {
	err := &APIError{}

	if err.Code != "" {
		t.Errorf("Code = %q, want empty", err.Code)
	}

	if err.Message != "" {
		t.Errorf("Message = %q, want empty", err.Message)
	}
}
