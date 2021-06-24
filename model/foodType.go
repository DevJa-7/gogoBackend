package model

import (
	"gopkg.in/mgo.v2/bson"
)

// FoodType is a food type	ex) Most Popular, Breakfast Meals ...
type FoodType struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty" description:"Object ID"`
	Code        int           `json:"code"`
	Name        string        `json:"name" description:"food type name"`
	Description string        `json:"description"`
	Image       string        `json:"image"`
	BusinessID  bson.ObjectId `json:"businessId" bson:"businessId"`
	Business    *Business     `json:"business,omitempty" bson:"business,omitempty"`
	FoodOption  []*FoodOption `json:"foodOptions" bson:"foodOptions"`
	CreatedAt   int64         `json:"createdAt" bson:"createdAt" description:"Created date."`
	UpdatedAt   int64         `json:"updatedAt" bson:"updatedAt" description:"Updated date."`
}

// FoodOption is the struct of food extra option
type FoodOption struct {
	Number     string `json:"number"`
	Title      string `json:"title"`
	IsRequired bool   `json:"isRequired" bson:"isRequired"`
	IsSingle   bool   `json:"isSingle" bson:"isSingle"` // false:multi selection, true:single option
	Enabled    bool   `json:"enabled"`
	Options    []*struct {
		Name    string  `json:"name"`
		Price   float64 `json:"price"`
		Enabled bool    `json:"enabled"`
		Default bool    `json:"default"`
	} `json:"options"`
}
