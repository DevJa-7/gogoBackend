package colorService

import (
	"errors"

	"../../../db"
	"../../../model"
	"../../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func colorCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("color"), session
}

// InitService inits service
func InitService() {
	// create script
	CreateColor(&model.Color{Name: "Red", Value: "#ff0000"})
	CreateColor(&model.Color{Name: "Black", Value: "#000000"})
	CreateColor(&model.Color{Name: "Yellow", Value: "#ffff00"})
	CreateColor(&model.Color{Name: "Orange", Value: "#ff8000"})
}

// CreateColor creates color
func CreateColor(color *model.Color) (*model.Color, error) {
	colorCollection, session := colorCollection()
	defer session.Close()

	// Check if url is existed already
	result := &model.Color{}
	if err := colorCollection.Find(bson.M{"name": color.Name}).One(result); err == nil && result.ID != "" {
		return nil, errors.New("This color is registered already")
	}
	// Create url with intialize data
	color.ID = bson.NewObjectId()
	color.CreatedAt = timeHelper.GetCurrentTime()
	color.UpdatedAt = timeHelper.GetCurrentTime()
	// Insert Data
	err := colorCollection.Insert(color)

	return color, err
}

// ReadColor returns color with object id
func ReadColor(objid bson.ObjectId) (*model.Color, error) {
	colorCollection, session := colorCollection()
	defer session.Close()

	color := &model.Color{}
	// Find admin with object id
	err := colorCollection.FindId(objid).One(color)
	return color, err
}

// UpdateColor updates color
func UpdateColor(objid bson.ObjectId, color *model.Color) (*model.Color, error) {
	colorCollection, session := colorCollection()
	defer session.Close()

	color.UpdatedAt = timeHelper.GetCurrentTime()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"name":      color.Name,
			"value":     color.Value,
			"updatedAt": color.UpdatedAt,
		}},
		ReturnNew: true,
	}
	// Update color
	_, err := colorCollection.FindId(objid).Apply(change, color)
	return color, err
}

// DeleteColor deletes color with object id
func DeleteColor(objid bson.ObjectId) error {
	colorCollection, session := colorCollection()
	defer session.Close()

	err := colorCollection.RemoveId(objid)
	return err
}

// ReadColors return colors after search query
func ReadColors() ([]*model.Color, error) {
	colorCollection, session := colorCollection()
	defer session.Close()

	colors := []*model.Color{}
	err := colorCollection.Find(bson.M{}).All(&colors)

	return colors, err
}
