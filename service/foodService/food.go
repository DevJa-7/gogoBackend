package foodService

import (
	"errors"

	"../../db"
	"../../model"
	"../../service/base/dietaryService"
	"../../service/base/mealKindService"
	"../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func foodCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("food"), session
}

// InitService inits service
func InitService() {

}

// CreateFood creates food
func CreateFood(food *model.Food) (*model.Food, error) {
	foodCollection, session := foodCollection()
	defer session.Close()

	// Check if url is existed already
	if c, _ := foodCollection.Find(bson.M{"name": food.Name, "businessId": food.BusinessID}).Count(); c > 0 {
		return nil, errors.New("This food is registered already")
	}
	// Create url with intialize data
	food.ID = bson.NewObjectId()
	food.CreatedAt = timeHelper.GetCurrentTime()
	food.UpdatedAt = timeHelper.GetCurrentTime()
	// Insert Data
	err := foodCollection.Insert(food)

	return food, err
}

// ReadFood returns food with object id
func ReadFood(objid bson.ObjectId) (*model.Food, error) {
	foodCollection, session := foodCollection()
	defer session.Close()

	food := &model.Food{}
	// Find food with object id
	err := foodCollection.FindId(objid).One(food)
	return food, err
}

// UpdateFood updates food
func UpdateFood(objid bson.ObjectId, food *model.Food) (*model.Food, error) {
	foodCollection, session := foodCollection()
	defer session.Close()

	food.UpdatedAt = timeHelper.GetCurrentTime()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"foodType":      food.FoodType,
			"mealKindCodes": food.MealKindCodes,
			"dietaryCodes":  food.DietaryCodes,
			"image":         food.Image,
			"name":          food.Name,
			"description":   food.Description,
			"price":         food.Price,
			"soldOut":       food.SoldOut,
			"enabled":       food.Enabled,
			"recommend":     food.Recommend,
			"mostPopular":   food.MostPopular,
			"freeDelivery":  food.FreeDelivery,
			"updatedAt":     food.UpdatedAt,
		}},
		ReturnNew: true,
	}
	// Update food
	_, err := foodCollection.FindId(objid).Apply(change, food)
	return food, err
}

// DeleteFood deletes food with object id
func DeleteFood(objid bson.ObjectId) error {
	foodCollection, session := foodCollection()
	defer session.Close()

	err := foodCollection.RemoveId(objid)
	return err
}

// ReadFoods return foods after search query
func ReadFoods(query string, offset int, count int, field string, sort int, businessID bson.ObjectId, mealKindCode int) ([]*model.Food, int, error) {
	foodCollection, session := foodCollection()
	defer session.Close()

	foods := []*model.Food{}
	totalCount := 0
	pipe := []bson.M{}

	if businessID != "" {
		pipe = append(pipe, bson.M{"$match": bson.M{"businessId": businessID}})
	}

	if query != "" {
		// Search foods by query
		param := bson.M{"$or": []interface{}{
			bson.M{"name": bson.RegEx{Pattern: query, Options: "i"}},
		}}
		pipe = append(pipe, bson.M{"$match": param})
	}
	// get total count of collection with initial query
	totalCount = db.GetCountOfCollection(foodCollection, &pipe)

	// add sort feature
	if field != "" && sort != 0 {
		pipe = append(pipe, bson.M{"$sort": bson.M{field: sort}})
	} else {
		pipe = append(pipe, bson.M{"$sort": bson.M{"code": 1}})
	}
	// add page feature
	if offset == 0 && count == 0 {
	} else {
		pipe = append(pipe, bson.M{"$skip": offset})
		pipe = append(pipe, bson.M{"$limit": count})
	}
	err := foodCollection.Pipe(pipe).All(&foods)

	return foods, totalCount, err
}

// ReadMostPopularFoods returns most popular foods of business
func ReadMostPopularFoods(businessID bson.ObjectId, mealKindCode int) ([]*model.Food, error) {
	foodCollection, session := foodCollection()
	defer session.Close()

	pipe := []bson.M{
		{"$match": bson.M{
			"businessId":    businessID,
			"mealKindCodes": mealKindCode,
			"mostPopular":   true,
			"enabled":       true,
			"soldOut":       false,
		}},
	}
	foods := []*model.Food{}
	err := foodCollection.Pipe(pipe).All(&foods)

	return foods, err
}

// ReadRecommendFoods returns recommend foods of business
func ReadRecommendFoods(businessID bson.ObjectId, mealKindCode int) ([]*model.Food, error) {
	foodCollection, session := foodCollection()
	defer session.Close()

	pipe := []bson.M{
		{"$match": bson.M{
			"businessId":    businessID,
			"mealKindCodes": mealKindCode,
			"recommend":     true,
			"enabled":       true,
			"soldOut":       false,
		}},
	}
	foods := []*model.Food{}
	err := foodCollection.Pipe(pipe).All(&foods)

	return foods, err
}

// ReadFoodsByBusiness return foods after search query
func ReadFoodsByBusiness(query string, offset int, count int, field string, sort int, businessID bson.ObjectId, mealKindCode int) ([]*model.BusinessFood, int, error) {
	foodCollection, session := foodCollection()
	defer session.Close()

	businessFoods := []*model.BusinessFood{}
	totalCount := 0
	pipe := []bson.M{}

	if businessID != "" {
		pipe = append(pipe, bson.M{"$match": bson.M{"businessId": businessID}})
	}
	if mealKindCode != 0 {
		pipe = append(pipe, bson.M{"$match": bson.M{"mealKindCodes": mealKindCode}})
	}
	if query != "" {
		// Search foods by query
		param := bson.M{"$or": []interface{}{
			bson.M{"name": bson.RegEx{Pattern: query, Options: "i"}},
		}}
		pipe = append(pipe, bson.M{"$match": param})
	}
	// get total count of collection with initial query
	totalCount = db.GetCountOfCollection(foodCollection, &pipe)

	pipe = append(pipe, []bson.M{
		{"$group": bson.M{"_id": "$foodType.name",
			"foodType": bson.M{"$first": "$foodType.name"},
			"foods":    bson.M{"$push": "$$ROOT"}}},
	}...)
	// add sort feature
	if field != "" && sort != 0 {
		pipe = append(pipe, bson.M{"$sort": bson.M{field: sort}})
	} else {
		pipe = append(pipe, bson.M{"$sort": bson.M{"code": 1}})
	}
	// add page feature
	if offset == 0 && count == 0 {
	} else {
		pipe = append(pipe, bson.M{"$skip": offset})
		pipe = append(pipe, bson.M{"$limit": count})
	}
	err := foodCollection.Pipe(pipe).All(&businessFoods)

	return businessFoods, totalCount, err
}

// RetrieveFoodBaseStructure retrieve base datas
func RetrieveFoodBaseStructure(f *model.Food) {
	f.Dietaries = dietaryService.ReadDietariesWithCodes(f.DietaryCodes)
	f.MealKinds = mealKindService.ReadMealKindsWithCodes(f.MealKindCodes)
}

func ReadWithoutApprove() ([]*model.Food, error) {
	foodCollection, session := foodCollection()
	defer session.Close()

	foods := []*model.Food{}
	err := foodCollection.Pipe([]bson.M{
		{"$match": bson.M{"status": false}},
		{"$lookup": bson.M{
			"from":         "business",
			"localField":   "businessId",
			"foreignField": "_id",
			"as":           "business",
		}},
		{"$unwind": bson.M{
			"path": "$business",
			"preserveNullAndEmptyArrays": true,
		}},
	}).All(&foods)

	return foods, err
}
