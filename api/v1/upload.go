package v1

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"../../config"
	"../../model"
	"../../util/random"
	"../../util/timeHelper"
	"../response"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// InitUpload inits upload file apis
// @Title Upload
// @Description Upload's router group.
func InitUpload(parentRoute *echo.Group) {
	route := parentRoute.Group("/uploads")
	route.POST("/image", uploadImage)

	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))
}

// @Title uploadImage
// @Description Upload a image.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   path       		form   	string  true	"Path to upload image"
// @Param   file			form   	file 	true	"Buffer of image."
// @Success 200 {object} model.File 			"Returns uploaded image"
// @Failure 400 {object} response.BasicResponse "err.file.create"
// @Failure 400 {object} response.BasicResponse "err.file.open"
// @Failure 400 {object} response.BasicResponse "err.file.copy"
// @Resource /upload
// @Router /upload/image [post]
func uploadImage(c echo.Context) error {
	path := c.FormValue("path")
	path = "images/" + path + "/"

	// create directory to upload on local path
	if err := os.MkdirAll(config.UploadLocalPath()+path, 0777); err != nil {
		fmt.Println("here", err)
		return response.KnownErrJSON(c, "err.file.create", err)
	}
	// source
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return response.KnownErrJSON(c, "err.file.create", err)
	}
	// open
	src, err := file.Open()
	if err != nil {
		fmt.Println(err)
		return response.KnownErrJSON(c, "err.file.open", err)
	}
	defer src.Close()
	// destination
	path += random.GenerateRandomString(12) + filepath.Ext(file.Filename)
	fmt.Println(path)
	dst, err := os.Create(config.UploadLocalPath() + path)
	if err != nil {
		fmt.Println(err)
		return response.KnownErrJSON(c, "err.file.create", err)
	}
	defer dst.Close()
	// copy
	if _, err := io.Copy(dst, src); err != nil {
		fmt.Println(err)
		return response.KnownErrJSON(c, "err.file.copy", err)
	}

	fileInfo, err := os.Stat(config.UploadLocalPath() + path)
	if err != nil {
		fmt.Println(err)
		return response.KnownErrJSON(c, "err.file.open", err)
	}

	return response.SuccessInterface(c, &model.File{
		Name:      fileInfo.Name(),
		Extension: filepath.Ext(file.Filename),
		Path:      path,
		Size:      fileInfo.Size(),
		CreatedAt: timeHelper.GetCurrentTime(),
	})
}
