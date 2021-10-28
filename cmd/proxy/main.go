package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
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
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	go func() {
		for {
			select {
			case <-stopChan:
				fmt.Println("Прокси завершает работу")
				os.Exit(0)
			}
		}
	}()
	http.HandleFunc("/getall", GetAll)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/make_friends", MakeFriends)
	http.HandleFunc("/delete", DeleteUser)
	http.HandleFunc("/friends", GetFriends)
	http.HandleFunc("/newage", NewAge)
	if err := http.ListenAndServe("localhost:9000", nil); err != nil {
		fmt.Println(err.Error())
	}
}
