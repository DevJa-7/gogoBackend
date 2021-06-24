package v1

import (
	"errors"
	"fmt"
	"strconv"

	"../../config"
	"../../model"
	"../../service/authService/businessService"
	"../../service/authService/permission"
	"../response"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitBusinesses inits business CRUD apis
// @Title Businesses
// @Description Businesses's router group.
func InitBusinesses(parentRoute *echo.Group) {
	parentRoute.GET("/public/businesses", readBusinesses)

	route := parentRoute.Group("/businesses")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createBusiness))
	route.GET("/:id", permission.AuthRequired(readBusiness))
	route.PUT("/:id", permission.AuthRequired(updateBusiness))
	route.DELETE("/:id", permission.AuthRequired(deleteBusiness))

	route.GET("", readBusinesses)

	route.GET("/:id/dietary", permission.AuthRequired(readBusinessDietary))
	route.GET("/:id/mealKind", permission.AuthRequired(readBusinessMealKind))

	route.POST("/query", permission.AuthRequired(readQueryBusinesses))

	businessService.InitService()
}

// @Title createBusiness
// @Description Create a business.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   email       	form   	string  true	"Business Email."
// @Param   password		form   	string 	true	"Business Password."
// @Success 200 {object} model.PublicBusiness 		"Returns created business"
// @Failure 400 {object} response.BasicResponse "err.business.bind"
// @Failure 400 {object} response.BasicResponse "err.business.create"
// @Resource /businesses
// @Router /businesses [post]
func createBusiness(c echo.Context) error {
	business := &model.Business{}
	if err := c.Bind(business); err != nil {
		return response.KnownErrJSON(c, "err.business.bind", err)
	}
	// create business
	business, err := businessService.CreateBusiness(business)
	if err != nil {
		return response.KnownErrJSON(c, "err.business.create", err)
	}

	//	publicBusiness := &model.PublicBusiness{Business: business}
	return response.SuccessInterface(c, business)
}

// @Title readBusiness
// @Description Read a business.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Business ID."
// @Success 200 {object} model.PublicBusiness 		"Returns read business"
// @Failure 400 {object} response.BasicResponse "err.business.bind"
// @Failure 400 {object} response.BasicResponse "err.business.read"
// @Resource /businesses
// @Router /businesses/{id} [get]
func readBusiness(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.business.bind", errors.New("Retreived object id is invalid"))
	}

	objid := bson.ObjectIdHex(c.Param("id"))
	business, err := businessService.ReadBusiness(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.business.read", err)
	}

	//	publicBusiness := &model.PublicBusiness{Business: business}
	return response.SuccessInterface(c, business)
}

// @Title updateBusiness
// @Description Update a business.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   id				path   	string  true	"Business ID."
// @Param   email	    	form   	string  true	"Business Email"
// @Success 200 {object} model.PublicBusiness 		"Returns read business"
// @Failure 400 {object} response.BasicResponse "err.business.bind"
// @Failure 400 {object} response.BasicResponse "err.business.update"
// @Resource /businesses
// @Router /businesses/{id} [put]
func updateBusiness(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.business.bind", errors.New("Retreived object id is invalid"))
	}

	business := &model.Business{}
	if err := c.Bind(business); err != nil {
		return response.KnownErrJSON(c, "err.business.bind", err)
	}

	objid := bson.ObjectIdHex(c.Param("id"))
	business, err := businessService.UpdateBusiness(objid, business)
	if err != nil {
		fmt.Println(err)
		return response.KnownErrJSON(c, "err.business.update", err)
	}

	//	publicBusiness := &model.PublicBusiness{Business: business}
	return response.SuccessInterface(c, business)
}

