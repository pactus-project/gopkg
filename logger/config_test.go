package logger

import "testing"

func TestConfigBasicCheck(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid default config",
			cfg:  DefaultConfig(),
		},
		{
			name:    "no targets",
			cfg:     &Config{},
			wantErr: false,
		},
		{
			name: "invalid target",
			cfg: &Config{
				Targets: []string{
					"console",
					"database",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.BasicCheck()

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
