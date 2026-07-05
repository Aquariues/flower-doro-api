package database

import (
	"embed"
	"fmt"
	"path"
	"sort"
	"strings"

	"gorm.io/gorm"
)

//go:embed migrations/*.up.sql
var migrationFiles embed.FS

func RunMigrations(db *gorm.DB) error {
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version varchar(255) PRIMARY KEY,
		applied_at timestamptz NOT NULL DEFAULT now()
	)`).Error; err != nil {
		return err
	}

	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)

	return db.Transaction(func(tx *gorm.DB) error {
		for _, name := range names {
			version := strings.TrimSuffix(name, ".up.sql")
			var applied bool
			if err := tx.Raw(`SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = ?)`, version).Scan(&applied).Error; err != nil {
				return err
			}
			if applied {
				continue
			}

			sqlBytes, err := migrationFiles.ReadFile(path.Join("migrations", name))
			if err != nil {
				return err
			}
			if err := tx.Exec(string(sqlBytes)).Error; err != nil {
				return fmt.Errorf("apply migration %s: %w", version, err)
			}
			if err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, version).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
