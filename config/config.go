package config

import (
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App         AppConfig         `yaml:"app"`
	RateLimiter RateLimiterConfig `yaml:"rate_limiter"`
}

type AppConfig struct {
	Env string `yaml:"env"`
}

type RateLimiterConfig struct {
	Default   LimiterConfig            `yaml:"default"`
	RPCLimits map[string]LimiterConfig `yaml:"rpc_limits"`
}

type LimiterConfig struct {
	Type        string             `yaml:"type"`
	TokenBucket TokenBucketConfig  `yaml:"token_bucket"`
	FixedWindow FixedWindowConfig  `yaml:"fixed_window"`
	LeakyBucket LeakyBucketConfig  `yaml:"leaky_bucket"`
	Redis       RedisLimiterConfig `yaml:"redis"`
}

type TokenBucketConfig struct {
	Capacity   int     `yaml:"capacity"`
	RefillRate float64 `yaml:"refill_rate"`
}

type FixedWindowConfig struct {
	Limit    int `yaml:"limit"`
	WindowMs int `yaml:"window_ms"`
}

type LeakyBucketConfig struct {
	Capacity int `yaml:"capacity"`
	LeakRate int `yaml:"leak_rate"`
}

type RedisLimiterConfig struct {
	Addr       string  `yaml:"addr"`
	Password   string  `yaml:"password"`
	DB         int     `yaml:"db"`
	Capacity   int     `yaml:"capacity"`
	RefillRate float64 `yaml:"refill_rate"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadViperConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	_ = viper.ReadInConfig()
}
