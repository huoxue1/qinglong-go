package cache

import (
	"github.com/nutsdb/nutsdb"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	db   *nutsdb.DB
	once = sync.Once{}
)

func initCache() {
	once.Do(func() {
		log.Infoln("start init cache")
		var err error
		db, err = nutsdb.Open(nutsdb.DefaultOptions, nutsdb.WithDir("./data/cache"))
		if err != nil {
			log.Errorf("open the cache error: %s", err.Error())
			return
		}
	})
}

func Exists(key string) bool {
	initCache()
	var exists bool
	_ = db.View(func(tx *nutsdb.Tx) error {
		_, err := tx.Get("default", []byte(key))
		if err != nil {
			exists = false
		} else {
			exists = true
		}
		return nil
	})
	return exists
}

func Get(key string) (string, error) {
	initCache()
	var val string
	err := db.View(func(tx *nutsdb.Tx) error {
		entry, err := tx.Get("default", []byte(key))
		if err != nil {
			return err
		}
		val = string(entry.Value)
		return nil
	})
	return val, err
}

func GetBytes(key string) ([]byte, error) {
	initCache()
	var val []byte
	err := db.View(func(tx *nutsdb.Tx) error {
		entry, err := tx.Get("default", []byte(key))
		if err != nil {
			return err
		}
		val = entry.Value
		return nil
	})
	return val, err
}

func Set(key string, value string) error {
	initCache()
	return db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put("default", []byte(key), []byte(value), 0)
	})
}

func SetBytes(key string, value []byte) error {
	initCache()
	return db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put("default", []byte(key), value, 0)
	})
}
