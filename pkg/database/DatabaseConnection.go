package database

import (
	config "PocGo/internal/configuration"
	notify "PocGo/internal/domain/notification"
	sqlDB "database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"sync"
)

type Database struct {
	db *sqlDB.DB
	mu sync.RWMutex
}

func NewConnection(cfg *config.Config) (*Database, error) {
	if cfg.Database == nil {
		return nil, notify.CreateNotification(notify.ErrorConfigDb)
	}

	db, err := sqlDB.Open("sqlserver", cfg.Database.ConnectionString)
	if err != nil {
		return nil, notify.CreateSimpleNotification(notify.ErrorOpenConnection, err)
	}

	maxOpen := cfg.Database.MaxOpenConns
	maxIdle := cfg.Database.MaxIdleConns

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)

	if err = db.Ping(); err != nil {
		return nil, notify.CreateSimpleNotification(notify.ErrorTestConnection, err)
	}

	return &Database{
		db: db,
	}, nil
}

func (d *Database) GetConnection() *sqlDB.DB {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.db
}

func (d *Database) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.db.Close()
}
