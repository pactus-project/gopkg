package config_test

import (
	"os"
	"testing"

	"github.com/pactus-project/gopkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Config struct {
	Key1 string `toml:"key1" yaml:"key1"`
	Key2 string `toml:"key2" yaml:"key2"`
}

func (c *Config) BasicCheck() error {
	if c.Key1 == "" || c.Key2 == "" {
		return assert.AnError
	}

	return nil
}

func TestYAMLConfig(t *testing.T) {
	loadYAMLConfig := func(t *testing.T, content string, opts ...config.Option,
	) (*Config, error) {
		t.Helper()

		dir := t.TempDir()
		err := os.WriteFile(dir+"/config.yaml", []byte(content), 0o600)
		require.NoError(t, err)

		cfg := &Config{}

		return cfg, config.LoadFromYAML(cfg, dir+"/config.yaml", opts...)
	}

	t.Run("TestYAMLSuccessfulLoad", func(t *testing.T) {
		configContent := `
key1: value1
key2: value2
`

		cfg, err := loadYAMLConfig(t, configContent, config.WithStrict(true))
		require.NoError(t, err)

		assert.Equal(t, &Config{Key1: "value1", Key2: "value2"}, cfg)
	})

	t.Run("TestYAMLStrict", func(t *testing.T) {
		configContent := `
key1: value1
key2: value2
key_unknown: value_unknown
`

		_, err := loadYAMLConfig(t, configContent, config.WithStrict(true))
		require.Error(t, err)

		_, err = loadYAMLConfig(t, configContent, config.WithStrict(false))
		require.NoError(t, err)
	})

	t.Run("TestYAMLBasicCheck", func(t *testing.T) {
		configContent := `
key1: value1
`
		_, err := loadYAMLConfig(t, configContent, config.WithStrict(true))
		require.Error(t, err)
	})

	t.Run("TestYAMLOverrideValues", func(t *testing.T) {
		configContent := `
key1: value1
key2: value2
`

		cfg, err := loadYAMLConfig(t, configContent,
			config.WithOverride(func(cfg *Config) {
				cfg.Key2 = "overridden2"
			}))
		require.NoError(t, err)

		assert.Equal(t, &Config{Key1: "value1", Key2: "overridden2"}, cfg)
	})
}

func TestTOMLConfig(t *testing.T) {
	loadTOMLConfig := func(t *testing.T, content string,
		opts ...config.Option,
	) (*Config, error) {
		t.Helper()

		dir := t.TempDir()
		err := os.WriteFile(dir+"/config.toml", []byte(content), 0o600)
		require.NoError(t, err)

		cfg := &Config{}

		return cfg, config.LoadFromTOML(cfg, dir+"/config.toml", opts...)
	}

	t.Run("TestTOMLSuccessfulLoad", func(t *testing.T) {
		configContent := `
key1 = 'value1'
key2 = 'value2'
`

		cfg, err := loadTOMLConfig(t, configContent, config.WithStrict(true))
		require.NoError(t, err)

		assert.Equal(t, &Config{Key1: "value1", Key2: "value2"}, cfg)
	})

	t.Run("TestTOMLStrict", func(t *testing.T) {
		configContent := `
key1 = 'value1'
key2 = 'value2'
key_unknown = 'value_unknown'
`

		_, err := loadTOMLConfig(t, configContent, config.WithStrict(true))
		require.Error(t, err)

		_, err = loadTOMLConfig(t, configContent, config.WithStrict(false))
		require.NoError(t, err)
	})

	t.Run("TestTOMLBasicCheck", func(t *testing.T) {
		configContent := `
key1 = 'value1'
`
		_, err := loadTOMLConfig(t, configContent, config.WithStrict(true))
		require.Error(t, err)
	})

	t.Run("TestTOMLOverrideValues", func(t *testing.T) {
		configContent := `
key1 = 'value1'
key2 = 'value2'
`

		cfg, err := loadTOMLConfig(t, configContent,
			config.WithOverride(func(cfg *Config) {
				cfg.Key2 = "overridden2"
			}))
		require.NoError(t, err)

		assert.Equal(t, &Config{Key1: "value1", Key2: "overridden2"}, cfg)
	})
}
