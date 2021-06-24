package v1

import (
	"net/http"
	"strconv"

	"../../api/response"
	"../../config"
	"../../model"
	"../../service/authService/permission"
	"../../service/documentService"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitDocument initialze doucment api
func InitDocument(parentRoute *echo.Group) {
	parentRoute.GET("/public/documents", readDocuments)

	route := parentRoute.Group("/documents")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createDocument))
	route.GET("/:id", permission.AuthRequired(readDocument))
	route.PUT("/:id", permission.AuthRequired(updateDocument))
	route.DELETE("/:id", permission.AuthRequired(deleteDocument))

	route.GET("", permission.AuthRequired(readDocuments))

	documentService.InitService()
}

//------------
// CRUD Handlers
//------------

// Create Document
func createDocument(c echo.Context) error {
	document := &model.Document{}
	if err := c.Bind(document); err != nil {
		return response.KnownErrJSON(c, "err.document.bind", err)
	}

	// Create document
	document, err := documentService.CreateDocument(document)
	if err != nil {
		return response.KnownErrJSON(c, "err.document.create", err)
	}
	return c.JSON(http.StatusOK, document)
}

// Read Document
func readDocument(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))

	// Retrieve document by id
	document, err := documentService.ReadDocument(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.document.read", err)
	}
	return response.SuccessInterface(c, document)
}

// Update Document
func updateDocument(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	document := &model.Document{}
	if err := c.Bind(document); err != nil {
		return response.KnownErrJSON(c, "err.document.bind", err)
	}

	// Update document
	document, err := documentService.UpdateDocument(objid, document)
	if err != nil {
		return response.KnownErrJSON(c, "err.document.update", err)
	}
	return response.SuccessInterface(c, document)
}

// Delete Document
func deleteDocument(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	// Remove document with object id
	err := documentService.DeleteDocument(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.document.delete", err)
	}
	return response.SuccessJSON(c, "Document is deleted correctly.")
}

// Read Documents
func readDocuments(c echo.Context) error {
	query := c.FormValue("query")
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))
	filter, _ := strconv.ParseBool(c.FormValue("is_filter")) // this variable is for only mobile app
	// Read documents with query
	documents, err := documentService.ReadDocuments(query, field, sort, filter)
	if err != nil {
		return response.KnownErrJSON(c, "err.document.read", err)
	}

	return response.SuccessInterface(c, documents)
}
