package service

import (
	"errors"
	"strconv"
	"strings"
)

func ValidateCoin(coin string) (string, error) {
	coin = strings.ToLower(strings.TrimSpace(coin))
	if coin == "" {
		return "", errors.New("coin param is required")
	}

	for _, r := range coin {
		if r < 'a' || r > 'z' {
			return "", errors.New("invalid coin name")
		}
	}
	return coin, nil
}

func ValidateLimit(raw string, def, min, max int) int {
	n, err := strconv.Atoi(raw)
	if err != nil || n < min {
		return def
	}
	if n > max {
		return max
	}
	return n
}
