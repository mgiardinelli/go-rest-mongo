package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func dbGetStudy(collection *mgo.Collection, id string) (interface{}, error) {
	var study study
	e := collection.FindId(bson.ObjectIdHex(id)).One(&study)

	if e != nil {
		return nil, e
	}

	return study, nil
}
