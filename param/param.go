package param

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mostafa-asg/pelican/store"
	"github.com/mostafa-asg/pelican/util"
)

var host string
var httpPort int
var tcpPort int
var enableHttp bool

var defaultExpire string
var defaultExpireStrategy string
var cleanUpInterval string

func Load() {
	flag.StringVar(&host, "host", getDefaultHost(), "Interface to bind")
	flag.IntVar(&httpPort, "http-port", getDefaultHttpPort(), "Http port to bind")
	flag.IntVar(&tcpPort, "tcp-port", getDefaultTcpPort(), "Tcp port to bind")
	flag.BoolVar(&enableHttp, "enable-http", false, "Eanble/Disable http server")

	flag.StringVar(&defaultExpire, "expire", getDefaultExpire(), "Default expiration time.")
	flag.StringVar(&defaultExpireStrategy, "strategy", getDefaultStrategy(), "Default expiration strategy")
	flag.StringVar(&cleanUpInterval, "cleanup", getDefaultCleanUp(), "At which interval should items evict from memory")

	flag.Parse()
}

func HttpBindAddress() string {
	return fmt.Sprintf("%s:%d", host, httpPort)
}

func TcpBindAddress() string {
	return fmt.Sprintf("%s:%d", host, tcpPort)
}

func DefaultExpire() (time.Duration, error) {
	return util.ToTimeDuration(defaultExpire)
}

func HttpEnabled() bool {
	return enableHttp
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
	return util.ToTimeDuration(cleanUpInterval)
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

func getDefaultHttpPort() int {
	port := os.Getenv("PELIKAN_HTTP_PORT")
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

func getDefaultTcpPort() int {
	port := os.Getenv("PELIKAN_TCP_PORT")
	if port == "" {
		port = "4051"
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil {
		panic(fmt.Sprintf("Invalid port number %s", port))
	} else {
		return portNumber
	}
}
