package flags

import (
    "flag"
    "os"
)

type Flags struct {
    Path string
}

var (
    configPath string
    parsed     bool
)

func Parse() *Flags {
    if !parsed {
        if os.Args == nil {
            os.Args = []string{"cmd"}
        }
        
        flag.StringVar(&configPath, "config", "./config/config.yaml", "path to config")
        flag.Parse()
        parsed = true
    }

    return &Flags{
        Path: configPath,
    }
}

func ResetForTesting() {
    flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
    configPath = ""
    parsed = false
}