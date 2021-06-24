package reasonService

import (
	"errors"

	"../../config"
	"../../db"
	"../../model"
	"../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func reasonCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("reason"), session
}

// InitService inits service
func InitService() {
	CreateReason(&model.Reason{Code: 300, Type: config.ReasonTripCancel, Status: true, Message: "Unable to find pickup"})
	CreateReason(&model.Reason{Code: 301, Type: config.ReasonTripCancel, Status: true, Message: "Oversized item"})
	CreateReason(&model.Reason{Code: 302, Type: config.ReasonTripCancel, Status: true, Message: "Excessive wait time"})
	CreateReason(&model.Reason{Code: 303, Type: config.ReasonTripCancel, Status: true, Message: "Other"})
	CreateReason(&model.Reason{Code: 304, Type: config.ReasonTripCancel, Status: true, Message: "Too far away"})
	CreateReason(&model.Reason{Code: 305, Type: config.ReasonTripCancel, Status: true, Message: "I don't want to do  delivery"})
}

// CreateReason creates reason
func CreateReason(reason *model.Reason) (*model.Reason, error) {
	reasonCollection, session := reasonCollection()
	defer session.Close()

	// Check if url is existed already
	result := &model.Reason{}
	if err := reasonCollection.Find(bson.M{"code": reason.Code}).One(&result); err == nil && result.ID != "" {
		return nil, errors.New("This reason is registered already")
	}
	// Create reason with intialize data
	reason.ID = bson.NewObjectId()
	if reason.Code == 0 {
		reason.Code = createCode()
	}
	reason.CreatedAt = timeHelper.GetCurrentTime()
	reason.UpdatedAt = timeHelper.GetCurrentTime()
	// Insert Data
	err := reasonCollection.Insert(reason)
	return reason, err
}

// ReadReason returns vehicle with object id
func ReadReason(objid bson.ObjectId) (*model.Reason, error) {
	reasonCollection, session := reasonCollection()
	defer session.Close()

	reason := &model.Reason{}
	err := reasonCollection.FindId(objid).One(reason)
	return reason, err
}

// UpdateReason updates reason
func UpdateReason(objid bson.ObjectId, reason *model.Reason) (*model.Reason, error) {
	reasonCollection, session := reasonCollection()
	defer session.Close()

	reason.UpdatedAt = timeHelper.GetCurrentTime()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"type":      reason.Type,
			"message":   reason.Message,
			"status":    reason.Status,
			"updatedAt": reason.UpdatedAt,
		}},
		ReturnNew: true,
	}

	// Update reason
	_, err := reasonCollection.FindId(objid).Apply(change, reason)
	return reason, err
}

// DeleteReason deletes reason with object id
func DeleteReason(objid bson.ObjectId) error {
	reasonCollection, session := reasonCollection()
	defer session.Close()

	err := reasonCollection.RemoveId(objid)
	return err
}

// ReadReasons return reasons after search query
func ReadReasons(query string, offset int, count int, field string, sort int, t int) ([]*model.Reason, int, error) {
	reasonCollection, session := reasonCollection()
	defer session.Close()

	reasons := []*model.Reason{}
	totalCount := 0
	pipe := []bson.M{}
	if t != 0 {
		pipe = append(pipe, bson.M{"$match": bson.M{"type": t}})
	}

	if query == "" {
		// Get all vehicles
		totalCount, _ = reasonCollection.Find(bson.M{}).Count()
	} else {
		// Search vehicle by query
		param := bson.M{"$or": []interface{}{
			bson.M{"message": bson.RegEx{Pattern: query, Options: "i"}},
		}}
		totalCount, _ = reasonCollection.Find(param).Count()
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

	err := reasonCollection.Pipe(pipe).All(&reasons)
	return reasons, totalCount, err
}

func createCode() int {
	reasonCollection, session := reasonCollection()
	defer session.Close()

	reason := &model.Reason{}
	if err := reasonCollection.Pipe([]bson.M{
		{"$sort": bson.M{"code": -1}},
		{"$limit": 1},
	}).One(&reason); err != nil {
		return 1
	}
	return reason.Code + 1
}
