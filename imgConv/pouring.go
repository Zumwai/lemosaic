package imgConv

import (
	//	"fmt"
	//"fmt"
	//"fmt"
	"image"
	"image/color"
	//"image/draw"
	"golang.org/x/image/draw"
	"sync"
)

/* pours chunk */
func pourColorImg(src Image, dst Image, dx, limitX, limitY, size int) {
	for x := dx; x < dx+limitX; x += size {
		for y := 0; y < limitY; y += size {
			col := CalcAverageChunk(x, y, size, src)
			//	fmt.Println(col)
			bounds := image.Rectangle{image.Point{x, y}, image.Point{x + size, y + size}}
			av := color.RGBA{uint8(col.R), uint8(col.G), uint8(col.B), uint8(col.A)}
			//	fmt.Println(av)
			draw.Draw(dst, bounds, &image.Uniform{av}, image.Point{0, 0}, draw.Src)
			//	draw.DrawMask(dst, bounds, &image.Uniform{av}, image.Point{0, 0}, nil, image.Point{0, 0}, draw.Src)
		}
	}
}

/* pours colors in memory */
func PreparePouring(src Image, chunk int, goCount int) Image {
	var wg sync.WaitGroup
	X, Y, goCount, goStep := caclulateNewLimits(src.Bounds().Max.X, src.Bounds().Max.Y, chunk, goCount)

	dst := GetEmptyPicture(X, Y)

	wg.Add(goCount)
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
