// @APIVersion 1.0.0
// @Title gogo API
// @Description gogo API usually works as expected. But sometimes its not true
// @Contact tiandage719@outlook.com
// @TermsOfServiceUrl http://google.com/
// @License BSD
// @LicenseUrl http://opensource.org/licenses/BSD-2-Clause
// @BasePath http://127.0.0.1:8000/api/v1
// @SubApi Auth management API [/]
// @SubApi Admins management API [/admins]
// @SubApi Businesses management API [/businesses]
// @SubApi Drivers management API [/drivers]
// @SubApi Users management API [/users]
// @SubApi Location management API [/locations]
// @SubApi Upload management API [/upload]

package api

import (
	"../config"
	"./v1"
	"./v1/auth"
	"./v1/base"

	"github.com/labstack/echo"
)

// RouteAPI contains router groups for API
func RouteAPI(parentRoute *echo.Echo) {
	route := parentRoute.Group(config.APIURL)
	{
		auth.Init(route)
		base.InitBrand(route)
		base.InitColor(route)
		base.InitDietary(route)
		base.InitMealKind(route)
		base.InitHelp(route)

		v1.InitDashboard(route)
		v1.InitAdmins(route)
		v1.InitDrivers(route)
		v1.InitVehicle(route)
		v1.InitUsers(route)
		v1.InitLocation(route)
		v1.InitDocument(route)
		v1.InitBusinesses(route)
		v1.InitOrder(route)
		v1.InitFoodType(route)
		v1.InitUpload(route)
		v1.InitRole(route)
		v1.InitURLGroup(route)
		v1.InitFood(route)
		v1.InitProblems(route)
		v1.InitReason(route)
		v1.InitAds(route)
	}
}
