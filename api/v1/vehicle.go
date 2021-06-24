package v1

import (
	"net/http"
	"strconv"

	"../../api/response"
	"../../config"
	"../../model"
	"../../service/authService/permission"
	"../../service/vehicleService"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitVehicle inits vehicle CRUD apis
// @Title Vehicles
// @Description Vehicle's router group.
func InitVehicle(parentRoute *echo.Group) {
	route := parentRoute.Group("/vehicles")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createVehicle))
	route.GET("/:id", permission.AuthRequired(readVehicle))
	route.PUT("/:id", permission.AuthRequired(updateVehicle))
	route.DELETE("/:id", permission.AuthRequired(deleteVehicle))

	route.GET("", permission.AuthRequired(readVehicles))
	route.GET("/active", permission.AuthRequired(readActiveVehicles))

	vehicleService.InitService()
}

//------------
// CRUD Handlers
//------------

// @Title createVehicle
// @Description Create a vehicle.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Success 200 {object} model.Vehicle	 		"Returns created vehicle"
// @Failure 400 {object} response.BasicResponse "err.vehicle.bind"
// @Failure 400 {object} response.BasicResponse "err.vehicle.create"
// @Resource /vehicles
// @Router /vehicles [post]
func createVehicle(c echo.Context) error {
	vehicle := &model.Vehicle{}
	if err := c.Bind(vehicle); err != nil {
		return response.KnownErrJSON(c, "err.vehicle.bind", err)
	}

	// Create vehicle
	publicVehicle, err := vehicleService.CreateVehicle(vehicle)
	if err != nil {
		return response.KnownErrJSON(c, "err.vehicle.create", err)
	}
	return c.JSON(http.StatusOK, publicVehicle)
}

// @Title readLocation
// @Description Read a location.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Success 200 {object} model.Location 		"Returns created location"
// @Failure 400 {object} response.BasicResponse "err.location.bind"
// @Failure 400 {object} response.BasicResponse "err.location.read"
// @Resource /vehicles
// @Router /vehicles/{id} [get]
func readVehicle(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))

	// Retrieve vehicle by id
	publicVehicle, err := vehicleService.ReadVehicle(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.vehicle.read", err)
	}
	return response.SuccessInterface(c, publicVehicle)
}

// @Title updateLocation
// @Description Update a location.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Success 200 {object} model.Location 		"Returns created location"
// @Failure 400 {object} response.BasicResponse "err.location.bind"
// @Failure 400 {object} response.BasicResponse "err.location.read"
// @Resource /vehicles
// @Router /vehicles/{id} [put]
func updateVehicle(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	vehicle := &model.Vehicle{}
	if err := c.Bind(vehicle); err != nil {
		return response.KnownErrJSON(c, "err.vehicle.bind", err)
	}

	// Update vehicle
	publicVehicle, err := vehicleService.UpdateVehicle(objid, vehicle)
	if err != nil {
		return response.KnownErrJSON(c, "err.vehicle.update", err)
	}
	return response.SuccessInterface(c, publicVehicle)
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
// @Resource /vehicles
// @Router /vehicles/{id} [delete]
func deleteVehicle(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	// Remove vehicle with object id
	err := vehicleService.DeleteVehicle(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.vehicle.delete", err)
	}
	return response.SuccessJSON(c, "Vehicle is deleted correctly.")
}

// @Title readLocations
// @Description Read vehicles with parameters.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string 	true	"Bearer {token}"
// @Param   query			form   	string  false	"Will search string."
// @Param   offset			form    int		false	"Offset for pagination."
// @Param   count 			form    int		false	"Count that will show per page."
// @Param   field			form    string  false	"Sort field."
// @Param   sort			form    int		false	"Sort direction. 0:default, 1:Ascending, -1:Descending"
// @Success 200 {object} model.ListForm 				"Retrieve all vehicles with parameters."
// @Failure 400 {object} response.BasicResponse "err.driver.read"
// @Resource /vehicles
// @Router /vehicles [get]
func readVehicles(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))

	// Read vehicles with query
	vehicles, total, err := vehicleService.ReadVehicles(query, offset, count, field, sort)
	if err != nil {
		return response.KnownErrJSON(c, "err.vehicle.read", err)
	}

	return response.SuccessInterface(c, model.ListForm{total, vehicles})
}

func readActiveVehicles(c echo.Context) error {
	// Read vehicles with query
	vehicles, err := vehicleService.ReadActiveVehicles()
	if err != nil {
		return response.KnownErrJSON(c, "err.vehicle.read", err)
	}

	return response.SuccessInterface(c, vehicles)
}
