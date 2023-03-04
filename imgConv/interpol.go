package imgConv

import (
	"errors"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"mosaic/config"
)

type Image interface {
	image.Image
	Set(x, y int, c color.Color)
}

/* main struct to hold info about image */
type ImgInfo struct {
	Av     Pixel
	Square Image
}

/* for calculating colors */
type Pixel struct {
	R uint32
	G uint32
	B uint32
	A uint32
}

/* copies image with resizing and required quality*/
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

type Union interface {
	*image.RGBA | *image.NRGBA
}

func GetEmptyPicture(sizeX, sizeY int) Image {

	ret := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	/*
	   f	ormat := config.FormatLookup()
	   	switch format {
	   	case "NRGBA":
	   		return image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	   	case "RGBA":
	   		return image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	   	default:
	   		return image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	   	}
	   	/*will never happen*/
	return ret
}
