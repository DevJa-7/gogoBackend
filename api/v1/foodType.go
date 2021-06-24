package v1

import (
	"errors"
	"strconv"

	"../../config"
	"../../model"
	"../../service/authService/permission"
	"../../service/foodTypeService"
	"../response"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitFoodType inits foodType CRUD apis
// @Title FoodTypes
// @Description FoodTypes's router group.
func InitFoodType(parentRoute *echo.Group) {
	route := parentRoute.Group("/foodTypes")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createFoodType))
	route.GET("/:id", permission.AuthRequired(readFoodType))
	route.PUT("/:id", permission.AuthRequired(updateFoodType))
	route.DELETE("/:id", permission.AuthRequired(deleteFoodType))

	route.GET("", permission.AuthRequired(readFoodTypes))
	route.GET("/by/business/:id", permission.AuthRequired(readFoodTypesByBusiness))

	// create spec
	route.POST("/by/spec/:id", permission.AuthRequired(createSpec))
	route.GET("/by/spec/:id/:number", permission.AuthRequired(readSpec))
	route.PUT("/by/spec/:id/:number", permission.AuthRequired(updateSpec))
	route.DELETE("/by/spec/:id/:number", permission.AuthRequired(deleteSpec))

	route.GET("/by/spec/:id", permission.AuthRequired(readSpecs))

	foodTypeService.InitService()
}

// @Title createFoodType
// @Description Create a foodType.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   name    	   	form   	string  true	"FoodType Name."
// @Param   image			form   	string 	true	"FoodType Image or Icon."
// @Success 200 {object} model.FoodType 		"Returns created foodType"
// @Failure 400 {object} response.BasicResponse "err.foodType.bind"
// @Failure 400 {object} response.BasicResponse "err.foodType.create"
// @Resource /foodTypes
// @Router /foodTypes [post]
func createFoodType(c echo.Context) error {
	foodType := &model.FoodType{}
	if err := c.Bind(foodType); err != nil {
		return response.KnownErrJSON(c, "err.foodType.bind", err)
	}

	// create foodType
	foodType, err := foodTypeService.CreateFoodType(foodType)
	if err != nil {
		return response.KnownErrJSON(c, "err.foodType.create", err)
	}

	return response.SuccessInterface(c, foodType)
}

// @Title readFoodType
// @Description Read a foodType.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"FoodType ID."
// @Success 200 {object} model.PublicFoodType 		"Returns read foodType"
// @Failure 400 {object} response.BasicResponse "err.foodType.bind"
// @Failure 400 {object} response.BasicResponse "err.foodType.read"
// @Resource /foodTypes
// @Router /foodTypes/{id} [get]
func readFoodType(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.foodType.bind", errors.New("Retreived object id is invalid"))
	}

	objid := bson.ObjectIdHex(c.Param("id"))
	foodType, err := foodTypeService.ReadFoodType(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.foodType.read", err)
	}

	return response.SuccessInterface(c, foodType)
}

// @Title updateFoodType
// @Description Update a foodType.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   id				path   	string  true	"FoodType ID."
// @Param   avatar      	form   	string  true	"FoodType Avatar"
// @Param   firstname		form   	string  true	"FoodType Firstname"
// @Param   lastname		form   	string  true	"FoodType Lastname"
// @Param   email	    	form   	string  true	"FoodType Email"
// @Param   birth      		form   	Time   	true	"FoodType Birth"
// @Success 200 {object} model.PublicFoodType 		"Returns read foodType"
// @Failure 400 {object} response.BasicResponse "err.foodType.bind"
// @Failure 400 {object} response.BasicResponse "err.foodType.read"
// @Resource /foodTypes
// @Router /foodTypes/{id} [put]
func updateFoodType(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.foodType.bind", errors.New("Retreived object id is invalid"))
	}

	foodType := &model.FoodType{}
	if err := c.Bind(foodType); err != nil {
		return response.KnownErrJSON(c, "err.foodType.bind", err)
	}

	objid := bson.ObjectIdHex(c.Param("id"))
	foodType, err := foodTypeService.UpdateFoodType(objid, foodType)
	if err != nil {
		return response.KnownErrJSON(c, "err.foodType.read", err)
	}

	return response.SuccessInterface(c, foodType)
}

