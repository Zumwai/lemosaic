package main

import (
	"flag"
	"fmt"
	"mosaic/config"
	"mosaic/localMosaic"
	"mosaic/logRuntime"
	//"github.com/pkg/profile"
	//_ "net/http/pprof"
	//"github.com/google/pprof/profile"
	"mosaic/serve"
)

func main() {
	flag.Int("interpol", 0, "Choose one out of the for methods for interpolation")
	flag.Int("format", 0, "choose internal format [0-1] - RGBA or NRGBA")
	flag.Int("encoder", 0, "choose final image format [0-3] - jpeg, png, tiff, gif")
	flag.Int("chunk", 20, "size of a square in image")
	flag.Int("routine", 50, "number of gourutine in use, will  be modified AKA tolower depending on the image size")
	flag.Bool("normal", false, "pass this flag to normalize square(chunk) size instead of image size")        //TODO
	flag.Bool("bounds", false, "flag for ignoring potenintial cutted squares at the boundaries of the image") //TODO
	/*sets additional tasks for changing size and format of the images*/
	var pourTarget = flag.String("pour", "", "Pour this file into squares")
	var calcAverage = flag.String("average", "", "what are average colors of a picture")
	var mosaic = flag.String("mosaic", "", "mosaic dat image")
	var browser = flag.Bool("serve", false, "enable -serve flag in order ot start a server")
	flag.Parse()

	//p := profile.Start(profile.CPUProfile)
	//defer profile.Start(profile.CPUProfile).Stop()
	//configTask(*chunkSize, *methodFlag)

	config.RegDecoders()
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

	logRuntime.PrintMemory("at the end\n")
	//
	// p.Stop()
}
