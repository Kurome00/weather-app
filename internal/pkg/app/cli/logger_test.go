package cli

import (
    "bytes"
    "io"
    "os"
    "strings"
    "testing"
)

func TestConsoleLogger(t *testing.T) {
    oldStdout := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    logger := NewConsoleLogger(true)

    logger.Info("test info message")
    w.Close()

    var buf bytes.Buffer
    io.Copy(&buf, r)
    os.Stdout = oldStdout

    output := buf.String()
    if !strings.Contains(output, "test info message") {
        t.Errorf("Info output doesn't contain expected message: %s", output)
    }
}

func TestConsoleLoggerDebugMode(t *testing.T) {
    oldStdout := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    logger := NewConsoleLogger(true)
    logger.Debug("debug message")
    w.Close()

    var buf bytes.Buffer
    io.Copy(&buf, r)
    os.Stdout = oldStdout

    if !strings.Contains(buf.String(), "debug message") {
        t.Error("Debug message should be printed when debugMode=true")
    }
    r, w, _ = os.Pipe()
    os.Stdout = w

    logger = NewConsoleLogger(false)
    logger.Debug("debug message")
    w.Close()

    buf.Reset()
    io.Copy(&buf, r)
    os.Stdout = oldStdout

    if buf.String() != "" {
        t.Error("Debug message should not be printed when debugMode=false")
    }
}

func TestFileLogger(t *testing.T) {
    tmpfile, err := os.CreateTemp("", "test*.log")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpfile.Name())

    logger, err := NewFileLogger(tmpfile.Name(), true)
    if err != nil {
        t.Fatal(err)
    }
    defer logger.Close()

    logger.Info("test info")
    logger.Debug("test debug")
    logger.Error("test error")

    content, err := os.ReadFile(tmpfile.Name())
    if err != nil {
        t.Fatal(err)
    }

    contentStr := string(content)
    if !strings.Contains(contentStr, "[INFO] test info") {
        t.Error("File doesn't contain info message")
    }
    if !strings.Contains(contentStr, "[DEBUG] test debug") {
        t.Error("File doesn't contain debug message")
    }
    if !strings.Contains(contentStr, "[ERROR] test error") {
        t.Error("File doesn't contain error message")
    }
}

func TestJSONLogger(t *testing.T) {
    oldStdout := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    logger := NewJSONLogger(true)
    logger.Info("test json message")
    w.Close()

    var buf bytes.Buffer
    io.Copy(&buf, r)
    os.Stdout = oldStdout

    output := buf.String()
    if !strings.Contains(output, "test json message") {
        t.Error("JSON logger output doesn't contain message")
    }
    if !strings.Contains(output, "INFO") {
        t.Error("JSON logger output doesn't contain level")
    }
}

func TestMultiLogger(t *testing.T) {

    var buf1, buf2 bytes.Buffer

    logger1 := &testLogger{&buf1}
    logger2 := &testLogger{&buf2}

    multi := NewMultiLogger(logger1, logger2)

    multi.Info("test info")
    multi.Debug("test debug")
    multi.Error("test error")

    if !strings.Contains(buf1.String(), "test info") {
        t.Error("First logger didn't receive info message")
    }
    if !strings.Contains(buf2.String(), "test info") {
        t.Error("Second logger didn't receive info message")
    }
    if !strings.Contains(buf1.String(), "test debug") {
        t.Error("First logger didn't receive debug message")
    }
    if !strings.Contains(buf2.String(), "test error") {
        t.Error("Second logger didn't receive error message")
    }
}

type testLogger struct {
    buf *bytes.Buffer
}

func (t *testLogger) Info(msg string)  { t.buf.WriteString("INFO: " + msg + "\n") }
func (t *testLogger) Debug(msg string) { t.buf.WriteString("DEBUG: " + msg + "\n") }
func (t *testLogger) Error(msg string) { t.buf.WriteString("ERROR: " + msg + "\n") }