// @Title deleteFoodType
// @Description Delete a foodType.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"FoodType ID."
// @Success 200 {object} response.BasicResponse "FoodType is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.foodType.bind"
// @Failure 400 {object} response.BasicResponse "err.foodType.delete"
// @Resource /foodTypes
// @Router /foodTypes/{id} [delete]
func deleteFoodType(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.foodType.bind", errors.New("Retreived object id is invalid"))
	}

	// delete foodType with object id
	objid := bson.ObjectIdHex(c.Param("id"))
	err := foodTypeService.DeleteFoodType(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.foodType.delete", err)
	}

	return response.SuccessJSON(c, "FoodType is deleted correctly.")
}

// @Title readFoodTypes
// @Description Read foodTypes with parameters.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string 	true	"Bearer {token}"
// @Param   query			form   	string  false	"Will search string."
// @Param   offset			form    int		false	"Offset for pagination."
// @Param   count 			form    int		false	"Count that will show per page."
// @Param   field			form    string  false	"Sort field."
// @Param   sort			form    int		false	"Sort direction. 0:default, 1:Ascending, -1:Descending"
// @Success 200 {object} model.ListForm 				"FoodType is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.foodType.read"
// @Resource /foodTypes
// @Router /foodTypes [get]
func readFoodTypes(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))

	// read foodTypes with params
	foodTypes, total, err := foodTypeService.ReadFoodTypes(query, offset, count, field, sort)
	if err != nil {
		return response.KnownErrJSON(c, "err.foodType.read", err)
	}
	return response.SuccessInterface(c, &model.ListForm{total, foodTypes})
}

func readFoodTypesByBusiness(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.foodType.bind", errors.New("Retreived object id is invalid"))
	}
	// business object id
	objid := bson.ObjectIdHex(c.Param("id"))
	foodTypes, err := foodTypeService.ReadFoodTypesWithBusiness(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.foodType.read", err)
	}
	return response.SuccessInterface(c, foodTypes)
}

func createSpec(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.foodType.bind", errors.New("Retreived object id is invalid"))
	}
	// foodtype object id
	objid := bson.ObjectIdHex(c.Param("id"))

	spec := &model.FoodOption{}
	if err := c.Bind(&spec); err != nil {
		return response.KnownErrJSON(c, "err.foodOption.bind", err)
	}
	spec, err := foodTypeService.CreateSpec(objid, spec)
	if err != nil {
		return response.KnownErrJSON(c, "err.foodOption.create", err)
	}
	return response.SuccessInterface(c, spec)
}

func readSpec(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.foodType.bind", errors.New("Retreived object id is invalid"))
	}
	// foodtype object id
	objid := bson.ObjectIdHex(c.Param("id"))
	number := c.Param("number")

	spec, err := foodTypeService.ReadSpec(objid, number)
	if err != nil {
		return response.KnownErrJSON(c, "err.foodOption.read", err)
	}
	return response.SuccessInterface(c, spec)
}

func updateSpec(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.foodType.bind", errors.New("Retreived object id is invalid"))
	}
	// foodtype object id
	objid := bson.ObjectIdHex(c.Param("id"))
	number := c.Param("number")

	spec := &model.FoodOption{}
	if err := c.Bind(&spec); err != nil {
		return response.KnownErrJSON(c, "err.foodOption.bind", err)
	}
	spec, err := foodTypeService.UpdateSpec(objid, number, spec)
	if err != nil {
		return response.KnownErrJSON(c, "err.foodOption.create", err)
	}
	return response.SuccessInterface(c, spec)
}

func deleteSpec(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.foodType.bind", errors.New("Retreived object id is invalid"))
	}
	// foodtype object id
	objid := bson.ObjectIdHex(c.Param("id"))
	number := c.Param("number")

	err := foodTypeService.DeleteSpec(objid, number)
	if err != nil {
		return response.KnownErrJSON(c, "err.foodOption.create", err)
	}

	return response.SuccessJSON(c, "Spec is deleted correctly.")
}

func readSpecs(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.foodType.bind", errors.New("Retreived object id is invalid"))
	}
	// foodtype object id
	objid := bson.ObjectIdHex(c.Param("id"))

	specs, err := foodTypeService.ReadSpecs(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.foodOptions.read", err)
	}
	return response.SuccessInterface(c, specs)
}
