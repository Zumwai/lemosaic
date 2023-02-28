package imgConv

import (
	//	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"
)

/* steps over x by the amount of size*goroutine and iterates from top to bottom of y, converts average chunk size of original image to av color and  pours a square with resulted color */

func pourColorImg(src image.Image, dst *image.RGBA, dx, limitX, size int) {
	for x := dx; x < dx+limitX; x += size {
		for y := 0; y < src.Bounds().Max.Y; y += size {
			col := CalcAverageChunk(x, y, size, src)
			bounds := image.Rectangle{image.Point{x, y}, image.Point{x + size, y + size}}
			av := color.RGBA{uint8(col.R), uint8(col.G), uint8(col.B), uint8(col.A)}
			draw.Draw(dst, bounds, &image.Uniform{av}, image.Point{0, 0}, draw.Src)
		}
	}
}

/* pours colors in memory */
func PreparePouring(src image.Image, chunk int, goCount int) image.Image {
	var wg sync.WaitGroup

	X, Y := CorrectImageSize(src.Bounds().Max.X, chunk), CorrectImageSize(src.Bounds().Max.Y, chunk)
	dst := image.NewRGBA(image.Rect(0, 0, X, Y))
	goStep := X / goCount
	wg.Add(goCount)
	for i := 0; i < goCount; i++ {
		go func(i int) {
			defer wg.Done()
			pourColorImg(src, dst, i*goStep, goStep, chunk)
		}(i)
	}
	wg.Wait()
	return dst
}
