package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/mostafa-asg/pelican/param"
	"github.com/mostafa-asg/pelican/parser"
	"github.com/mostafa-asg/pelican/store"
	"github.com/mostafa-asg/pelican/util"
	bytesUtil "github.com/mostafa-asg/pelican/util/bytes"
)

func main() {
	param.Load()

	defaultExpire, err := param.DefaultExpire()
	util.HandleError(err, util.Fatal)

	expireStrategy, err := param.DefaultStrategy()
	util.HandleError(err, util.Fatal)

	cleanUpInterval, err := param.CleanUpInterval()
	util.HandleError(err, util.Fatal)

	kvStore := store.New(defaultExpire, expireStrategy, cleanUpInterval)

	listener, err := net.Listen("tcp", param.BindAddresss())
	util.HandleError(err, util.Fatal)
	log.Println(fmt.Sprintf("Pelikan is up, listening on [%s]", param.BindAddresss()))

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
			go handleConnection(con, kvStore)
		}
	}
}

func handleConnection(con net.Conn, kvStore *store.Store) {
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
			kvStore.Put(key, value)

		} else if string(command[0:3]) == "GET" {
			key, err := parser.ParseGet(command)
			if err != nil {
				log.Println("Error", err)
				break
			}
			value, found := kvStore.Get(key)
			if found {
				con.Write(bytesUtil.ToBytes(uint32(len(value.([]byte)))))
				con.Write(value.([]byte))
			} else {
				con.Write(make([]byte, 4))
			}
		} else if string(command[0:3]) == "DEL" {
			key, err := parser.ParseDel(command)
			if err != nil {
				log.Println("Error", err)
				break
			}
			kvStore.Del(key)
		} else {
			log.Println("Unknown command")
			break
		}
	}

	con.Close()
}
