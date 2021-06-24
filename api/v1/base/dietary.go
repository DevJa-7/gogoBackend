package base

import (
	"strconv"

	"../../../api/response"
	"../../../config"
	"../../../model"
	"../../../service/authService/permission"
	"../../../service/base/dietaryService"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitDietary inits dietary CRUD apis
// @Title Dietaries
// @Description Dietaries's router group.
func InitDietary(parentRoute *echo.Group) {
	parentRoute.GET("/public/dietaries", readDietaries)

	route := parentRoute.Group("/dietaries")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createDietary))
	route.GET("/:id", permission.AuthRequired(readDietary))
	route.PUT("/:id", permission.AuthRequired(updateDietary))
	route.DELETE("/:id", permission.AuthRequired(deleteDietary))

	route.GET("", permission.AuthRequired(readDietaries))

	dietaryService.InitService()
}

// @Title createDietary
// @Description Create a dietary.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   name       	form   	string  true	"Dietary name."
// @Success 200 {object} model.Dietary             "Returns created dietary"
// @Failure 400 {object} response.BasicResponse "err.dietary.bind"
// @Failure 400 {object} response.BasicResponse "err.dietary.create"
// @Resource /dietaries
// @Router /dietaries [post]
func createDietary(c echo.Context) error {
	dietary := &model.Dietary{}
	if err := c.Bind(dietary); err != nil {
		return response.KnownErrJSON(c, "err.dietary.bind", err)
	}
	// fmt.Printf("%+v", dietary)
	// Create dietary
	dietary, err := dietaryService.CreateDietary(dietary)
	if err != nil {
		return response.KnownErrJSON(c, "err.dietary.create", err)
	}
	return response.SuccessInterface(c, dietary)
}

// @Title readDietary
// @Description Read a dietary.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Dietary ID."
// @Success 200 {object} model.Dietary 		"Returns read dietary"
// @Failure 400 {object} response.BasicResponse "err.dietary.bind"
// @Failure 400 {object} response.BasicResponse "err.dietary.read"
// @Resource /dietaries
// @Router /dietaries/{id} [get]
func readDietary(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	// Retrieve rider by id
	dietary, err := dietaryService.ReadDietary(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.dietary.read", err)
	}
	return response.SuccessInterface(c, dietary)
}

// @Title updateDietary
// @Description Update dietary of users.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Dietary ID."
// @Param   name			form   	string  true	"Dietary Name"
// @Success 200 {object} model.Dietary	 		"Returns updated dietary"
// @Failure 400 {object} response.BasicResponse "err.dietary.bind"
// @Failure 400 {object} response.BasicResponse "err.dietary.update"
// @Resource /dietaries
// @Router /dietaries/{id} [put]
func updateDietary(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	dietary := &model.Dietary{}
	if err := c.Bind(dietary); err != nil {
		return response.KnownErrJSON(c, "err.dietary.bind", err)
	}

	// Update dietary
	dietary, err := dietaryService.UpdateDietary(objid, dietary)
	if err != nil {
		return response.KnownErrJSON(c, "err.dietary.update", err)
	}
	return response.SuccessInterface(c, dietary)
}

// @Title deleteDietary
// @Description Delete a dietary.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"User ID."
// @Success 200 {object} response.BasicResponse "User is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.user.bind"
// @Failure 400 {object} response.BasicResponse "err.user.delete"
// @Resource /dietaries
// @Router /dietaries/{id} [delete]
func deleteDietary(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))

	// Remove rider with object id
	err := dietaryService.DeleteDietary(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.dietary.delete", err)
	}
	return response.SuccessJSON(c, "Dietary is deleted correctly.")
}

// @Title readDietaries
// @Description Read a dietaries.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   query			form   	string  false	"Will search string."
// @Param   offset			form    int		false	"Offset for pagination."
// @Param   count 			form    int		false	"Count that will show per page."
// @Param   field			form    string  false	"Sort field."
// @Param   sort			form    int		false	"Sort direction. 0:default, 1:Ascending, -1:Descending"
// @Success 200 {object} response.BasicResponse "User is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.user.bind"
// @Failure 400 {object} response.BasicResponse "err.user.delete"
// @Resource /dietaries
// @Router /dietaries/{id} [delete]
func readDietaries(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))
	top, _ := strconv.Atoi(c.FormValue("top"))
	def, _ := strconv.ParseBool(c.FormValue("default"))

	// Read dietaries with query
	dietaries, total, err := dietaryService.ReadDietaries(query, offset, count, field, sort, top, def)
	if err != nil {
		return response.KnownErrJSON(c, "err.dietary.read", err)
	}

	return response.SuccessInterface(c, &model.ListForm{total, dietaries})
}
