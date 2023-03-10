package localMosaic

import (
	"fmt"
	//"image"
	"golang.org/x/image/draw"

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
func PopulateHashDir(dirName string) (map[string]imgConv.ImgInfo, error) {
	var wg sync.WaitGroup
	var average AvColors

	dirReader, err := os.ReadDir(dirName)
	if err != nil {
		return average.hash, nil
	}

	average.hash = make(map[string]imgConv.ImgInfo, len(dirReader))
	wg.Add(len(dirReader))
	for _, f := range dirReader {
		go func(name string) {
			defer wg.Done()
			tmp, err := CalcAverageSrcColours(dirName + "/" + name)
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

/* resizes image in memory if needed, then calcultes average colors */
func CalcAverageSrcColours(name string) (pic imgConv.ImgInfo, err error) {
	img, err := getUnformattedImage(name)
	if err != nil {
		return
	}
	size := config.ChunkLookup()
	checked, ok := img.(draw.Image)
	if ok && (img.Bounds().Max.X == size && img.Bounds().Max.Y == size) {
		pic.Square = checked
		pic.Av = imgConv.GetAveragePixel(pic.Square, 0, 0, pic.Square.Bounds().Max.X, pic.Square.Bounds().Max.Y)
		return
	}

	pic.Square, err = imgConv.ResizeInMemory(img, size, size)
	if err != nil {
		return pic, err
	}
	pic.Av = imgConv.GetAveragePixel(pic.Square, 0, 0, pic.Square.Bounds().Max.X, pic.Square.Bounds().Max.Y)
	return
}
