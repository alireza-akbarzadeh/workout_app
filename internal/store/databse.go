package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/alireza-akbarzadeh/fem_project/internal/constants"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf(constants.Red+"💥 db: open %w"+constants.Reset, err)
	}
	fmt.Println(constants.Green + "✅ Connected to Database 🎉" + constants.Reset)

	return db, nil
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		fmt.Printf("%s❌ Error setting dialect: %v%s\n", constants.Red, err, constants.Reset)
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		fmt.Printf("%s❌ Migration failed: %v%s\n", constants.Red, err, constants.Reset)
		return fmt.Errorf("goose up: %w", err)
	}

	fmt.Printf("%s✅ Migration successful! 👍%s\n", constants.Green, constants.Reset)
	return nil
}

func MigrateFS(db *sql.DB, migrationFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)

}
