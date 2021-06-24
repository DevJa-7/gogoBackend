package v1

import (
	"net/http"
	"strconv"

	"../../api/response"
	"../../config"
	"../../model"
	"../../service/authService/driverService"
	"../../service/authService/permission"
	"../../service/authService/userService"
	"../../service/driverLocationService"
	"../../service/notificationService"
	"../../service/orderService"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

type M map[string]interface{}

// InitOrder inits order CRUD apis
// @Title Orders
// @Description Orders's router group.
func InitOrder(parentRoute *echo.Group) {
	route := parentRoute.Group("/orders")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createOrder))
	route.GET("/:id", permission.AuthRequired(readOrder))
	route.PUT("/:id", permission.AuthRequired(updateOrder))
	route.DELETE("/:id", permission.AuthRequired(deleteOrder))

	route.GET("", permission.AuthRequired(readOrders))

	route.POST("/business/process", permission.AuthRequired(processOrderForBusiness))
	route.POST("/trip/process", permission.AuthRequired(processTripForDriver))
	route.GET("/trip/nearDrivers", permission.AuthRequired(readNearDrivers))
	route.POST("/trip/request", permission.AuthRequired(requestTripToDriver))
	route.POST("/trip/pickup", permission.AuthRequired(submitPickupScore))

	route.POST("/submit/by/user", permission.AuthRequired(submitOrderByUser))
	route.POST("/submit/by/business", permission.AuthRequired(submitOrderByBusiness))
	route.POST("/submit/by/driver", permission.AuthRequired(submitOrderByDriver))

	orderService.InitService()
}

// @Title createOrder
// @Description Create a order.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   name       	form   	string  true	"Order name."
// @Success 200 {object} model.Order             "Returns created order"
// @Failure 400 {object} response.BasicResponse "err.order.bind"
// @Failure 400 {object} response.BasicResponse "err.order.create"
// @Resource /orders
// @Router /orders [post]
func createOrder(c echo.Context) error {
	order := &model.Order{}
	if err := c.Bind(order); err != nil {
		return response.KnownErrJSON(c, "err.order.bind", err)
	}
	// Create order
	order, err := orderService.CreateOrder(order)
	if err != nil {
		return response.KnownErrJSON(c, "err.order.create", err)
	}
	// Send message to business via websocket
	data := M{
		"type":    config.OrderRequest,
		"orderId": order.ID,
	}
	go notificationService.PushWebsocketNotification(order.BusinessID.Hex(), data)

	return response.SuccessInterface(c, order)
}

// @Title readOrder
// @Description Read a order.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Order ID."
// @Success 200 {object} model.Order 		"Returns read order"
// @Failure 400 {object} response.BasicResponse "err.order.bind"
// @Failure 400 {object} response.BasicResponse "err.order.read"
// @Resource /orders
// @Router /orders/{id} [get]
func readOrder(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	// Retrieve rider by id
	order, err := orderService.ReadOrder(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.order.read", err)
	}
	return response.SuccessInterface(c, order)
}

// @Title updateOrder
// @Description Update order of orders.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Order ID."
// @Success 200 {object} model.Order	 		"Returns updated order"
// @Failure 400 {object} response.BasicResponse "err.order.bind"
// @Failure 400 {object} response.BasicResponse "err.order.update"
// @Resource /orders
// @Router /orders/{id} [put]
func updateOrder(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	order := &model.Order{}
	if err := c.Bind(order); err != nil {
		return response.KnownErrJSON(c, "err.order.bind", err)
	}

	// Update order
	order, err := orderService.UpdateOrder(objid, order)
	if err != nil {
		return response.KnownErrJSON(c, "err.order.update", err)
	}
	return response.SuccessInterface(c, order)
}

// @Title deleteOrder
// @Description Delete a order.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"User ID."
// @Success 200 {object} response.BasicResponse "User is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.order.bind"
// @Failure 400 {object} response.BasicResponse "err.order.delete"
// @Resource /orders
// @Router /orders/{id} [delete]
func deleteOrder(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))

	// Remove rider with object id
	err := orderService.DeleteOrder(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.order.delete", err)
	}
	return response.SuccessJSON(c, "Order is deleted correctly.")
}

