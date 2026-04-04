package config

import (
    "io"
    "time"
    "gopkg.in/yaml.v3"
)

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

type Config struct {
    Provider Provider `yaml:"provider"`
    Location Location `yaml:"location"`
    Cache    Cache    `yaml:"cache"`
}

func Parse(r io.Reader) (Config, error) {
    var cfg Config
    if err := yaml.NewDecoder(r).Decode(&cfg); err != nil {
        return Config{}, err
    }

    if cfg.Cache.TTL == 0 {
        cfg.Cache.TTL = 5 * time.Minute
    }

    return cfg, nil
}