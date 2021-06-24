package authService

import (
	"time"

	"../../config"
	"../../model"
	smtp "../../util/email"
	"../../util/random"
	"../../util/timeHelper"
	"../../util/twillo"
	"./adminService"
	"./businessService"
	"./driverService"
	"./userService"

	"github.com/jinzhu/now"
	"gopkg.in/mgo.v2/bson"
)

// SendVerifyCode handle client email to verify codesss
func SendVerifyCode(email, role, method string) bool {
	// generate verify code to reset password
	var fullname string
	var phone string
	verifyCode := random.GenerateRandomDigitString(4)
	ql := bson.M{
		"$set": bson.M{
			"verify.code":      verifyCode,
			"verify.createdAt": timeHelper.GetCurrentTime(),
			"updatedAt":        timeHelper.GetCurrentTime(),
		},
	}

	switch role {
	case config.RoleAdmin:
		if admin, err := adminService.ReadAdminByEmail(email); err == nil {
			adminService.UpdateVerifyCode(admin.ID, ql)
			fullname = admin.Fullname()
		} else {
			return false
		}
	case config.RoleBusiness:
		if business, err := businessService.ReadBusinessByEmail(email); err == nil {
			businessService.UpdateVerifyCode(business.ID, ql)
			fullname = business.Name
			phone = business.PhoneCode + business.Phone
		} else {
			return false
		}
	case config.RoleDriver:
		if driver, err := driverService.ReadDriverByEmail(email); err == nil {
			driverService.UpdateVerifyCode(driver.ID, ql)
			fullname = driver.Fullname()
			phone = driver.PhoneCode + driver.Phone
		} else {
			return false
		}
	case config.RoleUser:
		if user, err := userService.ReadUserByEmail(email); err == nil {
			userService.UpdateVerifyCode(user.ID, ql)
			fullname = user.Fullname()
			phone = user.PhoneCode + user.Phone
		} else {
			return false
		}
	}

	if method == config.EmailMethod {
		// send forgot email to user email
		go smtp.SendVerifyEmail(email, fullname, verifyCode)
	} else if method == config.TwilloMethod {
		go twillo.SendVerifySMS(phone, verifyCode)
	}
	return true
}

// CheckVerifyCode checks that exists email and verifyCode
func CheckVerifyCode(email string, verifyCode string, role string) (bson.ObjectId, interface{}, error) {
	// calculate over time
	timeup := time.Now()
	timeup = timeup.Add(-time.Minute)
	timeup = now.New(timeup).BeginningOfMinute()
	// check verify code with email
	ql := bson.M{
		"email":            email,
		"verify.code":      verifyCode,
		"verify.createdAt": bson.M{"$gte": timeup.Unix()},
	}

	var objid bson.ObjectId
	var result interface{}
	var err error

	switch role {
	case config.RoleAdmin:
		if admin, e := adminService.CheckVerifyCode(ql); e == nil {
			objid = admin.ID
			result = model.PublicAdmin{Admin: admin}
		} else {
			err = e
		}
	case config.RoleBusiness:
		if business, e := businessService.CheckVerifyCode(ql); e == nil {
			objid = business.ID
			result = business
		} else {
			err = e
		}
	case config.RoleDriver:
		if driver, e := driverService.CheckVerifyCode(ql); e == nil {
			objid = driver.ID
			result = model.PublicDriver{Driver: driver}
		} else {
			err = e
		}
	case config.RoleUser:
		if user, e := userService.CheckVerifyCode(ql); e == nil {
			objid = user.ID
			result = model.PublicUser{User: user}
		} else {
			err = e
		}
	}

	return objid, result, err
}

// ReadMe returns self profile
func ReadMe(objid bson.ObjectId, role string) (interface{}, error) {
	var result interface{}
	var err error

	switch role {
	case config.RoleAdmin:
		if admin, err := adminService.ReadAdmin(objid); err == nil {
			result = model.PublicAdmin{Admin: admin}
		}
	case config.RoleBusiness:
		/*		if business, err := businessService.ReadBusiness(objid); err == nil {
				//			result = model.PublicBusiness{Business: business}
			}*/
	case config.RoleDriver:
		if driver, err := driverService.ReadDriver(objid); err == nil {
			result = model.PublicDriver{Driver: driver}
		}
	case config.RoleUser:
		if user, err := userService.ReadUser(objid); err == nil {
			result = model.PublicUser{User: user}
		}
	}

	return result, err
}

// ResetPassword returns self profile after reset password
func ResetPassword(objid bson.ObjectId, role string, pwd *model.Password) (interface{}, error) {
	var result interface{}
	var err error

	switch role {
	case config.RoleAdmin:
		if admin, err := adminService.ResetPassowrd(objid, pwd); err == nil {
			result = model.PublicAdmin{Admin: admin}
		}
	case config.RoleBusiness:
		if business, err := businessService.ResetPassword(objid, pwd); err == nil {
			result = business
		}
	case config.RoleDriver:
		if driver, err := driverService.ResetPassword(objid, pwd); err == nil {
			result = model.PublicDriver{Driver: driver}
		}
	case config.RoleUser:
		if user, err := userService.ResetPassword(objid, pwd); err == nil {
			result = model.PublicUser{User: user}
		}
	}

	return result, err
}
