package logger

import (
    "bytes"
    "testing"
)

func TestLogger(t *testing.T) {
    // Create a buffer to capture log output
    var buf bytes.Buffer
    logger := New()
    logger.SetOutput(&buf)

    // Test each log level
    logger.Info("This is an info message")
    logger.Error("This is an error message")
    logger.Debug("This is a debug message")
    // Fatal won't be tested in this manner since it calls os.Exit(1)
    
    // Check if the logged messages are present in the buffer
    output := buf.String()
    if !contains(output, "INFO: This is an info message") {
        t.Errorf("Info message not logged correctly")
    }
    if !contains(output, "ERROR: This is an error message") {
        t.Errorf("Error message not logged correctly")
    }
    if !contains(output, "DEBUG: This is a debug message") {
        t.Errorf("Debug message not logged correctly")
    }
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
    return len(s) > 0 && len(substr) > 0 && s[len(s)-len(substr):] == substr
}