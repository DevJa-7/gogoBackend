package userService

import (
	"errors"

	"../../../db"
	"../../../model"
	"../../../util/crypto"
	"../../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var basePipe []bson.M

func userCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("user"), session
}

// InitService inits service
func InitService() {
	basePipe = []bson.M{}
}

// CreateUser creates a user
func CreateUser(user *model.User) (*model.User, error) {
	userCollection, session := userCollection()
	defer session.Close()
	// check duplicate email
	if c, _ := userCollection.Find(bson.M{"email": user.Email}).Count(); c > 0 {
		return nil, errors.New("Same email is registered already")
	}
	user.ID = bson.NewObjectId()
	user.Password = crypto.GenerateHash(user.Password)
	user.CreatedAt = timeHelper.GetCurrentTime()
	user.UpdatedAt = timeHelper.GetCurrentTime()

	// Insert Data
	err := userCollection.Insert(user)
	return user, err
}

// ReadUser reads a user
func ReadUser(objid bson.ObjectId) (*model.User, error) {
	userCollection, session := userCollection()
	defer session.Close()

	user := &model.User{}
	// Read Data
	err := userCollection.FindId(objid).One(&user)
	return user, err
}

// UpdateUser reads a user
func UpdateUser(objid bson.ObjectId, user *model.User) (*model.User, error) {
	userCollection, session := userCollection()
	defer session.Close()
	// check duplicate email
	if c, _ := userCollection.Find(bson.M{"_id": bson.M{"$ne": objid}, "email": user.Email}).Count(); c > 0 {
		return nil, errors.New("Same email is registered already")
	}
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"avatar":            user.Avatar,
			"firstname":         user.Firstname,
			"lastname":          user.Lastname,
			"email":             user.Email,
			"brith":             user.Birth,
			"countryCode":       user.CountryCode,
			"phoneCode":         user.PhoneCode,
			"phone":             user.Phone,
			"verify.isVerified": user.Verify.IsVerified,
			"status":            user.Status,
			"promoCode":         user.PromoCode,
			"homeLocation":      user.HomeLocation,
			"workLocation":      user.WorkLocation,
			"updatedAt":         timeHelper.GetCurrentTime(),
		}},
		ReturnNew: true,
	}
	_, err := userCollection.FindId(objid).Apply(change, user)
	return user, err
}

// DeleteUser deletes user with object id
func DeleteUser(objid bson.ObjectId) error {
	userCollection, session := userCollection()
	defer session.Close()

	err := userCollection.RemoveId(objid)
	return err
}

// ReadUsers return users after retreive with params
func ReadUsers(query string, offset int, count int, field string, sort int) ([]*model.User, int, error) {
	userCollection, session := userCollection()
	defer session.Close()

	users := []*model.User{}
	totalCount := 0
	pipe := []bson.M{}
	if query != "" {
		// Search user by query
		param := bson.M{"$or": []interface{}{
			bson.M{"email": bson.RegEx{Pattern: query, Options: ""}},
			bson.M{"firstname": bson.RegEx{Pattern: query, Options: ""}},
			bson.M{"lastname": bson.RegEx{Pattern: query, Options: ""}},
		}}
		pipe = append(pipe, bson.M{"$match": param})
	}
	// get total count of collection with initial query
	totalCount = db.GetCountOfCollection(userCollection, &pipe)

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

	err := userCollection.Pipe(pipe).All(&users)

	return users, totalCount, err
}

// ReadUserFavBusiness reads user's favorite businesses
func ReadUserFavBusiness(objid bson.ObjectId) ([]*model.PublicBusiness, error) {
	userCollection, session := userCollection()
	defer session.Close()

	businesses := []*model.PublicBusiness{}
	err := userCollection.Pipe([]bson.M{
		{"$match": bson.M{"_id": objid}},
		{"$unwind": "$favorites"},
		{"$lookup": bson.M{
			"from":         "business",
			"localField":   "favorites",
			"foreignField": "_id",
			"as":           "favorite",
		}},
		{"$unwind": "$favorite"},
		{"$replaceRoot": bson.M{"newRoot": "$favorite"}},
	}).All(&businesses)
	return businesses, err
}

// UpdateFavorite add or remove favorite business id
func UpdateFavorite(objid bson.ObjectId, businessID bson.ObjectId) error {
	userCollection, session := userCollection()
	defer session.Close()

	if c, err := userCollection.Find(bson.M{"_id": objid, "favorites": businessID}).Count(); err == nil && c > 0 {
		return userCollection.UpdateId(objid, bson.M{
			"$pull": bson.M{"favorites": businessID},
		})
	}

	return userCollection.UpdateId(objid, bson.M{
		"$addToSet": bson.M{"favorites": businessID},
	})
}

// ReadCounts reads total user count and available count
func ReadCounts() (int, int) {
	userCollection, session := userCollection()
	defer session.Close()

	total, _ := userCollection.Find(bson.M{}).Count()
	available, _ := userCollection.Find(bson.M{"status": true}).Count()
	return total, available
}
