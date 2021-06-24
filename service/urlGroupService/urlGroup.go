package urlGroupService

import (
	"errors"

	"../../db"
	"../../model"
	"../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func urlGroupCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("url_group"), session
}

// InitService inits service
func InitService() {
	// init database
	CreateURLGroup(&model.URLGroup{Name: "Dashboard", URL: "/dashboard"})
	CreateURLGroup(&model.URLGroup{Name: "Administrators", URL: "/admins"})
	CreateURLGroup(&model.URLGroup{Name: "Businesses", URL: "/businesses"})
	CreateURLGroup(&model.URLGroup{Name: "Drivers", URL: "/drivers"})
	CreateURLGroup(&model.URLGroup{Name: "users", URL: "/users"})
	CreateURLGroup(&model.URLGroup{Name: "Locations", URL: "/locations"})
	CreateURLGroup(&model.URLGroup{Name: "Vehicles", URL: "/vehicles"})
	CreateURLGroup(&model.URLGroup{Name: "Urls", URL: "/urls"})
	CreateURLGroup(&model.URLGroup{Name: "Roles", URL: "/roles"})
	CreateURLGroup(&model.URLGroup{Name: "Documents", URL: "/documents"})
	CreateURLGroup(&model.URLGroup{Name: "Brands", URL: "/brands"})
	CreateURLGroup(&model.URLGroup{Name: "Profile", URL: "/profile"})
	CreateURLGroup(&model.URLGroup{Name: "Configuration", URL: "/configuration"})
}

// CreateURLGroup creates urlGroup
func CreateURLGroup(urlGroup *model.URLGroup) (*model.URLGroup, error) {
	urlGroupCollection, session := urlGroupCollection()
	defer session.Close()

	// Check if url is existed already
	result := &model.URLGroup{}
	if err := urlGroupCollection.Find(bson.M{"url": urlGroup.URL}).One(result); err == nil && result.ID != "" {
		return nil, errors.New("This url is registered already")
	}
	// Create url with intialize data
	urlGroup.ID = bson.NewObjectId()
	urlGroup.CreatedAt = timeHelper.GetCurrentTime()
	urlGroup.UpdatedAt = timeHelper.GetCurrentTime()

	// Insert Data
	err := urlGroupCollection.Insert(urlGroup)
	return urlGroup, err
}

// ReadURLGroup return URLGroup with object id
func ReadURLGroup(objid bson.ObjectId) (*model.URLGroup, error) {
	urlGroupCollection, session := urlGroupCollection()
	defer session.Close()

	urlGroup := &model.URLGroup{}
	err := urlGroupCollection.FindId(objid).One(urlGroup)
	return urlGroup, err
}

// UpdateURLGroup updates urlGroup
func UpdateURLGroup(objid bson.ObjectId, urlGroup *model.URLGroup) (*model.URLGroup, error) {
	urlGroupCollection, session := urlGroupCollection()
	defer session.Close()

	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"name":      urlGroup.Name,
			"url":       urlGroup.URL,
			"updatedAt": timeHelper.GetCurrentTime(),
		}},
		ReturnNew: true,
	}

	// Update urlGroup
	_, err := urlGroupCollection.FindId(objid).Apply(change, &urlGroup)
	return urlGroup, err
}

// DeleteURLGroup deletes urlGroup by object id
func DeleteURLGroup(objid bson.ObjectId) error {
	urlGroupCollection, session := urlGroupCollection()
	defer session.Close()

	err := urlGroupCollection.RemoveId(objid)
	return err
}

// ReadURLGroupsByGroup returns urls with user role
func ReadURLGroupsByGroup(groups []bson.ObjectId) interface{} {
	urlGroupCollection, session := urlGroupCollection()
	defer session.Close()

	result := bson.M{}
	err := urlGroupCollection.Pipe([]bson.M{
		{"$match": bson.M{"_id": bson.M{"$in": groups}}},
		{"$group": bson.M{"_id": nil, "urls": bson.M{"$push": "$url"}}},
	}).One(&result)
	if err == nil {
		return result["urls"]
	}
	return []string{}
}

// ReadURLGroups return url groups with search query
func ReadURLGroups(query string, offset int, count int, field string, sort int) ([]*model.URLGroup, int, error) {
	urlGroupCollection, session := urlGroupCollection()
	defer session.Close()

	urls := []*model.URLGroup{}
	totalCount := 0
	pipe := []bson.M{}
	if query == "" {
		// Get all riders
		totalCount, _ = urlGroupCollection.Find(bson.M{}).Count()
	} else {
		// Search rider by query
		param := bson.M{"$or": []interface{}{
			bson.M{"name": bson.RegEx{Pattern: query, Options: "i"}},
			bson.M{"url": bson.RegEx{Pattern: query, Options: "i"}},
		}}
		totalCount, _ = urlGroupCollection.Find(param).Count()
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
	err := urlGroupCollection.Pipe(pipe).All(&urls)

	return urls, totalCount, err
}
