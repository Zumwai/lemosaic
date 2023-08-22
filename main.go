package main

import (
	"flag"
	"fmt"
	"github.com/pkg/profile"
	"mosaic/config"
	"mosaic/localMosaic"
	"mosaic/logRuntime"
	"mosaic/serve"
)

func main() {
	flag.Int("interpol", 0, "Choose one out of the for methods for interpolation")
	flag.Int("format", 0, "choose internal format [0-5] - RGBA, RGBA64, NRGBA, NRGBA64, GRAY, CMYK")
	flag.Int("encoder", 0, "choose final image format [0-3] - jpeg, png, tiff, gif")
	flag.Int("chunk", 20, "size of a square in image")
	flag.Int("routine", 1000, "number of gourutine in use, will  be modified AKA tolower depending on the image size")
	flag.Bool("normal", false, "pass this to round down img size so that x/chunk % 0 and y/chynk % 0. pass -normal=true to avoid")
	flag.Int("qual", 50, "quality for jpeg output [1-100], will be adjusted to correct value")
	flag.String("source", "./pics", "dir to use for source image for mosaic")
	flag.Bool("unmax", false, "declares that arbitrary limit of 10mb is no longer needed")
	flag.Int("size", 0, "decided on final x size of the image")
	debug := flag.Bool("debug", false, "use this to enable memory tracker and generating pprof")
	var pourTarget = flag.String("pour", "", "Pour this file into squares")
	var calcAverage = flag.String("average", "", "what are average colors of a picture")
	var mosaic = flag.String("mosaic", "", "mosaic dat image")
	var browser = flag.Bool("serve", false, "enable -serve flag in order ot start a server")
	flag.Parse()

	/*actually useles*/
	config.RegDecoders()
	if *debug {
		defer profile.Start(profile.CPUProfile).Stop()
	}
	if *pourTarget != "" {
		err := localMosaic.ExecutePouring(*pourTarget)
		if err != nil {
			fmt.Println(err)
		}
	}
	if *calcAverage != "" {
		av, err := localMosaic.CalcAverageSrcColours(*calcAverage)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(av)
		}
	}
	if *mosaic != "" {
		err := localMosaic.ExecuteMosaic(*mosaic)
		if err != nil {
			fmt.Println(err)
		}
	}
	if *browser {
		serve.StartServer()
	}
	if *debug {
		logRuntime.PrintMemory("at the end\n")
	}
	//
	// p.Stop()
}
