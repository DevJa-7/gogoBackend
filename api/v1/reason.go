package v1

import (
	"net/http"
	"strconv"

	"../../api/response"
	"../../config"
	"../../model"
	"../../service/authService/permission"
	"../../service/reasonService"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitReason inits reason CRUD apis
// @Title Reasons
// @Description Reasons's router group.
func InitReason(parentRoute *echo.Group) {
	route := parentRoute.Group("/reasons")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createReason))
	route.GET("/:id", permission.AuthRequired(readReason))
	route.PUT("/:id", permission.AuthRequired(updateReason))
	route.DELETE("/:id", permission.AuthRequired(deleteReason))

	route.GET("", permission.AuthRequired(readReasons))

	reasonService.InitService()
}

//------------
// CRUD Handlers
//------------

// @Title createReason
// @Description Create a reason.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Success 200 {object} model.Reason	 	"Returns created reason"
// @Failure 400 {object} response.BasicResponse "err.reason.bind"
// @Failure 400 {object} response.BasicResponse "err.reason.create"
// @Resource /reasons
// @Router /reasons [post]
func createReason(c echo.Context) error {
	reason := &model.Reason{}
	if err := c.Bind(reason); err != nil {
		return response.KnownErrJSON(c, "err.reason.bind", err)
	}

	// Create reason
	publicReason, err := reasonService.CreateReason(reason)
	if err != nil {
		return response.KnownErrJSON(c, "err.reason.create", err)
	}
	return c.JSON(http.StatusOK, publicReason)
}

// @Title readReason
// @Description Read a reason.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Success 200 {object} model.Reason 		"Returns created reason"
// @Failure 400 {object} response.BasicResponse "err.reason.bind"
// @Failure 400 {object} response.BasicResponse "err.reason.read"
// @Resource /reasons
// @Router /reasons/{id} [get]
func readReason(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))

	// Retrieve reason by id
	publicReason, err := reasonService.ReadReason(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.reason.read", err)
	}
	return response.SuccessInterface(c, publicReason)
}

// @Title updateReason
// @Description Update a reason.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Success 200 {object} model.Location 		"Returns created reason"
// @Failure 400 {object} response.BasicResponse "err.reason.bind"
// @Failure 400 {object} response.BasicResponse "err.reason.read"
// @Resource /reasons
// @Router /reasons/{id} [put]
func updateReason(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	reason := &model.Reason{}
	if err := c.Bind(reason); err != nil {
		return response.KnownErrJSON(c, "err.reason.bind", err)
	}

	// Update reason
	publicReason, err := reasonService.UpdateReason(objid, reason)
	if err != nil {
		return response.KnownErrJSON(c, "err.reason.update", err)
	}
	return response.SuccessInterface(c, publicReason)
}

// @Title deleteReason
// @Description Delete a reason.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Location ID."
// @Success 200 {object} response.BasicResponse "Location is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.reason.bind"
// @Failure 400 {object} response.BasicResponse "err.reason.read"
// @Resource /reasons
// @Router /reasons/{id} [delete]
func deleteReason(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	// Remove reason with object id
	err := reasonService.DeleteReason(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.reason.delete", err)
	}
	return response.SuccessJSON(c, "Reason is deleted correctly.")
}

// @Title readReasons
// @Description Read reasons with parameters.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string 	true	"Bearer {token}"
// @Param   query			form   	string  false	"Will search string."
// @Param   offset			form    int		false	"Offset for pagination."
// @Param   count 			form    int		false	"Count that will show per page."
// @Param   field			form    string  false	"Sort field."
// @Param   sort			form    int		false	"Sort direction. 0:default, 1:Ascending, -1:Descending"
// @Success 200 {object} model.ListForm 				"Retrieve all reasons with parameters."
// @Failure 400 {object} response.BasicResponse "err.driver.read"
// @Resource /reasons
// @Router /reasons [get]
func readReasons(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))
	t, _ := strconv.Atoi(c.FormValue("type"))

	// Read reasons with query
	reasons, total, err := reasonService.ReadReasons(query, offset, count, field, sort, t)
	if err != nil {
		return response.KnownErrJSON(c, "err.reason.read", err)
	}

	return response.SuccessInterface(c, model.ListForm{total, reasons})
}
