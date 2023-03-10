package imgConv

import (
	"errors"
	//"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"mosaic/config"
)

/* similar to draw.Image. Declared for convinience */
type Image interface {
	image.Image
	Set(x, y int, c color.Color)
}

/* main struct to hold info about images for mosaic */
type ImgInfo struct {
	Av     Pixel
	Square Image
}

/*uint32 values used to calculate average colors */
type Pixel struct {
	R uint32
	G uint32
	B uint32
	A uint32
}

/* copies image, depends on required size and quality*/
func ApplyInterpol(src image.Image, dst Image, newRect image.Rectangle) error {
	interpolMethod := config.InterpolLookup()

	switch interpolMethod {
	case "ApproxBiLinear":
		draw.ApproxBiLinear.Scale(dst, newRect, src, src.Bounds(), draw.Over, nil)
	case "CatmullRom":
		draw.CatmullRom.Scale(dst, newRect, src, src.Bounds(), draw.Over, nil)
	case "BiLinear":
		draw.BiLinear.Scale(dst, newRect, src, src.Bounds(), draw.Over, nil)
	case "NearestNeighbor":
		draw.NearestNeighbor.Scale(dst, newRect, src, src.Bounds(), draw.Over, nil)
	default:
		return errors.New("incorrect interpolation type")
	}
	return nil
}

/* returns empy Image interface with drawable image inside, type depends on config */
func GetEmptyPicture(sizeX, sizeY int) Image {
	format := config.FormatLookup()
	switch format {
	case "NRGBA":
		return image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	case "RGBA":
		return image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	default:
		return image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	}
}

func ConvertToDrawable(src image.Image) Image {
	ret, ok := src.(draw.Image)
	if !ok {
		tmpPtr, err := ResizeInMemory(src, src.Bounds().Max.X, src.Bounds().Max.Y)
		if err != nil {
			return nil
		}
		return tmpPtr
	}
	return ret
}
