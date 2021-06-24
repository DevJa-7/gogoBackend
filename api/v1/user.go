package v1

import (
	"errors"
	"strconv"

	"../../config"
	"../../model"
	"../../service/authService/businessService"
	"../../service/authService/permission"
	"../../service/authService/userService"
	"../response"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitUsers inits user CRUD apis
// @Title Users
// @Description Users's router group.
func InitUsers(parentRoute *echo.Group) {
	route := parentRoute.Group("/users")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createUser))
	route.GET("/:id", permission.AuthRequired(readUser))
	route.PUT("/:id", permission.AuthRequired(updateUser))
	route.DELETE("/:id", permission.AuthRequired(deleteUser))

	route.GET("", permission.AuthRequired(readUsers))

	// general api
	route.GET("/:id/favorite", permission.AuthRequired(readFavorites))
	route.PUT("/:id/favorite", permission.AuthRequired(updateFavorite))

	userService.InitService()
}

// @Title createUser
// @Description Create a user.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   email       	form   	string  true	"User Email."
// @Param   password		form   	string 	true	"User Password."
// @Success 200 {object} model.PublicUser 		"Returns created user"
// @Failure 400 {object} response.BasicResponse "err.user.bind"
// @Failure 400 {object} response.BasicResponse "err.user.create"
// @Resource /users
// @Router /users [post]
func createUser(c echo.Context) error {
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		return response.KnownErrJSON(c, "err.user.bind", err)
	}

	// create user
	user, err := userService.CreateUser(user)
	if err != nil {
		return response.KnownErrJSON(c, "err.user.create", err)
	}

	publicUser := &model.PublicUser{User: user}
	return response.SuccessInterface(c, publicUser)
}

// @Title readUser
// @Description Read a user.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"User ID."
// @Success 200 {object} model.PublicUser 		"Returns read user"
// @Failure 400 {object} response.BasicResponse "err.user.bind"
// @Failure 400 {object} response.BasicResponse "err.user.read"
// @Resource /users
// @Router /users/{id} [get]
func readUser(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.user.bind", errors.New("Retreived object id is invalid"))
	}

	objid := bson.ObjectIdHex(c.Param("id"))
	user, err := userService.ReadUser(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.user.read", err)
	}

	publicUser := &model.PublicUser{User: user}
	return response.SuccessInterface(c, publicUser)
}

// @Title updateUser
// @Description Update a user.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   id				path   	string  true	"User ID."
// @Param   avatar      	form   	string  true	"User Avatar"
// @Param   firstname		form   	string  true	"User Firstname"
// @Param   lastname		form   	string  true	"User Lastname"
// @Param   email	    	form   	string  true	"User Email"
// @Param   birth      		form   	Time   	true	"User Birth"
// @Success 200 {object} model.PublicUser 		"Returns read user"
// @Failure 400 {object} response.BasicResponse "err.user.bind"
// @Failure 400 {object} response.BasicResponse "err.user.read"
// @Resource /users
// @Router /users/{id} [put]
func updateUser(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.user.bind", errors.New("Retreived object id is invalid"))
	}

	user := &model.User{}
	if err := c.Bind(user); err != nil {
		return response.KnownErrJSON(c, "err.user.bind", err)
	}

	objid := bson.ObjectIdHex(c.Param("id"))
	user, err := userService.UpdateUser(objid, user)
	if err != nil {
		return response.KnownErrJSON(c, "err.user.read", err)
	}

	publicUser := &model.PublicUser{User: user}
	return response.SuccessInterface(c, publicUser)
}

// @Title deleteUser
// @Description Delete a user.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"User ID."
// @Success 200 {object} response.BasicResponse "User is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.user.bind"
// @Failure 400 {object} response.BasicResponse "err.user.delete"
// @Resource /users
// @Router /users/{id} [delete]
func deleteUser(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.user.bind", errors.New("Retreived object id is invalid"))
	}

	// delete user with object id
	objid := bson.ObjectIdHex(c.Param("id"))
	err := userService.DeleteUser(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.user.delete", err)
	}

	return response.SuccessJSON(c, "User is deleted correctly.")
}

// @Title readUsers
// @Description Read users with parameters.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string 	true	"Bearer {token}"
// @Param   query			form   	string  false	"Will search string."
// @Param   offset			form    int		false	"Offset for pagination."
// @Param   count 			form    int		false	"Count that will show per page."
// @Param   field			form    string  false	"Sort field."
// @Param   sort			form    int		false	"Sort direction. 0:default, 1:Ascending, -1:Descending"
// @Success 200 {object} UsersForm 				"Returned list users."
// @Failure 400 {object} response.BasicResponse "err.user.read"
// @Resource /users
// @Router /users [get]
func readUsers(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))

	// read users with params
	users, total, err := userService.ReadUsers(query, offset, count, field, sort)
	if err != nil {
		return response.KnownErrJSON(c, "err.user.read", err)
	}

	// retreive by public format
	publicUsers := []*model.PublicUser{}
	for _, user := range users {
		publicUser := &model.PublicUser{User: user}
		publicUsers = append(publicUsers, publicUser)
	}

	return response.SuccessInterface(c, &model.ListForm{total, publicUsers})
}

func readFavorites(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.user.bind", errors.New("Retreived object id is invalid"))
	}
	// get user id
	objid := bson.ObjectIdHex(c.Param("id"))

	businesses, err := userService.ReadUserFavBusiness(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.favorite.read", err)
	}
	for _, b := range businesses {
		businessService.RetrieveBusinessBaseStructure(b)
	}

	return response.SuccessInterface(c, businesses)
}

func updateFavorite(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.user.bind", errors.New("Retreived object id is invalid"))
	}
	// get user id
	objid := bson.ObjectIdHex(c.Param("id"))
	// get business id
	var businessID bson.ObjectId
	if err := c.Bind(&businessID); err != nil {
		return response.KnownErrJSON(c, "err.business.bind", err)
	}

	if err := userService.UpdateFavorite(objid, businessID); err != nil {
		return response.KnownErrJSON(c, "err.user.udpate", err)
	}

	return response.SuccessJSON(c, "Favorite is modified correctly.")
}
