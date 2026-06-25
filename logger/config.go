package logger

import "fmt"

// Config defines parameters for the logger module.
type Config struct {
	Colorful           bool              `toml:"colorful"              yaml:"colorful"`
	Filename           string            `toml:"filename"              yaml:"filename"`
	MaxSize            int               `toml:"max_size"              yaml:"max_size"`
	MaxBackups         int               `toml:"max_backups"           yaml:"max_backups"`
	RotateLogAfterDays int               `toml:"rotate_log_after_days" yaml:"rotate_log_after_days"`
	Compress           bool              `toml:"compress"              yaml:"compress"`
	Targets            []string          `toml:"targets"               yaml:"targets"`
	Levels             map[string]string `toml:"levels"                yaml:"levels"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	conf := &Config{
		Colorful:           true,
		Filename:           "app.log",
		MaxSize:            10,
		MaxBackups:         0,
		RotateLogAfterDays: 1,
		Compress:           true,
		Targets:            []string{"console", "file"},
		Levels:             make(map[string]string),
	}

	conf.Levels["default"] = "info"

	return conf
}

// BasicCheck performs basic checks on the configuration.
func (c *Config) BasicCheck() error {
	validTargets := map[string]bool{
		"console": true,
		"file":    true,
	}

	for _, target := range c.Targets {
		if _, ok := validTargets[target]; !ok {
			return fmt.Errorf("invalid logging target %q (must be 'console' or 'file')", target)
		}
	}

	return nil
}
