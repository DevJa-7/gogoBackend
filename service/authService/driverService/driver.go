package driverService

import (
	"errors"

	"../../../db"
	"../../../model"
	"../../../util/crypto"
	"../../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var basePipe []bson.M

func driverCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("driver"), session
}

// InitService inits service
func InitService() {
	basePipe = []bson.M{}
}

// CreateDriver creates a driver
func CreateDriver(driver *model.Driver) (*model.Driver, error) {
	driverCollection, session := driverCollection()
	defer session.Close()
	// check duplicate email
	if c, _ := driverCollection.Find(bson.M{"email": driver.Email}).Count(); c > 0 {
		return nil, errors.New("Same email is registered already")
	}
	driver.ID = bson.NewObjectId()
	driver.Password = crypto.GenerateHash(driver.Password)
	driver.CreatedAt = timeHelper.GetCurrentTime()
	driver.UpdatedAt = timeHelper.GetCurrentTime()

	// Insert Data
	err := driverCollection.Insert(driver)
	return driver, err
}

// ReadDriver reads a driver
func ReadDriver(objid bson.ObjectId) (*model.Driver, error) {
	driverCollection, session := driverCollection()
	defer session.Close()

	driver := &model.Driver{}
	// Read Data
	err := driverCollection.FindId(objid).One(&driver)
	return driver, err
}

// UpdateDriver reads a driver
func UpdateDriver(objid bson.ObjectId, driver *model.Driver) (*model.Driver, error) {
	driverCollection, session := driverCollection()
	defer session.Close()
	// check duplicate email
	if c, _ := driverCollection.Find(bson.M{"_id": bson.M{"$ne": objid}, "email": driver.Email}).Count(); c > 0 {
		return nil, errors.New("Same email is registered already")
	}
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"avatar":            driver.Avatar,
			"email":             driver.Email,
			"firstname":         driver.Firstname,
			"lastname":          driver.Lastname,
			"countryCode":       driver.CountryCode,
			"phoneCode":         driver.PhoneCode,
			"phone":             driver.Phone,
			"verify.isVerified": driver.Verify.IsVerified,
			"status":            driver.Status,
			"updatedAt":         timeHelper.GetCurrentTime(),
		}},
		ReturnNew: true,
	}
	_, err := driverCollection.FindId(objid).Apply(change, driver)
	return driver, err
}

// DeleteDriver deletes driver with object id
func DeleteDriver(objid bson.ObjectId) error {
	driverCollection, session := driverCollection()
	defer session.Close()

	err := driverCollection.RemoveId(objid)
	return err
}

// ReadDrivers return drivers after retreive with params
func ReadDrivers(query string, offset int, count int, field string, sort int) ([]*model.Driver, int, error) {
	driverCollection, session := driverCollection()
	defer session.Close()

	drivers := []*model.Driver{}
	totalCount := 0
	pipe := []bson.M{}
	if query != "" {
		// Search driver by query
		param := bson.M{"$or": []interface{}{
			bson.M{"email": bson.RegEx{Pattern: query, Options: ""}},
			bson.M{"name": bson.RegEx{Pattern: query, Options: ""}},
			bson.M{"description": bson.RegEx{Pattern: query, Options: ""}},
		}}
		pipe = append(pipe, bson.M{"$match": param})
	}
	// get total count of collection with initial query
	totalCount = db.GetCountOfCollection(driverCollection, &pipe)

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
	pipe = append(pipe, basePipe...)

	err := driverCollection.Pipe(pipe).All(&drivers)

	return drivers, totalCount, err
}

// CreateDriverVehicle creates vehicle of driver
func CreateDriverVehicle(objid bson.ObjectId, driverVehicle *model.DriverVehicle) error {
	driverCollection, session := driverCollection()
	defer session.Close()

	err := driverCollection.Update(bson.M{"_id": objid}, bson.M{"$push": bson.M{"driverVehicles": driverVehicle}})
	return err
}

// UpdateDriverVehicle update vehicle with vehicle number
func UpdateDriverVehicle(objid bson.ObjectId, number string, driverVehicle *model.DriverVehicle) error {
	driverCollection, session := driverCollection()
	defer session.Close()

	err := driverCollection.Update(
		bson.M{"_id": objid, "driverVehicles.number": number},
		bson.M{"$set": bson.M{
			"driverVehicles.$.brand":     driverVehicle.Brand,
			"driverVehicles.$.model":     driverVehicle.Model,
			"driverVehicles.$.color":     driverVehicle.Color,
			"driverVehicles.$.year":      driverVehicle.Year,
			"driverVehicles.$.number":    driverVehicle.Number,
			"driverVehicles.$.vehicleId": driverVehicle.VehicleID,
			"driverVehicles.$.documents": driverVehicle.Documents,
		}})
	return err
}

// DeleteDriverVehicle removes vehicle from list
func DeleteDriverVehicle(objid bson.ObjectId, driverVehicle *model.DriverVehicle) error {
	driverCollection, session := driverCollection()
	defer session.Close()

	err := driverCollection.Update(bson.M{"_id": objid},
		bson.M{"$pull": bson.M{"driver_vehicles": bson.M{"vehicle_id": driverVehicle.VehicleID, "number": driverVehicle.Number}}})
	return err
}

// ReadDriverVehicles returns vehicles of driver
func ReadDriverVehicles(objid bson.ObjectId) ([]*model.DriverVehicle, error) {
	driverCollection, session := driverCollection()
	defer session.Close()

	driverVehicles := []*model.DriverVehicle{}
	// Find driver with object id
	err := driverCollection.Pipe([]bson.M{
		{"$match": bson.M{"_id": objid}},
		{"$unwind": bson.M{"path": "$driverVehicles",
			"preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         "vehicle",
			"localField":   "driverVehicles.vehicleId",
			"foreignField": "_id",
			"as":           "driverVehicles.vehicle"}},
		{"$unwind": "$driverVehicles.vehicle"},
		{"$replaceRoot": bson.M{"newRoot": "$driverVehicles"}},
	}).All(driverVehicles)
	return driverVehicles, err
}

// ReadCounts reads total user count and available count
func ReadCounts() (int, int) {
	driverCollection, session := driverCollection()
	defer session.Close()

	total, _ := driverCollection.Find(bson.M{}).Count()
	available, _ := driverCollection.Find(bson.M{"status": true}).Count()
	return total, available
}
