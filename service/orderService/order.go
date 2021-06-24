package orderService

import (
	"../../config"
	"../../db"
	"../../model"
	"../../util/random"
	"../../util/timeHelper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var basePipe []bson.M

func orderCollection() (*mgo.Collection, *mgo.Session) {
	mgoDB, session := db.MongoDB()
	return mgoDB.C("order"), session
}

// InitService inits service
func InitService() {
	basePipe = []bson.M{
		{"$lookup": bson.M{
			"from":         "user",
			"localField":   "userId",
			"foreignField": "_id",
			"as":           "user",
		}},
		{"$unwind": bson.M{
			"path": "$user",
			"preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         "business",
			"localField":   "businessId",
			"foreignField": "_id",
			"as":           "business",
		}},
		{"$unwind": bson.M{
			"path": "$business",
			"preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         "driver",
			"localField":   "driverId",
			"foreignField": "_id",
			"as":           "driver",
		}},
		{"$unwind": bson.M{
			"path": "$driver",
			"preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         "reason",
			"localField":   "reasonCode",
			"foreignField": "code",
			"as":           "reason",
		}},
		{"$unwind": bson.M{
			"path": "$reason",
			"preserveNullAndEmptyArrays": true}},
	}
}

// CreateOrder creates order
func CreateOrder(order *model.Order) (*model.Order, error) {
	orderCollection, session := orderCollection()
	defer session.Close()

	// Create url with intialize data
	order.ID = bson.NewObjectId()
	order.Number = random.GenerateRandomString(6)
	order.OrderStatus = config.OrderRequest
	order.CreatedAt = timeHelper.GetCurrentTime()
	order.UpdatedAt = timeHelper.GetCurrentTime()
	order.StatusAt = map[string]int64{
		config.OrderRequest: timeHelper.GetCurrentTime(),
	}
	// Insert Data
	err := orderCollection.Insert(order)

	return order, err
}

// ReadOrder returns order with object id
func ReadOrder(objid bson.ObjectId) (*model.Order, error) {
	orderCollection, session := orderCollection()
	defer session.Close()

	pipe := []bson.M{{"$match": bson.M{
		"_id": objid,
	}}}
	pipe = append(pipe, basePipe...)
	order := &model.Order{}
	// Find order with object id
	err := orderCollection.Pipe(pipe).One(order)

	return order, err
}

// UpdateOrder updates order
func UpdateOrder(objid bson.ObjectId, order *model.Order) (*model.Order, error) {
	orderCollection, session := orderCollection()
	defer session.Close()

	order.UpdatedAt = timeHelper.GetCurrentTime()
	// Create change info
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"updatedAt": order.UpdatedAt,
		}},
		ReturnNew: true,
	}
	// Update order
	_, err := orderCollection.FindId(objid).Apply(change, order)
	return order, err
}

// UpdateOrderProcess update order status
func UpdateOrderProcess(objid bson.ObjectId, order *model.Order) (*model.Order, error) {
	orderCollection, session := orderCollection()
	defer session.Close()

	order.UpdatedAt = timeHelper.GetCurrentTime()
	set := bson.M{
		"updatedAt": order.UpdatedAt,
	}

	if len(order.OrderStatus) > 0 {
		set["orderStatus"] = order.OrderStatus
		set["statusAt."+order.OrderStatus] = order.UpdatedAt
	}
	if len(order.TripStatus) > 0 {
		set["tripStatus"] = order.TripStatus
		set["statusAt."+order.TripStatus] = order.UpdatedAt
	}
	if order.ReasonCode > 0 {
		set["reasonCode"] = order.ReasonCode
	}
	if len(order.DriverID) > 0 {
		set["driverId"] = order.DriverID
	}
	change := mgo.Change{
		Update:    bson.M{"$set": set},
		ReturnNew: true,
	}
	// Update order
	_, err := orderCollection.FindId(objid).Apply(change, order)
	return order, err
}

