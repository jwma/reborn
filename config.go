package reborn

import (
	"fmt"
	"strconv"
)

type kvPairs map[string]interface{}

type Config struct {
	kvPairs
}

func NewConfig() *Config {
	return &Config{kvPairs{}}
}

func (c Config) valueToString(value interface{}) (string, error) {
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
	default:
		return "", UnsupportedValueTypeError{value}
	}
	return v, nil
}

func (c Config) SetValue(key string, value interface{}) error {
	v, err := c.valueToString(value)
	if err != nil {
		return err
	}
	c.kvPairs[key] = v
	return nil
}

func (c Config) GetValue(key string, defaults string) string {
	if value, ok := c.kvPairs[key]; ok {
		return value.(string)
	}
	return defaults
}

func (c Config) GetIntValue(key string, defaults int) int {
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

func (c Config) GetFloat64Value(key string, defaults float64) float64 {
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

func (c Config) GetBoolValue(key string, defaults bool) bool {
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
