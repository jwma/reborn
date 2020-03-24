<h1 align="center">
  <br>Reborn<br>
</h1>

<p align="center"><em>A redis-based configuration library developed using Go, easy to use.</em></p>
<p align="center">
  <a href="https://img.shields.io/github/license/mashape/apistatus.svg" target="_blank">
      <img src="https://img.shields.io/github/license/mashape/apistatus.svg" alt="license">
  </a>
</p>

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

