package reborn

import (
	"github.com/go-redis/redis"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

var client *redis.Client
var defaultDbIdx = os.Getenv("REDIS_DB")
var redisHost = os.Getenv("REDIS_HOST")
var redisPassword = os.Getenv("REDIS_PASSWORD")

func getRedisClient() *redis.Client {
	dbIdx, err := strconv.Atoi(defaultDbIdx)
	if err != nil {
		panic("please set REDIS_DB env")
	}

	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:        redisHost,
			Password:    redisPassword,
			DB:          dbIdx,
			IdleTimeout: -1,
		})
	}
	return client
}

func init() {
	c := getRedisClient()
	pong := c.Ping()
	if pong.Err() != nil {
		panic(pong.Err())
	}
}

func TestBasics(t *testing.T) {
	// 准备工作，先清理可能存在的数据
	configName := "config"
	client.Del(configName)

	reborn, err := New(getRedisClient(), configName)
	if err != nil {
		t.Errorf("failed to get Reborn instance, error: %v\n", err)
		return
	}

	// 使用 SetValue() 设置的值只会保存在当前实例中，还不会保存到数据库
	err = reborn.SetValue("websiteTitle", "Reborn")
	if err != nil {
		t.Errorf("failed to set websiteTitle, error: %v\n", err)
	}

	err = reborn.SetValue("toggle", false)
	if err != nil {
		t.Errorf("failed to set toggle, error: %v\n", err)
	}

	err = reborn.Persist()
	if err != nil {
		t.Errorf("failed to persist config to DB, error: %v\n", err)
	}
}

func TestDefaultsConfig(t *testing.T) {
	// 准备工作，先清理可能存在的数据
	configName := "config"
	client.Del(configName)

	var err error
	defaults := NewConfig()

	err = defaults.SetValue("websiteTitle", "Reborn")
	if err != nil {
		t.Errorf("failed to set websiteTitle, error: %v\n", err)
	}

	err = defaults.SetValue("toggle", false)
	if err != nil {
		t.Errorf("failed to set toggle, error: %v\n", err)
	}

	reborn, err := NewWithDefaults(getRedisClient(), configName, defaults)
	if err != nil {
		t.Errorf("failed to get Reborn instance, error: %v\n", err)
		return
	}
	err = reborn.Persist()
	if err != nil {
		t.Errorf("failed to persist config to DB, error: %v\n", err)
	}
}

func TestCompositeTypes(t *testing.T) {
	// 准备工作，先清理可能存在的数据
	configName := "config"
	client.Del(configName)

	reborn, err := New(getRedisClient(), configName)
	if err != nil {
		t.Errorf("failed to get Reborn instance, error: %v\n", err)
		return
	}

	intSlice := []int{1, 2, 3, 4, 5}
	err = reborn.SetValue("intSlice", intSlice)
	if err != nil {
		t.Errorf("failed to set []int value, error: %v\n", err)
	}
	if i := reborn.GetIntSliceValue("intSlice", make([]int, 0)); !reflect.DeepEqual(intSlice, i) {
		t.Errorf("intSlice expected: %v, got: %v\n", intSlice, i)
	}

	stringSlice := []string{"anmuji.com", "Reborn"}
	err = reborn.SetValue("stringSlice", stringSlice)
	if err != nil {
		t.Errorf("failed to set []string value, error: %v\n", err)
	}
	if j := reborn.GetStringSliceValue("stringSlice", make([]string, 0)); !reflect.DeepEqual(stringSlice, j) {
		t.Errorf("stringSlice expected: %v, got: %v\n", stringSlice, j)
	}

	stringIntMap := map[string]int{"mj": 1, "anmuji": 2}
	err = reborn.SetValue("stringIntMap", stringIntMap)
	if err != nil {
		t.Errorf("failed to set map[string]int value, error: %v\n", err)
	}
	if k := reborn.GetStringIntMapValue("stringIntMap", make(map[string]int)); !reflect.DeepEqual(stringIntMap, k) {
		t.Errorf("stringIntMap expected: %v, got: %v\n", stringIntMap, k)
	}

	stringStringMap := map[string]string{"mj": "MJ.MA", "anmuji": "安木鸡"}
	err = reborn.SetValue("stringStringMap", stringStringMap)
	if err != nil {
		t.Errorf("failed to set map[string]string value, error: %v\n", err)
	}
	if k := reborn.GetStringStringMapValue("stringStringMap", make(map[string]string)); !reflect.DeepEqual(stringStringMap, k) {
		t.Errorf("stringStringMap expected: %v, got: %v\n", stringStringMap, k)
	}

	err = reborn.Persist()
	if err != nil {
		t.Errorf("failed to persist config to DB, error: %v\n", err)
	}
}

func TestAutoReload(t *testing.T) {
	// 准备工作，先清理可能存在的数据
	configName := "config"
	client.Del(configName)

	reborn, err := New(getRedisClient(), configName)
	if err != nil {
		t.Errorf("failed to get Reborn instance, error: %v\n", err)
		return
	}

	key := "autoload"
	expected := "ok"
	reborn.SetAutoReloadDuration(time.Millisecond)
	reborn.StartAutoReload()

	reborn.SetValue("k1", "v1")
	time.AfterFunc(time.Millisecond*2, func() {
		reborn.Persist()
	})

	time.AfterFunc(time.Millisecond, func() {
		client.HSet(configName, key, expected) // 模拟数据库配置被修改
	})

	time.AfterFunc(time.Millisecond*7, func() {
		reborn.StopAuthReload()
		v := reborn.GetValue(key, "")
		if v != expected {
			t.Errorf("v expected: %s, got: %s\n", expected, v)
		}
	})

	time.Sleep(time.Millisecond * 8)
}
