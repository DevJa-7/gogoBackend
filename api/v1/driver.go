package v1

import (
	"errors"
	"fmt"
	"strconv"

	"../../config"
	"../../model"
	"../../service/authService/driverService"
	"../../service/authService/permission"
	"../../service/driverLocationService"
	"../../service/locationService"
	"../response"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitDrivers inits driver CRUD apis
// @Title Drivers
// @Description Drivers's router group.
func InitDrivers(parentRoute *echo.Group) {
	route := parentRoute.Group("/drivers")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createDriver))
	route.GET("/:id", permission.AuthRequired(readDriver))
	route.GET("/detail/:id", permission.AuthRequired(readDetailDriver))
	route.PUT("/:id", permission.AuthRequired(updateDriver))
	route.DELETE("/:id", permission.AuthRequired(deleteDriver))

	route.GET("", permission.AuthRequired(readDrivers))

	route.POST("/update/location", permission.AuthRequired(updateDriverLocation))

	route.POST("/vehicle/:id", permission.AuthRequired(createDriverVehicle))
	route.PUT("/vehicle/:id/:number", permission.AuthRequired(updateDriverVehicle))
	route.DELETE("/vehicle/:id/:number", permission.AuthRequired(deleteDriverVehicle))
	route.GET("/vehicles/:id", permission.AuthRequired(readDriverVehicles))

	driverService.InitService()
}

// @Title createDriver
// @Description Create a driver.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   email       	form   	string  true	"Driver Email."
// @Param   password		form   	string 	true	"Driver Password."
// @Success 200 {object} model.PublicDriver 		"Returns created driver"
// @Failure 400 {object} response.BasicResponse "err.driver.bind"
// @Failure 400 {object} response.BasicResponse "err.driver.create"
// @Resource /drivers
// @Router /drivers [post]
func createDriver(c echo.Context) error {
	driver := &model.Driver{}
	if err := c.Bind(driver); err != nil {
		return response.KnownErrJSON(c, "err.driver.bind", err)
	}

	// create driver
	driver, err := driverService.CreateDriver(driver)
	if err != nil {
		return response.KnownErrJSON(c, "err.driver.create", err)
	}

	publicDriver := &model.PublicDriver{Driver: driver}
	return response.SuccessInterface(c, publicDriver)
}

// @Title readDriver
// @Description Read a driver.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Driver ID."
// @Success 200 {object} model.PublicDriver 		"Returns read driver"
// @Failure 400 {object} response.BasicResponse "err.driver.bind"
// @Failure 400 {object} response.BasicResponse "err.driver.read"
// @Resource /drivers
// @Router /drivers/{id} [get]
func readDriver(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.driver.bind", errors.New("Retreived object id is invalid"))
	}

	objid := bson.ObjectIdHex(c.Param("id"))
	driver, err := driverService.ReadDriver(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.driver.read", err)
	}

	publicDriver := &model.PublicDriver{Driver: driver}
	return response.SuccessInterface(c, publicDriver)
}

// @Title updateDriver
// @Description Update a driver.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   id				path   	string  true	"Driver ID."
// @Param   avatar      	form   	string  true	"Driver Avatar"
// @Param   firstname		form   	string  true	"Driver Firstname"
// @Param   lastname		form   	string  true	"Driver Lastname"
// @Param   email	    	form   	string  true	"Driver Email"
// @Param   birth      		form   	Time   	true	"Driver Birth"
// @Success 200 {object} model.PublicDriver 		"Returns read driver"
// @Failure 400 {object} response.BasicResponse "err.driver.bind"
// @Failure 400 {object} response.BasicResponse "err.driver.read"
// @Resource /drivers
// @Router /drivers/{id} [put]
func updateDriver(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.driver.bind", errors.New("Retreived object id is invalid"))
	}

	driver := &model.Driver{}
	if err := c.Bind(driver); err != nil {
		return response.KnownErrJSON(c, "err.driver.bind", err)
	}

	objid := bson.ObjectIdHex(c.Param("id"))
	driver, err := driverService.UpdateDriver(objid, driver)
	if err != nil {
		return response.KnownErrJSON(c, "err.driver.read", err)
	}

	publicDriver := &model.PublicDriver{Driver: driver}
	return response.SuccessInterface(c, publicDriver)
}

// @Title deleteDriver
// @Description Delete a driver.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Driver ID."
// @Success 200 {object} response.BasicResponse "Driver is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.driver.bind"
// @Failure 400 {object} response.BasicResponse "err.driver.delete"
// @Resource /drivers
// @Router /drivers/{id} [delete]
func deleteDriver(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.driver.bind", errors.New("Retreived object id is invalid"))
	}

	// delete driver with object id
	objid := bson.ObjectIdHex(c.Param("id"))
	err := driverService.DeleteDriver(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.driver.delete", err)
	}

	return response.SuccessJSON(c, "Driver is deleted correctly.")
}

