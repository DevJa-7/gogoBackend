package model

// omit is the bool type for omitting a field of struct.
type omit bool

// Verify struct.
type Verify struct {
	Email      string `json:"email,omitempty" bson:"-"`
	Role       string `json:"role,omitempty" bson:"-"`
	IsVerified bool   `json:"isVerified" bson:"isVerified"`
	Code       string `json:"code"`
	CreatedAt  int64  `json:"createdAt" bson:"createdAt"`
}

// GeoJSON  struct
type GeoJSON struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// GeoLocation is geo struct of simple location
type GeoLocation struct {
	PlaceID string  `json:"placeId" bson:"placeId" `
	Address string  `json:"address"`
	GeoJSON GeoJSON `json:"geoJson" bson:"geoJson"`
}

// ListForm struct.
type ListForm struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

// Password struct for reset password
type Password struct {
	Old string `json:"old"`
	New string `json:"new"`
}

// Status is alias for status of every models
type Status int
