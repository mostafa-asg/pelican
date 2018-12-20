package param

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var host string
var port int

func Load() {
	flag.StringVar(&host, "host", getDefaultHost(), "Interface to bind")
	flag.IntVar(&port, "port", getDefaultPort(), "Port to bind")
	flag.Parse()
}

func BindAddresss() string {
	return fmt.Sprintf("%s:%d", host, port)
}

func getDefaultHost() string {
	host := os.Getenv("PELIKAN_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	return host
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
