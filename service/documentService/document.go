package documentService

import (
	"../../config"
	"../../db"
	"../../model"
	"../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func documentCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("document"), session
}

// InitService inits service
func InitService() {
}

// CreateDocument creates document
func CreateDocument(document *model.Document) (*model.Document, error) {
	documentCollection, session := documentCollection()
	defer session.Close()

	// Create document with intialize data
	document.ID = bson.NewObjectId()
	document.CreatedAt = timeHelper.GetCurrentTime()
	document.UpdatedAt = timeHelper.GetCurrentTime()
	// Insert Data
	err := documentCollection.Insert(document)
	return document, err
}

// ReadDocument returns document with object id
func ReadDocument(objid bson.ObjectId) (*model.Document, error) {
	documentCollection, session := documentCollection()
	defer session.Close()

	document := &model.Document{}
	// Find document with object id
	err := documentCollection.FindId(objid).One(document)
	return document, err
}

// UpdateDocument updates document
func UpdateDocument(objid bson.ObjectId, document *model.Document) (*model.Document, error) {
	documentCollection, session := documentCollection()
	defer session.Close()

	document.UpdatedAt = timeHelper.GetCurrentTime()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"type":       document.Type,
			"name":       document.Name,
			"isExpired":  document.IsExpired,
			"isRequired": document.IsRequired,
			"valid":      document.Valid,
			"updatedAt":  document.UpdatedAt,
		}},
		ReturnNew: true,
	}

	// Update document
	_, err := documentCollection.FindId(objid).Apply(change, document)
	return document, err
}

// DeleteDocument deletes document with object id
func DeleteDocument(objid bson.ObjectId) error {
	documentCollection, session := documentCollection()
	defer session.Close()

	return documentCollection.RemoveId(objid)
}

// ReadDocuments return documents after search query
func ReadDocuments(query string, field string, sort int, filter bool) (map[string][]*model.Document, error) {
	documentCollection, session := documentCollection()
	defer session.Close()

	type Documents struct {
		Type      int
		Documents []*model.Document
	}
	documents := []*Documents{}
	pipe := []bson.M{}
	// Search document by query
	if filter {
		pipe = append(pipe, bson.M{"$match": bson.M{"valid": true}})
	}
	param := bson.M{"$or": []interface{}{
		bson.M{"name": bson.RegEx{Pattern: query, Options: "i"}},
	}}
	pipe = append(pipe, bson.M{"$match": param})
	// add sort feature
	if field != "" && sort != 0 {
		pipe = append(pipe, bson.M{"$sort": bson.M{field: sort}})
	}
	pipe = append(pipe, bson.M{"$group": bson.M{"_id": "$type", "documents": bson.M{"$push": "$$ROOT"}}})
	pipe = append(pipe, bson.M{"$project": bson.M{"type": "$_id", "documents": 1}})
	err := documentCollection.Pipe(pipe).All(&documents)
	if err != nil {
		return nil, err
	}

	result := map[string][]*model.Document{}
	result["vehicle"] = []*model.Document{}
	result["driver"] = []*model.Document{}
	for _, d := range documents {
		if d.Type == config.DocTypeVehicle {
			result["vehicle"] = d.Documents
		} else if d.Type == config.DocTypeDriver {
			result["driver"] = d.Documents
		}
	}

	return result, err
}
