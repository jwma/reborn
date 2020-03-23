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
	for k, v := range defaults.kvPairs {
		err := r.Config.SetValue(k, v)
		if err != nil {
			return nil, err
		}
	}
	err = r.loadFromDB()
	return r, err
}

// 将 c.Config 的配置项目保存到数据库，已存在的配置项会被覆盖
func (r *Reborn) Save() error {
	_, err := r.rdb.HMSet(r.name, r.kvPairs).Result()
	return err
}

// 从数据库中加载配置项
func (r *Reborn) loadFromDB() error {
	rs, err := r.rdb.HGetAll(r.name).Result()
	if err != nil {
		return err
	}
	for k, v := range rs {
		r.kvPairs[k] = v // 这里不用 SetValue 是减少不必要的判断
	}
	return nil
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
