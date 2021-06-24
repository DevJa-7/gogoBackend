package base

import (
	"../../../api/response"
	"../../../config"
	"../../../model"
	"../../../service/authService/permission"
	"../../../service/base/colorService"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitColor initialze color api
func InitColor(parentRoute *echo.Group) {
	parentRoute.GET("/public/colors", readColors)

	route := parentRoute.Group("/colors")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createColor))
	route.GET("/:id", permission.AuthRequired(readColor))
	route.PUT("/:id", permission.AuthRequired(updateColor))
	route.DELETE("/:id", permission.AuthRequired(deleteColor))

	route.GET("", permission.AuthRequired(readColors))

	colorService.InitService()
}

// @Title createColor
// @Description Create a color.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   name       		form   	string  true	"Dietary name."
// @Success 200 {object} model.Dietary             	"Returns created color"
// @Failure 400 {object} response.BasicResponse 	"err.color.bind"
// @Failure 400 {object} response.BasicResponse 	"err.color.create"
// @Resource /colors
// @Router /colors [post]
func createColor(c echo.Context) error {
	color := &model.Color{}
	if err := c.Bind(color); err != nil {
		return response.KnownErrJSON(c, "err.color.bind", err)
	}
	// Create color
	color, err := colorService.CreateColor(color)
	if err != nil {
		return response.KnownErrJSON(c, "err.color.create", err)
	}
	return response.SuccessInterface(c, color)
}

// readColor
func readColor(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	// Retrieve rider by id
	color, err := colorService.ReadColor(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.color.read", err)
	}
	return response.SuccessInterface(c, color)
}

// updateColor
func updateColor(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	color := &model.Color{}
	if err := c.Bind(color); err != nil {
		return response.KnownErrJSON(c, "err.color.bind", err)
	}

	// Update color
	color, err := colorService.UpdateColor(objid, color)
	if err != nil {
		return response.KnownErrJSON(c, "err.color.update", err)
	}
	return response.SuccessInterface(c, color)
}

// deleteColor
func deleteColor(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))

	// Remove rider with object id
	err := colorService.DeleteColor(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.color.delete", err)
	}
	return response.SuccessJSON(c, "Color is deleted correctly.")
}

// readColors
func readColors(c echo.Context) error {
	// Read colors with query
	colors, err := colorService.ReadColors()
	if err != nil {
		return response.KnownErrJSON(c, "err.color.read", err)
	}

	return response.SuccessInterface(c, colors)
}
