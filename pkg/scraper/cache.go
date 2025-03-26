package scraper

import (
	"fmt"
	"os"
	"path/filepath"
)

var cacheDir = "cache"

func check_cache(slug string, req_func func(string) (*string, error)) (string, error) {
	if value, exist := get_cache(slug); exist {
		return value, nil
	}
	result, err := req_func(slug)

	if err != nil {
		return "", nil
	}
	if err := save_cache(slug, *result); err != nil {
		return "", fmt.Errorf("failed to save cache: %w", err)
	}
	return *result, err
}

func get_cache(slug string) (string, bool) {
	filename := filepath.Join(cacheDir, slug+".json")
	if _, err := os.Stat(filename); err == nil {
		body, err := os.ReadFile(filename)
		if err != nil {
			return "", false
		}
		return string(body), true
	} else if !os.IsNotExist(err) {
		fmt.Errorf("error checking cache file: %v", err)
		return "", false
	}
	return "", false
}

func save_cache(slug string, data string) error {
	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache dir: %w", err)
	}

	filename := filepath.Join(cacheDir, slug+".json")
	return os.WriteFile(filename, []byte(data), 0644)
}
