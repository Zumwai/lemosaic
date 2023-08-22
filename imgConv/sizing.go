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

func hcf(n1 int, n2 int) int {
	if n2 != 0 {
		return hcf(n2, n1%n2)
	} else {
		return n1
	}
}

/* calculates nearest integer(new <= num) divisible by div */
func CorrectImageSize(num, div int) int {
	if div > num {
		return num
	}
	var quotient = int(num / div)
	return int(div * quotient)
}

/*
corrects new image size/ number of go routine used and size of a step for them
according to size of a picture and squares
*/
func caclulateNewLimits(x, y int) Frame {
	var fr Frame
	//tmpRout := config.RoutineLookup()
	normal := config.NormalizeLookup()
	fr.Size = config.ChunkLookup()
	if !normal {
		fr.X, fr.Y = CorrectImageSize(x, fr.Size), CorrectImageSize(y, fr.Size)
	} else {
		fr.X, fr.Y = x, y

	}
	fr.Step = fr.Size
	fr.Routine = fr.X / fr.Step
	return fr
}
