package model

import (
	"gopkg.in/mgo.v2/bson"
)

// MealKind indicate meal-time or meal kind	ex) Breakfast, Dinner & Lunch, Chinese meal...
type MealKind struct {
	ID        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty" description:"Object ID"`
	Code      int           `json:"code"`
	Name      string        `json:"name" description:"food type name"`
	Image     string        `json:"image"`
	CreatedAt int64         `json:"createdAt" bson:"createdAt" description:"Created date."`
	UpdatedAt int64         `json:"updatedAt" bson:"updatedAt" description:"Updated date."`
}
