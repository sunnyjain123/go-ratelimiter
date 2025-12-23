package ratelimiter

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/sunnyjain123/go-ratelimiter/config"
)

func BuildLimiter(cfg config.LimiterConfig) (Limiter, error) {
	switch cfg.Type {

	case "noop":
		return NewNoopLimiter(), nil

	case "fixed_window":
		return NewFixedWindowLimiter(
			cfg.FixedWindow.Limit,
			time.Duration(cfg.FixedWindow.WindowMs)*time.Millisecond,
		), nil

	case "token_bucket":
		return NewTokenBucketLimiter(
			cfg.TokenBucket.Capacity,
			cfg.TokenBucket.RefillRate,
		), nil

	case "leaky_bucket":
		return NewLeakyBucketLimiter(
			cfg.LeakyBucket.Capacity,
			cfg.LeakyBucket.LeakRate,
		), nil

	}

	return nil, fmt.Errorf("unknown limiter type: %s", cfg.Type)
}

func BuildLimiterFromViper(prefix string) (Limiter, error) {
	typ := viper.GetString(prefix + ".type")

	switch typ {

	case "fixed_window":
		return NewFixedWindowLimiter(
			viper.GetInt(prefix+".fixed_window.limit"),
			time.Duration(viper.GetInt(prefix+".fixed_window.window_ms"))*time.Millisecond,
		), nil

	case "token_bucket":
		return NewTokenBucketLimiter(
			viper.GetInt(prefix+".token_bucket.capacity"),
			viper.GetFloat64(prefix+".token_bucket.refill_rate"),
		), nil

	case "leaky_bucket":
		return NewLeakyBucketLimiter(
			viper.GetInt(prefix+".leaky_bucket.capacity"),
			viper.GetInt(prefix+".leaky_bucket.leak_rate"),
		), nil

	}

	return nil, fmt.Errorf("unknown limiter type: %s", typ)
}
