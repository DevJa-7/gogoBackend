package model

import "gopkg.in/mgo.v2/bson"

// Document struct
type Document struct {
	ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Type       int           `json:"type"`
	Name       string        `json:"name"`
	IsExpired  bool          `json:"isExpired" bson:"isExpired"`
	IsRequired bool          `json:"isRequired" bson:"isRequired"`
	Valid      bool          `json:"valid"`
	CreatedAt  int64         `json:"createdAt" bson:"createdAt"`
	UpdatedAt  int64         `json:"updatedAt" bson:"updatedAt"`
}
