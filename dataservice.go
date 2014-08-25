package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func dbUpdateStudy(collection *mgo.Collection, study study) error {

	// Update the database
	e := collection.Update(bson.M{"_id": study.Id},
		bson.M{"studyname": study.StudyName,
			"_id":         study.Id,
			"description": study.Description,
		})
	if e != nil {
		return e
	}

	return nil
}
