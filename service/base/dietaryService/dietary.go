package dietaryService

import (
	"errors"

	"../../../db"
	"../../../model"
	"../../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func dietaryCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("dietary"), session
}

// InitService inits service
func InitService() {
	CreateDietary(&model.Dietary{Name: "Vegetarian"})
	CreateDietary(&model.Dietary{Name: "Vegan"})
	CreateDietary(&model.Dietary{Name: "Glueten-free"})
}

// CreateDietary creates dietary
func CreateDietary(dietary *model.Dietary) (*model.Dietary, error) {
	dietaryCollection, session := dietaryCollection()
	defer session.Close()

	// Check if dietary is existed already
	if c, _ := dietaryCollection.Find(bson.M{"name": dietary.Name}).Count(); c > 0 {
		return nil, errors.New("This dietary is registered already")
	}
	// Create url with intialize data
	dietary.ID = bson.NewObjectId()
	dietary.Code = createCode()
	dietary.CreatedAt = timeHelper.GetCurrentTime()
	dietary.UpdatedAt = timeHelper.GetCurrentTime()
	// Insert Data
	err := dietaryCollection.Insert(dietary)

	return dietary, err
}

// ReadDietary returns dietary with object id
func ReadDietary(objid bson.ObjectId) (*model.Dietary, error) {
	dietaryCollection, session := dietaryCollection()
	defer session.Close()

	dietary := &model.Dietary{}
	// Find dietary with object id
	err := dietaryCollection.FindId(objid).One(dietary)
	return dietary, err
}

// UpdateDietary updates dietary
func UpdateDietary(objid bson.ObjectId, dietary *model.Dietary) (*model.Dietary, error) {
	dietaryCollection, session := dietaryCollection()
	defer session.Close()

	// Check default dietary count
	// if dietary.Default {
	// 	if c, _ := dietaryCollection.Find(bson.M{"default": true}).Count(); c >= 3 {
	// 		return nil, errors.New("Default count is over")
	// 	}
	// }
	// Check if dietary is existed already
	if c, _ := dietaryCollection.Find(bson.M{"_id": bson.M{"$ne": dietary.ID}, "name": dietary.Name}).Count(); c > 0 {
		return nil, errors.New("This dietary is registered already")
	}

	dietary.UpdatedAt = timeHelper.GetCurrentTime()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"name":      dietary.Name,
			"image":     dietary.Image,
			"icon":      dietary.Icon,
			"top":       dietary.Top,
			"default":   dietary.Default,
			"updatedAt": dietary.UpdatedAt,
		}},
		ReturnNew: true,
	}
	// Update dietary
	_, err := dietaryCollection.FindId(objid).Apply(change, dietary)
	return dietary, err
}

// DeleteDietary deletes dietary with object id
func DeleteDietary(objid bson.ObjectId) error {
	dietaryCollection, session := dietaryCollection()
	defer session.Close()

	err := dietaryCollection.RemoveId(objid)
	return err
}

// ReadDietaries return dietaries after search query
func ReadDietaries(query string, offset int, count int, field string, sort int, top int, def bool) ([]*model.Dietary, int, error) {
	dietaryCollection, session := dietaryCollection()
	defer session.Close()

	dietaries := []*model.Dietary{}
	pipe := []bson.M{}
	if top < 2 {
		b := top != 0
		pipe = append(pipe, bson.M{"$match": bson.M{"top": b}})
	}
	if def {
		pipe = append(pipe, bson.M{"$match": bson.M{"default": def}})
	}

	if query != "" {
		// Search items by query
		param := bson.M{"$or": []interface{}{
			bson.M{"name": bson.RegEx{Pattern: query, Options: "i"}},
		}}
		pipe = append(pipe, bson.M{"$match": param})
	}
	// get total count of collection with initial query
	totalCount := db.GetCountOfCollection(dietaryCollection, &pipe)

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
	err := dietaryCollection.Pipe(pipe).All(&dietaries)

	return dietaries, totalCount, err
}

func createCode() int {
	dietaryCollection, session := dietaryCollection()
	defer session.Close()

	dietary := &model.Dietary{}
	if err := dietaryCollection.Pipe([]bson.M{
		{"$sort": bson.M{"code": -1}},
		{"$limit": 1},
	}).One(&dietary); err != nil {
		return 1
	}
	return dietary.Code + 1
}

// ReadDietariesWithCodes returns name array with codes
func ReadDietariesWithCodes(codes []int) []string {
	dietaryCollection, session := dietaryCollection()
	defer session.Close()

	type result struct {
		Results []string `json:"results"`
	}
	r := result{}
	dietaryCollection.Pipe([]bson.M{
		{"$match": bson.M{"code": bson.M{"$in": codes}}},
		{"$group": bson.M{"_id": 0, "results": bson.M{"$push": "$name"}}},
	}).One(&r)
	return r.Results
}
