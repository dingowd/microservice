package main

import (
	"fmt"
	"github.com/dingowd/microservice/internal/users"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	go func() {
		for {
			select {
			case <-stopChan:
				fmt.Println("Микросервис 2 завершает работу")
				os.Exit(0)
			}
		}
	}()
	u := new(users.DB)
	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	u.Collection = session.DB("users").C("users")
	http.HandleFunc("/getall", u.GetAll)
	http.HandleFunc("/create", u.Create)
	http.HandleFunc("/make_friends", u.MakeFriends)
	http.HandleFunc("/delete", u.DeleteUser)
	http.HandleFunc("/friends", u.GetFriends)
	http.HandleFunc("/newage", u.NewAge)
	if err := http.ListenAndServe("localhost:8082", nil); err != nil {
		fmt.Println(err.Error())
	}
}
