package v1

import (
	"errors"
	"fmt"
	"strconv"

	"../../api/response"
	"../../config"
	"../../model"
	"../../service/authService/permission"
	"../../service/locationService"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitLocation inits driver CRUD apis
// @Title Locations
// @Description Location's router group.
func InitLocation(parentRoute *echo.Group) {
	parentRoute.GET("/public/locations", readLocations)

	route := parentRoute.Group("/locations")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createLocation))
	route.GET("/:id", permission.AuthRequired(readLocation))
	route.PUT("/:id", permission.AuthRequired(updateLocation))
	route.DELETE("/:id", permission.AuthRequired(deleteLocation))

	route.GET("", permission.AuthRequired(readLocations))

	route.GET("/vehicles/:placeId", permission.AuthRequired(readLocationVehicles))
	route.PUT("/update/vehicleInfo/:id", permission.AuthRequired(updateVehicleInfos))

	locationService.InitService()
}

// @Title createLocation
// @Description Create a location.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Success 200 {object} model.Location 		"Returns created location"
// @Failure 400 {object} response.BasicResponse "err.location.bind"
// @Failure 400 {object} response.BasicResponse "err.location.create"
// @Resource /locations
// @Router /locations [post]
func createLocation(c echo.Context) error {
	location := &model.Location{}
	if err := c.Bind(location); err != nil {
		return response.KnownErrJSON(c, "err.location.bind", err)
	}
	// Check duplicate
	if _, err := locationService.ReadLocationWithPlaceID(location.PlaceID); err == nil {
		return response.KnownErrJSON(c, "err.location.duplicate", errors.New("This location is already registered"))
	}
	// Create location
	location, err := locationService.CreateLocation(location)
	if err != nil {
		return response.KnownErrJSON(c, "err.location.create", err)
	}
	return response.SuccessInterface(c, location)
}

// @Title readLocation
// @Description Read a location.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Success 200 {object} model.Location 		"Returns created location"
// @Failure 400 {object} response.BasicResponse "err.location.bind"
// @Failure 400 {object} response.BasicResponse "err.location.read"
// @Resource /locations
// @Router /locations/{id} [get]
func readLocation(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.location.bind", errors.New("Retreived object id is invalid"))
	}
	objid := bson.ObjectIdHex(c.Param("id"))

	// Retrieve location by id
	location, err := locationService.ReadLocation(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.location.read", err)
	}
	return response.SuccessInterface(c, location)
}

// @Title updateLocation
// @Description Update a location.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Success 200 {object} model.Location 		"Returns created location"
// @Failure 400 {object} response.BasicResponse "err.location.bind"
// @Failure 400 {object} response.BasicResponse "err.location.read"
// @Resource /locations
// @Router /locations/{id} [put]
func updateLocation(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.location.bind", errors.New("Retreived object id is invalid"))
	}
	objid := bson.ObjectIdHex(c.Param("id"))

	location := &model.Location{}
	if err := c.Bind(location); err != nil {
		fmt.Println("err.location.bind", err)
		return response.KnownErrJSON(c, "err.location.bind", err)
	}

	// Update location
	location, err := locationService.UpdateLocation(objid, location)
	if err != nil {
		return response.KnownErrJSON(c, "err.location.update", err)
	}
	return response.SuccessInterface(c, location)
}

// @Title deleteLocation
// @Description Delete a location.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Location ID."
// @Success 200 {object} response.BasicResponse "Location is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.location.bind"
// @Failure 400 {object} response.BasicResponse "err.location.read"
// @Resource /locations
// @Router /locations/{id} [delete]
func deleteLocation(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.location.bind", errors.New("Retreived object id is invalid"))
	}
	objid := bson.ObjectIdHex(c.Param("id"))

	// Remove location with object id
	err := locationService.DeleteLocation(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.location.delete", err)
	}
	return response.SuccessJSON(c, "Location is deleted correctly.")
}

// @Title readLocations
// @Description Read locations with parameters.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string 	true	"Bearer {token}"
// @Param   query			form   	string  false	"Will search string."
// @Param   offset			form    int		false	"Offset for pagination."
// @Param   count 			form    int		false	"Count that will show per page."
// @Param   field			form    string  false	"Sort field."
// @Param   sort			form    int		false	"Sort direction. 0:default, 1:Ascending, -1:Descending"
// @Success 200 {object} model.ListForm 				"Retrieve all locations with parameters."
// @Failure 400 {object} response.BasicResponse "err.driver.read"
// @Resource /locations
// @Router /locations [get]
func readLocations(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))

	// Read locations with query
	locations, total, err := locationService.ReadLocations(query, offset, count, field, sort)
	if err != nil {
		return response.KnownErrJSON(c, "err.location.read", err)
	}

	return response.SuccessInterface(c, model.ListForm{total, locations})
}

// @Title readLocationVehicles
// @Description Read vehicles of special location.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string 	true	"Bearer {token}"
// @Success 200 {object} 	array 					"Retrieve all locations with parameters."
// @Failure 400 {object} response.BasicResponse 	"err.location.read"
// @Resource /locations/vehicles
// @Router /locations/vehicles/{placeId} [get]
func readLocationVehicles(c echo.Context) error {
	placeID := c.Param("placeId")

	// Retrieve location by placeId
	location, err := locationService.ReadLocationWithPlaceID(placeID)
	if err != nil {
		return response.KnownErrJSON(c, "err.location.read", err)
	}
	return response.SuccessInterface(c, location.VehicleInfos)
}

// @Title updateVehicleInfos
// @Description Read locations with parameters.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string 	true	"Bearer {token}"
// @Success 200 {object} 	response.BasicResponse 	"Retrieve all locations with parameters."
// @Failure 400 {object} 	response.BasicResponse 	"err.location.read"
// @Resource /locations/update/vehicleInfo
// @Router /locations//update/vehicleInfo/{id} [put]
func updateVehicleInfos(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.location.bind", errors.New("Retreived object id is invalid"))
	}
	objid := bson.ObjectIdHex(c.Param("id"))

	location := &model.Location{}
	if err := c.Bind(location); err != nil {
		return response.KnownErrJSON(c, "err.location.bind", err)
	}
	// update vehicle infos
	err := locationService.UpdateVehicleInfos(objid, location.VehicleInfos)
	if err != nil {
		return response.KnownErrJSON(c, "err.vehicleInfo.update", err)
	}

	return response.SuccessInterface(c, "VehicleInfo is updated correctly.")
}
