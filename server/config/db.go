package config

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

// DBConnections holds the master and slave database connections
type DBConnections struct {
	Master *sqlx.DB
	Slave  *sqlx.DB
}

// NewDatabase initializes and returns master and slave database connections using sqlx
func NewDatabase(config *Config, log *logrus.Logger) *DBConnections {
	log.Info("Initializing database connections...")

	// Initialize master database connection
	master, err := initDBConnection(config.Data.MySQL.Master, "master", log)
	if err != nil {
		log.Fatalf("Error initializing database connection: %v", err)
	}

	// Initialize slave database connection
	slave, err := initDBConnection(config.Data.MySQL.Slave, "slave", log)
	if err != nil {
		err := master.Close()
		if err != nil {
			log.Fatalf("Error initializing closing connection: %v", err)
			return nil
		}
		log.Fatalf("Error initializing database connection: %v", err)
	}

	log.Info("Database connections initialized successfully")

	return &DBConnections{
		Master: master,
		Slave:  slave,
	}
}

// initDBConnection creates a new database connection with the given configuration
func initDBConnection(dbConfig DBConfig, dbType string, log *logrus.Logger) (*sqlx.DB, error) {
	log.Infof("Connecting to %s database...", dbType)

	// Open database connection
	db, err := sqlx.Open("mysql", dbConfig.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s database connection: %w", dbType, err)
	}

	// Parse connection max lifetime
	connMaxLifetime, err := time.ParseDuration(dbConfig.ConnMaxLifetime)
	if err != nil {
		err := db.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to parse connMaxLifetime for %s database: %w", dbType, err)
	}

	// Set connection pool settings
	db.SetMaxIdleConns(dbConfig.MaxIdle)
	db.SetMaxOpenConns(dbConfig.MaxOpen)
	db.SetConnMaxLifetime(connMaxLifetime)

	// Verify the connection is alive
	if err := db.Ping(); err != nil {
		err := db.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to ping %s database: %w", dbType, err)
	}

	log.Infof("%s database connection established successfully (MaxIdle: %d, MaxOpen: %d, ConnMaxLifetime: %s)",
		dbType, dbConfig.MaxIdle, dbConfig.MaxOpen, connMaxLifetime)

	return db, nil
}

// Close closes both master and slave database connections
func (dbc *DBConnections) Close() error {
	var masterErr, slaveErr error

	if dbc.Master != nil {
		masterErr = dbc.Master.Close()
	}

	if dbc.Slave != nil {
		slaveErr = dbc.Slave.Close()
	}

	if masterErr != nil {
		return fmt.Errorf("failed to close master database: %w", masterErr)
	}

	if slaveErr != nil {
		return fmt.Errorf("failed to close slave database: %w", slaveErr)
	}

	return nil
}
