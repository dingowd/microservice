package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var count bool = false

const firstservice string = "http://localhost:8080"
const secondservice string = "http://localhost:8082"

func GetAll(w http.ResponseWriter, r *http.Request) {
	url := GetURL("/getall")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	defer resp.Body.Close()
	text, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Fprint(w, err2)
		return
	}
	fmt.Fprint(w, string(text))
}

func post(w http.ResponseWriter, r *http.Request, url string) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	resp, err2 := http.Post(url, "text/plain", bytes.NewBuffer(req))
	if err2 != nil {
		fmt.Fprint(w, err)
		return
	}
	defer resp.Body.Close()
	req, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	if count {
		fmt.Fprint(w, "microservice 2\n", string(req))
		return
	}
	fmt.Fprint(w, "microservice 1\n", string(req))
}

func GetURL(url string) string {
	if count {
		count = !count
		return secondservice + url
	}
	count = !count
	return firstservice + url
}

func Create(w http.ResponseWriter, r *http.Request) {
	post(w, r, GetURL("/create"))
}

func MakeFriends(w http.ResponseWriter, r *http.Request) {
	post(w, r, GetURL("/make_friends"))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	post(w, r, GetURL("/delete"))
}

func GetFriends(w http.ResponseWriter, r *http.Request) {
	post(w, r, GetURL("/friends"))
}

func NewAge(w http.ResponseWriter, r *http.Request) {
	post(w, r, GetURL("/newage"))
}

func main() {
	//initial
	server := &http.Server{Addr: "localhost:9000", Handler: nil}
	//Graceful shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		//Shutdown signal with grace period of 30 seconds
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
	//request processing
	http.HandleFunc("/getall", GetAll)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/make_friends", MakeFriends)
	http.HandleFunc("/delete", DeleteUser)
	http.HandleFunc("/friends", GetFriends)
	http.HandleFunc("/newage", NewAge)
	//start http server
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}
