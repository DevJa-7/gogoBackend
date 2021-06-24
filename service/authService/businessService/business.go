package businessService

import (
	"errors"
	"fmt"
	"log"
	"time"

	"../../../config"
	"../../../db"
	"../../../model"
	"../../../service/base/dietaryService"
	"../../../service/base/mealKindService"
	"../../../util/crypto"
	"../../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var basePipe []bson.M

func businessCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("business"), session
}

// InitService inits service
func InitService() {
	businessCollection, session := businessCollection()
	defer session.Close()

	basePipe = []bson.M{}

	// indexing location
	index := mgo.Index{
		Key:        []string{"$2dsphere:geoLocation.geoJson"},
		Background: true,
	}
	err := businessCollection.EnsureIndex(index)
	if err != nil {
		log.Println(err)
	}
}

// CreateBusiness creates a business
func CreateBusiness(business *model.Business) (*model.Business, error) {
	businessCollection, session := businessCollection()
	defer session.Close()
	// check duplicate email
	if c, _ := businessCollection.Find(bson.M{"email": business.Email}).Count(); c > 0 {
		return nil, errors.New("Same email is registered already")
	}
	business.ID = bson.NewObjectId()
	business.Password = crypto.GenerateHash(business.Password)
	business.GeoLocation.GeoJSON.Type = "Point"
	business.GeoLocation.GeoJSON.Coordinates = []float64{0, 0}
	business.Schedules = []*model.Schedule{}
	for i := 0; i < 7; i++ {
		weekday := time.Weekday(i)
		schedule := &model.Schedule{
			Name:      weekday.String(),
			Weekday:   i,
			OpenTime:  config.ScheduleFromTime,
			CloseTime: config.ScheduleToTime,
			Enabled:   true,
		}
		business.Schedules = append(business.Schedules, schedule)
	}
	business.CreatedAt = timeHelper.GetCurrentTime()
	business.UpdatedAt = timeHelper.GetCurrentTime()

	// Insert Data
	err := businessCollection.Insert(business)
	return business, err
}

// ReadBusiness reads a business
func ReadBusiness(objid bson.ObjectId) (*model.Business, error) {
	businessCollection, session := businessCollection()
	defer session.Close()

	business := &model.Business{}
	// Read Data
	err := businessCollection.FindId(objid).One(&business)
	return business, err
}

// UpdateBusiness reads a business
func UpdateBusiness(objid bson.ObjectId, business *model.Business) (*model.Business, error) {
	businessCollection, session := businessCollection()
	defer session.Close()
	// check duplicate email
	if c, _ := businessCollection.Find(bson.M{"_id": bson.M{"$ne": objid}, "email": business.Email}).Count(); c > 0 {
		return nil, errors.New("Same email is registered already")
	}
	updateSet := bson.M{
		"email":             business.Email,
		"logo":              business.Logo,
		"identification":    business.Identification,
		"countryCode":       business.CountryCode,
		"phone":             business.Phone,
		"name":              business.Name,
		"description":       business.Description,
		"priceLevel":        business.PriceLevel,
		"preparationTime":   business.PreparationTime,
		"geoLocation":       business.GeoLocation,
		"bankInfo":          business.BankInfo,
		"dietaryCodes":      business.DietaryCodes,
		"mealKindCodes":     business.MealKindCodes,
		"website":           business.Website,
		"schedules":         business.Schedules,
		"closed":            business.Closed,
		"mostPopular":       business.MostPopular,
		"recommend":         business.Recommend,
		"verify.isVerified": business.Verify.IsVerified,
		"updatedAt":         timeHelper.GetCurrentTime(),
	}
	// Create change info
	change := mgo.Change{
		Update:    bson.M{"$set": updateSet},
		ReturnNew: true,
	}
	_, err := businessCollection.FindId(objid).Apply(change, business)
	return business, err
}

// DeleteBusiness deletes business with object id
func DeleteBusiness(objid bson.ObjectId) error {
	businessCollection, session := businessCollection()
	defer session.Close()

	err := businessCollection.RemoveId(objid)
	return err
}

// ReadBusinesses return businesses after retreive with params
func ReadBusinesses(query string, offset int, count int, field string, sort int) ([]*model.PublicBusiness, int, error) {
	businessCollection, session := businessCollection()
	defer session.Close()

	businesses := []*model.PublicBusiness{}
	pipe := []bson.M{}
	if query != "" {
		// Search business by query
		param := bson.M{"$or": []interface{}{
			bson.M{"email": bson.RegEx{Pattern: query, Options: ""}},
			bson.M{"name": bson.RegEx{Pattern: query, Options: ""}},
		}}
		pipe = append(pipe, bson.M{"$match": param})
	}
	// get total count of collection with initial query
	totalCount := db.GetCountOfCollection(businessCollection, &pipe)

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

	err := businessCollection.Pipe(pipe).All(&businesses)

	return businesses, totalCount, err
}

