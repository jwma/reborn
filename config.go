package reborn

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
)

type Config struct {
	kvPairs   sync.Map
	hasChange bool
	mux       sync.Mutex
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) valueToString(value interface{}) (string, error) {
	var v string
	switch value.(type) {
	case string:
		v = value.(string)
	case int:
		v = strconv.Itoa(value.(int))
		break
	case float64:
		v = fmt.Sprintf("%f", value.(float64))
		break
	case bool:
		v = fmt.Sprintf("%v", value.(bool))
		break
	case []int:
		j, _ := json.Marshal(value.([]int))
		v = string(j)
		break
	case []string:
		j, _ := json.Marshal(value.([]string))
		v = string(j)
		break
	case map[string]int:
		j, _ := json.Marshal(value.(map[string]int))
		v = string(j)
		break
	case map[string]string:
		j, _ := json.Marshal(value.(map[string]string))
		v = string(j)
		break
	default:
		return "", UnsupportedValueTypeError{value}
	}
	return v, nil
}

func (c *Config) set(key string, value string) {
	c.mux.Lock()
	c.hasChange = true
	c.mux.Unlock()
	c.kvPairs.Store(key, value)
}

func (c *Config) SetValue(key string, value interface{}) error {
	v, err := c.valueToString(value)
	if err != nil {
		return err
	}
	c.set(key, v)
	return nil
}

func (c *Config) GetValue(key string, defaults string) string {
	if value, ok := c.kvPairs.Load(key); ok {
		return value.(string)
	}
	return defaults
}

func (c *Config) GetIntValue(key string, defaults int) int {
	vStr := c.GetValue(key, "")
	if vStr == "" {
		return defaults
	}
	v, err := strconv.Atoi(vStr)
	if err != nil {
		return defaults
	}
	return v
}

// if key does not exists, return v
// if key exists and the value < v, return value
// if key exists and the value >= v, return v
func (c *Config) GetIntValueLT(key string, v int) int {
	value := c.GetIntValue(key, v)

	if value < v {
		return value
	}

	return v
}

// if key does not exists, return v
// if key exists and the value <= v, return value
// if key exists and the value > v, return v
func (c *Config) GetIntValueLTE(key string, v int) int {
	value := c.GetIntValue(key, v)

	if value <= v {
		return value
	}

	return v
}

// if key does not exists, return v
// if key exists and the value > v, return value
// if key exists and the value <= v, return v
func (c *Config) GetIntValueGT(key string, v int) int {
	value := c.GetIntValue(key, v)

	if value > v {
		return value
	}

	return v
}

// if key does not exists, return v
// if key exists and the value >= v, return value
// if key exists and the value < v, return v
func (c *Config) GetIntValueGTE(key string, v int) int {
	value := c.GetIntValue(key, v)

	if value >= v {
		return value
	}

	return v
}

func (c *Config) GetFloat64Value(key string, defaults float64) float64 {
	vStr := c.GetValue(key, "")
	if vStr == "" {
		return defaults
	}
	v, err := strconv.ParseFloat(vStr, 64)
	if err != nil {
		return defaults
	}
	return v
}

func (c *Config) GetBoolValue(key string, defaults bool) bool {
	vStr := c.GetValue(key, "")
	if vStr == "" {
		return defaults
	}
	v, err := strconv.ParseBool(vStr)
	if err != nil {
		return defaults
	}
	return v
}

func (c *Config) GetIntSliceValue(key string, defaults []int) []int {
	vStr := c.GetValue(key, "")
	if vStr == "" {
		return defaults
	}
	r := make([]int, 0)
	err := json.Unmarshal([]byte(vStr), &r)
	if err != nil {
		return defaults
	}
	return r
}

func (c *Config) GetStringSliceValue(key string, defaults []string) []string {
	vStr := c.GetValue(key, "")
	if vStr == "" {
		return defaults
	}
	r := make([]string, 0)
	err := json.Unmarshal([]byte(vStr), &r)
	if err != nil {
		return defaults
	}
	return r
}

func (c *Config) GetStringIntMapValue(key string, defaults map[string]int) map[string]int {
	vStr := c.GetValue(key, "")
	if vStr == "" {
		return defaults
	}
	r := make(map[string]int)
	err := json.Unmarshal([]byte(vStr), &r)
	if err != nil {
		return defaults
	}
	return r
}

func (c *Config) GetStringStringMapValue(key string, defaults map[string]string) map[string]string {
	vStr := c.GetValue(key, "")
	if vStr == "" {
		return defaults
	}
	r := make(map[string]string)
	err := json.Unmarshal([]byte(vStr), &r)
	if err != nil {
		return defaults
	}
	return r
}
