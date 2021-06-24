package v1

import (
	"fmt"
	"testing"

	"../../config"
	"../../model"
	"../../service/driverLocationService"
	"../../service/notificationService"
	"../../service/orderService"

	"gopkg.in/mgo.v2/bson"
)

var testOrderID = bson.ObjectIdHex("5a46b7f8a118744b8ca4e6b3")
var testDriverID = bson.ObjectIdHex("5a0c3abba1187403c4d4ed3d")

func TestOrderRequest(t *testing.T) {
	// test order
	order, _ := orderService.ReadOrder(testOrderID)
	data := M{
		"type":    config.OrderRequest,
		"orderId": order.ID,
	}
	notificationService.PushWebsocketNotification(order.BusinessID.Hex(), data)
}

func TestAcceptOrderFromBusiness(t *testing.T) {
	// business will call this api
	order := &model.Order{
		ID:          testOrderID,
		OrderStatus: config.OrderAccepted,
	}
	order, _ = orderService.UpdateOrderProcess(order.ID, order)
	runProcessOrderForBusiness(order)
}

func TestDeclineOrderFromBusiness(t *testing.T) {
	// business will call this api
	order := &model.Order{
		ID:          testOrderID,
		OrderStatus: config.OrderDeclined,
		ReasonCode:  300,
	}
	order, _ = orderService.UpdateOrderProcess(order.ID, order)
	runProcessOrderForBusiness(order)
}

func TestCancelOrderFromBusiness(t *testing.T) {
	// business will call this api
	order := &model.Order{
		ID:          testOrderID,
		OrderStatus: config.OrderCancelled,
	}
	order, _ = orderService.UpdateOrderProcess(order.ID, order)
	runProcessOrderForBusiness(order)
}

func TestPrepareOrderFromBusiness(t *testing.T) {
	// business will call this api
	order := &model.Order{
		ID:          testOrderID,
		OrderStatus: config.OrderPrepared,
	}
	order, _ = orderService.UpdateOrderProcess(order.ID, order)
	runProcessOrderForBusiness(order)
}

func TestCompleteOrderFromBusiness(t *testing.T) {
	// business will call this api
	order := &model.Order{
		ID:          testOrderID,
		OrderStatus: config.OrderCompleted,
	}
	order, _ = orderService.UpdateOrderProcess(order.ID, order)
	runProcessOrderForBusiness(order)
}

func TestReadNearDrivers(t *testing.T) {
	lat := 1.1
	lng := 2.2

	driverLocations, _ := driverLocationService.GetNearbyDrivers(lat, lng, 20)
	fmt.Println(driverLocations)
}

func TestRequestTripToDriver(t *testing.T) {
	order := &model.Order{
		ID:         testOrderID,
		DriverID:   testDriverID,
		TripStatus: config.TripRequest,
	}
	order, _ = orderService.UpdateOrderProcess(order.ID, order)
	runRequestTripToDriver(order)
}

func TestAcceptTripFromDriver(t *testing.T) {
	// business will call this api
	order := &model.Order{
		ID:         testOrderID,
		DriverID:   testDriverID,
		TripStatus: config.TripAccepted,
	}
	order, _ = orderService.UpdateOrderProcess(order.ID, order)
	runProcessTripForDriver(order)
}

func TestDeclineTripFromDriver(t *testing.T) {
	// business will call this api
	order := &model.Order{
		ID:         testOrderID,
		DriverID:   testDriverID,
		TripStatus: config.TripDeclined,
	}
	order, _ = orderService.UpdateOrderProcess(order.ID, order)
	runProcessTripForDriver(order)
}

func TestCancelTripFromDriver(t *testing.T) {
	// business will call this api
	order := &model.Order{
		ID:         testOrderID,
		DriverID:   testDriverID,
		TripStatus: config.TripCancelled,
		ReasonCode: 300,
	}
	order, _ = orderService.UpdateOrderProcess(order.ID, order)
	runProcessTripForDriver(order)
}

func TestConfirmTripFromDriver(t *testing.T) {
	// business will call this api
	order := &model.Order{
		ID:         testOrderID,
		TripStatus: config.TripConfirmed,
	}
	order, _ = orderService.UpdateOrderProcess(order.ID, order)
	runProcessTripForDriver(order)
}

func TestStartTripFromDriver(t *testing.T) {
	// business will call this api
	order := &model.Order{
		ID:         testOrderID,
		TripStatus: config.TripStarted,
	}
	order, _ = orderService.UpdateOrderProcess(order.ID, order)
	runProcessTripForDriver(order)
}

func TestArriveTripFromDriver(t *testing.T) {
	// business will call this api
	order := &model.Order{
		ID:         testOrderID,
		TripStatus: config.TripArrived,
	}
	order, _ = orderService.UpdateOrderProcess(order.ID, order)
	runProcessTripForDriver(order)
}

func TestDropTripFromDriver(t *testing.T) {
	// business will call this api
	order := &model.Order{
		ID:         testOrderID,
		TripStatus: config.TripDropped,
	}
	order, _ = orderService.UpdateOrderProcess(order.ID, order)
	runProcessTripForDriver(order)
}

func TestCompleteTripFromDriver(t *testing.T) {
	// business will call this api
	order := &model.Order{
		ID:         testOrderID,
		TripStatus: config.TripCompleted,
	}
	order, _ = orderService.UpdateOrderProcess(order.ID, order)
	runProcessTripForDriver(order)
}
