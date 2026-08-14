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

func (c *Config) SetDefaults() {
	c.Key1 = "default1"
}

func (c *Config) Override() error {
	val2 := os.Getenv("KEY2_OVERRIDE")
	if val2 != "" {
		c.Key2 = val2
	}

	return nil
}

func (c *Config) BasicCheck() error {
	if c.Key1 == "" || c.Key2 == "" {
		return assert.AnError
	}

	return nil
}

func loadConfig(t *testing.T, content string) (*Config, error) {
	t.Helper()

	dir := t.TempDir()
	if err := os.WriteFile(dir+"/config.yaml", []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{}

	return cfg, config.LoadFromYAML(cfg, dir+"/config.yaml")
}

func TestSuccessfulLoad(t *testing.T) {
	configContent := `
key1: value1
key2: value2
`

	cfg, err := loadConfig(t, configContent)
	require.NoError(t, err)

	assert.Equal(t, &Config{Key1: "value1", Key2: "value2"}, cfg)
}

func TestBasicCheck(t *testing.T) {
	configContent := `
key1: value1
`
	_, err := loadConfig(t, configContent)
	require.Error(t, err)
}

func TestDefaultValues(t *testing.T) {
	configContent := `
key2: value2
`

	cfg, err := loadConfig(t, configContent)
	require.NoError(t, err)

	assert.Equal(t, &Config{Key1: "default1", Key2: "value2"}, cfg)
}

func TestOverrideValues(t *testing.T) {
	configContent := `
key1: value1
key2: value2
`

	t.Setenv("KEY2_OVERRIDE", "overridden2")

	cfg, err := loadConfig(t, configContent)
	require.NoError(t, err)

	assert.Equal(t, &Config{Key1: "value1", Key2: "overridden2"}, cfg)
}

func TestUnknownField(t *testing.T) {
	configContent := `
invalid: yaml
`

	_, err := loadConfig(t, configContent)
	require.Error(t, err)
}
