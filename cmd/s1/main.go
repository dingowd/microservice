package main

import (
	"context"
	"fmt"
	"github.com/dingowd/microservice/internal/users"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//Initial
	args := os.Args[1:]
	server := &http.Server{Addr: args[0], Handler: nil}
	//Graceful shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()
		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()
	//Connect to mongodb
	u := new(users.DB)
	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	u.Collection = session.DB("users").C("users")
	//request processing
	http.HandleFunc("/getall", u.GetAll)
	http.HandleFunc("/create", u.Create)
	http.HandleFunc("/make_friends", u.MakeFriends)
	http.HandleFunc("/delete", u.DeleteUser)
	http.HandleFunc("/friends", u.GetFriends)
	http.HandleFunc("/newage", u.NewAge)
	//start http server
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}
