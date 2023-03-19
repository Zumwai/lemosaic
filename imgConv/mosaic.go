package imgConv

import (
	"golang.org/x/image/draw"
	"image"
	"sync"
)

func SqrtHDU32(x uint32) uint32 {
	var t, b, r uint
	t = uint(x)
	p := uint(1 << 30)
	for p > t {
		p >>= 2
	}
	for ; p != 0; p >>= 2 {
		b = r | p
		r >>= 1
		if t >= b {
			t -= b
			r |= p
		}
	}
	return uint32(r)
}

/* calculates euclid distance bewtween two 3-tuple of uint32 type */
func leEuclidCoordinates(target, src Pixel) uint32 {
	return SqrtHDU32((target.R-src.R)*(target.R-src.R) +
		(target.G-src.G)*(target.G-src.G) +
		(target.B-src.B)*(target.B-src.B))
}

/*
	calculates euclid distance between two 3-tuple. Ignores A for now. Slightly slower the le version

	func euclidCoordinates(target Pixel, src Pixel) float64 {
		return math.Sqrt(math.Pow(float64(target.R-src.R), 2) + math.Pow(float64(target.G-src.G), 2) + math.Pow(float64(target.B-src.B), 2))
	}

/* iterates over map of available squares and returns nearest image
*/
func calculateNearestPic(col Pixel, source map[string]ImgInfo) Image {
	var min uint32 = 90000
	var new Image

	for _, f := range source {
		tmp := leEuclidCoordinates(col, f.Av)
		if min > tmp {
			min, new = tmp, f.Square
		}
	}
	return new
}

/*
steps over x by the amount of (size of a square)*(number of goroutine) and  iterates from top to bottom of y,
converts average chunk size of original image to av color and substitutes it with nearest available image-square
*/
func mosaicDatImg(src Image, dst Image, dx, limitX, limitY, size int, source map[string]ImgInfo) {
	for x := dx; x < dx+limitX; x += size {
		for y := 0; y < limitY; y += size {
			col := CalcAverageChunk(x, y, size, src)
			bounds := image.Rectangle{image.Point{x, y}, image.Point{x + size, y + size}}
			av := calculateNearestPic(col, source)
			draw.Draw(dst, bounds, av, image.Point{0, 0}, draw.Over)
		}
	}
}

/* prepares mosaic image in memory */
func PrepareMosaic(src Image, source map[string]ImgInfo) (ret Image) {
	var wg sync.WaitGroup
	fr := caclulateNewLimits(src.Bounds().Max.X, src.Bounds().Max.Y)
	wg.Add(fr.Routine)

	dst := GetEmptyPicture(fr.X, fr.Y)

	for i := 0; i < fr.Routine; i++ {
		go func(i int) {
			defer wg.Done()
			mosaicDatImg(src, dst, i*fr.Step, fr.Step, fr.Y, fr.Size, source)
		}(i)
	}
	wg.Wait()
	return dst
}
