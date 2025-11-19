package migration

import (
	"embed"
	"path/filepath"
	"sort"
)

//go:embed *.sql
var migrationFiles embed.FS

type Migration struct {
	Content []byte
	Name    string
}

func GetMigrationFiles() ([]Migration, error) {
	migrationsDir := "."
	entries, err := migrationFiles.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}
		name := entry.Name()
		fullPath := filepath.Join(migrationsDir, name)
		content, err := migrationFiles.ReadFile(fullPath)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, Migration{Name: name, Content: content})
	}

	sort.Slice(migrations, func(i, j int) bool { return migrations[i].Name < migrations[j].Name })
	return migrations, nil
}
