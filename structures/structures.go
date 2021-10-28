package structures

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID      bson.ObjectId   `bson:"_id"`
	Name    string          `bson:"name"`
	Age     int             `bson:"age"`
	Friends []bson.ObjectId `bson:"friends"`
}

type InterUser struct {
	Name    string
	Age     int
	Friends []bson.ObjectId
}

type Friends struct {
	Us1 bson.ObjectId `json:"us1"`
	Us2 bson.ObjectId `json:"us2"`
}

type UserToDo struct {
	Id bson.ObjectId `json:"target_id"`
}

type DB struct {
	Collection *mgo.Collection
}

type NewAge struct {
	Id     bson.ObjectId `json:"id"`
	NewAge int           `json:"age"`
}
