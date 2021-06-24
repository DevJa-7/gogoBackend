package model

import "gopkg.in/mgo.v2/bson"

// Ads model
type Ads struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty" description:"Object ID"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Image       string        `json:"image"`
	Status      bool          `json:"status"`
	CreatedAt   int64         `json:"createdAt" bson:"createdAt" description:"Created date"`
	UpdatedAt   int64         `json:"updatedAt" bson:"updatedAt" description:"Updated date. This field will be updated when any update operation will be occured"`
}
