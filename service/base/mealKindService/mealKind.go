package mealKindService

import (
	"errors"

	"../../../db"
	"../../../model"
	"../../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func mealKindCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("meal_kind"), session
}

// InitService inits service
func InitService() {
	CreateMealKind(&model.MealKind{Name: "Breakfast"})
	CreateMealKind(&model.MealKind{Name: "Lunch & Dinner"})
}

// CreateMealKind creates mealKind
func CreateMealKind(mealKind *model.MealKind) (*model.MealKind, error) {
	mealKindCollection, session := mealKindCollection()
	defer session.Close()

	// Check if url is existed already
	if c, _ := mealKindCollection.Find(bson.M{"name": mealKind.Name}).Count(); c > 0 {
		return nil, errors.New("This mealKind is registered already")
	}
	// Create url with intialize data
	mealKind.ID = bson.NewObjectId()
	mealKind.Code = createCode()
	mealKind.CreatedAt = timeHelper.GetCurrentTime()
	mealKind.UpdatedAt = timeHelper.GetCurrentTime()
	// Insert Data
	err := mealKindCollection.Insert(mealKind)

	return mealKind, err
}

// ReadMealKind returns mealKind with object id
func ReadMealKind(objid bson.ObjectId) (*model.MealKind, error) {
	mealKindCollection, session := mealKindCollection()
	defer session.Close()

	mealKind := &model.MealKind{}
	// Find admin with object id
	err := mealKindCollection.FindId(objid).One(mealKind)
	return mealKind, err
}

// UpdateMealKind updates mealKind
func UpdateMealKind(objid bson.ObjectId, mealKind *model.MealKind) (*model.MealKind, error) {
	mealKindCollection, session := mealKindCollection()
	defer session.Close()

	mealKind.UpdatedAt = timeHelper.GetCurrentTime()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"name":      mealKind.Name,
			"image":     mealKind.Image,
			"updatedAt": mealKind.UpdatedAt,
		}},
		ReturnNew: true,
	}
	// Update mealKind
	_, err := mealKindCollection.FindId(objid).Apply(change, mealKind)
	return mealKind, err
}

// DeleteMealKind deletes mealKind with object id
func DeleteMealKind(objid bson.ObjectId) error {
	mealKindCollection, session := mealKindCollection()
	defer session.Close()

	err := mealKindCollection.RemoveId(objid)
	return err
}

// ReadMealKinds return mealKinds after search query
func ReadMealKinds(query string, offset int, count int, field string, sort int) ([]*model.MealKind, int, error) {
	mealKindCollection, session := mealKindCollection()
	defer session.Close()

	mealKinds := []*model.MealKind{}
	pipe := []bson.M{}
	if query != "" {
		// Search business by query
		param := bson.M{"$or": []interface{}{
			bson.M{"name": bson.RegEx{Pattern: query, Options: ""}},
		}}
		pipe = append(pipe, bson.M{"$match": param})
	}
	// get total count of collection with initial query
	totalCount := db.GetCountOfCollection(mealKindCollection, &pipe)

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

	err := mealKindCollection.Pipe(pipe).All(&mealKinds)

	return mealKinds, totalCount, err
}

func createCode() int {
	mealKindCollection, session := mealKindCollection()
	defer session.Close()

	mealKind := &model.MealKind{}
	if err := mealKindCollection.Pipe([]bson.M{
		{"$sort": bson.M{"code": -1}},
		{"$limit": 1},
	}).One(&mealKind); err != nil {
		return 1
	}
	return mealKind.Code + 1
}

// ReadMealKindsWithCodes returns name array with codes
func ReadMealKindsWithCodes(codes []int) []string {
	mealKindCollection, session := mealKindCollection()
	defer session.Close()

	type result struct {
		Results []string `json:"results"`
	}
	r := result{}
	mealKindCollection.Pipe([]bson.M{
		{"$match": bson.M{"code": bson.M{"$in": codes}}},
		{"$group": bson.M{"_id": 0, "results": bson.M{"$push": "$name"}}},
	}).One(&r)
	return r.Results
}
