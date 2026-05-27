package xorm

import (
	"time"

	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string

	// DSN overrides all fields above when set.
	DSN string

	// Connection pool
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration

	// PG: "disable", "require", "verify-ca", "verify-full". Default "disable".
	SSLMode string
	// PG: e.g. "Asia/Shanghai". Default empty (server default).
	TimeZone string
	// MySQL: e.g. "utf8mb4". Default "utf8mb4".
	Charset string
	// MySQL: time.Location name for parseTime, e.g. "Local", "UTC". Default "Local".
	Loc string
}

func ConfigurePool(db *gorm.DB, cfg Config) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}
	return nil
}
