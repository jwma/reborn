<h1 align="center">
  <br>Reborn<br>
</h1>

<p align="center"><em>A redis-based configuration library developed using Go, ready to go, easy to use.</em></p>
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

## Install
```console
go get github.com/jwma/reborn
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jwma/reborn"
	"log"
)

func main() {
	// get your redis client here
	client := redis.NewClient(&redis.Options{
		Addr:        "127.0.0.1:6379",
		Password:    "",
		DB:          0,
		IdleTimeout: -1,
	})

	var err error

	// new empty config
	defaults := reborn.NewConfig()

	// set default config
	err = defaults.SetValue("websiteTitle", "Reborn")
	if err != nil {
		log.Printf("failed to set websiteTitle, error: %v\n", err)
	}

	err = defaults.SetValue("toggle", false)
	if err != nil {
		log.Printf("failed to set toggle, error: %v\n", err)
	}

	// create a Reborn instance using defaults config
	// this instance will store all the configs into YOUR_CONFIG_KEY Hash in Redis
	r, err := reborn.NewWithDefaults(client, "YOUR_CONFIG_KEY", defaults)

	// by passing the second parameter, 
	// you can get this default value when the key does not exists
	websiteTitle := r.GetValue("websiteTitle", "default value")
	fmt.Println(websiteTitle)                      // Output: Reborn
	fmt.Println(r.GetValue("noExistsKey", "oops")) // Output: oops

	toggle := r.GetBoolValue("toggle", false)
	fmt.Println(toggle)
}

```

## More usage

You can use Reborn like below.

### Create a Reborn instance without default configs

```go
r, _ := reborn.New(client, "YOUR_CONFIG_KEY")

// Set() can set the Key-Value in Reborn instance and Redis at the same time.
// each time the Set() is called, the Redis is requested once.
r.Set("websiteTitle", "Reborn")
r.Set("requestTimeout", 30)
```

### Save multiple Key-Value at one time

```go
r, _ := reborn.New(client, "YOUR_CONFIG_KEY")

// SetValue() only update the Reborn instance configs
r.SetValue("discount", 8.8)
r.SetValue("websiteTitle", "Promotion")

// after SetValue(), you can call Persist() save you Key-Value into Redis.
r.Persist()
```

## Supported types

Here are the types you can pass when you call `Set()` or `SetValue()`:
- `int`
- `float64`
- `string`
- `bool`
- `[]int`
- `[]string`
- `map[string]int`
- `map[string]string`

_⚠️ You will get the `UnsupportedValueTypeError` when you pass unsupported types._

Call the below functions to get different types value:
- `GetValue()` return `string`
- `GetIntValue()` return `int`
- `GetFloat64Value()` return `float64`
- `GetBoolValue()` return  `bool`
- `GetIntSliceValue()` return `[]int`
- `GetStringSliceValue()` return `[]string`
- `GetStringIntMapValue()` return `map[string]int`
- `GetStringStringMapValue()` return `map[string]string`

### How to support the other types?
You can parse the variable to `string`, then you can call `Set()` or `SetValue()`. When you want to get this config item,
you can call `GetValue()` to get `string` types variable, then you parse it back.

Not cool, but it works~