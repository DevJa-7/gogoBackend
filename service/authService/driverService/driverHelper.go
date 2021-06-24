package driverService

import (
	"errors"

	"../../../model"
	"../../../util/crypto"
	"../../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// LoginByInfo returns driver with email and password
func LoginByInfo(driver *model.Driver) (*model.Driver, error) {
	driverCollection, session := driverCollection()
	defer session.Close()

	err := driverCollection.Find(bson.M{"email": driver.Email, "password": driver.Password}).One(driver)
	return driver, err
}

// ReadDriverByEmail returns driver
func ReadDriverByEmail(email string) (*model.Driver, error) {
	driverCollection, session := driverCollection()
	defer session.Close()

	driver := &model.Driver{}
	err := driverCollection.Find(bson.M{"email": email}).Select(bson.M{}).One(driver)
	return driver, err
}

// UpdateVerifyCode update verify code for forgot password
func UpdateVerifyCode(objid bson.ObjectId, ql bson.M) error {
	driverCollection, session := driverCollection()
	defer session.Close()

	// update verify code with object id
	return driverCollection.UpdateId(objid, ql)
}

// CheckVerifyCode checks that exists email and verifyCode
func CheckVerifyCode(ql bson.M) (*model.Driver, error) {
	driverCollection, session := driverCollection()
	defer session.Close()

	driver := &model.Driver{}
	err := driverCollection.Find(ql).One(&driver)
	if err == nil {
		driverCollection.UpdateId(driver.ID, bson.M{
			"$set": bson.M{"verify.isVerified": true},
		})
	}
	return driver, err
}

// UpdateLoginInfo updates login info
func UpdateLoginInfo(objid bson.ObjectId, lastLogin *model.LastLogin, onesignalPlayerID string, platform string) error {
	driverCollection, session := driverCollection()
	defer session.Close()

	err := driverCollection.UpdateId(objid, bson.M{
		"$set": bson.M{
			"lastLogin":         lastLogin,
			"platform":          platform,
			"onesignalPlayerId": onesignalPlayerID,
			"updatedAt":         timeHelper.GetCurrentTime(),
		}},
	)
	return err
}

// ResetPassword resets password for admin
func ResetPassword(objid bson.ObjectId, pwd *model.Password) (*model.Driver, error) {
	driverCollection, session := driverCollection()
	defer session.Close()

	c, _ := driverCollection.Find(bson.M{"_id": objid, "password": crypto.GenerateHash(pwd.Old)}).Count()
	if c == 0 {
		return nil, errors.New("Business is not existed")
	}

	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"password":  crypto.GenerateHash(pwd.New),
			"updatedAt": timeHelper.GetCurrentTime(),
		}},
		ReturnNew: true,
	}

	driver := &model.Driver{}
	_, err := driverCollection.FindId(objid).Apply(change, driver)
	if err != nil {
		return nil, err
	}

	return driver, nil
}
