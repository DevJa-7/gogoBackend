package brandService

import (
	"errors"

	"../../../db"
	"../../../model"
	"../../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func brandCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("brand"), session
}

// InitService inits service
func InitService() {
	// create script
	CreateBrand(&model.Brand{Title: "Toyota", Models: []*model.Model{{1, "Prius"}, {2, "Tacoma"}, {3, "Corolla"}, {4, "Yaris"}, {5, "Camry Hybrid"}}})
	CreateBrand(&model.Brand{Title: "Hyundai", Models: []*model.Model{{1, "Elantra"}, {2, "Sonata"}, {3, "Azera"}}})
}

// CreateBrand creates brand
func CreateBrand(brand *model.Brand) (*model.Brand, error) {
	brandCollection, session := brandCollection()
	defer session.Close()

	// Check if url is existed already
	result := &model.Brand{}
	if err := brandCollection.Find(bson.M{"title": brand.Title}).One(result); err == nil && result.ID != "" {
		return nil, errors.New("This brand is registered already")
	}
	// Create url with intialize data
	brand.ID = bson.NewObjectId()
	brand.CreatedAt = timeHelper.GetCurrentTime()
	brand.UpdatedAt = timeHelper.GetCurrentTime()
	// Insert Data
	err := brandCollection.Insert(brand)

	return brand, err
}

// ReadBrand returns brand with object id
func ReadBrand(objid bson.ObjectId) (*model.Brand, error) {
	brandCollection, session := brandCollection()
	defer session.Close()

	brand := &model.Brand{}
	// Find admin with object id
	err := brandCollection.FindId(objid).One(brand)
	return brand, err
}

// UpdateBrand updates brand
func UpdateBrand(objid bson.ObjectId, brand *model.Brand) (*model.Brand, error) {
	brandCollection, session := brandCollection()
	defer session.Close()

	brand.UpdatedAt = timeHelper.GetCurrentTime()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"title":     brand.Title,
			"updatedAt": brand.UpdatedAt,
		}},
		ReturnNew: true,
	}
	// Update brand
	_, err := brandCollection.FindId(objid).Apply(change, brand)
	return brand, err
}

// DeleteBrand deletes brand with object id
func DeleteBrand(objid bson.ObjectId) error {
	brandCollection, session := brandCollection()
	defer session.Close()

	err := brandCollection.RemoveId(objid)
	return err
}

// ReadBrands return brands after search query
func ReadBrands() ([]*model.Brand, error) {
	brandCollection, session := brandCollection()
	defer session.Close()

	brands := []*model.Brand{}
	err := brandCollection.Find(bson.M{}).All(&brands)

	return brands, err
}

// UpdateModels updates models
func UpdateModels(objid bson.ObjectId, models []*model.Model) error {
	brandCollection, session := brandCollection()
	defer session.Close()

	// Create change models
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"models":    models,
			"updatedAt": timeHelper.GetCurrentTime(),
		}},
		ReturnNew: true,
	}

	// Update location
	_, err := brandCollection.FindId(objid).Apply(change, nil)
	return err
}
