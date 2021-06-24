package model

import (
	"gopkg.in/mgo.v2/bson"
)

// Model struct
type Model struct {
	Code  int    `json:"code"`
	Value string `json:"value"`
}

// Brand struct
type Brand struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title     string        `json:"title"`
	Models    []*Model      `json:"models"`
	CreatedAt int64         `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64         `json:"updatedAt" bson:"updatedAt"`
}
