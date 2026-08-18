// Package config provides a minimal configuration-loading utility: it reads a
// config file into a struct implementing Config, applying defaults, overrides,
// and basic validation.
package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// Config is the interface implemented by configuration structs that can be
// loaded from a file.
type Config interface {
	BasicCheck() error

	Override()
}

// parseOptions holds the configuration settings for the parsing operation.
type parseOptions struct {
	Strict bool
}

// Option defines a function signature used to configure parseOptions.
type Option func(*parseOptions)

// WithStrict configures whether the parser should run in strict mode.
func WithStrict(strict bool) Option {
	return func(o *parseOptions) {
		o.Strict = strict
	}
}

type loader = func(cfg Config, r io.Reader, opts *parseOptions) error

func loadFromFile(cfg Config, path string, loader loader, opts ...Option) error {
	parsOpts := &parseOptions{Strict: false} // default value
	for _, opt := range opts {
		opt(parsOpts)
	}

	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = loader(cfg, file, parsOpts); err != nil {
		return fmt.Errorf("failed to load config from file: %w", err)
	}

	cfg.Override()

	if err = cfg.BasicCheck(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	return nil
}

func loaderYAML(cfg Config, r io.Reader, opts *parseOptions) error {
	decoder := yaml.NewDecoder(r)

	if opts.Strict {
		decoder.KnownFields(true)
	}

	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("failed to decode config file: %w", err)
	}

	return nil
}

// LoadFromYAML loads cfg from the YAML file at path,
// applying overrides, and basic validation.
func LoadFromYAML(cfg Config, path string, opts ...Option) error {
	return loadFromFile(cfg, path, loaderYAML, opts...)
}

func loaderTOML(cfg Config, r io.Reader, opts *parseOptions) error {
	decoder := toml.NewDecoder(r)
	if opts.Strict {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("failed to decode config file: %w", err)
	}

	return nil
}

// LoadFromTOML loads cfg from the TOML file at path,
// applying overrides, and basic validation.
func LoadFromTOML(cfg Config, path string, opts ...Option) error {
	return loadFromFile(cfg, path, loaderTOML, opts...)
}
