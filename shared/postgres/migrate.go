package postgres

import (
	"database/sql"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

func RunMigrations(db *sql.DB, migrations fs.FS, dir string) error {
	entries, err := fs.ReadDir(migrations, dir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)

	for _, name := range files {
		path := dir + "/" + name
		sqlBytes, err := fs.ReadFile(migrations, path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		if _, err := db.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("apply migration %s: %w", name, err)
		}
	}

	return nil
}
