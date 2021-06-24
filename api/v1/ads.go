package v1

import (
	"errors"
	"fmt"
	"strconv"

	"../../config"
	"../../model"
	"../../service/adsService"
	"../../service/authService/permission"
	"../response"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitAds inits ads CRUD apis
// @Title Ads
// @Description Ads's router group.
func InitAds(parentRoute *echo.Group) {
	route := parentRoute.Group("/ads")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createAds))
	route.GET("/:id", permission.AuthRequired(readAds))
	route.PUT("/:id", permission.AuthRequired(updateAds))
	route.DELETE("/:id", permission.AuthRequired(deleteAds))

	route.GET("", readAllAds)

	adsService.InitService()
}

// @Title createAds
// @Description Create a ads.
// @Accept  json
// @Produce	json
// @Success 200 {object} model.PublicAds 		"Returns created ads"
// @Failure 400 {object} response.BasicResponse "err.ads.bind"
// @Failure 400 {object} response.BasicResponse "err.ads.create"
// @Resource /ads
// @Router /ads [post]
func createAds(c echo.Context) error {
	ads := &model.Ads{}
	if err := c.Bind(ads); err != nil {
		return response.KnownErrJSON(c, "err.ads.bind", err)
	}
	// create ads
	ads, err := adsService.CreateAds(ads)
	if err != nil {
		return response.KnownErrJSON(c, "err.ads.create", err)
	}

	//	publicAds := &model.PublicAds{Ads: ads}
	return response.SuccessInterface(c, ads)
}

// @Title readAds
// @Description Read a ads.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Ads ID."
// @Success 200 {object} model.PublicAds 		"Returns read ads"
// @Failure 400 {object} response.BasicResponse "err.ads.bind"
// @Failure 400 {object} response.BasicResponse "err.ads.read"
// @Resource /ads
// @Router /ads/{id} [get]
func readAds(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.ads.bind", errors.New("Retreived object id is invalid"))
	}

	objid := bson.ObjectIdHex(c.Param("id"))
	ads, err := adsService.ReadAds(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.ads.read", err)
	}

	//	publicAds := &model.PublicAds{Ads: ads}
	return response.SuccessInterface(c, ads)
}

// @Title updateAds
// @Description Update a ads.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   id				path   	string  true	"Ads ID."
// @Param   email	    	form   	string  true	"Ads Email"
// @Success 200 {object} model.PublicAds 		"Returns read ads"
// @Failure 400 {object} response.BasicResponse "err.ads.bind"
// @Failure 400 {object} response.BasicResponse "err.ads.update"
// @Resource /ads
// @Router /ads/{id} [put]
func updateAds(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.ads.bind", errors.New("Retreived object id is invalid"))
	}

	ads := &model.Ads{}
	if err := c.Bind(ads); err != nil {
		return response.KnownErrJSON(c, "err.ads.bind", err)
	}

	objid := bson.ObjectIdHex(c.Param("id"))
	ads, err := adsService.UpdateAds(objid, ads)
	if err != nil {
		fmt.Println(err)
		return response.KnownErrJSON(c, "err.ads.update", err)
	}

	return response.SuccessInterface(c, ads)
}

// @Title deleteAds
// @Description Delete a ads.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Ads ID."
// @Success 200 {object} response.BasicResponse "Ads is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.ads.bind"
// @Failure 400 {object} response.BasicResponse "err.ads.delete"
// @Resource /ads
// @Router /ads/{id} [delete]
func deleteAds(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.ads.bind", errors.New("Retreived object id is invalid"))
	}

	// delete ads with object id
	objid := bson.ObjectIdHex(c.Param("id"))
	err := adsService.DeleteAds(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.ads.delete", err)
	}

	return response.SuccessJSON(c, "Ads is deleted correctly.")
}

// @Title readAds
// @Description Read ads with parameters.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string 	true	"Bearer {token}"
// @Param   query			form   	string  false	"Will search string."
// @Param   offset			form    int		false	"Offset for pagination."
// @Param   count 			form    int		false	"Count that will show per page."
// @Param   field			form    string  false	"Sort field."
// @Param   sort			form    int		false	"Sort direction. 0:default, 1:Ascending, -1:Descending"
// @Success 200 {object} model.ListForm 				"Ads is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.ads.read"
// @Resource /ads
// @Router /ads [get]
func readAllAds(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))

	// read ads with params
	ads, total, err := adsService.ReadAllAds(query, offset, count, field, sort)
	if err != nil {
		return response.KnownErrJSON(c, "err.ads.read", err)
	}

	return response.SuccessInterface(c, &model.ListForm{total, ads})
}
