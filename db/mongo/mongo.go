package mongo

import (
	"github.com/dingowd/microservice/structures"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type toCheck struct {
	ID bson.ObjectId `bson:"_id"`
}

//создание нового пользователя в коллекции col
func Create(u structures.User, col *mgo.Collection) error {
	err := col.Insert(u)
	return err
}

//функция для проверки существования объекта с _id n в коллекции col
func IsExist(n bson.ObjectId, col *mgo.Collection) error {
	var mU toCheck
	query := bson.M{
		"_id": n,
	}
	err := col.Find(query).One(&mU)
	return err
}

//удаление пользователя с _id n в коллекции col
func RemoveUser(n bson.ObjectId, col *mgo.Collection) error {
	err := col.RemoveId(n)
	return err
}
