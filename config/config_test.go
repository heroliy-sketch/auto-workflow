package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("app_name: test-service\napiport: 9090\nenvironment: test\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if got.AppName != "test-service" || got.APIPort != 9090 || got.Environment != "test" {
		t.Fatalf("unexpected config: %+v", got)
	}
}

func TestLoadDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if got.AppName != "auto-workflow" || got.APIPort != 8080 || got.Environment != "development" {
		t.Fatalf("unexpected defaults: %+v", got)
	}
}
