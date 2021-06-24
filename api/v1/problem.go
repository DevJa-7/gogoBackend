package v1

import (
	"net/http"
	"strconv"

	"../../api/response"
	"../../config"
	"../../model"
	"../../service/authService/permission"
	"../../service/problemService"
	"../../service/vehicleService"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitProblems inits problem CRUD apis
// @Title Problems
// @Description Problems's router group.
func InitProblems(parentRoute *echo.Group) {
	route := parentRoute.Group("/problems")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createProblem))
	route.GET("/:id", permission.AuthRequired(readProblem))
	route.PUT("/:id", permission.AuthRequired(updateProblem))
	route.DELETE("/:id", permission.AuthRequired(deleteProblem))

	route.GET("", permission.AuthRequired(readProblems))
	route.GET("/business/:id", permission.AuthRequired(readMerchantProblems))

	route.GET("/resolved", permission.AuthRequired(readResolveProblems))

	vehicleService.InitService()
}

//------------
// CRUD Handlers
//------------

// @Title createProblem
// @Description Create a problem.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Success 200 {object} model.Problem	 	"Returns created problem"
// @Failure 400 {object} response.BasicResponse "err.problem.bind"
// @Failure 400 {object} response.BasicResponse "err.problem.create"
// @Resource /problems
// @Router /problems [post]
func createProblem(c echo.Context) error {
	problem := &model.Problem{}
	if err := c.Bind(problem); err != nil {
		return response.KnownErrJSON(c, "err.problem.bind", err)
	}

	// Create problem
	publicProblem, err := problemService.CreateProblem(problem)
	if err != nil {
		return response.KnownErrJSON(c, "err.problem.create", err)
	}
	return c.JSON(http.StatusOK, publicProblem)
}

// @Title readProblem
// @Description Read a problem.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Success 200 {object} model.Problem 		"Returns created problem"
// @Failure 400 {object} response.BasicResponse "err.problem.bind"
// @Failure 400 {object} response.BasicResponse "err.problem.read"
// @Resource /problems
// @Router /problems/{id} [get]
func readProblem(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))

	// Retrieve problem by id
	publicProblem, err := problemService.ReadProblem(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.problem.read", err)
	}
	return response.SuccessInterface(c, publicProblem)
}

// @Title updateProblem
// @Description Update a problem.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Success 200 {object} model.Location 		"Returns created problem"
// @Failure 400 {object} response.BasicResponse "err.problem.bind"
// @Failure 400 {object} response.BasicResponse "err.problem.read"
// @Resource /problems
// @Router /problems/{id} [put]
func updateProblem(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	problem := &model.Problem{}
	if err := c.Bind(problem); err != nil {
		return response.KnownErrJSON(c, "err.problem.bind", err)
	}

	// Update problem
	publicProblem, err := problemService.UpdateProblem(objid, problem)
	if err != nil {
		return response.KnownErrJSON(c, "err.problem.update", err)
	}
	return response.SuccessInterface(c, publicProblem)
}

// @Title deleteProblem
// @Description Delete a problem.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Location ID."
// @Success 200 {object} response.BasicResponse "Location is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.problem.bind"
// @Failure 400 {object} response.BasicResponse "err.problem.read"
// @Resource /problems
// @Router /problems/{id} [delete]
func deleteProblem(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	// Remove problem with object id
	err := problemService.DeleteProblem(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.problem.delete", err)
	}
	return response.SuccessJSON(c, "Problem is deleted correctly.")
}

// @Title readProblems
// @Description Read problems with parameters.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string 	true	"Bearer {token}"
// @Param   query			form   	string  false	"Will search string."
// @Param   offset			form    int		false	"Offset for pagination."
// @Param   count 			form    int		false	"Count that will show per page."
// @Param   field			form    string  false	"Sort field."
// @Param   sort			form    int		false	"Sort direction. 0:default, 1:Ascending, -1:Descending"
// @Success 200 {object} model.ListForm 				"Retrieve all problems with parameters."
// @Failure 400 {object} response.BasicResponse "err.driver.read"
// @Resource /problems
// @Router /problems [get]
func readProblems(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))

	// Read problems with query
	problems, total, err := problemService.ReadProblems(query, offset, count, field, sort)
	if err != nil {
		return response.KnownErrJSON(c, "err.problem.read", err)
	}

	return response.SuccessInterface(c, model.ListForm{total, problems})
}

func readMerchantProblems(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))

	// Read problems with query
	problems, err := problemService.ReadMerchantProblems(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.problem.read", err)
	}

	return response.SuccessInterface(c, problems)
}

func readResolveProblems(c echo.Context) error {
	// Read problems with query
	problems, err := problemService.ReadResolveProblems()
	if err != nil {
		return response.KnownErrJSON(c, "err.problem.read", err)
	}

	return response.SuccessInterface(c, problems)
}
