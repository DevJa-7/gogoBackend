package model

import "gopkg.in/mgo.v2/bson"

// Food is a restaurant model
type Food struct {
	ID            bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty" description:"Object ID"`
	BusinessID    bson.ObjectId `json:"businessId" bson:"businessId"`
	Business      *Business     `json:"business,omitempty" bson:"business,omitempty"`
	FoodType      *FoodType     `json:"foodType,omitempty" bson:"foodType,omitempty"`
	MealKindCodes []int         `json:"mealKindCodes" bson:"mealKindCodes"` // breakfast, lunch
	MealKinds     []string      `json:"mealKinds,omitempty" bson:"mealKinds,omitempty"`
	DietaryCodes  []int         `json:"dietaryCodes" bson:"dietaryCodes"` // veterian, meat
	Dietaries     []string      `json:"dietaries,omitempty" bson:"dietaries,omitempty"`
	Image         string        `json:"image"`
	Name          string        `json:"name" description:"Restaurant name"`
	Description   string        `json:"description" description:"Restaurant description"`
	Price         float64       `json:"price"`
	// extra variable
	SoldOut      bool  `json:"soldOut" bson:"soldOut"`
	Enabled      bool  `json:"enabled"`
	Recommend    bool  `json:"recommend"`
	MostPopular  bool  `json:"mostPopular" bson:"mostPopular"`
	Status       bool  `json:"status"`
	FreeDelivery bool  `json:"freeDelivery" bson:"freeDelivery"`
	CreatedAt    int64 `json:"createdAt" bson:"createdAt" description:"Created date."`
	UpdatedAt    int64 `json:"updatedAt" bson:"updatedAt" description:"Updated date."`
}

// BusinessFood is for business query
type BusinessFood struct {
	FoodType string  `json:"foodType" bson:"foodType"`
	Foods    []*Food `json:"foods"`
}

// OrderFood is ordered food struct
type OrderFood struct {
	Food        *Food   `json:"food"`
	Price       float64 `json:"price"`
	Instruction string  `json:"instruction"`
}
