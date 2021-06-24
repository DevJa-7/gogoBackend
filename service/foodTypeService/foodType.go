package foodTypeService

import (
	"errors"

	"../../db"
	"../../model"
	"../../util/random"
	"../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var basePipe []bson.M

func foodTypeCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("food_type"), session
}

// InitService inits service
func InitService() {
	basePipe = []bson.M{
		{"$lookup": bson.M{
			"from":         "business",
			"localField":   "businessId",
			"foreignField": "_id",
			"as":           "business",
		}},
		{"$unwind": bson.M{
			"path": "$business",
			"preserveNullAndEmptyArrays": true}},
	}
}

// CreateFoodType creates foodType
func CreateFoodType(foodType *model.FoodType) (*model.FoodType, error) {
	foodTypeCollection, session := foodTypeCollection()
	defer session.Close()

	// check duplicate email
	if c, _ := foodTypeCollection.Find(bson.M{
		"name":       foodType.Name,
		"businessId": foodType.BusinessID}).Count(); c > 0 {
		return nil, errors.New("Same foodType is registered already in same business")
	}

	// Create foodtype with intialize data
	foodType.ID = bson.NewObjectId()
	foodType.Code = createCode()
	foodType.CreatedAt = timeHelper.GetCurrentTime()
	foodType.UpdatedAt = timeHelper.GetCurrentTime()
	// Insert Data
	err := foodTypeCollection.Insert(foodType)

	return foodType, err
}

// ReadFoodType returns foodType with object id
func ReadFoodType(objid bson.ObjectId) (*model.FoodType, error) {
	foodTypeCollection, session := foodTypeCollection()
	defer session.Close()

	foodType := &model.FoodType{}
	// Find admin with object id
	err := foodTypeCollection.FindId(objid).One(foodType)
	return foodType, err
}

// UpdateFoodType updates foodType
func UpdateFoodType(objid bson.ObjectId, foodType *model.FoodType) (*model.FoodType, error) {
	foodTypeCollection, session := foodTypeCollection()
	defer session.Close()
	// check duplicate name
	if c, _ := foodTypeCollection.Find(bson.M{
		"_id":        bson.M{"$ne": objid},
		"name":       foodType.Name,
		"businessId": foodType.BusinessID}).Count(); c > 0 {
		return nil, errors.New("Same foodType is registered already in same business")
	}

	foodType.UpdatedAt = timeHelper.GetCurrentTime()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"name":        foodType.Name,
			"description": foodType.Description,
			"image":       foodType.Image,
			"updatedAt":   foodType.UpdatedAt,
		}},
		ReturnNew: true,
	}
	// Update foodType
	_, err := foodTypeCollection.FindId(objid).Apply(change, foodType)
	return foodType, err
}

// DeleteFoodType deletes foodType with object id
func DeleteFoodType(objid bson.ObjectId) error {
	foodTypeCollection, session := foodTypeCollection()
	defer session.Close()

	err := foodTypeCollection.RemoveId(objid)
	return err
}

// ReadFoodTypes return foodTypes after search query
func ReadFoodTypes(query string, offset int, count int, field string, sort int) ([]*model.FoodType, int, error) {
	foodTypeCollection, session := foodTypeCollection()
	defer session.Close()

	foodTypes := []*model.FoodType{}
	pipe := []bson.M{}
	if query != "" {
		// Search business by query
		param := bson.M{"$or": []interface{}{
			bson.M{"name": bson.RegEx{Pattern: query, Options: ""}},
		}}
		pipe = append(pipe, bson.M{"$match": param})
	}
	// get total count of collection with initial query
	totalCount := db.GetCountOfCollection(foodTypeCollection, &pipe)

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

	err := foodTypeCollection.Pipe(pipe).All(&foodTypes)

	return foodTypes, totalCount, err
}

func createCode() int {
	foodTypeCollection, session := foodTypeCollection()
	defer session.Close()

	foodType := &model.FoodType{}
	if err := foodTypeCollection.Pipe([]bson.M{
		{"$sort": bson.M{"code": -1}},
		{"$limit": 1},
	}).One(&foodType); err != nil {
		return 1
	}
	return foodType.Code + 1
}

// ReadFoodTypesWithBusiness returns name array with codes
func ReadFoodTypesWithBusiness(businessID bson.ObjectId) ([]*model.FoodType, error) {
	foodTypeCollection, session := foodTypeCollection()
	defer session.Close()

	foodTypes := []*model.FoodType{}
	err := foodTypeCollection.Pipe([]bson.M{
		{"$match": bson.M{"businessId": businessID}},
		{"$sort": bson.M{"name": 1}},
	}).All(&foodTypes)
	return foodTypes, err
}

// CreateSpec creates food option in foodtype
func CreateSpec(objid bson.ObjectId, spec *model.FoodOption) (*model.FoodOption, error) {
	foodTypeCollection, session := foodTypeCollection()
	defer session.Close()

	number := random.GenerateRandomString(6)
	spec.Number = number
	err := foodTypeCollection.UpdateId(objid, bson.M{"$addToSet": bson.M{"foodOptions": spec}})

	return spec, err
}

// ReadSpec reads food option in foodtype
func ReadSpec(objid bson.ObjectId, number string) (*model.FoodOption, error) {
	foodTypeCollection, session := foodTypeCollection()
	defer session.Close()

	spec := &model.FoodOption{}
	err := foodTypeCollection.Pipe([]bson.M{
		{"$match": bson.M{"_id": objid}},
		{"$unwind": "$foodOptions"},
		{"$replaceRoot": bson.M{"newRoot": "$foodOptions"}},
		{"$match": bson.M{"number": number}},
	}).One(&spec)
	return spec, err
}

// UpdateSpec update food option in foodtype
func UpdateSpec(objid bson.ObjectId, number string, spec *model.FoodOption) (*model.FoodOption, error) {
	foodTypeCollection, session := foodTypeCollection()
	defer session.Close()

	err := foodTypeCollection.Update(bson.M{"_id": objid, "foodOptions.number": number},
		bson.M{"$set": bson.M{"foodOptions.$.options": spec.Options}})

	return spec, err
}

// DeleteSpec delete food option in foodtype
func DeleteSpec(objid bson.ObjectId, number string) error {
	foodTypeCollection, session := foodTypeCollection()
	defer session.Close()

	err := foodTypeCollection.Update(bson.M{"_id": objid},
		bson.M{"$pull": bson.M{"foodOptions": bson.M{"number": number}}})

	return err
}

// ReadSpecs reads food option in foodtype
func ReadSpecs(objid bson.ObjectId) ([]*model.FoodOption, error) {
	foodTypeCollection, session := foodTypeCollection()
	defer session.Close()

	specs := []*model.FoodOption{}
	err := foodTypeCollection.Pipe([]bson.M{
		{"$match": bson.M{"_id": objid}},
		{"$unwind": "$foodOptions"},
		{"$replaceRoot": bson.M{"newRoot": "$foodOptions"}},
	}).All(&specs)
	return specs, err
}
