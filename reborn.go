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
	c := NewConfig()
	for k, v := range r.kvPairs {
		if _, ok := cfgFromDB[k]; !ok {
			c.kvPairs[k] = v
		}
	}
	if len(c.kvPairs) > 0 {
		_, err := r.rdb.HMSet(r.name, c.kvPairs).Result()
		if err != nil {
			return SyncDefaultsToDBError{err}
		}
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
		r.kvPairs[k] = v
	}
}

func (r *Reborn) Persist() error {
	_, err := r.rdb.HMSet(r.name, r.kvPairs).Result()
	return err
}

func (r *Reborn) Set(key string, value interface{}) error {
	err := r.Config.SetValue(key, value)
	if err != nil {
		return err
	}

	_, err = r.rdb.HSet(r.name, key, r.kvPairs[key]).Result()
	if err != nil {
		return err
	}
	return nil
}
