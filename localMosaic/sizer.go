package localMosaic

import (
	"fmt"
	//"image"
	"mosaic/imgConv"
	"os"

	"sync"
)

/* need to do it in goroutine plus additional feature of resizing in place */
func resizeSrcImage(dirName, name string, size int) error {
	src, err := getDecodedFile(dirName + "/" + name)
	if err != nil {
		return err
	}
	dst, err := imgConv.ResizeInMemory(src, size, size)
	if err != nil {
		return err
	}

	err = encodeToFile("./smaller/", name, "_resized.png", dst)
	if err != nil {
		return err
	}
	return nil
}

/* changes size of the images in dir and squares them  to desirable size.. Puts in ./smaller dir. */
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
