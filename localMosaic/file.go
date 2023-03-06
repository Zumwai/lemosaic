package localMosaic

import (
	//"image"
	"fmt"
	"golang.org/x/image/draw"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mosaic/config"
	"mosaic/imgConv"
	"net/http"
	"os"
)

/* decodes given image */
func getDecodedFile(name string) (imgConv.Image, error) {
	tmp, err := getUnformattedImage(name)
	if err != nil {
		return nil, err
	}
	ret := convertToDrawable(tmp)
	return ret, nil
}

func convertToDrawable(src image.Image) imgConv.Image {
	ret, ok := src.(draw.Image)
	if !ok {
		tmpPtr, err := imgConv.ResizeInMemory(src, src.Bounds().Max.X, src.Bounds().Max.Y)
		if err != nil {
			return nil
		}
		return tmpPtr
	}
	return ret
}

/* encodes to local file */
func encodeToFile(path, name, suffix string, dst imgConv.Image) error {
	//format := config.FormatLookup()
	format := config.EncoderLookup()
	newFile, err := os.Create(path + "/" + name + suffix + "." + format)

	if err != nil {
		return err
	}
	defer newFile.Close()

	switch format {
	case "png":
		enc := png.Encoder{
			CompressionLevel: png.BestSpeed,
		}
		return enc.Encode(newFile, dst)
	case "jpeg":
		return jpeg.Encode(newFile, dst, &jpeg.Options{Quality: 50})
	case "gif":
		return gif.Encode(newFile, dst, nil)
	case "tiff":
		return tiff.Encode(newFile, dst, &tiff.Options{Compression: tiff.Deflate, Predictor: false})
	default:
		return fmt.Errorf("unrecognized- %s", format)
	}
}

/* inspects file for type, checks for boundaries. Returns error if file is too large, too unexpected or unavailable*/
func getUnformattedImage(name string) (image.Image, error) {
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
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	format := http.DetectContentType(buff)
	dst, err := DecodeByType(format, file)

	if err != nil {
		return nil, err
	}
	return dst, nil
}

/* calls corresponding decoder, depends on image format, returns image.Image with an undefined underlying actual type */
func DecodeByType(format string, file io.Reader) (dst image.Image, err error) {
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
