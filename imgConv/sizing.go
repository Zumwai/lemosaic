package imgConv

import (
	"image"
	"mosaic/config"
)

/*stores de-facto config for processing image */
type Frame struct {
	X       int
	Y       int
	Routine int
	Step    int
	Size    int
}

/* resizes image in memory and returns it in drawable image format*/
func ResizeInMemory(src image.Image, sizeX, sizeY int) (Image, error) {
	//new := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))
	new := GetEmptyPicture(sizeX, sizeY)
	err := ApplyInterpol(src, new, image.Rectangle{image.Point{0, 0}, image.Point{sizeX, sizeY}})
	if err != nil {
		return nil, err
	}
	return new, nil
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

/* slightly corrects colors */
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{r / 257, g / 257, b / 257, a / 257}
}

/* calculates average color of a given chunk. */
func GetAveragePixel(pic Image, dx, dy, maxx, maxy int) (av Pixel) {
	for x := dx; x < maxx; x++ {
		for y := dy; y < maxy; y++ {

			col := rgbaToPixel(pic.At(x, y).RGBA())
			av.R += col.R
			av.G += col.G
			av.B += col.B
			/*
				col = rgbaToPixel(pic.At(x, y+1).RGBA())
				av.R += col.R
				av.G += col.G
				av.B += col.B
				col = rgbaToPixel(pic.At(x, y+2).RGBA())
				av.R += col.R
				av.G += col.G
				av.B += col.B
				col = rgbaToPixel(pic.At(x, y+3).RGBA())
				av.R += col.R
				av.G += col.G
				av.B += col.B
			*/
			//av.A += col.A
		}
	}
	imgArea := uint32((maxx - dx) * (maxy - dy))
	av.R, av.G, av.B, av.A = av.R/imgArea, av.G/imgArea, av.B/imgArea, 255
	return av
}

func hcf(n1 int, n2 int) int {
	if n2 != 0 {
		return hcf(n2, n1%n2)
	} else {
		return n1
	}
}

/* calculates nearest integer(new <= num) divisible by div */
func CorrectImageSize(num, div int) int {
	var quotient = int(num / div)
	return int(div * quotient)
}

/*
corrects new image size/ number of go routine used and size of a step for them
according to size of a picture and squares
*/
func caclulateNewLimits(x, y int) Frame {
	var fr Frame
	tmpRout := config.RoutineLookup()
	normal := config.NormalizeLookup()
	fr.Size = config.ChunkLookup()

	if !normal {
		fr.X, fr.Y = CorrectImageSize(x, fr.Size), CorrectImageSize(y, fr.Size)
	} else {
		fr.X, fr.Y = x, y
	}

	fr.Routine = hcf(fr.X/fr.Size, tmpRout)
	fr.Step = fr.X / fr.Routine
	return fr
}
