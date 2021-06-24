package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Fare struct.
type Fare struct {
	BaseFare   float32 `json:"baseFare" bson:"baseFare"`
	MinFare    float32 `json:"minFare" bson:"minFare"` // limit fare
	PerKm      float32 `json:"perKm" bson:"perKm"`
	PerMinute  float32 `json:"perMinute" bson:"perMinute"`
	NightSurge float32 `json:"nightSurge" bson:"nightSurge"`
	BookingFee float32 `json:"bookingFee" bson:"bookingFee"`
}

// VehicleInfo struct.
type VehicleInfo struct {
	VehicleID bson.ObjectId `json:"vehicleId" bson:"vehicleId"`
	Vehicle   *Vehicle      `json:"vehicle,omitempty" bson:"vehicle,omitempty"`
	Fare      *Fare         `json:"fare"`

	EstimateTime float64 `json:"estimateTime,omitempty" bson:"-,omitempty"`
}

// Period struct.
type Period struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// Location struct.
type Location struct {
	ID           bson.ObjectId  `json:"id" bson:"_id,omitempty"`
	CountryCode  string         `json:"countryCode" bson:"countryCode"`
	City         string         `json:"city"`
	Latitude     float64        `json:"latitude"`
	Longitude    float64        `json:"longitude"`
	PlaceID      string         `json:"placeId" bson:"placeId"`
	VehicleInfos []*VehicleInfo `json:"vehicleInfos,omitempty" bson:"vehicleInfos,omitempty"`
	IsDayTime    bool           `json:"isDayTime" bson:"isDayTime"`
	DayTime      *Period        `json:"dayTime,omitempty" bson:"dayTime,omitempty"`
	IsNightTime  bool           `json:"isNightTime" bson:"isNightTime"`
	NightTime    *Period        `json:"nightTime,omitempty" bson:"nightTime,omitempty"`
	CreatedAt    int64          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    int64          `json:"updatedAt" bson:"updatedAt"`
}

// ContainVehicle returns true when location vehicleInfos contatin vehicle
func (location *Location) ContainVehicle(vehicleID bson.ObjectId) bool {
	for _, v := range location.VehicleInfos {
		if vehicleID == v.VehicleID {
			return true
		}
	}
	return false
}
