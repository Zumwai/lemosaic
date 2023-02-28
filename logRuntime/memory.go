package logRuntime

import (
	"fmt"
	"runtime"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// https://golang.org/pkg/runtime/#MemStats
func PrintMemory(message string) {
	var m runtime.MemStats
	fmt.Printf("breakpoint %s\t", message)
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v MiB\n", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB\n", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB\n", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
	fmt.Printf("\tMallocs = %v\n", bToMb(m.Mallocs))
	fmt.Printf("\tHeapSys = %v\n", bToMb(m.HeapSys))
}
