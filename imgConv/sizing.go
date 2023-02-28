package imgConv

import (
	"image"
	"math"
)

func ResizeInMemory(src image.Image, sizeX, sizeY int, interpolMethod string) (image.Image, error) {
	new := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))
	err := ApplyInterpol(src, new, interpolMethod)
	if err != nil {
		return nil, err
	}
	return new, nil
}

/* calculates nearest integer(new < num) divisible by div */

func CorrectImageSize(num, div int) int {
	quotient := num / div
	return div * quotient
}

/* calcs average color of a square (x or y) + size in a given image. May be redone into int calc for speed*/
func CalcAverageChunk(x, y, size int, img image.Image) (tuple Pixel) {
	var frameX int = x + size
	var frameY int = y + size
	//tmp := img.(*image.YCbCr)
	//start := uintptr(unsafe.Pointer(&tmp.Pix[0]))
	tmp := img.(*image.RGBA)
	//var tmp image.Image
	/*	if _, ok := img.(*image.RGBA); ok {
			tmp = img.(*image.RGBA)
		} else if _, ok := img.(*image.YCbCr); ok {
			tmp = leconvert(img)
		}
	*/
	for dx := x; dx < frameX; dx++ {
		for dy := y; dy < frameY; dy++ {

			offset := tmp.PixOffset(dx, y)
			oops := tmp.Pix[offset : offset+8 : offset+16]
			tuple.R += float64(oops[0])
			tuple.G += float64(oops[1])
			tuple.B += float64(oops[2])
			/*
				tuple.R += float64(*(*uint8)((unsafe.Pointer)(start + uintptr(dstPixOffset))))
				tuple.G += float64(*(*uint8)((unsafe.Pointer)(start + uintptr(dstPixOffset+8))))
				tuple.B += float64(*(*uint8)((unsafe.Pointer)(start + uintptr(dstPixOffset+16))))
			*/
		}
	}
	imgArea := float64(size * size)
	tuple.R = math.Round(tuple.R / imgArea)
	tuple.G = math.Round(tuple.G / imgArea)
	tuple.B = math.Round(tuple.B / imgArea)
	tuple.A = 255
	return tuple
}

func RgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{float64(r / 257), float64(g / 257), float64(b / 257), float64(a / 257)}
}