// @Title readOrders
// @Description Read a orders.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   query			form   	string  false	"Will search string."
// @Param   offset			form    int		false	"Offset for pagination."
// @Param   count 			form    int		false	"Count that will show per page."
// @Param   field			form    string  false	"Sort field."
// @Param   sort			form    int		false	"Sort direction. 0:default, 1:Ascending, -1:Descending"
// @Success 200 {object} response.BasicResponse "User is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.order.bind"
// @Failure 400 {object} response.BasicResponse "err.order.delete"
// @Resource /orders
// @Router /orders/{id} [delete]
func readOrders(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))
	var businessID bson.ObjectId
	if bson.IsObjectIdHex(c.FormValue("businessId")) {
		businessID = bson.ObjectIdHex(c.FormValue("businesId"))
	}
	var userID bson.ObjectId
	if bson.IsObjectIdHex(c.FormValue("userId")) {
		userID = bson.ObjectIdHex(c.FormValue("userId"))
	}
	var driverID bson.ObjectId
	if bson.IsObjectIdHex(c.FormValue("driverId")) {
		driverID = bson.ObjectIdHex(c.FormValue("driverId"))
	}
	status := c.FormValue("status")
	past, _ := strconv.ParseBool(c.FormValue("is_past"))
	upcoming, _ := strconv.ParseBool(c.FormValue("is_upcoming"))
	rated, _ := strconv.Atoi(c.FormValue("rated"))

	// Read orders with query
	orders, total, err := orderService.ReadOrders(query, offset, count, field, sort, businessID, userID, driverID, status, past, upcoming, rated)
	if err != nil {
		return response.KnownErrJSON(c, "err.order.read", err)
	}
	// return response.SuccessInterface(c, &model.ListForm{total, orders})
	return c.JSON(http.StatusOK, &model.ListForm{total, orders})
}

func processOrderForBusiness(c echo.Context) error {
	order := &model.Order{}
	if err := c.Bind(order); err != nil {
		return response.KnownErrJSON(c, "err.order.bind", err)
	}
	// Update order
	order, err := orderService.UpdateOrderProcess(order.ID, order)
	if err != nil {
		return response.KnownErrJSON(c, "err.order.update", err)
	}
	// run process
	go runProcessOrderForBusiness(order)

	return response.SuccessInterface(c, order)
}

func runProcessOrderForBusiness(order *model.Order) {
	user, _ := userService.ReadUser(order.UserID)

	notification := &model.OneSignalNotification{}
	notification.AppID = config.UserAppID
	notification.PlayerIds = []string{user.OneSignalPlayerID}
	notification.Title = M{"en": config.UserAppName}

	data := M{}
	data["type"] = order.OrderStatus
	data["orderId"] = order.ID
	notification.Data = data

	switch order.OrderStatus {
	case config.OrderAccepted:
		notification.Message = M{"en": "Your order is accepted."}
	case config.OrderDeclined:
		notification.Message = M{"en": "Your order is declined."}
	case config.OrderCancelled:
		notification.Message = M{"en": "Your order is cancelled."}
	case config.OrderPrepared:
		notification.Message = M{"en": "Your order is preparing."}
	case config.OrderCompleted:
		notification.Message = M{"en": "Your order is completed."}
		notificationService.PushWebsocketNotification(order.BusinessID.Hex(), data)
		return
	}

	notificationService.PushOneSignalNotification(notification, config.UserAPIKey)
}

func processTripForDriver(c echo.Context) error {
	order := &model.Order{}
	if err := c.Bind(order); err != nil {
		return response.KnownErrJSON(c, "err.order.bind", err)
	}
	// Update order
	order, err := orderService.UpdateOrderProcess(order.ID, order)
	if err != nil {
		return response.KnownErrJSON(c, "err.order.update", err)
	}
	// run process
	go runProcessTripForDriver(order)

	return response.SuccessInterface(c, order)
}

