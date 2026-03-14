package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/Kurome00/weather-app.git/internal/pkg/app/cli"
)

func main() {
    var (
        loggerType = flag.String("logger", "console", "Тип логгера: console, file, json, multi")
        debugMode  = flag.Bool("debug", false, "Режим отладки")
        logFile    = flag.String("logfile", "app.log", "Файл для логов")
        latitude   = flag.Float64("lat", 53.6688, "Широта")
        longitude  = flag.Float64("lon", 23.8223, "Долгота")
    )
    flag.Parse()

    config := cli.Config{
        Latitude:  *latitude,
        Longitude: *longitude,
        DebugMode: *debugMode,
    }

    var logger cli.Logger

    switch *loggerType {
    case "console":
        logger = cli.NewConsoleLogger(*debugMode)
        fmt.Println("Используется консольный логгер")

    case "file":
        fileLogger, err := cli.NewFileLogger(*logFile, *debugMode)
        if err != nil {
            fmt.Printf("Ошибка создания файлового логгера: %s\n", err)
            os.Exit(1)
        }
        defer fileLogger.Close()
        logger = fileLogger
        fmt.Printf("Используется файловый логгер: %s\n", *logFile)

    case "json":
        logger = cli.NewJSONLogger(*debugMode)
        fmt.Println("Используется JSON логгер")

    case "multi":
        consoleLogger := cli.NewConsoleLogger(*debugMode)
        fileLogger, err := cli.NewFileLogger(*logFile, *debugMode)
        if err != nil {
            fmt.Printf("Ошибка создания файлового логгера: %s\n", err)
            os.Exit(1)
        }
        defer fileLogger.Close()
        logger = cli.NewMultiLogger(consoleLogger, fileLogger)
        fmt.Println("Используется многоцелевой логгер (консоль + файл)")

    default:
        fmt.Printf("Неизвестный тип логгера: %s\n", *loggerType)
        os.Exit(1)
    }

    app := cli.New(logger, config)
    
    if err := app.Run(); err != nil {
        logger.Error(fmt.Sprintf("Ошибка приложения: %s", err.Error()))
        os.Exit(1)
    }
}