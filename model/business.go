package model

import "gopkg.in/mgo.v2/bson"

// BankInfo struct for bank
type BankInfo struct {
	Name    string `json:"name"`
	Account string `json:"account"`
	Number  string `json:"number"`
}

// Schedule is for open, close of restaurant
type Schedule struct {
	Name      string `json:"name"`
	Weekday   int    `json:"weekday"`
	OpenTime  int64  `json:"openTime"`
	CloseTime int64  `json:"closeTime"`
	Enabled   bool   `json:"enabled"`
}

// Business struct.
type Business struct {
	ID                bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty" description:"Object ID"`
	Email             string        `json:"email" description:"Business email address"`
	Password          string        `json:"password"`
	ParentID          bson.ObjectId `json:"parentId,omitempty" bson:"parentId,omitempty"`
	Parent            *Business     `json:"parent,omitempty" bson:"parent,omitempty"`
	SubBusiness       []*Business   `json:"subBusiness" bson:"subBusiness"`
	CountryCode       string        `json:"countryCode" bson:"countryCode"`
	PhoneCode         string        `json:"phoneCode" bson:"phoneCode"`
	Phone             string        `json:"phone"`
	Logo              string        `json:"logo"`
	Identification    string        `json:"identification"`
	Name              string        `json:"name" description:"Restaurant name"`
	Description       string        `json:"description" description:"Restaurant description"`
	Website           string        `json:"website"`
	BankInfo          BankInfo      `json:"bankInfo" bson:"bankInfo"`
	GeoLocation       GeoLocation   `json:"geoLocation" bson:"geoLocation"`
	PriceLevel        int           `json:"priceLevel" bson:"priceLevel"`
	DietaryCodes      []int         `json:"dietaryCodes" bson:"dietaryCodes"`
	Dietaries         []string      `json:"dietaries,omitempty" bson:"dietaries,omitempty"`
	MealKindCodes     []int         `json:"mealKindCodes" bson:"mealKindCodes"`
	MealKinds         []string      `json:"mealKinds,omitempty" bson:"mealKinds,omitempty"`
	Schedules         []*Schedule   `json:"schedules"`
	PreparationTime   int           `json:"preparationTime" bson:"preparationTime"`
	Closed            bool          `json:"closed"`
	Recommend         bool          `json:"recommend"`
	MostPopular       bool          `json:"mostPopular" bson:"mostPopular"`
	Verify            *Verify       `json:"verify,omitempty" bson:"verify,omitempty"`
	Updated           bool          `json:"isUpdated"`
	Distance          float64       `json:"distance,omitempty" bson:"distance,omitempty"`
	OneSignalPlayerID string        `json:"onesignalPlayerId,omitempty" bson:"onesignalPlayerId,omitempty"`
	CreatedAt         int64         `json:"createdAt" bson:"createdAt" description:"Created date."`
	UpdatedAt         int64         `json:"updatedAt" bson:"updatedAt" description:"Updated date."`
}

// PublicBusiness struct.
type PublicBusiness struct {
	ID              bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty" description:"Object ID"`
	Email           string        `json:"email" description:"Business email address"`
	Phone           string        `json:"phone"`
	Logo            string        `json:"logo"`
	Identification  string        `json:"identification"`
	Name            string        `json:"name" description:"Restaurant name"`
	Description     string        `json:"description" description:"Restaurant description"`
	Website         string        `json:"website"`
	BankInfo        BankInfo      `json:"bankInfo" bson:"bankInfo"`
	GeoLocation     GeoLocation   `json:"geoLocation" bson:"geoLocation"`
	PriceLevel      int           `json:"priceLevel" bson:"priceLevel"`
	DietaryCodes    []int         `json:"dietaryCodes" bson:"dietaryCodes"`
	Dietaries       []string      `json:"dietaries,omitempty" bson:"dietaries,omitempty"`
	MealKindCodes   []int         `json:"mealKindCodes" bson:"mealKindCodes"`
	MealKinds       []string      `json:"mealKinds,omitempty" bson:"mealKinds,omitempty"`
	Schedules       []*Schedule   `json:"schedules"`
	PreparationTime int           `json:"preparationTime" bson:"preparationTime"`
	Verify          *Verify       `json:"verify,omitempty" bson:"verify,omitempty"`
	Closed          bool          `json:"closed"`
	Recommend       bool          `json:"recommend"`
	MostPopular     bool          `json:"mostPopular" bson:"mostPopular"`
	CreatedAt       int64         `json:"createdAt" bson:"createdAt" description:"Created date."`
	UpdatedAt       int64         `json:"updatedAt" bson:"updatedAt" description:"Updated date."`
}

// ListBusiness is strucut for genear user search
type ListBusiness struct {
	QueryCode int               `json:"queryCode"`
	IsSlide   bool              `json:"isSlide"`
	Items     []*PublicBusiness `json:"items"`
}

// QueryBusiness is struct for query business
type QueryBusiness struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	PlaceID string  `json:"placeId"`
	Sort    int     `json:"sort"` // 0:all, 1:recommend, 2:most popular, 3:delivery time
	Price   []int   `json:"price"`
	Dietary []int   `json:"dietary"`
}
