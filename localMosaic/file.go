package localMosaic

import (
	"fmt"

	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"golang.org/x/image/tiff"

	// "golang.org/x/image/webp"
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
	x := config.SizeLookup()
	if x == 0 {
		return imgConv.ConvertToDrawable(tmp), nil
	}
	return imgConv.ConvertWithResizing(tmp, x), nil
}

/* encodes to local file, type of resulting image depends on config */
func EncodeToFile(path, name, suffix string, dst imgConv.Image) error {
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
			//BufferPool:       pool,
		}
		return enc.Encode(newFile, dst)
	case "jpeg":
		return jpeg.Encode(newFile, dst, &jpeg.Options{Quality: config.JpegQualityLookup()})
	case "gif":
		return gif.Encode(newFile, dst, nil)
	case "tiff":
		return tiff.Encode(newFile, dst, &tiff.Options{Compression: tiff.Deflate, Predictor: false})
	case "webp":
		options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 80)
		options.Lossless = true
		if err != nil {
			fmt.Println(err)
		}
		return webp.Encode(newFile, dst, options)
	default:
		return fmt.Errorf("unrecognized format - %s", format)
	}
}

/* opens given filename and inspects file for type. Returns error if file is too large, too unexpected or unavailable or if something went wrong*/
func getUnformattedImage(name string) (image.Image, error) {
	stat, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	if config.UnmaxLookup() {
		size := stat.Size()
		if size > 1e+7 {
			return nil, fmt.Errorf("file is too large %d, max size is 10 mb, use -unmax flag to unset", size)
		}
	}
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	dst, err := DecodeByType(file)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

/* calls corresponding decoder, depends on image format, returns image.Image with an undefined underlying actual type */
func DecodeByType(file io.ReadSeeker) (dst image.Image, err error) {
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return nil, err
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	format := http.DetectContentType(buff)
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
		return webp.Decode(file, &decoder.Options{})
	default:
		return nil, fmt.Errorf("unrecognized - %s", format)
	}
}
