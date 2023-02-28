package imgConv

import (
	"errors"
	"golang.org/x/image/draw"
	"image"
)

type ImgInfo struct {
	Av     Pixel
	Square image.Image
}

type Pixel struct {
	R float64
	G float64
	B float64
	A float64
}

func ApplyInterpol(src image.Image, dst *image.RGBA, interpolMethod string) error {
	switch interpolMethod {
	case "ApproxBiLinear":
		draw.ApproxBiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	case "CatmullRom":
		draw.CatmullRom.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	case "BiLinear":
		draw.BiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	case "NearestNeighbor":
		draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	default:
		return errors.New("incorrect interpolation type")
	}
	return nil
}
