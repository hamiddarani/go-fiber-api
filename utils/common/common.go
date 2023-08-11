package common

import (
	"encoding/json"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ExpireTime time.Duration
	Digits     int
	Limiter    time.Duration
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func GenerateOtp(cfg *Config) string {
	rand.Seed(time.Now().UnixNano())
	min := int(math.Pow(10, float64(cfg.Digits-1)))
	max := int(math.Pow(10, float64(cfg.Digits)) - 1)

	var num = rand.Intn(max-min) + min
	return strconv.Itoa(num)
}

func TypeConverter[T any](data any) (*T, error) {
	var result T
	dataJson, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataJson, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
