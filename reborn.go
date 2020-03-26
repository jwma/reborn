package reborn

import (
	"github.com/go-redis/redis"
)

type Reborn struct {
	*Config
	rdb  *redis.Client
	name string
}

func New(rdb *redis.Client, name string) (*Reborn, error) {
	return NewWithDefaults(rdb, name, NewConfig())
}

func NewWithDefaults(rdb *redis.Client, name string, defaults *Config) (*Reborn, error) {
	var err error
	r := &Reborn{
		Config: defaults,
		rdb:    rdb,
		name:   name,
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
}

func (r *Reborn) Persist() error {
	pipeline := r.rdb.Pipeline()
	r.kvPairs.Range(func(k, v interface{}) bool {
		pipeline.HSet(r.name, k.(string), v)
		return true
	})
	_, err := pipeline.Exec()
	return err
}

func (r *Reborn) Set(key string, value interface{}) error {
	err := r.Config.SetValue(key, value)
	if err != nil {
		return err
	}

	v, _ := r.kvPairs.Load(key)
	_, err = r.rdb.HSet(r.name, key, v).Result()
	return err
}
