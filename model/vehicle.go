package model

import "gopkg.in/mgo.v2/bson"

// Vehicle struct.
type Vehicle struct {
	ID           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Photo        string        `json:"photo"`                    // png file
	Image        string        `json:"image"`                    // drawn image
	MenuIcon     string        `json:"menuIcon" bson:"menuIcon"` // icon that will be shown in menu
	MapIcon      string        `json:"mapIcon" bson:"mapIcon"`   // icon that will be shown in map
	Title        string        `json:"title"`
	Detail       string        `json:"detail"`
	Description  string        `json:"description"`
	MaxSeat      int           `json:"maxSeat" bson:"maxSeat"`
	SearchRadius float32       `json:"searchRadius" bson:"searchRadius"`
	Status       bool          `json:"status"`
	Enabled      bool          `json:"enabled"`
	CreatedAt    int64         `json:"createdAt" bson:"createdAt"`
	UpdatedAt    int64         `json:"updatedAt" bson:"updatedAt"`
}
