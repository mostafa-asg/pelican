package param

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mostafa-asg/pelican/store"
)

var host string
var port int

var defaultExpire string
var defaultExpireStrategy string
var cleanUpInterval string

func Load() {
	flag.StringVar(&host, "host", getDefaultHost(), "Interface to bind")
	flag.IntVar(&port, "port", getDefaultPort(), "Port to bind")

	flag.StringVar(&defaultExpire, "expire", getDefaultExpire(), "Default expiration time.")
	flag.StringVar(&defaultExpireStrategy, "strategy", getDefaultStrategy(), "Default expiration strategy")
	flag.StringVar(&cleanUpInterval, "cleanup", getDefaultCleanUp(), "At which interval should items evict from memory")

	flag.Parse()
}

func BindAddresss() string {
	return fmt.Sprintf("%s:%d", host, port)
}

func DefaultExpire() (time.Duration, error) {
	return toTimeDuration(defaultExpire)
}

func DefaultStrategy() (store.Strategy, error) {
	switch defaultExpireStrategy {
	case "sliding":
		return store.Sliding, nil
	case "absolute":
		return store.Absolute, nil
	default:
		return -1, errors.New("Invalid strategy")
	}
}

func CleanUpInterval() (time.Duration, error) {
	return toTimeDuration(cleanUpInterval)
}

func toTimeDuration(s string) (time.Duration, error) {
	number, err := strconv.Atoi(s[:len(s)-1])
	if err != nil {
		return 0, err
	}

	duration := s[len(s)-1:]
	switch duration {
	case "s":
		return time.Duration(number) * time.Second, nil
	case "m":
		return time.Duration(number) * time.Minute, nil
	case "h":
		return time.Duration(number) * time.Hour, nil
	default:
		return 0, errors.New("invalid time duration")
	}

}

func getDefaultHost() string {
	host := os.Getenv("PELIKAN_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	return host
}

func getDefaultCleanUp() string {
	expire := os.Getenv("PELIKAN_CLEAN_UP")
	if expire == "" {
		expire = "30m" // 30 minutes
	}

	return expire
}

func getDefaultExpire() string {
	expire := os.Getenv("PELIKAN_ITEM_EXPIRE")
	if expire == "" {
		expire = "30m" // 30 minutes
	}

	return expire
}

func getDefaultStrategy() string {
	strategy := os.Getenv("PELIKAN_EXPIRE_STRATEGY")
	if strategy == "" {
		strategy = "sliding"
	}

	return strategy
}

func getDefaultPort() int {
	port := os.Getenv("PELIKAN_PORT")
	if port == "" {
		port = "4050"
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil {
		panic(fmt.Sprintf("Invalid port number %s", port))
	} else {
		return portNumber
	}
}