func runProcessTripForDriver(order *model.Order) {
	user, _ := userService.ReadUser(order.UserID)

	notification := &model.OneSignalNotification{}
	notification.AppID = config.UserAppID
	notification.PlayerIds = []string{user.OneSignalPlayerID}
	notification.Title = M{"en": config.UserAppName}

	data := M{}
	data["type"] = order.TripStatus
	data["orderId"] = order.ID
	notification.Data = data

	switch order.TripStatus {
	case config.TripAccepted:
		notification.Message = M{"en": "Your order trip is accepted."}
		notificationService.PushWebsocketNotification(order.BusinessID.Hex(), data)
		return
	case config.TripDeclined:
		notification.Message = M{"en": "Your order trip is declined."}
		notificationService.PushWebsocketNotification(order.BusinessID.Hex(), data)
		return
	case config.TripCancelled:
		notification.Message = M{"en": "Your order trip is cancelled."}
		notificationService.PushWebsocketNotification(order.BusinessID.Hex(), data)
		return
	case config.TripConfirmed:
		notification.Message = M{"en": "Your order trip is confirmed."}
	case config.TripStarted:
		notification.Message = M{"en": "Your order trip is started."}
	case config.TripArrived:
		notification.Message = M{"en": "Your order trip is arrived."}
	case config.TripDropped:
		notification.Message = M{"en": "Your order trip is dropped."}
	case config.TripCompleted:
		notification.Message = M{"en": "Your order trip is completed."}
	}

	notificationService.PushOneSignalNotification(notification, config.UserAPIKey)
	notificationService.PushWebsocketNotification(order.BusinessID.Hex(), data)
}

func readNearDrivers(c echo.Context) error {
	lat, _ := strconv.ParseFloat(c.FormValue("lat"), 64)
	lng, _ := strconv.ParseFloat(c.FormValue("lng"), 64)

	driverLocations, err := driverLocationService.GetNearbyDrivers(lat, lng, 20)
	if err != nil {
		return response.KnownErrJSON(c, "err.drivers.read", err)
	}
	return response.SuccessInterface(c, driverLocations)
}

func requestTripToDriver(c echo.Context) error {
	order := &model.Order{}
	if err := c.Bind(order); err != nil {
		return response.KnownErrJSON(c, "err.order.bind", err)
	}

	// Update order
	order, err := orderService.UpdateOrderProcess(order.ID, order)
	if err != nil {
		return response.KnownErrJSON(c, "err.order.update", err)
	}

	go runRequestTripToDriver(order)

	return response.SuccessInterface(c, order)
}

func runRequestTripToDriver(order *model.Order) {
	driver, _ := driverService.ReadDriver(order.DriverID)

	notification := &model.OneSignalNotification{}
	notification.AppID = config.DriverAppID
	notification.PlayerIds = []string{driver.OneSignalPlayerID}
	notification.Title = M{"en": config.DriverAppName}

	data := M{}
	data["type"] = config.TripRequest
	data["orderId"] = order.ID
	notification.Data = data

	notification.Message = M{"en": "New trip is requested."}

	notificationService.PushOneSignalNotification(notification, config.DriverAPIKey)
}

func submitPickupScore(c echo.Context) error {
	order := &model.Order{}
	if err := c.Bind(order); err != nil {
		return response.KnownErrJSON(c, "err.order.bind", err)
	}
	// update order pickup score
	order, err := orderService.UpdateOrderPickupScore(order.ID, order)
	if err != nil {
		return response.KnownErrJSON(c, "err.order.update", err)
	}

	return response.SuccessInterface(c, order)
}

func submitOrderByUser(c echo.Context) error {
	order := &model.Order{}
	if err := c.Bind(order); err != nil {
		return response.KnownErrJSON(c, "err.order.bind", err)
	}
	// update feedback
	err := orderService.UpdateOrderRate(order)
	if err != nil {
		return response.KnownErrJSON(c, "err.order.update", err)
	}

	return response.SuccessInterface(c, order)
}

func submitOrderByBusiness(c echo.Context) error {

	return nil
}

func submitOrderByDriver(c echo.Context) error {

	return nil
}
