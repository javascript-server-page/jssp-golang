package db

import (
	"database/sql"
	"sync"
)

type DBCache struct {
	// lock sync.RWMutex
	// data map[string]*sql.DB
	data *sync.Map
}

func NewDBCache() *DBCache {
	return &DBCache{new(sync.Map)}
}

func (c *DBCache) GetDB(driverName, dataSourceName string) (*sql.DB, error) {
	val, ok := c.data.Load(dataSourceName)
	if ok {
		db := val.(*sql.DB)
		if db.Ping() != nil {
			db, err := sql.Open(driverName, dataSourceName)
			if err != nil {
				return nil, err
			}
			c.data.Store(dataSourceName, db)
			return db, nil
		}
		return db, nil
	}
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	c.data.Store(dataSourceName, db)
	return db, nil
}
