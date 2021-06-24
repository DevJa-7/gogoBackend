package model

import "gopkg.in/mgo.v2/bson"

// Reason struct for decline or cancel
type Reason struct {
	ID        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty" description:"Object ID"`
	Code      int           `json:"code"`
	Type      int           `json:"type"`
	Message   string        `json:"message"`
	Status    bool          `json:"status"`
	CreatedAt int64         `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64         `json:"updatedAt" bson:"updatedAt"`
}
