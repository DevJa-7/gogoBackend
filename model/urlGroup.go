package model

import "gopkg.in/mgo.v2/bson"

// URLGroup struct.
type URLGroup struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name      string        `json:"name"`
	URL       string        `json:"url"`
	CreatedAt int64         `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64         `json:"updatedAt" bson:"updatedAt"`
}
