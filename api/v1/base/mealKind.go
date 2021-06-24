package base

import (
	"strconv"

	"../../../api/response"
	"../../../config"
	"../../../model"
	"../../../service/authService/permission"
	"../../../service/base/mealKindService"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitMealKind initialze mealKind api
func InitMealKind(parentRoute *echo.Group) {
	parentRoute.GET("/public/mealKinds", readMealKinds)

	route := parentRoute.Group("/mealKinds")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createMealKind))
	route.GET("/:id", permission.AuthRequired(readMealKind))
	route.PUT("/:id", permission.AuthRequired(updateMealKind))
	route.DELETE("/:id", permission.AuthRequired(deleteMealKind))

	route.GET("", permission.AuthRequired(readMealKinds))

	mealKindService.InitService()
}

//--------------
// CRUD Handlers
//--------------

// Create Url MealKind
func createMealKind(c echo.Context) error {
	mealKind := &model.MealKind{}
	if err := c.Bind(mealKind); err != nil {
		return response.KnownErrJSON(c, "err.mealKind.bind", err)
	}
	// Create mealKind
	mealKind, err := mealKindService.CreateMealKind(mealKind)
	if err != nil {
		return response.KnownErrJSON(c, "err.mealKind.create", err)
	}
	return response.SuccessInterface(c, mealKind)
}

// readMealKind
func readMealKind(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	// Retrieve rider by id
	mealKind, err := mealKindService.ReadMealKind(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.mealKind.read", err)
	}
	return response.SuccessInterface(c, mealKind)
}

// updateMealKind
func updateMealKind(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	mealKind := &model.MealKind{}
	if err := c.Bind(mealKind); err != nil {
		return response.KnownErrJSON(c, "err.mealKind.bind", err)
	}

	// Update mealKind
	mealKind, err := mealKindService.UpdateMealKind(objid, mealKind)
	if err != nil {
		return response.KnownErrJSON(c, "err.mealKind.update", err)
	}
	return response.SuccessInterface(c, mealKind)
}

// @Title readMealKinds
// @Description Read mealKinds with parameters.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string 	true	"Bearer {token}"
// @Param   query			form   	string  false	"Will search string."
// @Param   offset			form    int		false	"Offset for pagination."
// @Param   count 			form    int		false	"Count that will show per page."
// @Param   field			form    string  false	"Sort field."
// @Param   sort			form    int		false	"Sort direction. 0:default, 1:Ascending, -1:Descending"
// @Success 200 {object} model.ListForm 			"Return queried mealkinds"
// @Failure 400 {object} response.BasicResponse "err.foodType.read"
// @Resource /mealKinds
// @Router /mealKinds [get]
func deleteMealKind(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))

	// Remove rider with object id
	err := mealKindService.DeleteMealKind(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.mealKind.delete", err)
	}
	return response.SuccessJSON(c, "MealKind is deleted correctly.")
}

// readMealKinds
func readMealKinds(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))

	// Read mealKinds with query
	mealKinds, total, err := mealKindService.ReadMealKinds(query, offset, count, field, sort)
	if err != nil {
		return response.KnownErrJSON(c, "err.mealKind.read", err)
	}

	return response.SuccessInterface(c, &model.ListForm{total, mealKinds})
}
