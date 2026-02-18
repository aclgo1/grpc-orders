package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (

	//go:embed *.up.sql
	engineMigrationsFS embed.FS
	appMigrationsFS    fs.FS
	DB                 *sqlx.DB
)

type migrationEntry struct {
	id       string
	filename string
	fsys     fs.FS
}

func SetAppMigrations(db *sqlx.DB, fs fs.FS) error {

	DB = db
	appMigrationsFS = fs

	return nil
}

func checkTableExists(tx *sql.Tx) (bool, error) {
	const query = `SELECT count(*) FROM information_schema.tables
	WHERE table_schema = 'public' AND table_name = 'schema_migrations';`

	var count int
	if err := tx.QueryRow(query).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to check if schema_migrations table exists: %w", err)
	}
	return count > 0, nil
}

func createMigrationsTable(tx *sql.Tx) error {
	const query = `CREATE TABLE IF NOT EXISTS schema_migrations 
	(id TEXT PRIMARY KEY, applied_at TEXT DEFAULT CURRENT_TIMESTAMP)`

	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	return nil
}

// get all
func getAppliedMigrations(tx *sql.Tx) (map[string]bool, error) {
	const query = `SELECT id FROM schema_migrations`

	rows, err := tx.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}

	defer rows.Close()

	applied := make(map[string]bool)

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan migration id: %w", err)
		}

		applied[id] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating applied migrations: %w", err)
	}

	return applied, nil
}

// insert
func recordMigration(id string, tx *sql.Tx) error {
	const query = `INSERT INTO schema_migrations (id) VALUES ($1)`

	if _, err := tx.Exec(query, id); err != nil {
		return fmt.Errorf("failed to record migration %q: %w", id, err)
	}

	return nil
}
func getMigrationMaxTx(tx *sql.Tx) (int, error) {
	const query = `SELECT MAX(version) FROM schema_migrations`

	var max sql.NullInt64
	if err := tx.QueryRow(query).Scan(&max); err != nil {
		return 0, fmt.Errorf("failed to get max migration version: %w", err)
	}

	if !max.Valid {
		return 0, nil
	}

	return int(max.Int64), nil
}

// validate
func parseMigrationFilename(filename string) (string, error) {
	if !strings.HasSuffix(filename, ".up.sql") {
		return "", fmt.Errorf("invalid migration filename %q: must end with .up.sql", filename)
	}

	id := strings.TrimSuffix(filename, ".up.sql")

	if len(id) < 6 {
		return "", fmt.Errorf("invalid migration filename %q: too short", filename)
	}

	prefix := id[:4]

	for _, c := range prefix {
		if c < '0' || c > '9' {
			return "", fmt.Errorf("invalid migration filename %q: must start with 4-digit prefix", filename)
		}
	}

	if id[4] != '_' {
		return "", fmt.Errorf("invalid migration filename %q: digit prefix must be followed by underscore", filename)
	}

	return id, nil
}

func collectMigrations(fsys fs.FS) ([]migrationEntry, error) {
	if fsys == nil {
		return nil, nil
	}

	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []migrationEntry

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		id, err := parseMigrationFilename(name)
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, migrationEntry{
			id:       id,
			filename: name,
			fsys:     fsys,
		})
	}

	return migrations, nil
}

func Run() error {

	if DB == nil {
		return errors.New("DB *sqlx.DB is nil")
	}

	engineMigrations, err := collectMigrations(engineMigrationsFS)
	if err != nil {
		return fmt.Errorf("failed to collect engine migrations: %w", err)
	}

	appMigrations, err := collectMigrations(appMigrationsFS)
	if err != nil {
		return fmt.Errorf("failed to collect application migrations: %w", err)
	}

	allMigrations := append(engineMigrations, appMigrations...)

	seen := make(map[string]string)

	for _, m := range allMigrations {
		if existing, ok := seen[m.id]; ok {
			return fmt.Errorf("duplicate migration id %q: found in %q and %q", m.id, existing, m.filename)
		}

		seen[m.id] = m.filename
	}

	sort.Slice(allMigrations, func(i, j int) bool {
		return allMigrations[i].id < allMigrations[j].id
	})

	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transcation: %w", err)
	}

	defer func() {
		if tx != nil {
			rberr := tx.Rollback()
			if rberr != nil {
				log.Printf("failed to rollback transaction: %v", rberr)
			}
		}
	}()

	exists, err := checkTableExists(tx)
	if !exists {
		err = createMigrationsTable(tx)
		if err != nil {
			return fmt.Errorf("failed to ensure schema_migrations table exists: %w", err)
		}
	}

	//get all ids
	applied, err := getAppliedMigrations(tx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedCount := 0

	for _, m := range allMigrations {
		if applied[m.id] {
			continue
		}
		content, err := fs.ReadFile(m.fsys, m.filename)
		if err != nil {
			return fmt.Errorf("failed to read migration file %q: %w", m.filename, err)
		}

		log.Printf("applying migrations: %s", m.id)

		//exec sql contein in files
		if _, err := tx.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to apply migration %q: %w", m.id, err)
		}

		//insert id database
		if err := recordMigration(m.id, tx); err != nil {
			return err
		}

		appliedCount++
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	tx = nil

	if appliedCount == 0 {
		log.Printf("no new migrations to apply")
	} else {
		log.Printf("applied %d migration(s)", appliedCount)
	}

	return nil
}
