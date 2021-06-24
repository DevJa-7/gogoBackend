package model

import "gopkg.in/mgo.v2/bson"

type LastLogin struct {
	Date int64  `json:"date"`
	IP   string `json:"IP"`
}

// User is a user model
type User struct {
	ID                bson.ObjectId   `json:"id,omitempty" bson:"_id,omitempty" description:"Object ID"`
	Email             string          `json:"email" description:"User email address"`
	Password          string          `json:"password"`
	Avatar            string          `json:"avatar" description:"User avatar url"`
	Firstname         string          `json:"firstname" description:"User firstname"`
	Lastname          string          `json:"lastname" description:"User lastname"`
	FacebookUserID    string          `json:"facebookUserId" bson:"facebookUserId" description:"When user login or register signup with Facebook, this field will be updated"`
	Birth             int64           `json:"birth" bson:"birth,omitempty" description:"User birth"`
	CountryCode       string          `json:"countryCode" bson:"countryCode"`
	PhoneCode         string          `json:"phoneCode" bson:"phoneCode"`
	Phone             string          `json:"phone"`
	PromoCode         string          `json:"promoCode" bson:"promoCode"`
	Rating            float32         `json:"rating"`
	Status            bool            `json:"status"`
	Verify            *Verify         `json:"verify,omitempty" bson:"verify,omitempty"`
	OneSignalPlayerID string          `json:"onesignalPlayerId,omitempty" bson:"onesignalPlayerId,omitempty"`
	LastLogin         *LastLogin      `json:"lastLogin" bson:"lastLogin"`
	Platform          string          `json:"platform"`
	Favorites         []bson.ObjectId `json:"favorites"`
	RecentLocation    GeoLocation     `json:"recentLocation" bson:"recentLocation"`
	HomeLocation      GeoLocation     `json:"homeLocation" bson:"homeLocation"`
	WorkLocation      GeoLocation     `json:"workLocation" bson:"workLocation"`
	CreatedAt         int64           `json:"createdAt" bson:"createdAt" description:"User created date"`
	UpdatedAt         int64           `json:"updatedAt" bson:"updatedAt" description:"User updated date. This field will be updated when any update operation will be occured"`
}

// PublicUser struct.
type PublicUser struct {
	*User
	Password          omit `json:"password,omitempty"`
	OneSignalPlayerID omit `json:"onesignalPlayerId,omitempty"`
}

// Fullname method
func (p *User) Fullname() string {
	return p.Firstname + " " + p.Lastname
}
