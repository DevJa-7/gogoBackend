package driverLocationService

import (
	"log"

	"../../config"
	"../../db"
	"../../model"
	"../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func driverLocationCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("driver_location"), session
}

// InitService inits service
func InitService() {
	driverLocationCollection, session := driverLocationCollection()
	defer session.Close()

	// indexing location
	index := mgo.Index{
		Key:        []string{"$2dsphere:location"},
		Background: true,
	}
	err := driverLocationCollection.EnsureIndex(index)
	if err != nil {
		log.Println(err)
	}
	// db.shops.ensureIndex({location:"2dsphere"})
	// db.getCollection('driver_location').aggregate([
	// {$geoNear: {
	//         "near": { "type": "Point", "coordinates":  [ -73.99279 , 40.719296 ] },
	//         "distanceField": "calculated",
	//         "maxDistance": 2,
	//         "query":{"vehicle_id" : ObjectId("58d94b7abab1f81cfd0b9cd4")},
	//         "includeLocs": "location",
	//         "num": 5,
	//         "spherical": true
	//      }}
	// ])
}

// UpdateDriverLocation updates driverLocation
func UpdateDriverLocation(driverLocation *model.DriverLocation) (*model.DriverLocation, error) {
	driverLocationCollection, session := driverLocationCollection()
	defer session.Close()

	driverLocation.Location.Type = "Point"
	driverLocation.UpdatedAt = timeHelper.GetCurrentTime()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"driverId":  driverLocation.DriverID,
			"vehicleId": driverLocation.VehicleID,
			"location":  driverLocation.Location,
			"status":    driverLocation.Status,
			"updatedAt": driverLocation.UpdatedAt,
		}},
		Upsert:    true,
		ReturnNew: true,
	}

	// Update driverLocation
	_, err := driverLocationCollection.Find(bson.M{"driverId": driverLocation.DriverID}).Apply(change, driverLocation)
	return driverLocation, err
}

// GetNearbyDrivers returns nearby drivers
func GetNearbyDrivers(lat, lng float64, limit int) ([]*model.DriverLocation, error) {
	driverLocationCollection, session := driverLocationCollection()
	defer session.Close()

	var err error
	driverLocations := []*model.DriverLocation{}
	i := 0
	searchRadius := config.DefaultSearchRadius
	for {
		err := driverLocationCollection.Pipe([]bson.M{
			{"$geoNear": bson.M{
				"near":          bson.M{"type": "Point", "coordinates": []float64{lng, lat}},
				"distanceField": "distance",
				"maxDistance":   searchRadius,
				"query":         bson.M{"status": config.Online},
				"includeLocs":   "location",
				"num":           limit,
				"spherical":     true,
			}},
			{"$lookup": bson.M{
				"from":         "driver",
				"localField":   "driverId",
				"foreignField": "_id",
				"as":           "driver",
			}},
			{"$unwind": bson.M{
				"path": "$driver",
				"preserveNullAndEmptyArrays": true}},
		}).All(&driverLocations)

		if err == nil && len(driverLocations) > 0 {
			break
		}
		if i > 3 {
			break
		}
		searchRadius += config.DefaultSearchRadius
		i++
	}
	return driverLocations, err
}

// UpdateDriverStatus update driver's status
func UpdateDriverStatus(objid bson.ObjectId, status int) error {
	driverLocationCollection, session := driverLocationCollection()
	defer session.Close()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"status":    status,
			"updatedAt": timeHelper.GetCurrentTime(),
		}},
		ReturnNew: true,
	}
	// Update driverLocation
	_, err := driverLocationCollection.Find(bson.M{"driverId": objid}).Apply(change, nil)
	return err
}

// ReadDriverLocation reads driver location with driver object id
func ReadDriverLocation(driverID bson.ObjectId) (*model.DriverLocation, error) {
	driverLocationCollection, session := driverLocationCollection()
	defer session.Close()

	driverLocation := &model.DriverLocation{}
	err := driverLocationCollection.Find(bson.M{"driverId": driverID}).One(&driverLocation)
	return driverLocation, err
}
