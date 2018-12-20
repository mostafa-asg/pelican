package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"

	bytesUtil "github.com/mostafa-asg/pelican/bytes/util"
	"github.com/mostafa-asg/pelican/parser"
)

var kvStore sync.Map

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
	reader := bufio.NewReader(con)

	for {
		bytes, err := bytesUtil.ReadBytes(4, reader)
		if err != nil {
			log.Println("Error", err)
			break
		}

		commandSize := bytesUtil.ToUint32(bytes)
		if commandSize == 0 {
			break
		}

		command, err := bytesUtil.ReadBytes(int(commandSize), reader)
		if err != nil {
			log.Println("Error", err)
			break
		}

		if string(command[0:3]) == "PUT" {
			key, value, err := parser.ParsePut(command)
			if err != nil {
				log.Println("Error", err)
				break
			}
			kvStore.Store(key, value)

		} else if string(command[0:3]) == "GET" {
			key, err := parser.ParseGet(command)
			if err != nil {
				log.Println("Error", err)
				break
			}
			value, found := kvStore.Load(key)
			if found {
				con.Write(value.([]byte))
			}
		} else if string(command[0:3]) == "DEL" {
			key, err := parser.ParseDel(command)
			if err != nil {
				log.Println("Error", err)
				break
			}
			kvStore.Delete(key)
		} else {
			log.Println("Unknown command")
			break
		}
	}

	con.Close()
}