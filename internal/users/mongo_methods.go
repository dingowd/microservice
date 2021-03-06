package users

import (
	"encoding/json"
	"fmt"
	"github.com/dingowd/microservice/db/mongo"
	"github.com/dingowd/microservice/structures"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

//type DB structures.DB

type MongoDB struct {
	Collection *mgo.Collection
	Session    *mgo.Session
}

func (d *MongoDB) Init() {
	var err error
	d.Session, err = mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	d.Collection = d.Session.DB("users").C("users")
}

func (d *MongoDB) Close() {
	d.Session.Close()
}

// Function to create user in mongo DB by json
func (d *MongoDB) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
		u := new(structures.InterUser)
		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		udb := &structures.User{ID: bson.NewObjectId(), Name: u.Name, Age: u.Age, Friends: u.Friends}
		err = mongo.Create(*udb, d.Collection)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(w, "User"+"\nName:"+u.Name+
			"\nAge:"+strconv.Itoa(u.Age)+"\nwas created")
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "Method isn`t POST")
}

// Function to delete user in mongo DB by json
func (d *MongoDB) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
		var El structures.UserToDo
		if err := json.Unmarshal(content, &El); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
		if err = mongo.RemoveUser(El.Id, d.Collection); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
		fmt.Fprint(w, "User ", El.Id, " was deleted")
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "Method isn`t POST")
}

// Function to display all users in mongo DB
func (d *MongoDB) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		query := bson.M{}
		users := []structures.User{}
		if err := d.Collection.Find(query).All(&users); err != nil {
			fmt.Fprint(w, err)
			return
		}
		for _, u := range users {
			fmt.Fprint(w, u, "\n")
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "Method isn`t GET")
}

// Function to make friends from user1 and user2 in mongo DB by json
func (d *MongoDB) MakeFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
		//?????????????????? id ?????????????????????????? ???? ??????????????, ?????????????? ???????? ??????????????????
		var Fr structures.Friends
		if err := json.Unmarshal(content, &Fr); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}
		n1 := Fr.Us1
		n2 := Fr.Us2
		var mU1, mU2 structures.User
		//???????????????? ?????????????????????????? ???????????????????????? 1
		if err = mongo.IsExist(Fr.Us1, d.Collection); err != nil {
			fmt.Fprint(w, "User ", Fr.Us1, " ", err)
			return
		}
		query1 := bson.M{
			"_id": n1,
		}
		_ = d.Collection.Find(query1).One(&mU1)
		//???????????????? ?????????????????????????? ???????????????????????? 2
		if err = mongo.IsExist(Fr.Us2, d.Collection); err != nil {
			fmt.Fprint(w, "User ", Fr.Us2, " ", err)
			return
		}
		query2 := bson.M{
			"_id": n2,
		}
		_ = d.Collection.Find(query2).One(&mU2)
		//???????????????? ???? ????, ?????? ???????????????????????? ?????? ????????????
		for _, val := range mU1.Friends {
			if val == n2 {
				fmt.Fprint(w, "User\n", mU1, "\nand User\n", mU2, "\nalready friends")
				return
			}
		}
		//???????????????????? ???????????????????????? 2 ?? ???????????? ?? ???????????????????????? 1
		mU1.Friends = append(mU1.Friends, n2)
		err = d.Collection.Update(bson.M{"_id": n1}, bson.M{"$set": bson.M{"friends": mU1.Friends}})
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		//???????????????????? ???????????????????????? 1 ?? ???????????? ?? ???????????????????????? 2
		mU2.Friends = append(mU2.Friends, n1)
		err = d.Collection.Update(bson.M{"_id": n2}, bson.M{"$set": bson.M{"friends": mU2.Friends}})
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, Fr.Us1, " ", Fr.Us2, " are friends")
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "Method isn`t POST")
}

// Function to display user's friends by json
func (d *MongoDB) GetFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
		var El structures.UserToDo
		if err := json.Unmarshal(content, &El); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}
		//???????????????? ?????????????????????????? ????????????????????????
		if err = mongo.IsExist(El.Id, d.Collection); err != nil {
			fmt.Fprint(w, "User ", El.Id, " ", err)
			return
		}
		var mU structures.User
		query := bson.M{
			"_id": El.Id,
		}
		_ = d.Collection.Find(query).One(&mU)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Friends of User with ID ", El.Id, " are ", mU.Friends)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "Method isn`t POST")
}

// Function to set new user age in mongo DB by json
func (d *MongoDB) NewAge(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
		var El structures.NewAge
		if err := json.Unmarshal(content, &El); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}
		//???????????????? ?????????????????????????? ????????????????????????
		if err = mongo.IsExist(El.Id, d.Collection); err != nil {
			fmt.Fprint(w, "User ", El.Id, " ", err)
			return
		}
		//???????????? ?????????????? ????????????????????????
		err = d.Collection.Update(bson.M{"_id": El.Id}, bson.M{"$set": bson.M{"age": El.NewAge}})
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		fmt.Fprint(w, "Age of User with ID ", El.Id, " is corrected to ", El.NewAge)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "Method isn`t POST")
}
