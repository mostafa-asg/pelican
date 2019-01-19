package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	api "github.com/mostafa-asg/pelican/http"
	"github.com/mostafa-asg/pelican/param"
	"github.com/mostafa-asg/pelican/socket"
	"github.com/mostafa-asg/pelican/store"
	"github.com/mostafa-asg/pelican/util"
	"github.com/prometheus/client_golang/prometheus"
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

	listener, err := net.Listen("tcp", param.TcpBindAddress())
	util.HandleError(err, util.Fatal)
	log.Println(fmt.Sprintf("Listening on [%s]", param.TcpBindAddress()))

	go func() {
		if param.HttpEnabled() {
			r := mux.NewRouter()
			r.HandleFunc("/{key}", api.PutHandler(kvStore)).Methods("put")
			r.HandleFunc("/{key}", api.GetHandler(kvStore)).Methods("get")
			r.HandleFunc("/{key}", api.DelHandler(kvStore)).Methods("delete")
			r.HandleFunc("/counter/inc/{key}/{value}", api.IncCounter(kvStore)).Methods("put")
			r.HandleFunc("/counter/dec/{key}/{value}", api.DecCounter(kvStore)).Methods("put")
			r.Handle("/get/metrics", prometheus.Handler())
			log.Println(fmt.Sprintf("Http server will listen on [%s]", param.HttpBindAddress()))
			http.ListenAndServe(param.HttpBindAddress(), r)
		}
	}()

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
			go socket.HandleConnection(con, kvStore)
		}
	}
}
