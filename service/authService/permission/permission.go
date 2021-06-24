package permission

import (
	"errors"
	"net/http"
	"time"

	"../../../api/response"
	"../../../config"
	"../../../util/timeHelper"
	"../adminService"
	"../businessService"
	"../driverService"
	"../userService"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"gopkg.in/mgo.v2/bson"
)

// GenerateToken returns token after generate with user
func GenerateToken(objid bson.ObjectId, role string) (string, error) {
	// create token
	token := jwt.New(jwt.SigningMethodHS256)

	// set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["idx"] = objid.Hex()
	claims["exp"] = timeHelper.FewDaysLater(config.AuthTokenExpirationDay)
	claims["role"] = role

	// generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.AuthTokenKey))
	return t, err
}

// InfoFromToken returns idx from token
func InfoFromToken(c echo.Context) (bson.ObjectId, string) {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	var objid bson.ObjectId
	var role string
	// retrieve object id of client
	if claims["idx"] != nil {
		var idx string
		idx = claims["idx"].(string)
		if bson.IsObjectIdHex(idx) {
			objid = bson.ObjectIdHex(idx)
		}

	}
	// retrieve role: admin, user, driver, business
	if claims["role"] != nil {
		role = claims["role"].(string)
	}

	return objid, role
}

func expiredFromToken(c echo.Context) bool {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	exp, err := time.Parse(time.RFC3339, claims["exp"].(string))
	if err != nil {
		return false
	}

	return timeHelper.IsExpired(exp)
}

// AuthRequired run function when user logged in.
func AuthRequired(f func(c echo.Context) error) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// checking expire date
		if expiredFromToken(c) {
			log.Error("Token is expired.")
			return response.KnownErrorJSON(c, http.StatusUnauthorized, "error.token.expire", errors.New("Token is expired"))
		}

		// checking client validation
		{
			objid, role := InfoFromToken(c)
			var err error
			switch role {
			case config.RoleAdmin:
				_, err = adminService.ReadAdmin(objid)
			case config.RoleBusiness:
				_, err = businessService.ReadBusiness(objid)
			case config.RoleDriver:
				_, err = driverService.ReadDriver(objid)
			case config.RoleUser:
				_, err = userService.ReadUser(objid)
			}
			if err != nil {
				log.Error("Auth failed.")
				return response.KnownErrorJSON(c, http.StatusUnauthorized, "error.auth.fail", errors.New("Auth failed"))
			}
		}

		f(c)

		return nil
	}
}
