package imgConv

import (
	//	"fmt"
	//"fmt"
	//"fmt"
	"image"
	"image/color"
	//"image/draw"
	"sync"

	"golang.org/x/image/draw"
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
func PreparePouring(src Image) Image {
	var wg sync.WaitGroup
	var fr Frame
	fr = caclulateNewLimits(src.Bounds().Max.X, src.Bounds().Max.Y)

	dst := GetEmptyPicture(fr.X, fr.Y)

	wg.Add(fr.Routine)
	for i := 0; i < fr.Routine; i++ {
		go func(i int) {
			defer wg.Done()
			pourColorImg(src, dst, i*fr.Step, fr.Step, fr.Y, fr.Size)
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
