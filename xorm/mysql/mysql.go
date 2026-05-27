package mysql

import (
	"fmt"
	"strings"

	"github.com/xxxryan/go-infra/xorm"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open(cfg xorm.Config) (*gorm.DB, error) {
	dsn := buildDSN(cfg)
	db, err := gorm.Open(gormMySQL.Open(dsn), &gorm.Config{})
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
	sys.DBName = ""
	sys.DSN = ""

	dsn := buildDSN(sys)
	db, err := gorm.Open(gormMySQL.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("xorm/mysql: connect to system db: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", quoteIdent(target))).Error
	if err != nil {
		return fmt.Errorf("xorm/mysql: create database %q: %w", target, err)
	}
	return nil
}

func TruncateAll(db *gorm.DB) error {
	var tables []string
	err := db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE() AND table_type = 'BASE TABLE'").Scan(&tables).Error
	if err != nil {
		return fmt.Errorf("xorm/mysql: list tables: %w", err)
	}
	if len(tables) == 0 {
		return nil
	}
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return fmt.Errorf("xorm/mysql: disable fk checks: %w", err)
	}
	for _, t := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", quoteIdent(t))).Error; err != nil {
			db.Exec("SET FOREIGN_KEY_CHECKS = 1")
			return fmt.Errorf("xorm/mysql: truncate table %q: %w", t, err)
		}
	}
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	return nil
}

func buildDSN(cfg xorm.Config) string {
	if cfg.DSN != "" {
		return cfg.DSN
	}
	port := cfg.Port
	if port == 0 {
		port = 3306
	}
	charset := cfg.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	loc := cfg.Loc
	if loc == "" {
		loc = "Local"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
		cfg.User, cfg.Password, cfg.Host, port, cfg.DBName, charset, loc)
}

func quoteIdent(s string) string {
	return "`" + strings.ReplaceAll(s, "`", "``") + "`"
}
