package v1

import (
	"strconv"

	"../../api/response"
	"../../config"
	"../../model"
	"../../service/authService/permission"
	"../../service/urlGroupService"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitURLGroup initialze url group api
// @Title URLGroups
// @Description URLGroups's router group.
func InitURLGroup(parentRoute *echo.Group) {
	route := parentRoute.Group("/urls")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	// group crud service
	route.POST("", permission.AuthRequired(createURLGroup))
	route.GET("/:id", permission.AuthRequired(readURLGroup))
	route.PUT("/:id", permission.AuthRequired(updateURLGroup))
	route.DELETE("/:id", permission.AuthRequired(deleteURLGroup))

	route.GET("", permission.AuthRequired(readURLGroups))

	urlGroupService.InitService()
}

//--------------
// CRUD Handlers
//--------------

// @Title createURLGroup
// @Description creates urlGroup
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   name       		form   	string  true	"URL Name"
// @Param   url				form   	string 	true	"URL"
// @Success 200 {object} model.URLGroup 		"Returns created url group"
// @Failure 400 {object} response.BasicResponse "err.urlGroup.bind"
// @Failure 400 {object} response.BasicResponse "err.urlGroup.create"
// @Resource /urls
// @Router /urls [post]
func createURLGroup(c echo.Context) error {
	urlGroup := &model.URLGroup{}
	if err := c.Bind(urlGroup); err != nil {
		return response.KnownErrJSON(c, "err.urlGroup.bind", err)
	}
	// Insert Data
	urlGroup, err := urlGroupService.CreateURLGroup(urlGroup)
	if err != nil {
		return response.KnownErrJSON(c, "err.urlGroup.create", err)
	}

	return response.SuccessInterface(c, urlGroup)
}

// @Title readURLGroup
// @Description reads urlGroup
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"URLGroup ID."
// @Success 200 {object} model.URLGroup 		"Returns read URLGroup"
// @Failure 400 {object} response.BasicResponse "err.urlGroup.bind"
// @Failure 400 {object} response.BasicResponse "err.urlGroup.read"
// @Resource /urls
// @Router /urls/{id} [get]
func readURLGroup(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	// read url group with obj id
	urlGroup, err := urlGroupService.ReadURLGroup(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.urlGroup.read", err)
	}

	return response.SuccessInterface(c, urlGroup)
}

// @Title updateURLGroup
// @Description updates URLGroup
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   id				path   	string  true	"URLGroup ID."
// @Param   name	      	form   	string  true	"URLGroup Name"
// @Param   url				form   	string  true	"URLGroup URL"
// @Success 200 {object} model.URLGroup 		"Returns updated url"
// @Failure 400 {object} response.BasicResponse "err.urlGroup.bind"
// @Failure 400 {object} response.BasicResponse "err.urlGroup.update"
// @Resource /urls
// @Router /urls/{id} [put]
func updateURLGroup(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	// bind data
	urlGroup := &model.URLGroup{}
	if err := c.Bind(urlGroup); err != nil {
		return response.KnownErrJSON(c, "err.urlGroup.bind", err)
	}
	// update data
	urlGroup, err := urlGroupService.UpdateURLGroup(objid, urlGroup)
	if err != nil {
		return response.KnownErrJSON(c, "err.urlGroup.update", err)
	}

	return response.SuccessInterface(c, urlGroup)
}

// @Title deleteURLGroup
// @Description deletes URLGroup
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"User ID."
// @Success 200 {object} response.BasicResponse "URLGroup is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.urlGroup.bind"
// @Failure 400 {object} response.BasicResponse "err.urlGroup.delete"
// @Resource /urls
// @Router /urls/{id} [delete]
func deleteURLGroup(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	err := urlGroupService.DeleteURLGroup(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.urlGroup.delete", err)
	}

	return response.SuccessJSON(c, "URLGroup is deleted correctly.")
}

// @Title readURLGroups
// @Description returns all url groups
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string 	true	"Bearer {token}"
// @Param   query			form   	string  false	"Will search string."
// @Param   offset			form    int		false	"Offset for pagination."
// @Param   count 			form    int		false	"Count that will show per page."
// @Param   field			form    string  false	"Sort field."
// @Param   sort			form    int		false	"Sort direction. 0:default, 1:Ascending, -1:Descending"
// @Success 200 {object} model.ListForm 				"Returned list URLGroups."
// @Failure 400 {object} response.BasicResponse "err.urlGroup.read"
// @Resource /urls
// @Router /urls [get]
func readURLGroups(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))

	// Read urlGroups with query
	urlGroups, total, err := urlGroupService.ReadURLGroups(query, offset, count, field, sort)
	if err != nil {
		return response.KnownErrJSON(c, "err.urlGroup.read", err)
	}

	return response.SuccessInterface(c, &model.ListForm{total, urlGroups})
}
