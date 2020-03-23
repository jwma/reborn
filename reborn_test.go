package reborn

import (
	"github.com/go-redis/redis"
	"os"
	"strconv"
	"testing"
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
		t.Errorf("failed to get Reborn instance, %v\n", err)
	}

	// 测试设置配置
	err = reborn.Set("websiteTitle", "Reborn")
	if err != nil {
		t.Errorf("failed to set string value, %v\n", err)
	}
	err = reborn.Set("requestTimeout", 30)
	if err != nil {
		t.Errorf("failed to set int value, %v\n", err)
	}
	err = reborn.Set("discount", 5.5)
	if err != nil {
		t.Errorf("failed to set float64 value, %v\n", err)
	}
	err = reborn.Set("toggle", false)
	if err != nil {
		t.Errorf("failed to set bool value, %v\n", err)
	}

	// 测试获取配置
	websiteTitle := reborn.GetValue("websiteTitle", "")
	if websiteTitle != "Reborn" {
		t.Errorf("websiteTitle expected: %s, got: %s\n", "Reborn", websiteTitle)
	}
}
