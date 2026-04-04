package flags

import (
    "testing"
	"os"
)

func TestParseFlags(t *testing.T) {
    ResetForTesting()
    
    oldArgs := os.Args
    defer func() { os.Args = oldArgs }()

    os.Args = []string{"cmd", "-config", "./custom/config.yaml"}

    flags := Parse()

    if flags.Path != "./custom/config.yaml" {
        t.Errorf("Expected path './custom/config.yaml', got '%s'", flags.Path)
    }
}

func TestParseFlagsDefault(t *testing.T) {
    ResetForTesting()
    
    oldArgs := os.Args
    defer func() { os.Args = oldArgs }()

    os.Args = []string{"cmd"}

    flags := Parse()

    if flags.Path != "./config/config.yaml" {
        t.Errorf("Expected default path './config/config.yaml', got '%s'", flags.Path)
    }
}

func TestParseFlagsEmpty(t *testing.T) {
    ResetForTesting()
    
    oldArgs := os.Args
    defer func() { os.Args = oldArgs }()

    os.Args = []string{"cmd"}

    flags := Parse()

    if flags == nil {
        t.Error("Parse() returned nil")
    }
}

func TestParseFlagsMultipleTimes(t *testing.T) {
    ResetForTesting()
    
    oldArgs := os.Args
    defer func() { os.Args = oldArgs }()

    os.Args = []string{"cmd", "-config", "./test.yaml"}

    flags1 := Parse()
    flags2 := Parse()

    if flags1.Path != flags2.Path {
        t.Error("Multiple Parse() calls should return same result")
    }
}

func TestParseFlagsDifferentPaths(t *testing.T) {
    testCases := []struct {
        name     string
        args     []string
        expected string
    }{
        {"default", []string{"cmd"}, "./config/config.yaml"},
        {"custom path", []string{"cmd", "-config", "./myconfig.yaml"}, "./myconfig.yaml"},
        {"absolute path", []string{"cmd", "-config", "/etc/app/config.yaml"}, "/etc/app/config.yaml"},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            ResetForTesting()
            
            oldArgs := os.Args
            defer func() { os.Args = oldArgs }()

            os.Args = tc.args

            flags := Parse()
            if flags.Path != tc.expected {
                t.Errorf("Expected '%s', got '%s'", tc.expected, flags.Path)
            }
        })
    }
}