// RetrieveBusinessBaseStructure retrieve base datas
func RetrieveBusinessBaseStructure(b *model.PublicBusiness) {
	b.Dietaries = dietaryService.ReadDietariesWithCodes(b.DietaryCodes)
	b.MealKinds = mealKindService.ReadMealKindsWithCodes(b.MealKindCodes)
}

// ReadNearbyPopularBusiness returns nearby popular business
func ReadNearbyPopularBusiness(lat, lng float64) ([]*model.PublicBusiness, error) {
	businessCollection, session := businessCollection()
	defer session.Close()

	businesses := []*model.PublicBusiness{}
	err := businessCollection.Pipe([]bson.M{
		{"$geoNear": bson.M{
			"near":          bson.M{"type": "Point", "coordinates": []float64{lng, lat}},
			"distanceField": "distance",
			"maxDistance":   10000,
			"query":         bson.M{"closed": false, "mostPopular": true},
			"includeLocs":   "geoLocation.geoJson",
			"num":           10,
			"spherical":     true,
		}},
	}).All(&businesses)
	return businesses, err
}

// ReadRecommendBusiness read recommend business
func ReadRecommendBusiness(lat, lng float64) ([]*model.PublicBusiness, error) {
	businessCollection, session := businessCollection()
	defer session.Close()

	businesses := []*model.PublicBusiness{}
	err := businessCollection.Pipe([]bson.M{
		{"$geoNear": bson.M{
			"near":          bson.M{"type": "Point", "coordinates": []float64{lng, lat}},
			"distanceField": "distance",
			"maxDistance":   10000,
			"query":         bson.M{"closed": false, "recommend": true},
			"includeLocs":   "geoLocation.geoJson",
			"num":           10,
			"spherical":     true,
		}},
	}).All(&businesses)
	return businesses, err
}

// ReadUnder30Business returns businesses under 30mins
func ReadUnder30Business(lat, lng float64) ([]*model.PublicBusiness, error) {
	businessCollection, session := businessCollection()
	defer session.Close()

	businesses := []*model.PublicBusiness{}
	err := businessCollection.Pipe([]bson.M{
		{"$geoNear": bson.M{
			"near":          bson.M{"type": "Point", "coordinates": []float64{lng, lat}},
			"distanceField": "distance",
			"maxDistance":   10000,
			"query":         bson.M{"closed": false, "mostPopular": false, "recommend": false},
			"includeLocs":   "geoLocation.geoJson",
			"num":           10,
			"spherical":     true,
		}},
	}).All(&businesses)
	return businesses, err
}

// ReadMoreBusiness returns other
func ReadMoreBusiness(lat, lng float64, excepts []bson.ObjectId) ([]*model.PublicBusiness, error) {
	businessCollection, session := businessCollection()
	defer session.Close()

	businesses := []*model.PublicBusiness{}
	err := businessCollection.Pipe([]bson.M{
		{"$geoNear": bson.M{
			"near":          bson.M{"type": "Point", "coordinates": []float64{lng, lat}},
			"distanceField": "distance",
			"maxDistance":   10000,
			"query":         bson.M{"_id": bson.M{"$nin": excepts}},
			"includeLocs":   "geoLocation.geoJson",
			"num":           100,
			"spherical":     true,
		}},
	}).All(&businesses)
	return businesses, err
}

// ReadQueryBusiness returns other
func ReadQueryBusiness(queryBusiness *model.QueryBusiness) ([]*model.PublicBusiness, error) {
	businessCollection, session := businessCollection()
	defer session.Close()

	businesses := []*model.PublicBusiness{}
	query := bson.M{}
	if queryBusiness.Sort == config.SortPopular {
		query["mostPopular"] = true
	}

	if len(queryBusiness.Price) > 0 {
		query["priceLevel"] = bson.M{"$in": queryBusiness.Price}
	}
	if len(queryBusiness.Dietary) > 0 {
		query["dietaryCodes"] = bson.M{"$in": queryBusiness.Dietary}
	}
	fmt.Println(query)

	pipe := []bson.M{
		{"$geoNear": bson.M{
			"near":          bson.M{"type": "Point", "coordinates": []float64{queryBusiness.Lng, queryBusiness.Lat}},
			"distanceField": "distance",
			"maxDistance":   10000,
			"query":         query,
			"includeLocs":   "geoLocation.geoJson",
			"num":           100,
			"spherical":     true,
		}},
	}
	if queryBusiness.Sort == config.SortDeliveryTime {
		pipe = append(pipe, bson.M{
			"$sort": bson.M{"distance": 1},
		})
	}

	err := businessCollection.Pipe(pipe).All(&businesses)
	return businesses, err
}

// ReadCounts reads total user count and available count
func ReadCounts() (int, int) {
	businessCollection, session := businessCollection()
	defer session.Close()

	total, _ := businessCollection.Find(bson.M{}).Count()
	available, _ := businessCollection.Find(bson.M{"verify.isVerified": true}).Count()
	return total, available
}
