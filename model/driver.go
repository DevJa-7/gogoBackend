package model

import "gopkg.in/mgo.v2/bson"

// DriverLocation struct
type DriverLocation struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	DriverID  bson.ObjectId `json:"driverId,omitempty" bson:"driverId,omitempty"`
	Driver    *Driver       `json:"driver,omitempty"  bson:"driver,omitempty"`
	VehicleID bson.ObjectId `json:"vehicleId,omitempty" bson:"vehicleId,omitempty"`
	Vehicle   *Vehicle      `json:"vehicle,omitempty" bson:"vehicle,omitempty"`
	Location  GeoJSON       `json:"location"`
	PlaceID   string        `json:"placeId" bson:"placeId"`
	Status    Status        `json:"status"` //100:offline, 101:online, 102:ongoing
	Distance  float64       `json:"distance,omitempty" bson:"distance,omitempty"`
	UpdatedAt int64         `json:"updatedAt" bson:"updatedAt"`
}

// VerifyDocument struct is that for verify document to enable driver's vehicle
type VerifyDocument struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Document   *Document     `json:"document,omitempty" bson:"document,omitempty"`
	Status     Status        `json:"status"` //0:not uploaded 1:pending 2:verifed
	ExpireDate int64         `json:"expireDate,omitempty" bson:"expireDate,omitempty"`
	Image      string        `json:"image"`
}

// DriverVehicle struct
type DriverVehicle struct {
	Number    string            `json:"number"`
	VehicleID bson.ObjectId     `json:"vehicleId" bson:"vehicleId"`
	Vehicle   *Vehicle          `json:"vehicle,omitempty" bson:"vehicle,omitempty"`
	Brand     string            `json:"brand"`
	Model     string            `json:"model"`
	Color     string            `json:"color"`
	Year      int               `json:"year"`
	Documents []*VerifyDocument `json:"documents"`
}

// Driver model
type Driver struct {
	ID                bson.ObjectId    `json:"id,omitempty" bson:"_id,omitempty" description:"Object ID"`
	Email             string           `json:"email" description:"Business email address"`
	Password          string           `json:"password"`
	Avatar            string           `json:"avatar" description:"Driver avatar url"`
	Firstname         string           `json:"firstname" description:"Driver firstname"`
	Lastname          string           `json:"lastname" description:"Driver lastname"`
	Verify            *Verify          `json:"verify,omitempty" bson:"verify,omitempty"`
	Status            bool             `json:"status"` //0:pending 2:activated
	CountryCode       string           `json:"countryCode" bson:"countryCode"`
	PhoneCode         string           `json:"phoneCode" bson:"phoneCode"`
	Phone             string           `json:"phone"`
	Rating            float32          `json:"rating"`
	DriverVehicles    []*DriverVehicle `json:"driverVehicles" bson:"driverVehicles"`
	LocationPlaceID   string           `json:"locationPlaceId,omitempty" bson:"locationPlaceId,omitempty"`
	OneSignalPlayerID string           `json:"onesignalPlayerId,omitempty" bson:"onesignalPlayerId,omitempty"`
	LastLogin         *LastLogin       `json:"lastLogin" bson:"lastLogin"`
	Platform          string           `json:"platform"`
	CreatedAt         int64            `json:"createdAt" bson:"createdAt" description:"Driver created date"`
	UpdatedAt         int64            `json:"updatedAt" bson:"updatedAt" description:"Driver updated date. This field will be updated when any update operation will be occured"`
}

// PublicDriver struct.
type PublicDriver struct {
	*Driver
	Password          omit `json:"password,omitempty"`
	OneSignalPlayerID omit `json:"onesignalPlayerId,omitempty"`
}

// Fullname method
func (p *Driver) Fullname() string {
	return p.Firstname + " " + p.Lastname
}