// @Title readDrivers
// @Description Read drivers with parameters.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string 	true	"Bearer {token}"
// @Param   query			form   	string  false	"Will search string."
// @Param   offset			form    int		false	"Offset for pagination."
// @Param   count 			form    int		false	"Count that will show per page."
// @Param   field			form    string  false	"Sort field."
// @Param   sort			form    int		false	"Sort direction. 0:default, 1:Ascending, -1:Descending"
// @Success 200 {object} model.ListForm 				"Retrieve all drivers with parameters."
// @Failure 400 {object} response.BasicResponse "err.driver.read"
// @Resource /drivers
// @Router /drivers [get]
func readDrivers(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))

	// read drivers with params
	drivers, total, err := driverService.ReadDrivers(query, offset, count, field, sort)
	if err != nil {
		return response.KnownErrJSON(c, "err.driver.read", err)
	}

	// retreive by public format
	publicDrivers := []*model.PublicDriver{}
	for _, driver := range drivers {
		publicDriver := &model.PublicDriver{Driver: driver}
		publicDrivers = append(publicDrivers, publicDriver)
	}

	return response.SuccessInterface(c, &model.ListForm{total, publicDrivers})
}

// readDetailDriver returns details info for driver profile
func readDetailDriver(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	// get driver information
	publicDriver, err := driverService.ReadDriver(objid)
	if err != nil {
		fmt.Println("err.driver.read - ", err)
	}
	// get city name
	location, err := locationService.ReadLocationWithPlaceID(publicDriver.LocationPlaceID)
	if err != nil {
		fmt.Println("err.location.read - ", err)
	}
	// // get recent 8 trips of driver
	// trips, err := tripService.ReadTripsHistory("driver", publicDriver.ID, 0, 8)
	// if err != nil {
	// 	fmt.Println("err.trips.read - ", err)
	// }

	// cancelledCount := tripService.ReadTripsCount("driver", publicDriver.ID, config.TripCancelled)
	// completedCount := tripService.ReadTripsCount("driver", publicDriver.ID, config.TripCompleted)
	return response.SuccessInterface(c, map[string]interface{}{
		"driver":   publicDriver,
		"location": location.City,
		// "recent_trips":    trips,
		// "cancelled_count": cancelledCount,
		// "completed_count": completedCount,
	})
}

// updateDriverLocation updates driver's location
func updateDriverLocation(c echo.Context) error {
	// bind param
	driverLocation := &model.DriverLocation{}
	if err := c.Bind(driverLocation); err != nil {
		return response.KnownErrJSON(c, "err.driverLocation.bind", err)
	}

	driverLocation, err := driverLocationService.UpdateDriverLocation(driverLocation)
	if err != nil {
		return response.KnownErrJSON(c, "err.driverLocation.update", err)
	}

	return response.SuccessInterface(c, driverLocation)
}

// createDriverVehicle creates driver vehicle
func createDriverVehicle(c echo.Context) error {
	driverVehicle := &model.DriverVehicle{}
	if err := c.Bind(driverVehicle); err != nil {
		return response.KnownErrJSON(c, "err.driverVehicle.bind", err)
	}
	objid := bson.ObjectIdHex(c.Param("id")) // driver id

	for _, document := range driverVehicle.Documents {
		if document.Status == config.DocumentNone {
			document.Status = config.DocumentPending
		}
	}

	if err := driverService.CreateDriverVehicle(objid, driverVehicle); err != nil {
		return response.KnownErrJSON(c, "err.driverVehicle.create", err)
	}

	return response.SuccessInterface(c, "Driver vehicle is created correctly.")
}

// updateDriverVehicle updates driver vehicle
func updateDriverVehicle(c echo.Context) error {
	driverVehicle := &model.DriverVehicle{}
	if err := c.Bind(driverVehicle); err != nil {
		return response.KnownErrJSON(c, "err.driverVehicle.bind", err)
	}
	objid := bson.ObjectIdHex(c.Param("id")) //driver id
	number := c.Param("number")

	for _, document := range driverVehicle.Documents {
		if document.Status == config.DocumentNone {
			document.Status = config.DocumentPending
		}
	}

	if err := driverService.UpdateDriverVehicle(objid, number, driverVehicle); err != nil {
		return response.KnownErrJSON(c, "err.driverVehicle.update", err)
	}

	return response.SuccessInterface(c, "Driver vehicle is updated correctly.")
}

// deleteDriverVehicle deletes driver vehicle
func deleteDriverVehicle(c echo.Context) error {
	driverVehicle := &model.DriverVehicle{}
	if err := c.Bind(driverVehicle); err != nil {
		return response.KnownErrJSON(c, "err.driverVehicle.bind", err)
	}
	objid := bson.ObjectIdHex(c.Param("id")) //driver id

	if err := driverService.DeleteDriverVehicle(objid, driverVehicle); err != nil {
		return response.KnownErrJSON(c, "err.driverVehicle.delete", err)
	}

	return response.SuccessInterface(c, "Driver vehicle is delete correctly.")
}

func readDriverVehicles(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id")) //driver id

	driverVehicles, err := driverService.ReadDriverVehicles(objid)
	if err != nil {
		return response.SuccessInterface(c, []interface{}{})
	}
	return response.SuccessInterface(c, driverVehicles)
}
