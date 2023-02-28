package localMosaic

import (
	"fmt"
	"mosaic/imgConv"
	"os"
	"sync"
)

type AvColors struct {
	mu   sync.Mutex
	hash map[string]imgConv.ImgInfo
}

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
