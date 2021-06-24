package model

import "gopkg.in/mgo.v2/bson"

// Order is a order model
type Order struct {
	ID         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty" description:"Object ID"`
	Number     string        `json:"number" description:"Order number"`
	UserID     bson.ObjectId `json:"userId,omitempty" bson:"userId,omitempty"`
	User       *User         `json:"user,omitempty" bson:"user,omitempty"`
	BusinessID bson.ObjectId `json:"businessId,omitempty" bson:"businessId,omitempty"`
	Business   *Business     `json:"business,omitempty" bson:"business,omitempty"`
	DriverID   bson.ObjectId `json:"driverId,omitempty" bson:"driverId,omitempty"`
	Driver     *Driver       `json:"driver,omitempty" bson:"driver,omitempty"`
	Foods      []*struct {
		Food  *OrderFood `json:"food"`
		Count int        `json:"count"`
		Price float64    `json:"price"`
		Note  string     `json:"note"`
	} `json:"foods"`
	Tax              float64          `json:"tax"`
	BookingFee       float64          `json:"bookingFee" bson:"bookingFee"`
	Price            float64          `json:"price"`
	Instruction      string           `json:"instruction"`
	DeliveryLocation *GeoLocation     `json:"deliveryLocation" bson:"deliveryLocation"`
	OrderStatus      string           `json:"orderStatus" bson:"orderStatus"`
	TripStatus       string           `json:"tripStatus" bson:"tripStatus"`
	StatusAt         map[string]int64 `json:"statusAt" bson:"statusAt"`
	ReasonCode       int              `json:"reasonCode,omitempty" bson:"reasonCode,omitempty"`
	Reason           *Reason          `json:"reason"`
	PickupScore      int              `json:"pickupScore" bson:"pickupScore"`
	Recipient        string           `json:"recipient"`
	Rated            bool             `json:"rated"`
	BusinessRate     float32          `json:"businessRate" bson:"businessRate"`
	BusinessFeedback string           `json:"businessFeedback" bson:"businessFeedback"`
	DriverRate       float32          `json:"driverRate" bson:"driverRate"`
	DriverFeedback   string           `json:"driverFeedback" bson:"driverFeedback"`
	UserRate         float32          `json:"userRate" bson:"userRate"`
	CreatedAt        int64            `json:"createdAt" bson:"createdAt" description:"Created date."`
	UpdatedAt        int64            `json:"updatedAt" bson:"updatedAt" description:"Updated date."`
}
