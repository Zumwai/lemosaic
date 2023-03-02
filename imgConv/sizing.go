package imgConv

import (
	//"fmt"
	"image"
	//"mosaic/imgConv"
)

/* resizes image in memory and returns it in PNG format */
func ResizeInMemory(src image.Image, sizeX, sizeY int, interpolMethod string) (*image.NRGBA, error) {
	new := image.NewNRGBA(image.Rect(0, 0, sizeX, sizeY))
	err := ApplyInterpol(src, new, interpolMethod)
	if err != nil {
		return nil, err
	}
	return new, nil
}

/* calculates nearest integer(new < num) divisible by div */
func CorrectImageSize(num, div int) int {
	var quotient int
	quotient = int(num / div)
	return int(div * quotient)
}

/* corrects limit size in case of overflow */

func CalcAverageChunk(x, y, size int, img *image.NRGBA) Pixel {
	var limitX int
	var limitY int

	if x+size > img.Rect.Max.X {
		limitX = img.Rect.Max.X
	} else {
		limitX = x + size
	}
	if y+size > img.Rect.Max.Y {
		limitY = img.Rect.Max.Y
	} else {
		limitY = y + size
	}
	return GetAveragePixel(img, x, y, limitX, limitY)
}

/*
	func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
		return Pixel{r / 257, g / 257, b / 257, a / 257}
	}
*/

/* calculates average color of a given chunk. Needs testing for разворачивание цикла*/

func GetAveragePixel(pic *image.NRGBA, dx, dy, maxx, maxy int) (av Pixel) {
	for x := dx; x < maxx; x++ {
		//		srcPixOffset := pic.PixOffset(x, 0)
		//  	dstPixOffset := dst.PixOffset(x, 0)
		for y := dy; y < maxy; y++ {
			offset := pic.PixOffset(x, y)
			col := pic.Pix[offset : offset+4 : offset+4]
			av.R += uint32(col[0])
			av.G += uint32(col[1])
			av.B += uint32(col[2])
			/*
				offset = pic.PixOffset(x, y+1)
				col = pic.Pix[offset : offset+4 : offset+4]
				av.R += uint32(col[0])
				av.G += uint32(col[1])
				av.B += uint32(col[2])
				offset = pic.PixOffset(x, y+2)
				col = pic.Pix[offset : offset+4 : offset+4]
				av.R += uint32(col[0])
				av.G += uint32(col[1])
				av.B += uint32(col[2])
				offset = pic.PixOffset(x, y+3)
				col = pic.Pix[offset : offset+4 : offset+4]
				av.R += uint32(col[0])
				av.G += uint32(col[1])
				av.B += uint32(col[2])
			*/
		}
	}
	//fmt.Println(av)
	imgArea := uint32((maxx - dx) * (maxy - dy))
	//fmt.Println(imgArea)
	av.R, av.G, av.B, av.A = av.R/imgArea, av.G/imgArea, av.B/imgArea, 255
	//fmt.Println(av)
	return av
}
