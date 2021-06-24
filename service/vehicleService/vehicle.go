package vehicleService

import (
	"../../db"
	"../../model"
	"../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func vehicleCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("vehicle"), session
}

// InitService inits service
func InitService() {
}

// CreateVehicle creates vehicle
func CreateVehicle(vehicle *model.Vehicle) (*model.Vehicle, error) {
	vehicleCollection, session := vehicleCollection()
	defer session.Close()

	// Create vehicle with intialize data
	vehicle.ID = bson.NewObjectId()
	vehicle.CreatedAt = timeHelper.GetCurrentTime()
	vehicle.UpdatedAt = timeHelper.GetCurrentTime()
	// Insert Data
	err := vehicleCollection.Insert(vehicle)
	return vehicle, err
}

// ReadVehicle returns vehicle with object id
func ReadVehicle(objid bson.ObjectId) (*model.Vehicle, error) {
	vehicleCollection, session := vehicleCollection()
	defer session.Close()

	vehicle := &model.Vehicle{}
	err := vehicleCollection.FindId(objid).One(vehicle)
	return vehicle, err
}

// UpdateVehicle updates vehicle
func UpdateVehicle(objid bson.ObjectId, vehicle *model.Vehicle) (*model.Vehicle, error) {
	vehicleCollection, session := vehicleCollection()
	defer session.Close()

	vehicle.UpdatedAt = timeHelper.GetCurrentTime()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"photo":        vehicle.Photo,
			"image":        vehicle.Image,
			"menuIcon":     vehicle.MenuIcon,
			"mapIcon":      vehicle.MapIcon,
			"title":        vehicle.Title,
			"detail":       vehicle.Detail,
			"description":  vehicle.Description,
			"maxSeat":      vehicle.MaxSeat,
			"searchRadius": vehicle.SearchRadius,
			"status":       vehicle.Status,
			"updatedAt":    vehicle.UpdatedAt,
		}},
		ReturnNew: true,
	}

	// Update vehicle
	_, err := vehicleCollection.FindId(objid).Apply(change, vehicle)
	return vehicle, err
}

// DeleteVehicle deletes vehicle with object id
func DeleteVehicle(objid bson.ObjectId) error {
	vehicleCollection, session := vehicleCollection()
	defer session.Close()

	err := vehicleCollection.RemoveId(objid)
	return err
}

// ReadVehicles return vehicles after search query
func ReadVehicles(query string, offset int, count int, field string, sort int) ([]*model.Vehicle, int, error) {
	vehicleCollection, session := vehicleCollection()
	defer session.Close()

	vehicles := []*model.Vehicle{}
	totalCount := 0
	pipe := []bson.M{}
	if query == "" {
		// Get all vehicles
		totalCount, _ = vehicleCollection.Find(bson.M{}).Count()
	} else {
		// Search vehicle by query
		param := bson.M{"$or": []interface{}{
			bson.M{"title": bson.RegEx{Pattern: query, Options: "i"}},
			bson.M{"description": bson.RegEx{Pattern: query, Options: "i"}},
		}}
		totalCount, _ = vehicleCollection.Find(param).Count()
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

	err := vehicleCollection.Pipe(pipe).All(&vehicles)
	return vehicles, totalCount, err
}

// ReadActiveVehicles return active vehicles
func ReadActiveVehicles() ([]*model.Vehicle, error) {
	vehicleCollection, session := vehicleCollection()
	defer session.Close()

	vehicles := []*model.Vehicle{}
	err := vehicleCollection.Find(bson.M{"status": true}).All(&vehicles)
	return vehicles, err
}
