package config

import (
    "io"
    "time"
    "gopkg.in/yaml.v3"
)

type ConfigFile struct {
    Service ServiceConfig `yaml:"service"`
}

type Provider struct {
    Type string `yaml:"type"`
}

type Location struct {
    Lat  float64 `yaml:"lat"`
    Long float64 `yaml:"long"`
}

type Cache struct {
    Enabled bool          `yaml:"enabled"`
    TTL     time.Duration `yaml:"ttl"`
}

type ServiceConfig struct {
    P Provider `yaml:"provider"`
    L Location `yaml:"location"`
    C Cache    `yaml:"cache"`
}

type Config struct {
    P Provider `yaml:"provider"`
    L Location `yaml:"location"`
    C Cache    `yaml:"cache"`
}

func Parse(r io.Reader) (Config, error) {
    var cf ConfigFile
    if err := yaml.NewDecoder(r).Decode(&cf); err != nil {
        return Config{}, err
    }
    
    return Config{
        P: cf.Service.P,
        L: cf.Service.L,
        C: cf.Service.C,
    }, nil
}