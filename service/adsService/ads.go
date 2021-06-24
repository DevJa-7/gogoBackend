package adsService

import (
	"errors"

	"../../db"
	"../../model"
	"../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func adsCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("ads"), session
}

// InitService inits service
func InitService() {

}

// CreateAds creates ads
func CreateAds(ads *model.Ads) (*model.Ads, error) {
	adsCollection, session := adsCollection()
	defer session.Close()

	// Check if url is existed already
	if c, _ := adsCollection.Find(bson.M{"name": ads.Name}).Count(); c > 0 {
		return nil, errors.New("This ads is registered already")
	}
	// Create url with intialize data
	ads.ID = bson.NewObjectId()
	ads.CreatedAt = timeHelper.GetCurrentTime()
	ads.UpdatedAt = timeHelper.GetCurrentTime()
	// Insert Data
	err := adsCollection.Insert(ads)

	return ads, err
}

// ReadAds returns ads with object id
func ReadAds(objid bson.ObjectId) (*model.Ads, error) {
	adsCollection, session := adsCollection()
	defer session.Close()

	ads := &model.Ads{}
	// Find ads with object id
	err := adsCollection.FindId(objid).One(ads)
	return ads, err
}

// UpdateAds updates ads
func UpdateAds(objid bson.ObjectId, ads *model.Ads) (*model.Ads, error) {
	adsCollection, session := adsCollection()
	defer session.Close()

	ads.UpdatedAt = timeHelper.GetCurrentTime()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"name":        ads.Name,
			"description": ads.Description,
			"image":       ads.Image,
			"status":      ads.Status,
			"updatedAt":   ads.UpdatedAt,
		}},
		ReturnNew: true,
	}
	// Update ads
	_, err := adsCollection.FindId(objid).Apply(change, ads)
	return ads, err
}

// DeleteAds deletes ads with object id
func DeleteAds(objid bson.ObjectId) error {
	adsCollection, session := adsCollection()
	defer session.Close()

	err := adsCollection.RemoveId(objid)
	return err
}

// ReadAllAds return adss after search query
func ReadAllAds(query string, offset int, count int, field string, sort int) ([]*model.Ads, int, error) {
	adsCollection, session := adsCollection()
	defer session.Close()

	adss := []*model.Ads{}
	totalCount := 0
	pipe := []bson.M{}

	if query != "" {
		// Search adss by query
		param := bson.M{"$or": []interface{}{
			bson.M{"name": bson.RegEx{Pattern: query, Options: "i"}},
		}}
		pipe = append(pipe, bson.M{"$match": param})
	}
	// get total count of collection with initial query
	totalCount = db.GetCountOfCollection(adsCollection, &pipe)

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
	err := adsCollection.Pipe(pipe).All(&adss)

	return adss, totalCount, err
}
