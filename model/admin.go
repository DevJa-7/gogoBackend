package model

import "gopkg.in/mgo.v2/bson"

// Admin model
type Admin struct {
	ID        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty" description:"Object ID"`
	Email     string        `json:"email" description:"Business email address"`
	Password  string        `json:"password"`
	Avatar    string        `json:"avatar" description:"Avatar url"`
	Firstname string        `json:"firstname" description:"Firstname"`
	Lastname  string        `json:"lastname" description:"Lastname"`
	RoleCode  int           `json:"roleCode" bson:"roleCode"`
	Role      *Role         `json:"role,omitempty" bson:"role,omitempty"`
	Verify    *Verify       `json:"verify,omitempty" bson:"verify,omitempty"`
	Status    bool          `json:"status"`
	CreatedAt int64         `json:"createdAt" bson:"createdAt" description:"Created date"`
	UpdatedAt int64         `json:"updatedAt" bson:"updatedAt" description:"Updated date. This field will be updated when any update operation will be occured"`
}

// PublicAdmin struct.
type PublicAdmin struct {
	*Admin
	Password omit `json:"password,omitempty"`
	// Verify   omit `json:"verify,omitempty"`
}

// Fullname method
func (p *Admin) Fullname() string {
	return p.Firstname + " " + p.Lastname
}
