<h1 align="center">
  <br>Reborn<br>
</h1>

<p align="center"><em>Reborn 是使用 Go 开发的，基于 Redis 存储的配置库，简单配置，易于使用。</em></p>
<p align="center">
  <a href="https://github.com/jwma/reborn/workflows/Go/badge.svg?branch=master" target="_blank">
    <img src="https://github.com/jwma/reborn/workflows/Go/badge.svg?branch=master" alt="ci">
  </a>
  <a href="https://img.shields.io/github/license/mashape/apistatus.svg" target="_blank">
      <img src="https://img.shields.io/github/license/mashape/apistatus.svg" alt="license">
  </a>
</p>

---

[中文](README.zh_cn.md "中文") | [English](README.md "English")

---

## 安装
```console
go get github.com/jwma/reborn
```

## 快速开始

```go
package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jwma/reborn"
	"log"
)

func main() {
	// 获得一个 redis client
	client := redis.NewClient(&redis.Options{
		Addr:        "127.0.0.1:6379",
		Password:    "",
		DB:          0,
		IdleTimeout: -1,
	})

	var err error

	// 获得一个空的配置，将会作为 Reborn 实例的默认配置
	defaults := reborn.NewConfig()

	// 设置配置项
	err = defaults.SetValue("websiteTitle", "Reborn")
	if err != nil {
		log.Printf("failed to set websiteTitle, error: %v\n", err)
	}

	err = defaults.SetValue("toggle", false)
	if err != nil {
		log.Printf("failed to set toggle, error: %v\n", err)
	}

	// 通过默认配置获得一个 Reborn 实例，该实例的配置会保存在 Redis 名为 YOUR_CONFIG_KEY 的 Hash 中
	r, err := reborn.NewWithDefaults(client, "YOUR_CONFIG_KEY", defaults)

	// 通过key 获取配置，第二个参数为获取的 key 不存的时候所使用的默认值
	websiteTitle := r.GetValue("websiteTitle", "默认值")
	fmt.Println(websiteTitle)  // 输出：Reborn
	fmt.Println(r.GetValue("noExistsKey", "oops")) // 输出：oops

	toggle := r.GetBoolValue("toggle", false)
	fmt.Println(toggle)
}
```

## 更多用法

除了可以使用默认配置项获得 Reborn 实例外，还可以有如下用法。

### 获得没有任何配置项的 Reborn 实例

```go
r, _ := reborn.New(client, "YOUR_CONFIG_KEY")

// 可以在程序运行时，根据需要设置配置项，此操作会同时更新 Reborn 实例以及数据库的配置项
// 每一次 Set 的调用，都会请求一次数据库
r.Set("websiteTitle", "Reborn")
r.Set("requestTimeout", 30)
```

### 一次数据库请求设置多个配置

```go
r, _ := reborn.New(client, "YOUR_CONFIG_KEY")

// SetValue 只会更新 Reborn 实例的配置项而不会影响数据库的配置项
r.SetValue("discount", 8.8)
r.SetValue("websiteTitle", "Promotion")

// 通过调用 Persist 方法，对当前 Reborn 实例的配置进行保存到数据库的操作
r.Persist()
```

## 支持的数据类型

你可以在使用 `Set()` 或 `SetValue()` 时，传递如下的数据类型：
- `int`
- `float64`
- `string`
- `bool`
- `[]int`
- `[]string`
- `map[string]int`
- `map[string]string`

_⚠️ 当你 `Set()` 或 `SetValue()` 传递了不支持的数据类型时，你会得到一个 `UnsupportedValueTypeError`。_

你可以通过如下的方法获取不同类型的配置项：
- `GetValue()` 获取 `string`
- `GetIntValue()` 获取 `int`
- `GetFloat64Value()` 获取 `float64`
- `GetBoolValue()` 获取  `bool`
- `GetIntSliceValue()` 获取 `[]int`
- `GetStringSliceValue()` 获取 `[]string`
- `GetStringIntMapValue()` 获取 `map[string]int`
- `GetStringStringMapValue()` 获取 `map[string]string`

### 想要支持其他数据类型？
如果你的配置项使用的数据类型不在支持列表中，你可以在 `Set()` 或 `SetValue()` 时，传入已经转换为 `string` 的值，在读取配置项时，
可以使用 `GetValue()` 获取，最后再转换为原本的数据类型。
