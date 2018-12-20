package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	host := "0.0.0.0"
	port := 8001

	listener, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, os.Kill)

		<-c
		log.Println("Got interrupted ...")
		listener.Close()
		log.Println("Exit the application")
		os.Exit(0)
	}()

	for {
		log.Println("Waiting for new connections ...")
		con, err := listener.Accept()
		if err != nil {
			log.Println("Error in accepting connecton", err)
		} else {
			log.Println("Accept new connection from", con.RemoteAddr())
			go handleConnection(con)
		}

	}
}

func handleConnection(con net.Conn) {
	con.Write([]byte("Hello"))
	con.Close()
}
