package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvironmentFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	content := "# local configuration\nDOTENV_TEST_URL=https://example.com\nDOTENV_TEST_TOKEN='secret value'\nDOTENV_TEST_QUOTED=\"line\\nvalue\"\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DOTENV_TEST_URL", "https://process.example.com")
	t.Cleanup(func() {
		os.Unsetenv("DOTENV_TEST_TOKEN")
		os.Unsetenv("DOTENV_TEST_QUOTED")
	})

	if err := loadEnvironmentFile(path); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("DOTENV_TEST_URL"); got != "https://process.example.com" {
		t.Fatalf("process environment was overridden: %q", got)
	}
	if got := os.Getenv("DOTENV_TEST_TOKEN"); got != "secret value" {
		t.Fatalf("single-quoted value = %q", got)
	}
	if got := os.Getenv("DOTENV_TEST_QUOTED"); got != "line\nvalue" {
		t.Fatalf("double-quoted value = %q", got)
	}
}

func TestLoadEnvironmentFileRejectsInvalidEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("INVALID KEY=value\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := loadEnvironmentFile(path); err == nil {
		t.Fatal("invalid environment key was accepted")
	}
}

func TestLoadEnvironmentFileReplacesEmptyButNotConfiguredValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("DOTENV_EMPTY=fallback\nDOTENV_CONFIGURED=file-value\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DOTENV_EMPTY", "")
	t.Setenv("DOTENV_CONFIGURED", "process-value")
	if err := loadEnvironmentFile(path); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("DOTENV_EMPTY"); got != "fallback" {
		t.Fatalf("empty environment was not filled: %q", got)
	}
	if got := os.Getenv("DOTENV_CONFIGURED"); got != "process-value" {
		t.Fatalf("configured environment was overridden: %q", got)
	}
}
