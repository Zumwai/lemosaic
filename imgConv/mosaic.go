package imgConv

import (
	//"fmt"
	"image"
	//"image/draw"
	"golang.org/x/image/draw"
	//"math"
	"sync"
)

func wikiSqrt(n uint32) uint32 {
	var c uint32 = 0
	var d uint32 = 1 << 30
	var x uint32 = n
	for d > n {
		d >>= 2
	}
	for d != 0 {
		if x >= c+d {
			x -= c + d
			c = (c >> 1) + d
		} else {
			c >>= 1
		}
		d >>= 2
	}
	return c
}

func leEuclidCoordinates(target, src Pixel) uint32 {
	return wikiSqrt((target.R-src.R)*(target.R-src.R) +
		(target.G-src.G)*(target.G-src.G) +
		(target.B-src.B)*(target.B-src.B))
}

/* calculates euclid distance between two 3-tuple. Ignores A for now. Slightly slower the le version
func euclidCoordinates(target Pixel, src Pixel) float64 {
	return math.Sqrt(math.Pow(float64(target.R-src.R), 2) + math.Pow(float64(target.G-src.G), 2) + math.Pow(float64(target.B-src.B), 2))
}
*/
/* iterates over map of available squares and returns nearest image  */
func calculateNearestPic(col Pixel, source map[string]ImgInfo) Image {
	var min uint32 = 90000
	var new Image

	for _, f := range source {
		tmp := leEuclidCoordinates(col, f.Av)
		//	fmt.Println(tmp)
		if min > tmp {
			min, new = tmp, f.Square
		}
	}
	//	fmt.Println(namepls)
	return new
}

/* steps over x by the amount of size*goroutine and iterates from top to bottom of y, converts average chunk size of original image to av color and  substitutes it with nearest available image-square */
func mosaicDatImg(src Image, dst Image, dx, limitX, limitY, size int, source map[string]ImgInfo) {
	for x := dx; x < dx+limitX; x += size {
		for y := 0; y < limitY; y += size {
			col := CalcAverageChunk(x, y, size, src)
			bounds := image.Rectangle{image.Point{x, y}, image.Point{x + size, y + size}}
			av := calculateNearestPic(col, source)
			draw.Draw(dst, bounds, av, image.Point{0, 0}, draw.Over)
			//draw.DrawMask(dst, bounds, av, image.Point{0, 0}, nil, image.Point{0, 0}, draw.Over)
		}
	}
}

/* prepares mosaic image in memory */
func PrepareMosaic(src Image, chunk, goCount int, source map[string]ImgInfo) Image {
	var wg sync.WaitGroup
	sizeX, sizeY, goCount, goStep := caclulateNewLimits(src.Bounds().Max.X, src.Bounds().Max.Y, chunk, goCount)

	dst := GetEmptyPicture(sizeX, sizeY)

	wg.Add(goCount)

	/* possible error */
	limitY := src.Bounds().Max.Y

	for i := 0; i < goCount; i++ {
		go func(i int) {
			defer wg.Done()
			mosaicDatImg(src, dst, i*goStep, goStep, limitY, chunk, source)
		}(i)
	}
	wg.Wait()
	return dst
}
