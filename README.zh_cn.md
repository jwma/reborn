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

* [安装](#安装)
* [快速开始](#快速开始)
* [更多用法](#更多用法)
    * [获得没有任何配置项的 Reborn 实例](#获得没有任何配置项的-reborn-实例)
    * [一次保存多个配置项](#一次保存多个配置项)
    * [开启自动重载配置](#开启自动重载配置)
* [支持的数据类型](#支持的数据类型)
    * [想要支持其他数据类型？](#想要支持其他数据类型)
* [使用须知](#使用须知)
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
```

### 一次保存多个配置项

```go
r, _ := reborn.New(client, "YOUR_CONFIG_KEY")

// SetValue 只会更新 Reborn 实例的配置项而不会影响数据库的配置项
r.SetValue("discount", 8.8)
r.SetValue("websiteTitle", "Promotion")

// 通过调用 Persist 方法，对当前 Reborn 实例的配置进行保存到数据库的操作
r.Persist()
```

### 开启自动重载配置

开启自动重载后，Reborn 会定时使用数据库的配置去覆盖 Reborn 实例的配置项，但当 Reborn 实例的数据有变动的情况下，自动重载功能会
跳过重载操作，进入下一轮等待，直到你调用 `Persist()` 之后，自动重载将会继续进行重载配置的工作。

```go
r, _ := reborn.New(client, "YOUR_CONFIG_KEY")
r.SetAutoReloadDuration(time.Second * 10)  // 默认间隔: 5s
r.StartAutoReload()

// 停止自动重载配置
r.StopAutoReload()
```

## 支持的数据类型

你可以在使用 `SetValue()` 时，传递如下的数据类型：
- `int`
- `float64`
- `string`
- `bool`
- `[]int`
- `[]string`
- `map[string]int`
- `map[string]string`

_⚠️ 当你 `SetValue()` 传递了不支持的数据类型时，你会得到一个 `UnsupportedValueTypeError`。_

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
如果你的配置项使用的数据类型不在支持列表中，你可以在 `SetValue()` 时，传入已经转换为 `string` 的值，在读取配置项时，
可以使用 `GetValue()` 获取，最后再转换为原本的数据类型。

## 使用须知

**⚠️ 记得调用 `Persist()`。**

当你调用 `SetValue()` 设置了配置项后，不要忘记调用 `Persist()`，忘记调用的后果：

1. 你的配置当然不会被保存到数据库；
2. 如果你开启了自动重载配置的话，忘记调用 `Persist()` 会影响配置的重载，因为 Reborn 认为你有改动未保存，所以不会使用数据库配置对 
Reborn 实例的配置进行覆盖。 
