package users

import "net/http"

type DataBase interface {
	Init()
	Create(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	MakeFriends(w http.ResponseWriter, r *http.Request)
	GetFriends(w http.ResponseWriter, r *http.Request)
	NewAge(w http.ResponseWriter, r *http.Request)
	Close()
}
