package base

import (
	"fmt"

	"../../../api/response"
	"../../../config"
	"../../../model"
	"../../../service/authService/permission"
	"../../../service/base/brandService"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitBrand inits brand of vehicle CRUD apis
// @Title Brands
// @Description Brand's router group.
func InitBrand(parentRoute *echo.Group) {
	parentRoute.GET("/public/brands", readBrands)

	route := parentRoute.Group("/brands")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createBrand))
	route.GET("/:id", permission.AuthRequired(readBrand))
	route.PUT("/:id", permission.AuthRequired(updateBrand))
	route.DELETE("/:id", permission.AuthRequired(deleteBrand))

	route.GET("", permission.AuthRequired(readBrands))
	route.PUT("/update/model/:id", permission.AuthRequired(updateModels))

	brandService.InitService()
}

// @Title createBrand
// @Description Create a brand.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   name       		form   	string  true	"Dietary name."
// @Success 200 {object} model.Dietary             	"Returns created dietary"
// @Failure 400 {object} response.BasicResponse 	"err.dietary.bind"
// @Failure 400 {object} response.BasicResponse 	"err.dietary.create"
// @Resource /dietaries
// @Router /dietaries [post]
func createBrand(c echo.Context) error {
	brand := &model.Brand{}
	if err := c.Bind(brand); err != nil {
		return response.KnownErrJSON(c, "err.brand.bind", err)
	}
	fmt.Println(brand)
	// Create brand
	brand, err := brandService.CreateBrand(brand)
	if err != nil {
		return response.KnownErrJSON(c, "err.brand.create", err)
	}
	return response.SuccessInterface(c, brand)
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
func readBrand(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	// Retrieve rider by id
	brand, err := brandService.ReadBrand(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.brand.read", err)
	}
	return response.SuccessInterface(c, brand)
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
func updateBrand(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	brand := &model.Brand{}
	if err := c.Bind(brand); err != nil {
		return response.KnownErrJSON(c, "err.brand.bind", err)
	}

	// Update brand
	brand, err := brandService.UpdateBrand(objid, brand)
	if err != nil {
		return response.KnownErrJSON(c, "err.brand.update", err)
	}
	return response.SuccessInterface(c, brand)
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
func deleteBrand(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))

	// Remove rider with object id
	err := brandService.DeleteBrand(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.brand.delete", err)
	}
	return response.SuccessJSON(c, "Brand is deleted correctly.")
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
func readBrands(c echo.Context) error {
	// Read brands with query
	brands, err := brandService.ReadBrands()
	if err != nil {
		return response.KnownErrJSON(c, "err.brand.read", err)
	}

	return response.SuccessInterface(c, brands)
}

// @Title updateModels
// @Description Upddate models of the brand.
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
func updateModels(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	brand := &model.Brand{}
	if err := c.Bind(brand); err != nil {
		return response.KnownErrJSON(c, "err.brand.bind", err)
	}

	err := brandService.UpdateModels(objid, brand.Models)
	if err != nil {
		return response.KnownErrJSON(c, "err.brand.update", err)
	}

	return response.SuccessInterface(c, "Model is updated correctly.")
}
