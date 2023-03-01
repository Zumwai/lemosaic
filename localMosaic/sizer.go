package localMosaic

import (
	"fmt"
	"image"
	"image/png"
	"mosaic/config"
	"mosaic/imgConv"
	"os"
	"strings"
	"sync"
)

/* need to do it in goroutine plus additional feature of resizing in place */

func resizeSrcImage(dirName, name string, size int) error {
	src, err := getDecodedFile(dirName + "/" + name)
	if err != nil {
		return err
	}
	resized, err := os.Create("./smaller/" + strings.TrimRight(name, ".png") + "_resized.png")
	if err != nil {
		return err
	}
	defer resized.Close()

	dst := image.NewNRGBA(image.Rect(0, 0, size, size))
	err = imgConv.ApplyInterpol(src, dst, config.InterpolLookup())
	if err != nil {
		return err
	}

	err = png.Encode(resized, dst)
	if err != nil {
		return err
	}
	return nil
}

/* changes size of the images in dir to desirable size. Format  x/y get twisted for now */
func ChangeSrcsSize(targetDir string, size int) error {
	var wg sync.WaitGroup

	dirSrc, err := os.ReadDir(targetDir)
	if err != nil {
		return err
	}

	wg.Add(len(dirSrc))
	for _, f := range dirSrc {
		go func(name string) {
			defer wg.Done()
			err = resizeSrcImage(targetDir, name, size)
			if err != nil {
				fmt.Printf("failed because %s: %v\n", name, err)
			}
		}(f.Name())
	}
	wg.Wait()
	return nil
}
