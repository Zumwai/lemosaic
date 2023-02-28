package localMosaic

import (
	//"bytes"
	"fmt"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
)

const seekStart = 0 // local const identical in it's meaning to io.SeekStart. Does no need package for this

/* inspects file for type, checks for boundaries. Calls to  Convert by type and returns image/png. Returns error if file is too large or too unexpected*/
func ConvertCorrectType(name string) (image.Image, error) {
	stat, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	size := stat.Size()
	if size > 1e+7 {
		return nil, fmt.Errorf("file is too large %d, max size is 10 mb", size)
	}

	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return nil, err
	}

	// resets reader pointer for decoder
	_, err = file.Seek(seekStart, seekStart)
	if err != nil {
		return nil, err
	}
	format := http.DetectContentType(buff)
	dst, err := decodeByType(format, file)
	if format != "image/png" {
		dst = leconvert(dst)
	}
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func decodeByType(format string, file *os.File) (dst image.Image, err error) {
	switch format {
	case "image/png":
		return png.Decode(file)
	case "image/jpeg":
		return jpeg.Decode(file)
	case "image/gif":
		return gif.Decode(file)
	case "image/tiff":
		return tiff.Decode(file)
	case "image/webp":
		return webp.Decode(file)
	default:
		return nil, fmt.Errorf("unrecognized- %s", format)
	}
}
