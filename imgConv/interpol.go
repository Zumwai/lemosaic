package imgConv

import (
	"errors"
	//	"fmt"
	"golang.org/x/image/draw"
	"image"
	// "mosaic/config"
)

type ImgInfo struct {
	Av     Pixel
	Square *image.NRGBA
}

type Pixel struct {
	R uint32
	G uint32
	B uint32
	A uint32
}

/* copies image with resizing and required quality*/
func ApplyInterpol(src image.Image, dst *image.NRGBA, interpolMethod string) error {
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
