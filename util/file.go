package util

import "mime/multipart"

type FileUpload struct {
	File   multipart.File        `json:"file"`
	Header *multipart.FileHeader `json:"header"`
	Ext    string                `json:"extension"`
}

func ValidateImage(ext string) bool {

	if ext == ".jpg" ||
		ext == ".jpeg" ||
		ext == ".png" ||
		ext == ".JPG" ||
		ext == ".JPEG" ||
		ext == ".PNG" {
		return true
	}
	return false

}
