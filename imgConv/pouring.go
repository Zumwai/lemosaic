package imgConv

import (
	//	"fmt"
	//"fmt"
	"fmt"
	"image"
	"image/color"
	//"image/draw"
	"golang.org/x/image/draw"
	"sync"
)

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func caclulateNewLimits(x, y, chunk, routine int) (nx, ny, ncount int) {
	nx, ny = CorrectImageSize(x, chunk), CorrectImageSize(y, chunk)
	tmp := nx / chunk
	ncount = GCD(tmp, ncount)
	return nx, ny, ncount
}

func pourColorImg(src *image.NRGBA, dst *image.NRGBA, dx, limitX, limitY, size int) {
	for x := dx; x < dx+limitX; x += size {
		for y := 0; y < limitY; y += size {
			col := CalcAverageChunk(x, y, size, src)
			bounds := image.Rectangle{image.Point{x, y}, image.Point{x + size, y + size}}
			av := color.RGBA{uint8(col.R), uint8(col.G), uint8(col.B), uint8(col.A)}
			draw.Draw(dst, bounds, &image.Uniform{av}, image.Point{x, y}, draw.Src)
			//			draw.DrawMask(dst, bounds, &image.Uniform{av}, image.Point{0, 0}, nil, image.Point{0, 0}, draw.Src)

		}
	}
}

/* pours colors in memory */
func PreparePouring(src *image.NRGBA, chunk int, goCount int) *image.NRGBA {
	var wg sync.WaitGroup
	X, Y, goCount := caclulateNewLimits(src.Bounds().Max.X, src.Bounds().Max.Y, chunk, goCount)

	fmt.Println(X)
	goStep := X / goCount
	fmt.Println(goStep, goCount)
	//X = CorrectImageSize(X, goCount)
	//fmt.Println(X)
	//goStep = CorrectImageSize(goStep, goCount)
	//X = CorrectImageSize(X, goStep)
	dst := image.NewNRGBA(image.Rect(0, 0, X, Y))
	//tmpPtr := src.(*image.NRGBA)

	//1365
	//1353
	wg.Add(goCount)
	//fmt.Println(goCount, X, Y, goStep)
	for i := 0; i < goCount; i++ {
		go func(i int) {
			defer wg.Done()
			//fmt.Println("dx", i*goStep, "limitX", goStep, "true limit", i*goStep+goStep, "size", chunk, "i am", i, "routine!")
			pourColorImg(src, dst, i*goStep, goStep, Y, chunk)
		}(i)
	}
	wg.Wait()
	return dst
}

/*
fmt.Println(goStep, i*goStep, Y, goCount, "i am", i, " routine!")
0 routine - has dx as 0, limitx with goStep 945
1 routine - has dx as  945, dx as 945,
*/
