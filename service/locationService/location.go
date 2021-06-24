package locationService

import (
	"../../db"
	"../../model"
	"../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var basePipe []bson.M

func locationCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("location"), session
}

// InitService inits service
func InitService() {
	basePipe = []bson.M{
		{"$unwind": bson.M{
			"path": "$vehicleInfos",
			"preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         "vehicle",
			"localField":   "vehicleInfos.vehicleId",
			"foreignField": "_id",
			"as":           "vehicleInfos.vehicle",
		}},
		{"$unwind": bson.M{
			"path": "$vehicleInfos.vehicle",
			"preserveNullAndEmptyArrays": true}},
		{"$group": bson.M{
			"_id":          "$_id",
			"countryCode":  bson.M{"$first": "$countryCode"},
			"city":         bson.M{"$first": "$city"},
			"placeId":      bson.M{"$first": "$placeId"},
			"vehicleInfos": bson.M{"$push": "$$ROOT.vehicleInfos"}}},
	}
}

// CreateLocation creates location
func CreateLocation(location *model.Location) (*model.Location, error) {
	locationCollection, session := locationCollection()
	defer session.Close()

	// Create location with intialize data
	location.ID = bson.NewObjectId()
	location.CreatedAt = timeHelper.GetCurrentTime()
	location.UpdatedAt = timeHelper.GetCurrentTime()
	// Insert Data
	err := locationCollection.Insert(location)
	return location, err
}

// ReadLocation returns location with object id
func ReadLocation(objid bson.ObjectId) (*model.Location, error) {
	locationCollection, session := locationCollection()
	defer session.Close()

	location := &model.Location{}
	// Find location with object id
	pipe := []bson.M{
		{"$match": bson.M{"_id": objid}},
	}
	pipe = append(pipe, basePipe...)
	err := locationCollection.Pipe(pipe).One(location)

	// check empty array for admin panel
	if len(location.VehicleInfos) > 0 && len(location.VehicleInfos[0].VehicleID) == 0 {
		location.VehicleInfos = nil
	}
	return location, err
}

// ReadLocationWithPlaceID returns location with place id
func ReadLocationWithPlaceID(placeID string) (*model.Location, error) {
	locationCollection, session := locationCollection()
	defer session.Close()

	location := &model.Location{}
	// Find location with object id
	pipe := []bson.M{
		{"$match": bson.M{"placeId": placeID}},
	}
	pipe = append(pipe, basePipe...)
	err := locationCollection.Pipe(pipe).One(location)
	return location, err
}

// UpdateLocation updates location
func UpdateLocation(objid bson.ObjectId, location *model.Location) (*model.Location, error) {
	locationCollection, session := locationCollection()
	defer session.Close()

	location.UpdatedAt = timeHelper.GetCurrentTime()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"countryCode": location.CountryCode,
			"city":        location.City,
			"latitude":    location.Latitude,
			"longitude":   location.Longitude,
			"placeId":     location.PlaceID,
			"isDayTime":   location.IsDayTime,
			"dayTime":     location.DayTime,
			"isNightTime": location.IsNightTime,
			"nightTime":   location.NightTime,
			"updatedAt":   location.UpdatedAt,
		}},
		ReturnNew: true,
	}

	// Update location
	_, err := locationCollection.FindId(objid).Apply(change, location)
	return location, err
}

// DeleteLocation deletes location with object id
func DeleteLocation(objid bson.ObjectId) error {
	locationCollection, session := locationCollection()
	defer session.Close()

	err := locationCollection.RemoveId(objid)
	return err
}

// ReadLocations return locations after search query
func ReadLocations(query string, offset int, count int, field string, sort int) ([]*model.Location, int, error) {
	locationCollection, session := locationCollection()
	defer session.Close()

	locations := []*model.Location{}
	totalCount := 0
	pipe := []bson.M{}
	if query == "" {
		// Get all locations
		totalCount, _ = locationCollection.Find(bson.M{}).Count()
	} else {
		// Search location by query
		param := bson.M{"$or": []interface{}{
			bson.M{"geoLocation.address": bson.RegEx{Pattern: query, Options: "i"}},
		}}
		totalCount, _ = locationCollection.Find(param).Count()
		pipe = append(pipe, bson.M{"$match": param})
	}
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

	err := locationCollection.Pipe(pipe).All(&locations)
	return locations, totalCount, err
}
