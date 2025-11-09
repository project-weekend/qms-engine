package config

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/project-weekend/qms-engine/server/config"

	_ "github.com/go-sql-driver/mysql"
)

// NewDatabase initializes and returns master and slave database connections using sqlx
func NewDatabase(appCfg *config.Config, logger *slog.Logger) *sqlx.DB {
	logger.Info("Initializing database connections...")
	username := appCfg.Database.Username
	password := appCfg.Database.Password
	host := appCfg.Database.Host
	port := appCfg.Database.Port
	database := appCfg.Database.Name
	idleConnection := appCfg.Database.Pool.Idle
	maxConnection := appCfg.Database.Pool.Max
	maxLifeTimeConnection := appCfg.Database.Pool.Lifetime

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		log.Fatal(err)
	}

	db.SetMaxIdleConns(idleConnection)
	db.SetMaxOpenConns(maxConnection)
	db.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	if err != nil {
		logger.Error("Failed to connect database", "error", err)
		log.Fatalf("failed to connect database: %v", err)
	}

	return db
}
