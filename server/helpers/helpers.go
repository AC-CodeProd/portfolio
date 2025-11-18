package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func GetRootPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Failed to get current working directory: %v", err))
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "config.yaml")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	panic("Could not find project root directory (config.yaml not found while traversing up the directory tree)")
}

func BoolPtr(b bool) *bool {
	return &b
}

func RandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)[:n]
}

func MkdirIfNotExists(path string) error {
	if path == "" {
		return errors.New("path cannot be empty")
	}
	path = filepath.Clean(path)
	if filepath.Ext(path) != "" {
		path = filepath.Dir(path)
	}

	path = filepath.Join(GetRootPath(), path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}
