package imgConv

import (
	"errors"
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

/* returns Image interface with drawable image inside, type depends on config */
func GetEmptyPicture(sizeX, sizeY int) Image {
	format := config.FormatLookup()
	switch format {
	case "NRGBA":
		return image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	case "RGBA":
		return image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	case "GRAY":
		return image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	case "CMYK":
		return image.NewCMYK(image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	case "RGBA64":
		return image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	case "NRGBA64":
		return image.NewNRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	default:
		return image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	}
}

/*checks if underlying image is drawable, if not - replaces it with one */
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

/* corrects limit size in case of overflow */
func CalcAverageChunk(x, y, size int, img Image) Pixel {
	var limitX int
	var limitY int

	if x+size > img.Bounds().Max.X {
		limitX = img.Bounds().Max.X
	} else {
		limitX = x + size
	}
	if y+size > img.Bounds().Max.Y {

		limitY = img.Bounds().Max.Y
	} else {
		limitY = y + size
	}
	return GetAveragePixel(img, x, y, limitX, limitY)
}

/* calculates average color of a given chunk.*/
func GetAveragePixel(pic Image, dx, dy, maxx, maxy int) (av Pixel) {
	for x := dx; x < maxx; x++ {
		for y := dy; y < maxy; y++ {
			r, g, b, _ := pic.At(x, y).RGBA()
			av.R += r / 255
			av.G += g / 255
			av.B += b / 255
		}
	}
	imgArea := uint32((maxx - dx) * (maxy - dy))
	av.R, av.G, av.B, av.A = av.R/imgArea, av.G/imgArea, av.B/imgArea, 255
	return av
}
