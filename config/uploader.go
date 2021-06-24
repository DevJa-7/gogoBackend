package config

import (
	"path/filepath"
	"runtime"
	"time"
)

// Constants for uploader.
const (
	UploadS3ImagePath = "images/"
	//	UploadTarget = LOCAL | S3
	UploadTarget = "LOCAL"
	// UploadTarget  = "S3"
	// UploadBucket = TEST | PRODUCTION
	UploadBucket = "TEST"

	UploadTimeout = 30 * time.Second
)

// UploadLocalPath return upload local path
func UploadLocalPath() string {
	_, currentFile, _, _ := runtime.Caller(0)
	path := filepath.Dir(currentFile)
	path += "/../../upload/"

	return path
}