// UpdateOrderPickupScore updates order pickup score
func UpdateOrderPickupScore(objid bson.ObjectId, order *model.Order) (*model.Order, error) {
	orderCollection, session := orderCollection()
	defer session.Close()

	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"score":     order.PickupScore,
			"updatedAt": timeHelper.GetCurrentTime(),
		}},
		ReturnNew: true,
	}
	// Update order
	_, err := orderCollection.FindId(objid).Apply(change, order)
	return order, err
}

// DeleteOrder deletes order with object id
func DeleteOrder(objid bson.ObjectId) error {
	orderCollection, session := orderCollection()
	defer session.Close()

	err := orderCollection.RemoveId(objid)
	return err
}

// ReadOrders return orders after search query
func ReadOrders(query string, offset int, count int, field string, sort int, businessID, userID, driverID bson.ObjectId, status string, past, upcoming bool, rated int) ([]*model.Order, int, error) {
	orderCollection, session := orderCollection()
	defer session.Close()

	orders := []*model.Order{}
	totalCount := 0
	pipe := []bson.M{}

	if businessID != "" {
		pipe = append(pipe, bson.M{"$match": bson.M{"businessId": businessID}})
	}
	if userID != "" {
		pipe = append(pipe, bson.M{"$match": bson.M{"userId": userID}})
	}
	if driverID != "" {
		pipe = append(pipe, bson.M{"$match": bson.M{"driverId": driverID}})
	}
	if status != "" {
		pipe = append(pipe, bson.M{"$match": bson.M{"orderStatus": status}})
	}
	if upcoming {
		pipe = append(pipe, bson.M{"$match": bson.M{"orderStatus": bson.M{"$in": []string{
			config.OrderRequest, config.OrderAccepted, config.OrderPrepared,
		}}}})
	}
	if past {
		pipe = append(pipe, bson.M{"$match": bson.M{"orderStatus": bson.M{"$in": []string{
			config.OrderDeclined, config.OrderCancelled, config.OrderCompleted,
		}}}})
	}
	if rated > 0 {
		if rated == config.Rated {
			pipe = append(pipe, bson.M{"$match": bson.M{"rated": true}})
		} else if rated == config.NoRated {
			pipe = append(pipe, bson.M{"$match": bson.M{"rated": false}})
		}
	}

	pipe = append(pipe, basePipe...)

	if query == "" {
		// Get all riders
		totalCount, _ = orderCollection.Find(bson.M{}).Count()
	} else {
		// Search rider by query
		param := bson.M{"$or": []interface{}{
			bson.M{"number": bson.RegEx{Pattern: query, Options: "i"}},
		}}
		totalCount, _ = orderCollection.Find(param).Count()
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
	err := orderCollection.Pipe(pipe).All(&orders)

	return orders, totalCount, err
}

// UpdateOrderRate updates order rate by user
func UpdateOrderRate(order *model.Order) error {
	orderCollection, session := orderCollection()
	defer session.Close()

	err := orderCollection.Update(bson.M{"_id": order.ID, "userId": order.UserID}, bson.M{
		"$set": bson.M{
			"rated":            true,
			"orderStatus":      config.OrderCompleted,
			"businessRate":     order.BusinessRate,
			"businessFeedback": order.BusinessFeedback,
			"driverRate":       order.DriverRate,
			"driverFeedback":   order.DriverFeedback,
		}},
	)
	return err
}

// NoRatedOrders read not rated orders of user
func NoRatedOrders(objid bson.ObjectId) ([]*model.Order, error) {
	orderCollection, session := orderCollection()
	defer session.Close()

	orders := []*model.Order{}
	t := timeHelper.FewDaysLater(-7)
	pipe := []bson.M{
		{"$match": bson.M{
			"userId":    objid,
			"rated":     false,
			"updatedAt": bson.M{"$ge": t},
		}},
	}
	pipe = append(pipe, basePipe...)
	err := orderCollection.Pipe(pipe).All(&orders)

	return orders, err
}
