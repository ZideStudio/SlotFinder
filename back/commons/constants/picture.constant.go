package constants

type PictureFormat string

const (
	PICTURE_FORMAT_JPEG PictureFormat = "jpeg"
	PICTURE_FORMAT_PNG  PictureFormat = "png"
	PICTURE_FORMAT_WEBP PictureFormat = "webp"
)

var ALLOWED_PICTURE_FORMATS = []PictureFormat{
	PICTURE_FORMAT_JPEG,
	PICTURE_FORMAT_PNG,
	PICTURE_FORMAT_WEBP,
}
