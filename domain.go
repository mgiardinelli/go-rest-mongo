package main

import (
	"gopkg.in/mgo.v2/bson"
)

// study model
type study struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	StudyName   string        `json:"studyname"`
	Description string        `json:"description"`
	Levels      []level       `json:"levels"`
	CreatedBy   string        `json:"createdby"`
}

// level model
type level struct {
	LevelOrder int     `json:"levelorder"`
	Values     []value `json:"values"`
	LevelName  string  `json:"levelname"`
	CreatedBy  string  `json:"createdby"`
}

// value model
type value struct {
	ValueOrder int    `json:"valueorder"`
	ValueName  string `json:"valuename"`
	CreatedBy  string `json:"createdby"`
}
