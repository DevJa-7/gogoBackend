package model

// File struct.
type File struct {
	Name      string `json:"name" description:"File name"`
	Extension string `json:"extension" description:"File name"`
	Path      string `json:"path"  description:"File stored path(url)"`
	Size      int64  `json:"size"  description:"File size(byte)"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt"  description:"File created time"`
}
