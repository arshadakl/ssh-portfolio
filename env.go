package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// loadEnvironment loads local development values without overriding variables
// already provided by Docker, systemd, a shell, or another process manager.
func loadEnvironment() error {
	paths := []string{".env"}
	if executable, err := os.Executable(); err == nil {
		executableEnv := filepath.Join(filepath.Dir(executable), ".env")
		if absolute, err := filepath.Abs(".env"); err != nil || !strings.EqualFold(absolute, executableEnv) {
			paths = append(paths, executableEnv)
		}
	}

	var loadErrors []error
	for _, path := range paths {
		if err := loadEnvironmentFile(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			loadErrors = append(loadErrors, err)
		}
	}
	return errors.Join(loadErrors...)
}

func loadEnvironmentFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(strings.TrimPrefix(scanner.Text(), "\ufeff"))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		key, value, found := strings.Cut(line, "=")
		key = strings.TrimSpace(key)
		if !found || !validEnvironmentKey(key) {
			return fmt.Errorf("%s:%d: invalid environment entry", path, lineNumber)
		}
		if existing, exists := os.LookupEnv(key); exists && strings.TrimSpace(existing) != "" {
			continue
		}
		value, err = parseEnvironmentValue(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("%s:%d: %w", path, lineNumber, err)
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("%s:%d: set %s: %w", path, lineNumber, key, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	return nil
}

func validEnvironmentKey(key string) bool {
	if key == "" {
		return false
	}
	for index, r := range key {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_' || (index > 0 && r >= '0' && r <= '9') {
			continue
		}
		return false
	}
	return true
}

func parseEnvironmentValue(value string) (string, error) {
	if len(value) < 2 {
		return value, nil
	}
	if value[0] == '\'' && value[len(value)-1] == '\'' {
		return value[1 : len(value)-1], nil
	}
	if value[0] == '"' {
		if value[len(value)-1] != '"' {
			return "", errors.New("unterminated quoted value")
		}
		decoded, err := strconv.Unquote(value)
		if err != nil {
			return "", fmt.Errorf("invalid quoted value: %w", err)
		}
		return decoded, nil
	}
	return value, nil
}
