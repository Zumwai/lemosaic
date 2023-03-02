package localMosaic

import (
	"fmt"
	//"image"
	"mosaic/config"
	"mosaic/imgConv"
	"os"
	"sync"
)

type AvColors struct {
	mu   sync.Mutex
	hash map[string]imgConv.ImgInfo
}

/* method for concurrently add entity to map */
func (m *AvColors) add(name string, hash imgConv.ImgInfo) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hash[name] = hash
}

/* opens dir and puts all its files into map, containing name and average colors */
func PopulateHashDir(dirName string, size int) (map[string]imgConv.ImgInfo, error) {
	var wg sync.WaitGroup
	var average AvColors

	dirReader, err := os.ReadDir(dirName)
	if err != nil {
		return average.hash, nil
	}

	average.hash = make(map[string]imgConv.ImgInfo, len(dirReader))
	//average.hash = make(map[string]imgConv.ImgInfo)
	wg.Add(len(dirReader))
	for _, f := range dirReader {
		go func(name string) {
			defer wg.Done()
			tmp, err := CalcAverageColours(dirName+"/"+name, size)
			if err != nil {
				fmt.Println(name, ":\t", err)
			} else {
				average.add(name, tmp)
			}
		}(f.Name())
	}
	wg.Wait()
	return average.hash, nil
}

/* calcultes average colors of given file, resized it in memory if requested*/
func CalcAverageColours(name string, size int) (pic imgConv.ImgInfo, err error) {
	img, err := getDecodedFile(name)
	if err != nil {
		return
	}
	if size == 0 {
		pic.Square = img
	} else {
		pic.Square, err = imgConv.ResizeInMemory(img, size, size, config.InterpolLookup())
		if err != nil {
			return pic, err
		}
	}
	//tmpPtr := pic.Square.(*image.NRGBA)
	//pic.Av = imgConv.GetAveragePixel(tmpPtr, tmpPtr.Rect)
	pic.Av = imgConv.GetAveragePixel(pic.Square, 0, 0, pic.Square.Rect.Max.X, pic.Square.Rect.Max.Y)

	return
}
