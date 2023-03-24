package imgConv

import (
	"image"
	"sync"

	"golang.org/x/image/draw"
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
func calcEuclidCoordinates(target, src Pixel) uint32 {
	return SqrtHDU32((target.R-src.R)*(target.R-src.R) +
		(target.G-src.G)*(target.G-src.G) +
		(target.B-src.B)*(target.B-src.B))
}

func semiEuclidCoordinates(target, src Pixel) uint32 {
	return (target.R-src.R)*(target.R-src.R) +
		(target.G-src.G)*(target.G-src.G) +
		(target.B-src.B)*(target.B-src.B)
}

func expandedSemiEuclid(target, src Pixel) uint32 {
	//x1 := (target.R*target.R + target.G*target.G + target.B*target.B + src.G*src.G + src.B*src.B + src.R*src.R - 2*(target.R*src.R + target.B*src.B + target.G*src.G)
	return (target.R+target.G+target.B)*(target.R+target.G+target.B) +
		(src.R+src.G+src.B)*(src.R+src.G+src.B) -
		2*(target.R+target.G+target.B)*(src.R+src.G+src.B)
}

/*
	func nonEuclidCoordinates(target, src Pixel) uint32 {
		return (target.R - src.R) + (target.G - src.G) + (target.B - src.B)
	}
*/

func calculateNearestPic(col Pixel, source []ImgInfo) Image {
	var min uint32 = 90000
	var ret Image

	for i, _ := range source {
		tmp := semiEuclidCoordinates(col, source[i].Av)
		if min > tmp {
			min, ret = tmp, source[i].Square
		}
	}

	return ret
}

/*
steps over x by the amount of (size of a square)*(number of goroutine) and  iterates from top to bottom of y,
converts average chunk size of original image to av color and substitutes it with nearest available image-square
*/
func mosaicDatImg(src Image, dst Image, dx, limitX, limitY, size int, source []ImgInfo) {
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
func PrepareMosaic(src Image, source []ImgInfo) (ret Image) {
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
