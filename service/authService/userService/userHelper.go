package userService

import (
	"errors"

	"../../../model"
	"../../../util/crypto"
	"../../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// LoginByInfo returns user with email and password
func LoginByInfo(user *model.User) (*model.User, error) {
	userCollection, session := userCollection()
	defer session.Close()

	err := userCollection.Find(bson.M{"email": user.Email, "password": user.Password}).One(user)
	return user, err
}

// ReadUserByEmail returns user
func ReadUserByEmail(email string) (*model.User, error) {
	userCollection, session := userCollection()
	defer session.Close()

	user := &model.User{}
	err := userCollection.Find(bson.M{"email": email}).One(&user)
	return user, err
}

// ReadUserByFacebookID returns user with facebook id
func ReadUserByFacebookID(facebookID string) (*model.User, error) {
	userCollection, session := userCollection()
	defer session.Close()

	user := &model.User{}
	err := userCollection.Find(bson.M{"facebookUserId": facebookID}).One(&user)
	return user, err
}

// UpdateVerifyCode update verify code for forgot password
func UpdateVerifyCode(objid bson.ObjectId, ql bson.M) error {
	userCollection, session := userCollection()
	defer session.Close()

	// update verify code with object id
	return userCollection.UpdateId(objid, ql)
}

// CheckVerifyCode checks that exists email and verifyCode
func CheckVerifyCode(ql bson.M) (*model.User, error) {
	userCollection, session := userCollection()
	defer session.Close()

	user := &model.User{}
	err := userCollection.Find(ql).One(&user)
	if err == nil {
		userCollection.UpdateId(user.ID, bson.M{
			"$set": bson.M{"verify.isVerified": true},
		})
	}
	return user, err
}

// UpdateLoginInfo updates login info
func UpdateLoginInfo(objid bson.ObjectId, lastLogin *model.LastLogin, onesignalPlayerID string, platform string) error {
	userCollection, session := userCollection()
	defer session.Close()

	err := userCollection.UpdateId(objid, bson.M{
		"$set": bson.M{
			"lastLogin":         lastLogin,
			"platform":          platform,
			"onesignalPlayerId": onesignalPlayerID,
			"updatedAt":         timeHelper.GetCurrentTime(),
		}},
	)
	return err
}

// ResetPassword resets password for user
func ResetPassword(objid bson.ObjectId, pwd *model.Password) (*model.User, error) {
	userCollection, session := userCollection()
	defer session.Close()

	c, _ := userCollection.Find(bson.M{"_id": objid, "password": crypto.GenerateHash(pwd.Old)}).Count()
	if c == 0 {
		return nil, errors.New("User is not existed")
	}

	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"password":  crypto.GenerateHash(pwd.New),
			"updatedAt": timeHelper.GetCurrentTime(),
		}},
		ReturnNew: true,
	}

	user := &model.User{}
	_, err := userCollection.FindId(objid).Apply(change, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
