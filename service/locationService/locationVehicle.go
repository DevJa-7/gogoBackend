package locationService

import (
	"../../model"
	"../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UpdateVehicleInfos updates location
func UpdateVehicleInfos(objid bson.ObjectId, vehicleInfos []*model.VehicleInfo) error {
	locationCollection, session := locationCollection()
	defer session.Close()

	for _, info := range vehicleInfos {
		info.Vehicle = nil
	}
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"vehicleInfos": vehicleInfos,
			"updatedAt":    timeHelper.GetCurrentTime(),
		}},
		ReturnNew: true,
	}

	// Update location
	_, err := locationCollection.FindId(objid).Apply(change, nil)
	return err
}

// ReadVehicleInfoFare returns fare of vehicel in the location
func ReadVehicleInfoFare(placeID string, objid bson.ObjectId) (*model.Fare, error) {
	locationCollection, session := locationCollection()
	defer session.Close()

	vehicleInfo := &model.VehicleInfo{}
	err := locationCollection.Pipe([]bson.M{
		{"$match": bson.M{"geoLocation.placeId": placeID}},
		{"$unwind": "$vehicleInfos"},
		{"$match": bson.M{"vehicleInfos.vehicleId": objid}},
		{"$replaceRoot": bson.M{"newRoot": "$vehicleInfos"}},
	}).One(vehicleInfo)
	return vehicleInfo.Fare, err
}