// @Title deleteBusiness
// @Description Delete a business.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Business ID."
// @Success 200 {object} response.BasicResponse "Business is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.business.bind"
// @Failure 400 {object} response.BasicResponse "err.business.delete"
// @Resource /businesses
// @Router /businesses/{id} [delete]
func deleteBusiness(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.business.bind", errors.New("Retreived object id is invalid"))
	}

	// delete business with object id
	objid := bson.ObjectIdHex(c.Param("id"))
	err := businessService.DeleteBusiness(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.business.delete", err)
	}

	return response.SuccessJSON(c, "Business is deleted correctly.")
}

// @Title readBusinesses
// @Description Read businesses with parameters.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string 	true	"Bearer {token}"
// @Param   query			form   	string  false	"Will search string."
// @Param   offset			form    int		false	"Offset for pagination."
// @Param   count 			form    int		false	"Count that will show per page."
// @Param   field			form    string  false	"Sort field."
// @Param   sort			form    int		false	"Sort direction. 0:default, 1:Ascending, -1:Descending"
// @Success 200 {object} model.ListForm 				"Business is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.business.read"
// @Resource /businesses
// @Router /businesses [get]
func readBusinesses(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))

	// read businesses with params
	businesses, total, err := businessService.ReadBusinesses(query, offset, count, field, sort)
	if err != nil {
		return response.KnownErrJSON(c, "err.business.read", err)
	}

	return response.SuccessInterface(c, &model.ListForm{total, businesses})
}

func readBusinessDietary(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.business.bind", errors.New("Retreived object id is invalid"))
	}
	// get business objectid with object id
	objid := bson.ObjectIdHex(c.Param("id"))
	dietaries, err := businessService.ReadBusinessDietary(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.business.dietary.read", err)
	}

	return response.SuccessInterface(c, dietaries)
}

func readBusinessMealKind(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.business.bind", errors.New("Retreived object id is invalid"))
	}
	// get business objectid with object id
	objid := bson.ObjectIdHex(c.Param("id"))
	mealKinds, err := businessService.ReadBusinessMealKind(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.business.mealKind.read", err)
	}

	return response.SuccessInterface(c, mealKinds)
}

func readQueryBusinesses(c echo.Context) error {
	queryBusiness := &model.QueryBusiness{}
	if err := c.Bind(queryBusiness); err != nil {
		return response.KnownErrJSON(c, "err.query.bind", err)
	}

	lat := queryBusiness.Lat
	lng := queryBusiness.Lng
	if queryBusiness.Sort == config.SortRecommend && len(queryBusiness.Price) == 0 && len(queryBusiness.Dietary) == 0 {
		// search popular
		popular, _ := businessService.ReadNearbyPopularBusiness(lat, lng)
		excepts := []bson.ObjectId{}
		for _, b := range popular {
			excepts = append(excepts, b.ID)
			businessService.RetrieveBusinessBaseStructure(b)
		}
		// serach recommend
		recommend, _ := businessService.ReadRecommendBusiness(lat, lng)
		for _, b := range recommend {
			excepts = append(excepts, b.ID)
			businessService.RetrieveBusinessBaseStructure(b)
		}
		// serarch unders
		under, _ := businessService.ReadUnder30Business(lat, lng)
		for _, b := range under {
			excepts = append(excepts, b.ID)
			businessService.RetrieveBusinessBaseStructure(b)
		}
		// read more
		more, _ := businessService.ReadMoreBusiness(lat, lng, excepts)
		for _, b := range more {
			businessService.RetrieveBusinessBaseStructure(b)
		}

		return response.SuccessInterface(c, []*model.ListBusiness{
			{config.QueryPopular, true, popular},
			{config.QueryRecommend, true, recommend},
			{config.QueryUnder, true, under},
			{config.QueryAll, false, more},
		})
	}

	businesses, err := businessService.ReadQueryBusiness(queryBusiness)
	if err != nil {
		return response.KnownErrJSON(c, "err.business.mealKind.read", err)
	}
	for _, b := range businesses {
		businessService.RetrieveBusinessBaseStructure(b)
	}

	return response.SuccessInterface(c, []*model.ListBusiness{
		{config.QueryAll, false, businesses},
	})
}
