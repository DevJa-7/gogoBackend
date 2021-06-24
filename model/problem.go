package model

import (
	"gopkg.in/mgo.v2/bson"
)

// Problem struct
type Problem struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty" description:"Object ID"`
	BusinessID  bson.ObjectId `json:"businessId" bson:"businessId"`
	Business    *Business     `json:"business,omitempty" bson:"business,omitempty"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      Status        `json:"status"` // 0:pending, 1:resolved, 2:in process
	Answer      []*struct {
		Description string `json:"description"`
	} `json:"answer"`
	CreatedAt int64 `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64 `json:"updatedAt" bson:"updatedAt"`
}

type PublicProblem struct {
	*Problem
}
