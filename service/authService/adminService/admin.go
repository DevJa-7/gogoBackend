package adminService

import (
	"errors"

	"../../../config"
	"../../../db"
	"../../../model"
	"../../../util/crypto"
	"../../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var basePipe []bson.M

func adminCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("admin"), session
}

// InitService inits service
func InitService() {
	basePipe = []bson.M{
		{"$lookup": bson.M{
			"from":         "role",
			"localField":   "roleCode",
			"foreignField": "code",
			"as":           "role",
		}},
		{"$unwind": bson.M{
			"path": "$role",
			"preserveNullAndEmptyArrays": true}},
	}

	CreateAdmin(&model.Admin{
		Firstname: "",
		Lastname:  "Administrator",
		Email:     "admin@gogo.cr",
		Password:  "admin1234!",
		RoleCode:  config.AdminCode,
		Verify: &model.Verify{
			IsVerified: true,
		},
		Status: true,
	})
}

// CreateAdmin creates a admin
func CreateAdmin(admin *model.Admin) (*model.Admin, error) {
	adminCollection, session := adminCollection()
	defer session.Close()
	// check duplicate email
	if c, _ := adminCollection.Find(bson.M{"email": admin.Email}).Count(); c > 0 {
		return nil, errors.New("Same email is registered already")
	}
	admin.ID = bson.NewObjectId()
	admin.Password = crypto.GenerateHash(admin.Password)
	admin.CreatedAt = timeHelper.GetCurrentTime()
	admin.UpdatedAt = timeHelper.GetCurrentTime()

	// Insert Data
	err := adminCollection.Insert(admin)
	return admin, err
}

// ReadAdmin reads a admin
func ReadAdmin(objid bson.ObjectId) (*model.Admin, error) {
	adminCollection, session := adminCollection()
	defer session.Close()

	admin := &model.Admin{}
	// Read Data
	err := adminCollection.FindId(objid).One(&admin)
	return admin, err
}

// UpdateAdmin reads a admin
func UpdateAdmin(objid bson.ObjectId, admin *model.Admin) (*model.Admin, error) {
	adminCollection, session := adminCollection()
	defer session.Close()
	// check duplicate email
	if c, _ := adminCollection.Find(bson.M{"_id": bson.M{"$ne": objid}, "email": admin.Email}).Count(); c > 0 {
		return nil, errors.New("Same email is registered already")
	}
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"avatar":            admin.Avatar,
			"email":             admin.Email,
			"firstname":         admin.Firstname,
			"lastname":          admin.Lastname,
			"roleCode":          admin.RoleCode,
			"verify.isVerified": admin.Verify.IsVerified,
			"status":            admin.Status,
			"updatedAt":         timeHelper.GetCurrentTime(),
		}},
		ReturnNew: true,
	}
	_, err := adminCollection.FindId(objid).Apply(change, admin)
	return admin, err
}

// DeleteAdmin deletes admin with object id
func DeleteAdmin(objid bson.ObjectId) error {
	adminCollection, session := adminCollection()
	defer session.Close()

	err := adminCollection.RemoveId(objid)
	return err
}

// ReadAdmins return admins after retreive with params
func ReadAdmins(query string, offset int, count int, field string, sort int) ([]*model.Admin, int, error) {
	adminCollection, session := adminCollection()
	defer session.Close()

	admins := []*model.Admin{}
	totalCount := 0
	pipe := []bson.M{}
	if query != "" {
		// Search admin by query
		param := bson.M{"$or": []interface{}{
			bson.M{"email": bson.RegEx{Pattern: query, Options: ""}},
			bson.M{"firstname": bson.RegEx{Pattern: query, Options: ""}},
			bson.M{"lastname": bson.RegEx{Pattern: query, Options: ""}},
		}}
		pipe = append(pipe, bson.M{"$match": param})
	}
	// get total count of collection with initial query
	totalCount = db.GetCountOfCollection(adminCollection, &pipe)

	// add sort feature
	if field != "" && sort != 0 {
		pipe = append(pipe, bson.M{"$sort": bson.M{field: sort}})
	}
	// add page feature
	if offset == 0 && count == 0 {
	} else {
		pipe = append(pipe, bson.M{"$skip": offset})
		pipe = append(pipe, bson.M{"$limit": count})
	}
	pipe = append(pipe, basePipe...)

	err := adminCollection.Pipe(pipe).All(&admins)

	return admins, totalCount, err
}

// ReadCounts reads total user count and available count
func ReadCounts() (int, int) {
	adminCollection, session := adminCollection()
	defer session.Close()

	total, _ := adminCollection.Find(bson.M{}).Count()
	available, _ := adminCollection.Find(bson.M{"status": true}).Count()
	return total, available
}
