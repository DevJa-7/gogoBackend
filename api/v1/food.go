package v1

import (
	"strconv"

	"../../config"
	"../../model"
	"../../service/authService/permission"
	"../../service/foodService"
	"../response"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitFood inits food CRUD apis
// @Title Foods
// @Description Foods's router group.
func InitFood(parentRoute *echo.Group) {
	route := parentRoute.Group("/foods")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createFood))
	route.GET("/:id", permission.AuthRequired(readFood))
	route.PUT("/:id", permission.AuthRequired(updateFood))
	route.DELETE("/:id", permission.AuthRequired(deleteFood))

	route.GET("", permission.AuthRequired(readFoods))
	route.GET("/by/business", permission.AuthRequired(readFoodsByBusiness))
	route.GET("/require/approved", permission.AuthRequired(requireApproved))

	foodService.InitService()
}

// @Title createFood
// @Description Create a food.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   name       	form   	string  true	"Food name."
// @Success 200 {object} model.Food             "Returns created food"
// @Failure 400 {object} response.BasicResponse "err.food.bind"
// @Failure 400 {object} response.BasicResponse "err.food.create"
// @Resource /foods
// @Router /foods [post]
func createFood(c echo.Context) error {
	food := &model.Food{}
	if err := c.Bind(food); err != nil {
		return response.KnownErrJSON(c, "err.food.bind", err)
	}
	// Create food
	food, err := foodService.CreateFood(food)
	if err != nil {
		return response.KnownErrJSON(c, "err.food.create", err)
	}
	return response.SuccessInterface(c, food)
}

// @Title readFood
// @Description Read a food.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Food ID."
// @Success 200 {object} model.Food 		"Returns read food"
// @Failure 400 {object} response.BasicResponse "err.food.bind"
// @Failure 400 {object} response.BasicResponse "err.food.read"
// @Resource /foods
// @Router /foods/{id} [get]
func readFood(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	// Retrieve rider by id
	food, err := foodService.ReadFood(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.food.read", err)
	}
	return response.SuccessInterface(c, food)
}

// @Title updateFood
// @Description Update food of users.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Food ID."
// @Success 200 {object} model.Food	 		"Returns updated food"
// @Failure 400 {object} response.BasicResponse "err.food.bind"
// @Failure 400 {object} response.BasicResponse "err.food.update"
// @Resource /foods
// @Router /foods/{id} [put]
func updateFood(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	food := &model.Food{}
	if err := c.Bind(food); err != nil {
		return response.KnownErrJSON(c, "err.food.bind", err)
	}

	// Update food
	food, err := foodService.UpdateFood(objid, food)
	if err != nil {
		return response.KnownErrJSON(c, "err.food.update", err)
	}
	return response.SuccessInterface(c, food)
}

// @Title deleteFood
// @Description Delete a food.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"User ID."
// @Success 200 {object} response.BasicResponse "User is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.user.bind"
// @Failure 400 {object} response.BasicResponse "err.user.delete"
// @Resource /users
// @Router /users/{id} [delete]
func deleteFood(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))

	// Remove rider with object id
	err := foodService.DeleteFood(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.food.delete", err)
	}
	return response.SuccessJSON(c, "Food is deleted correctly.")
}

// @Title readFoods
// @Description Read a foods.
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
// @Resource /users
// @Router /users/{id} [delete]
func readFoods(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))
	mealKindCode, _ := strconv.Atoi(c.FormValue("mealKindCode"))
	var businessID bson.ObjectId
	if bson.IsObjectIdHex(c.FormValue("businessId")) {
		businessID = bson.ObjectIdHex(c.FormValue("businessId"))
	}

	// Read foods with query
	foods, total, err := foodService.ReadFoods(query, offset, count, field, sort, businessID, mealKindCode)
	if err != nil {
		return response.KnownErrJSON(c, "err.food.read", err)
	}
	// retreive base data
	for _, food := range foods {
		foodService.RetrieveFoodBaseStructure(food)
	}

	return response.SuccessInterface(c, &model.ListForm{total, foods})
}

func readFoodsByBusiness(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))
	mealKindCode, _ := strconv.Atoi(c.FormValue("mealKindCode"))
	var businessID bson.ObjectId
	if bson.IsObjectIdHex(c.FormValue("businessId")) {
		businessID = bson.ObjectIdHex(c.FormValue("businessId"))
	}

	results := []*model.BusinessFood{}
	// Most popular foods
	if foods, err := foodService.ReadMostPopularFoods(businessID, mealKindCode); err == nil && len(foods) > 0 {
		results = append(results, &model.BusinessFood{
			"Most Popular",
			foods,
		})
	}
	// Recommend foods
	if foods, err := foodService.ReadRecommendFoods(businessID, mealKindCode); err == nil && len(foods) > 0 {
		results = append(results, &model.BusinessFood{
			"Recommend",
			foods,
		})
	}
	// Read foods with query
	businessFoods, total, err := foodService.ReadFoodsByBusiness(query, offset, count, field, sort, businessID, mealKindCode)
	if err != nil {
		return response.KnownErrJSON(c, "err.food.read", err)
	}
	results = append(results, businessFoods...)

	// retreive base data
	for _, businessFood := range results {
		for _, food := range businessFood.Foods {
			foodService.RetrieveFoodBaseStructure(food)
		}
	}

	return response.SuccessInterface(c, &model.ListForm{total, results})
}

func requireApproved(c echo.Context) error {
	// Read Foods without approved
	foods, err := foodService.ReadWithoutApprove()
	if err != nil {
		return response.KnownErrJSON(c, "err.food.read", err)
	}

	return response.SuccessInterface(c, foods)
}
