// Package config provides a minimal configuration-loading utility: it reads a
// config file into a struct implementing Config, applying defaults, overrides,
// and basic validation.
package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config is the interface implemented by configuration structs that can be
// loaded from a file.
type Config interface {
	SetDefaults()
	Override() error
	BasicCheck() error
}

type loader = func(cfg Config, r io.Reader) error

func loadFromFile(cfg Config, path string, loader loader) error {
	cfg.SetDefaults()

	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = loader(cfg, file); err != nil {
		return fmt.Errorf("failed to load config from file: %w", err)
	}

	if err = cfg.Override(); err != nil {
		return fmt.Errorf("failed to override config: %w", err)
	}

	if err = cfg.BasicCheck(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	return nil
}

func loaderYAML(cfg Config, r io.Reader) error {
	decoder := yaml.NewDecoder(r)
	decoder.KnownFields(true)

	if decodeErr := decoder.Decode(cfg); decodeErr != nil {
		return fmt.Errorf("failed to decode config file: %w", decodeErr)
	}

	return nil
}

// LoadFromYAML loads cfg from the YAML file at path, applying defaults,
// environment overrides, and basic validation.
func LoadFromYAML(cfg Config, path string) error {
	return loadFromFile(cfg, path, loaderYAML)
}
