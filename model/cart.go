package model

import "gopkg.in/mgo.v2/bson"

// Cart is a cart model
type Cart struct {
	ID        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty" description:"Object ID"`
	CreatedAt int64         `json:"createdAt" bson:"createdAt" description:"Created date."`
	UpdatedAt int64         `json:"updatedAt" bson:"updatedAt" description:"Updated date."`
}
