package imgConv

import (
	"image"
	"image/draw"
	"math"
	"sync"
)

/* calculates euclid distance between two 3-tuple. Ignores A for now */
func euclidCoordinates(target Pixel, src Pixel) float64 {
	return math.Sqrt(math.Pow(float64(target.R-src.R), 2) +
		math.Pow(float64(target.G-src.G), 2) +
		math.Pow(float64(target.B-src.B), 2))
}

/* iterates over map of available squares and returns nearest image  */
func CalculateNearestPic(col Pixel, source map[string]ImgInfo) image.Image {
	var min float64 = 9000
	var new image.Image
	for _, f := range source {
		tmp := euclidCoordinates(col, f.Av)
		if min > tmp {
			min, new = tmp, f.Square
		}
	}
	return new
}

/* steps over x by the amount of size*goroutine and iterates from top to bottom of y, converts average chunk size of original image to av color and  substitutes it with nearest available image-square */
func mosaicDatImg(src image.Image, dst *image.RGBA, dx, limitX, size int, source map[string]ImgInfo) {
	for x := dx; x < dx+limitX; x += size {
		for y := 0; y < src.Bounds().Max.Y; y += size {
			col := CalcAverageChunk(x, y, size, src)
			bounds := image.Rectangle{image.Point{x, y}, image.Point{x + size, y + size}}
			av := CalculateNearestPic(col, source)
			if av == nil {
				/*costyl. breaks from finishing the image */
				return
			}
			draw.Draw(dst, bounds, av, image.Point{0, 0}, draw.Src)
		}
	}
}

/* prepares mosaic image in memory */
func PrepareMosaic(src image.Image, chunk, goCount int, source map[string]ImgInfo) image.Image {
	var wg sync.WaitGroup

	sizeX, sizeY := CorrectImageSize(src.Bounds().Max.X, chunk), CorrectImageSize(src.Bounds().Max.Y, chunk)
	dst := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))
	goStep := sizeX / goCount
	wg.Add(goCount)
	for i := 0; i < goCount; i++ {
		go func(i int) {
			defer wg.Done()
			mosaicDatImg(src, dst, i*goStep, goStep, chunk, source)
		}(i)
	}
	wg.Wait()
	return dst
}
