package businessService

import (
	"errors"

	"../../../model"
	"../../../util/crypto"
	"../../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// LoginByInfo returns business with email and password
func LoginByInfo(business *model.Business) (*model.PublicBusiness, error) {
	businessCollection, session := businessCollection()
	defer session.Close()

	b := &model.PublicBusiness{}
	err := businessCollection.Find(bson.M{"email": business.Email}).One(&b)
	return b, err
}

// ReadBusinessByEmail returns business
func ReadBusinessByEmail(email string) (*model.Business, error) {
	businessCollection, session := businessCollection()
	defer session.Close()

	business := &model.Business{}
	err := businessCollection.Find(bson.M{"email": email}).Select(bson.M{}).One(business)
	return business, err
}

// UpdateVerifyCode update verify code for forgot password
func UpdateVerifyCode(objid bson.ObjectId, ql bson.M) error {
	businessCollection, session := businessCollection()
	defer session.Close()

	// update verify code with object id
	return businessCollection.UpdateId(objid, ql)
}

// CheckVerifyCode checks that exists email and verifyCode
func CheckVerifyCode(ql bson.M) (*model.PublicBusiness, error) {
	businessCollection, session := businessCollection()
	defer session.Close()

	business := &model.PublicBusiness{}
	err := businessCollection.Find(ql).One(&business)
	if err == nil {
		businessCollection.UpdateId(business.ID, bson.M{
			"$set": bson.M{"verify.isVerified": true},
		})
	}
	return business, err
}

// ReadBusinessDietary returns dietaries of business
func ReadBusinessDietary(objid bson.ObjectId) ([]*model.Dietary, error) {
	businessCollection, session := businessCollection()
	defer session.Close()

	dietaries := []*model.Dietary{}
	err := businessCollection.Pipe([]bson.M{
		{"$match": bson.M{"_id": objid}},
		{"$unwind": "$dietaryCodes"},
		{"$sort": bson.M{"dietaryCodes": 1}},
		{"$lookup": bson.M{
			"from":         "dietary",
			"localField":   "dietaryCodes",
			"foreignField": "code",
			"as":           "dietary",
		}},
		{"$unwind": bson.M{
			"path": "$dietary",
			"preserveNullAndEmptyArrays": true}},
		{"$replaceRoot": bson.M{"newRoot": "$dietary"}},
	}).All(&dietaries)

	return dietaries, err
}

// ReadBusinessMealKind returns meal kinds of business
func ReadBusinessMealKind(objid bson.ObjectId) ([]*model.MealKind, error) {
	businessCollection, session := businessCollection()
	defer session.Close()

	mealKinds := []*model.MealKind{}
	err := businessCollection.Pipe([]bson.M{
		{"$match": bson.M{"_id": objid}},
		{"$unwind": "$mealKindCodes"},
		{"$sort": bson.M{"mealKindCodes": 1}},
		{"$lookup": bson.M{
			"from":         "meal_kind",
			"localField":   "mealKindCodes",
			"foreignField": "code",
			"as":           "mealKind",
		}},
		{"$unwind": bson.M{
			"path": "$mealKind",
			"preserveNullAndEmptyArrays": true}},
		{"$replaceRoot": bson.M{"newRoot": "$mealKind"}},
	}).All(&mealKinds)

	return mealKinds, err
}

// ReadBusinessFoodType returns food types of business
func ReadBusinessFoodType(objid bson.ObjectId) ([]*model.FoodType, error) {
	businessCollection, session := businessCollection()
	defer session.Close()

	foodTypes := []*model.FoodType{}
	err := businessCollection.Pipe([]bson.M{
		{"$match": bson.M{"_id": objid}},
		{"$unwind": "$foodTypeCodes"},
		{"$sort": bson.M{"foodTypeCodes": 1}},
		{"$lookup": bson.M{
			"from":         "food_type",
			"localField":   "foodTypeCodes",
			"foreignField": "code",
			"as":           "foodType",
		}},
		{"$unwind": bson.M{
			"path": "$foodType",
			"preserveNullAndEmptyArrays": true}},
		{"$replaceRoot": bson.M{"newRoot": "$foodType"}},
	}).All(&foodTypes)

	return foodTypes, err
}

// ResetPassword resets password for admin
func ResetPassword(objid bson.ObjectId, pwd *model.Password) (*model.PublicBusiness, error) {
	businessCollection, session := businessCollection()
	defer session.Close()

	c, _ := businessCollection.Find(bson.M{"_id": objid, "password": crypto.GenerateHash(pwd.Old)}).Count()
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

	business := &model.PublicBusiness{}
	_, err := businessCollection.FindId(objid).Apply(change, business)
	if err != nil {
		return nil, err
	}

	return business, nil
}
