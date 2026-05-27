package pg

import (
	"fmt"
	"strings"

	"github.com/xxxryan/go-infra/xorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(cfg xorm.Config) (*gorm.DB, error) {
	dsn := buildDSN(cfg)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := xorm.ConfigurePool(db, cfg); err != nil {
		return nil, err
	}
	return db, nil
}

func CreateDatabaseIfNotExists(cfg xorm.Config) error {
	target := cfg.DBName
	sys := cfg
	sys.DBName = "postgres"
	sys.DSN = ""

	dsn := buildDSN(sys)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("xorm/pg: connect to system db: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	var exists bool
	err = db.Raw("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)", target).Scan(&exists).Error
	if err != nil {
		return fmt.Errorf("xorm/pg: check database existence: %w", err)
	}
	if exists {
		return nil
	}

	err = db.Exec(fmt.Sprintf(`CREATE DATABASE %s`, quoteIdent(target))).Error
	if err != nil {
		return fmt.Errorf("xorm/pg: create database %q: %w", target, err)
	}
	return nil
}

func TruncateAll(db *gorm.DB) error {
	var tables []string
	err := db.Raw(`SELECT tablename FROM pg_tables WHERE schemaname = 'public'`).Scan(&tables).Error
	if err != nil {
		return fmt.Errorf("xorm/pg: list tables: %w", err)
	}
	if len(tables) == 0 {
		return nil
	}
	quoted := make([]string, len(tables))
	for i, t := range tables {
		quoted[i] = quoteIdent(t)
	}
	err = db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", strings.Join(quoted, ", "))).Error
	if err != nil {
		return fmt.Errorf("xorm/pg: truncate tables: %w", err)
	}
	return nil
}

func buildDSN(cfg xorm.Config) string {
	if cfg.DSN != "" {
		return cfg.DSN
	}
	port := cfg.Port
	if port == 0 {
		port = 5432
	}
	sslmode := cfg.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}
	parts := []string{
		fmt.Sprintf("host=%s", cfg.Host),
		fmt.Sprintf("port=%d", port),
		fmt.Sprintf("user=%s", cfg.User),
		fmt.Sprintf("password=%s", cfg.Password),
		fmt.Sprintf("dbname=%s", cfg.DBName),
		fmt.Sprintf("sslmode=%s", sslmode),
	}
	if cfg.TimeZone != "" {
		parts = append(parts, fmt.Sprintf("TimeZone=%s", cfg.TimeZone))
	}
	return strings.Join(parts, " ")
}

func quoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}
