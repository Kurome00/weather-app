package cli

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

// Logger определяет интерфейс для логирования
type Logger interface {
    Info(msg string)
    Debug(msg string)
    Error(msg string)
}

// ==================== 1. Консольный логгер ====================
// ConsoleLogger выводит логи в консоль с цветами
type ConsoleLogger struct {
    debugMode bool
}

func NewConsoleLogger(debugMode bool) *ConsoleLogger {
    return &ConsoleLogger{
        debugMode: debugMode,
    }
}

func (l *ConsoleLogger) Info(msg string) {
    fmt.Printf("\033[32m[INFO]\033[0m %s\n", msg) // Зеленый цвет
}

func (l *ConsoleLogger) Debug(msg string) {
    if l.debugMode {
        fmt.Printf("\033[36m[DEBUG]\033[0m %s\n", msg) // Голубой цвет
    }
}

func (l *ConsoleLogger) Error(msg string) {
    fmt.Printf("\033[31m[ERROR]\033[0m %s\n", msg) // Красный цвет
}

// ==================== 2. Файловый логгер ====================
// FileLogger записывает логи в файл
type FileLogger struct {
    file      *os.File
    debugMode bool
    logger    *log.Logger
}

func NewFileLogger(filename string, debugMode bool) (*FileLogger, error) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    
    return &FileLogger{
        file:      file,
        debugMode: debugMode,
        logger:    log.New(file, "", log.LstdFlags),
    }, nil
}

func (l *FileLogger) Info(msg string) {
    l.logger.Printf("[INFO] %s", msg)
}

func (l *FileLogger) Debug(msg string) {
    if l.debugMode {
        l.logger.Printf("[DEBUG] %s", msg)
    }
}

func (l *FileLogger) Error(msg string) {
    l.logger.Printf("[ERROR] %s", msg)
}

func (l *FileLogger) Close() error {
    return l.file.Close()
}

// ==================== 3. Структурированный логгер (JSON) ====================
// JSONLogger выводит логи в JSON формате
type JSONLogger struct {
    debugMode bool
}

func NewJSONLogger(debugMode bool) *JSONLogger {
    return &JSONLogger{
        debugMode: debugMode,
    }
}

type logEntry struct {
    Level   string    `json:"level"`
    Message string    `json:"message"`
    Time    time.Time `json:"time"`
}

func (l *JSONLogger) Info(msg string) {
    entry := logEntry{
        Level:   "INFO",
        Message: msg,
        Time:    time.Now(),
    }
    fmt.Printf("%s\n", toJSON(entry))
}

func (l *JSONLogger) Debug(msg string) {
    if l.debugMode {
        entry := logEntry{
            Level:   "DEBUG",
            Message: msg,
            Time:    time.Now(),
        }
        fmt.Printf("%s\n", toJSON(entry))
    }
}

func (l *JSONLogger) Error(msg string) {
    entry := logEntry{
        Level:   "ERROR",
        Message: msg,
        Time:    time.Now(),
    }
    fmt.Printf("%s\n", toJSON(entry))
}

// Вспомогательная функция для JSON
func toJSON(v interface{}) string {
    bytes, _ := json.Marshal(v)
    return string(bytes)
}

// ==================== 4. Многоцелевой логгер ====================
// MultiLogger пишет в несколько логгеров одновременно
type MultiLogger struct {
    loggers []Logger
}

func NewMultiLogger(loggers ...Logger) *MultiLogger {
    return &MultiLogger{
        loggers: loggers,
    }
}

func (m *MultiLogger) Info(msg string) {
    for _, logger := range m.loggers {
        logger.Info(msg)
    }
}

func (m *MultiLogger) Debug(msg string) {
    for _, logger := range m.loggers {
        logger.Debug(msg)
    }
}

func (m *MultiLogger) Error(msg string) {
    for _, logger := range m.loggers {
        logger.Error(msg)
    }
}