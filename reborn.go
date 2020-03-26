package reborn

import (
	"github.com/go-redis/redis"
	"time"
)

type Reborn struct {
	*Config
	rdb                *redis.Client
	name               string
	ready              bool
	autoReloadTicker   *time.Ticker
	autoReloadDuration time.Duration
}

func New(rdb *redis.Client, name string) (*Reborn, error) {
	return NewWithDefaults(rdb, name, NewConfig())
}

func NewWithDefaults(rdb *redis.Client, name string, defaults *Config) (*Reborn, error) {
	var err error
	r := &Reborn{
		Config:             defaults,
		rdb:                rdb,
		name:               name,
		autoReloadDuration: time.Second * 5,
	}

	cfgFromDB, err := r.loadFromDB()
	if err != nil {
		return nil, err
	}

	err = r.syncDefaultsToDB(cfgFromDB)
	if err != nil {
		return nil, err
	}

	r.overrideDefaultsWithDBConfig(cfgFromDB)
	r.ready = true

	return r, nil
}

// 把只存在于当前 reborn 实例的配置信息同步到数据库
func (r *Reborn) syncDefaultsToDB(cfgFromDB map[string]string) error {
	pipeline := r.rdb.Pipeline()
	r.kvPairs.Range(func(k, v interface{}) bool {
		kk := k.(string)
		if _, ok := cfgFromDB[kk]; !ok {
			pipeline.HSet(r.name, kk, v)
		}
		return true
	})
	_, err := pipeline.Exec()
	if err != nil {
		return SyncDefaultsToDBError{err}
	}
	return nil
}

func (r *Reborn) loadFromDB() (map[string]string, error) {
	rs, err := r.rdb.HGetAll(r.name).Result()
	if err != nil {
		return make(map[string]string), LoadFromDBError{err}
	}
	return rs, nil
}

func (r *Reborn) overrideDefaultsWithDBConfig(cfgFromDB map[string]string) {
	for k, v := range cfgFromDB {
		r.set(k, v)
	}
	r.mux.Lock()
	r.hasChange = false
	r.mux.Unlock()
}

func (r *Reborn) Persist() error {
	pipeline := r.rdb.Pipeline()
	r.kvPairs.Range(func(k, v interface{}) bool {
		pipeline.HSet(r.name, k.(string), v)
		return true
	})
	_, err := pipeline.Exec()

	r.mux.Lock()
	r.hasChange = false
	r.mux.Unlock()

	return err
}

func (r *Reborn) SetAutoReloadDuration(d time.Duration) {
	r.autoReloadDuration = d
}

func (r *Reborn) StartAutoReload() {
	r.autoReloadTicker = time.NewTicker(r.autoReloadDuration)
	go func() {
		for {
			select {
			case <-r.autoReloadTicker.C:
				if r.ready && !r.hasChange {
					cfgFromDB, _ := r.loadFromDB()
					r.overrideDefaultsWithDBConfig(cfgFromDB)
				}
			}
		}
	}()
}

func (r *Reborn) StopAutoReload() {
	r.autoReloadTicker.Stop()
}
