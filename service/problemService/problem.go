package problemService

import (
	"../../db"
	"../../model"
	"../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func problemCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("problem"), session
}

// InitService inits service
func InitService() {
}

// CreateProblem creates problem
func CreateProblem(problem *model.Problem) (*model.Problem, error) {
	problemCollection, session := problemCollection()
	defer session.Close()

	// Create problem with intialize data
	problem.ID = bson.NewObjectId()
	problem.CreatedAt = timeHelper.GetCurrentTime()
	problem.UpdatedAt = timeHelper.GetCurrentTime()
	// Insert Data
	err := problemCollection.Insert(problem)
	return problem, err
}

// ReadProblem returns vehicle with object id
func ReadProblem(objid bson.ObjectId) (*model.Problem, error) {
	problemCollection, session := problemCollection()
	defer session.Close()

	problem := &model.Problem{}
	err := problemCollection.FindId(objid).One(problem)
	return problem, err
}

// UpdateProblem updates problem
func UpdateProblem(objid bson.ObjectId, problem *model.Problem) (*model.Problem, error) {
	problemCollection, session := problemCollection()
	defer session.Close()

	problem.UpdatedAt = timeHelper.GetCurrentTime()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"title":       problem.Title,
			"description": problem.Description,
			"status":      problem.Status,
			"updatedAt":   problem.UpdatedAt,
		}},
		ReturnNew: true,
	}

	// Update problem
	_, err := problemCollection.FindId(objid).Apply(change, problem)
	return problem, err
}

// DeleteProblem deletes problem with object id
func DeleteProblem(objid bson.ObjectId) error {
	problemCollection, session := problemCollection()
	defer session.Close()

	err := problemCollection.RemoveId(objid)
	return err
}

// ReadProblems return problems after search query
func ReadProblems(query string, offset int, count int, field string, sort int) ([]*model.Problem, int, error) {
	problemCollection, session := problemCollection()
	defer session.Close()

	problems := []*model.Problem{}
	totalCount := 0
	pipe := []bson.M{}
	if query == "" {
		// Get all vehicles
		totalCount, _ = problemCollection.Find(bson.M{}).Count()
	} else {
		// Search vehicle by query
		param := bson.M{"$or": []interface{}{
			bson.M{"title": bson.RegEx{Pattern: query, Options: "i"}},
			bson.M{"description": bson.RegEx{Pattern: query, Options: "i"}},
		}}
		totalCount, _ = problemCollection.Find(param).Count()
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

	err := problemCollection.Pipe(pipe).All(&problems)
	return problems, totalCount, err
}

// ReadMerchantProblems returns problems with object id
func ReadMerchantProblems(objid bson.ObjectId) ([]*model.Problem, error) {
	problemCollection, session := problemCollection()
	defer session.Close()

	problems := []*model.Problem{}
	err := problemCollection.Find(bson.M{"businessId": objid}).All(&problems)
	return problems, err
}

// ReadActiveVehicles return active problems
func ReadResolveProblems() ([]*model.Problem, error) {
	problemCollection, session := problemCollection()
	defer session.Close()

	problems := []*model.Problem{}
	err := problemCollection.Pipe([]bson.M{
		{"$match": bson.M{"status": 0}},
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
	}).All(&problems)

	return problems, err
}
