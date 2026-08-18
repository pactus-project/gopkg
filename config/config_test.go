package config_test

import (
	"os"
	"testing"

	"github.com/pactus-project/gopkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Config struct {
	Key1 string `yaml:"key1"`
	Key2 string `yaml:"key2"`
}

func (c *Config) Override() {
	val2 := os.Getenv("KEY2_OVERRIDE")
	if val2 != "" {
		c.Key2 = val2
	}
}

func (c *Config) BasicCheck() error {
	if c.Key1 == "" || c.Key2 == "" {
		return assert.AnError
	}

	return nil
}

func TestYAMLConfig(t *testing.T) {
	loadYAMLConfig := func(t *testing.T, content string, strict bool) (*Config, error) {
		t.Helper()

		dir := t.TempDir()
		if err := os.WriteFile(dir+"/config.yaml", []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}

		cfg := &Config{}

		return cfg, config.LoadFromYAML(cfg, dir+"/config.yaml", config.WithStrict(strict))
	}

	t.Run("TestYAMLSuccessfulLoad", func(t *testing.T) {
		configContent := `
key1: value1
key2: value2
`

		cfg, err := loadYAMLConfig(t, configContent, true)
		require.NoError(t, err)

		assert.Equal(t, &Config{Key1: "value1", Key2: "value2"}, cfg)
	})

	t.Run("TestYAMLStrict", func(t *testing.T) {
		configContent := `
key1: value1
key2: value2
key_unknown: value_unknown
`

		_, err := loadYAMLConfig(t, configContent, true)
		require.Error(t, err)

		_, err = loadYAMLConfig(t, configContent, false)
		require.NoError(t, err)
	})

	t.Run("TestYAMLBasicCheck", func(t *testing.T) {
		configContent := `
key1: value1
`
		_, err := loadYAMLConfig(t, configContent, true)
		require.Error(t, err)
	})

	t.Run("TestYAMLOverrideValues", func(t *testing.T) {
		configContent := `
key1: value1
key2: value2
`

		t.Setenv("KEY2_OVERRIDE", "overridden2")

		cfg, err := loadYAMLConfig(t, configContent, true)
		require.NoError(t, err)

		assert.Equal(t, &Config{Key1: "value1", Key2: "overridden2"}, cfg)
	})

	t.Run("TestYAMLUnknownField", func(t *testing.T) {
		configContent := `
invalid: yaml
`

		_, err := loadYAMLConfig(t, configContent, true)
		require.Error(t, err)
	})
}
