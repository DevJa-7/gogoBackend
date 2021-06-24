package adminService

import (
	"errors"

	"../../../model"
	"../../../util/crypto"
	"../../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// LoginByInfo returns admin with email and password
func LoginByInfo(admin *model.Admin) (*model.Admin, error) {
	adminCollection, session := adminCollection()
	defer session.Close()

	pipe := []bson.M{{"$match": bson.M{"email": admin.Email, "password": admin.Password}}}
	pipe = append(pipe, basePipe...)
	err := adminCollection.Pipe(pipe).One(admin)
	return admin, err
}

// ReadAdminByEmail returns admin
func ReadAdminByEmail(email string) (*model.Admin, error) {
	adminCollection, session := adminCollection()
	defer session.Close()

	admin := &model.Admin{}
	err := adminCollection.Find(bson.M{"email": email}).Select(bson.M{}).One(admin)
	return admin, err
}

// UpdateVerifyCode update verify code for forgot password
func UpdateVerifyCode(objid bson.ObjectId, ql bson.M) error {
	adminCollection, session := adminCollection()
	defer session.Close()

	// update verify code with object id
	return adminCollection.UpdateId(objid, ql)
}

// CheckVerifyCode checks that exists email and verifyCode
func CheckVerifyCode(ql bson.M) (*model.Admin, error) {
	adminCollection, session := adminCollection()
	defer session.Close()

	admin := &model.Admin{}
	err := adminCollection.Find(ql).One(&admin)
	if err == nil {
		adminCollection.UpdateId(admin.ID, bson.M{
			"$set": bson.M{"verify.isVerified": true},
		})
	}
	return admin, err
}

// ResetPassowrd resets password for admin
func ResetPassowrd(objid bson.ObjectId, pwd *model.Password) (*model.Admin, error) {
	adminCollection, session := adminCollection()
	defer session.Close()

	c, _ := adminCollection.Find(bson.M{"_id": objid, "password": crypto.GenerateHash(pwd.Old)}).Count()
	if c == 0 {
		return nil, errors.New("Admin is not existed")
	}

	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"password":  crypto.GenerateHash(pwd.New),
			"updatedAt": timeHelper.GetCurrentTime(),
		}},
		ReturnNew: true,
	}

	admin := &model.Admin{}
	_, err := adminCollection.FindId(objid).Apply(change, admin)
	if err != nil {
		return nil, err
	}

	return admin, nil
}
