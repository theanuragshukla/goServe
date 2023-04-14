package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func Usage() {
	log.Printf("Usage: %s [OPTIONS] DIRECTORY\n", os.Args[0])
	log.Printf("Starts a file server to serve files from DIRECTORY.\n\n")
	log.Printf("Options:\n")
	flag.PrintDefaults()
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func humanizeSize(size int64) string {
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
		TB = 1 << 40
	)

	switch {
	case size < KB:
		return strconv.FormatInt(size, 10) + " B"
	case size < MB:
		return strconv.FormatInt(size/KB, 10) + " KB"
	case size < GB:
		return strconv.FormatInt(size/MB, 10) + " MB"
	case size < TB:
		return strconv.FormatInt(size/GB, 10) + " GB"
	default:
		return strconv.FormatInt(size/TB, 10) + " TB"
	}
}
