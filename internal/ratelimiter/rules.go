package ratelimiter

type Rule struct {
	Limit      int
	RefillRate float64
}

type RuleSet map[string]Rule
