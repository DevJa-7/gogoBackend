package model

import "gopkg.in/mgo.v2/bson"

// Role struct.
type Role struct {
	ID        bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	Name      string          `json:"name"`
	Code      int             `json:"code"`
	URLGroup  []bson.ObjectId `json:"urlGroup" bson:"urlGroup"`
	CreatedAt int64           `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64           `json:"updatedAt" bson:"updatedAt"`
}
