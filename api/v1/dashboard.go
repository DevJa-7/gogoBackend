package v1

import (
	"../../config"
	"../../service/authService/adminService"
	"../../service/authService/businessService"
	"../../service/authService/driverService"
	"../../service/authService/permission"
	"../../service/authService/userService"
	"../response"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// InitDashboard initialze dashboard api
// @Title Dashboard
// @Description Dashboard router group.
func InitDashboard(parentRoute *echo.Group) {
	route := parentRoute.Group("/dashboard")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.GET("/general", permission.AuthRequired(readGeneralData))
}

// @Title readGeneralData
// @Description ScheduleDashboard
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Success 200 {object} 	interface				"Returns users, businesses, drivers count"
// @Resource /dashboard
// @Router /dashboard/general [get]
func readGeneralData(c echo.Context) error {
	// read processes
	admins, adminAvailables := adminService.ReadCounts()
	users, userAvailables := userService.ReadCounts()
	drivers, driverAvailables := driverService.ReadCounts()
	businesses, businessAvailables := businessService.ReadCounts()

	return response.SuccessInterface(c, map[string]interface{}{
		"admins":             admins,
		"adminAvaliables":    adminAvailables,
		"users":              users,
		"userAvailables":     userAvailables,
		"drivers":            drivers,
		"driverAvailables":   driverAvailables,
		"businesses":         businesses,
		"businessAvailables": businessAvailables,
	})
}
