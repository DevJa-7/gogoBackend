package model

import "gopkg.in/mgo.v2/bson"

// Dietary struct ex) American, burger...
type Dietary struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Code      int           `json:"code"`
	Name      string        `json:"name"`
	Icon      string        `json:"icon"`
	Image     string        `json:"image"`
	Top       bool          `json:"top"`
	Default   bool          `json:"default"`
	CreatedAt int64         `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64         `json:"updatedAt" bson:"updatedAt"`
}
