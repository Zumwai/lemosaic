package imgConv

import (
	//	"fmt"
	//
	"image"
	// "mosaic/imgConv"
	"fmt"
	"mosaic/config"
)

/* default opacity to store */
const opacity = 255

type Frame struct {
	X       int
	Y       int
	Routine int
	Step    int
	Size    int
}

/* resizes image in memory and returns it in Picture format
 */
func ResizeInMemory(src image.Image, sizeX, sizeY int) (Image, error) {
	//new := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))
	new := GetEmptyPicture(sizeX, sizeY)
	err := ApplyInterpol(src, new, image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	if err != nil {
		return nil, err
	}
	return new, nil
}

/* calculates nearest integer(new < num) divisible by div */
func CorrectImageSize(num, div int) int {
	var quotient = int(num / div)
	return int(div * quotient)
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

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{r / 257, g / 257, b / 257, a / 257}
}

/* calculates average color of a given chunk. Needs testing for разворачивание цикла*/
func GetAveragePixel(pic Image, dx, dy, maxx, maxy int) (av Pixel) {
	for x := dx; x < maxx; x++ {
		//		srcPixOffset := pic.PixOffset(x, 0)
		//  	dstPixOffset := dst.PixOffset(x, 0)
		for y := dy; y < maxy; y++ {

			col := rgbaToPixel(pic.At(x, y).RGBA())
			av.R += col.R
			av.G += col.G
			av.B += col.B
			/*

				var pic *image.RGBA

				//	fmt.Println(offset, col, x, y)
				av.R += uint32(col[0])
				av.G += uint32(col[1])
				av.B += uint32(col[2])

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
	av.R, av.G, av.B, av.A = av.R/imgArea, av.G/imgArea, av.B/imgArea, opacity
	//fmt.Println(av)
	return av
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

/* corrects new image size and number of go routine used accotding to size of a picture and squares */
func caclulateNewLimits(x, y int) Frame {
	var fr Frame
	tmpRout := config.RoutineLookup()
	fr.Size = config.ChunkLookup()

	fr.X, fr.Y = CorrectImageSize(x, fr.Size), CorrectImageSize(y, fr.Size)

	tmp := fr.X / fr.Size
	fr.Routine = gcd(tmp, tmpRout)

	fr.Step = fr.X / fr.Routine
	fmt.Printf("%+v\n", fr)
	return fr
}
