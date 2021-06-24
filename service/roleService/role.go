package roleService

import (
	"errors"
	"fmt"

	"../../config"
	"../../db"
	"../../model"
	"../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func roleCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("role"), session
}

// InitService inits service
func InitService() {

	// create script
	CreateRole(&model.Role{Name: "Administrator", Code: config.AdminCode})
	CreateRole(&model.Role{Name: "Developer", Code: config.DeveloperCode})
	CreateRole(&model.Role{Name: "Marketer", Code: config.MarketerCode})
	CreateRole(&model.Role{Name: "Reporter", Code: config.ReportsCode})
}

// CreateRole creates role
func CreateRole(role *model.Role) (*model.Role, error) {
	roleCollection, session := roleCollection()
	defer session.Close()

	// Check if url is existed already
	result := &model.Role{}
	if err := roleCollection.Find(bson.M{"code": role.Code}).One(result); err == nil && result.ID != "" {
		return nil, errors.New("This role is registered already")
	}
	// Create url with intialize data
	role.ID = bson.NewObjectId()
	if role.Code == 0 {
		role.Code = createCode()
	}
	role.CreatedAt = timeHelper.GetCurrentTime()
	role.UpdatedAt = timeHelper.GetCurrentTime()
	// Insert Data
	err := roleCollection.Insert(role)

	return role, err
}

// ReadRole returns role with object id
func ReadRole(objid bson.ObjectId) (*model.Role, error) {
	roleCollection, session := roleCollection()
	defer session.Close()

	role := &model.Role{}
	// Find role with object id
	err := roleCollection.FindId(objid).One(role)
	return role, err
}

// ReadRoleWithCode returns role with code
func ReadRoleWithCode(code int) (*model.Role, error) {
	roleCollection, session := roleCollection()
	defer session.Close()

	role := &model.Role{}
	// Find role with object id
	err := roleCollection.Find(bson.M{"code": code}).One(role)
	return role, err
}

// UpdateRole updates role
func UpdateRole(objid bson.ObjectId, role *model.Role) (*model.Role, error) {
	roleCollection, session := roleCollection()
	defer session.Close()

	fmt.Println(role.URLGroup)
	role.UpdatedAt = timeHelper.GetCurrentTime()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"name":      role.Name,
			"urlGroup":  role.URLGroup,
			"updatedAt": role.UpdatedAt,
		}},
		ReturnNew: true,
	}
	// Update role
	_, err := roleCollection.FindId(objid).Apply(change, role)
	return role, err
}

// DeleteRole deletes role with object id
func DeleteRole(objid bson.ObjectId) error {
	roleCollection, session := roleCollection()
	defer session.Close()

	err := roleCollection.RemoveId(objid)
	return err
}

// ReadRoles return roles after search query
func ReadRoles(query string, offset int, count int, field string, sort int) ([]*model.Role, int, error) {
	roleCollection, session := roleCollection()
	defer session.Close()

	roles := []*model.Role{}
	totalCount := 0
	pipe := []bson.M{}
	if query == "" {
		// Get all riders
		totalCount, _ = roleCollection.Find(bson.M{}).Count()
	} else {
		// Search rider by query
		param := bson.M{"$or": []interface{}{
			bson.M{"name": bson.RegEx{Pattern: query, Options: "i"}},
		}}
		totalCount, _ = roleCollection.Find(param).Count()
		pipe = append(pipe, bson.M{"$match": param})
	}

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
	err := roleCollection.Pipe(pipe).All(&roles)

	return roles, totalCount, err
}

func createCode() int {
	roleCollection, session := roleCollection()
	defer session.Close()

	role := &model.Role{}
	err := roleCollection.Find(bson.M{}).Sort("-code").Limit(1).One(role)
	if err != nil || role.Code < config.MinRoleCode {
		return config.MinRoleCode
	}
	return role.Code + 1
}

// ReadURLGroupsWithCode returns urlgroups of role code
func ReadURLGroupsWithCode(code int) ([]*model.URLGroup, error) {
	roleCollection, session := roleCollection()
	defer session.Close()

	urlGroups := []*model.URLGroup{}
	err := roleCollection.Pipe([]bson.M{
		{"$match": bson.M{"code": code}},
		{"$unwind": "$urlGroup"},
		{"$lookup": bson.M{
			"from":         "url_group",
			"foreignField": "_id",
			"localField":   "urlGroup",
			"as":           "urlGroup",
		}},
		{"$unwind": "$urlGroup"},
		{"$replaceRoot": bson.M{"newRoot": "$urlGroup"}},
	}).All(&urlGroups)
	return urlGroups, err
}
