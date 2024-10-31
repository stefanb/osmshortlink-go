package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/stefanb/osmshortlink-go"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: osmshortlink [latitude] [longitude] [zoom]")
		os.Exit(1)
	}

	lat, err := strconv.ParseFloat(os.Args[1], 32)
	if err != nil {
		fmt.Printf("Invalid latitude %q: %s\n", os.Args[1], err.Error())
		os.Exit(1)
	}

	lon, err := strconv.ParseFloat(os.Args[2], 32)
	if err != nil {
		fmt.Printf("Invalid longitude %q: %s\n", os.Args[2], err.Error())
		os.Exit(1)
	}

	zoom, err := strconv.ParseUint(os.Args[3], 10, 5)
	if err != nil {
		fmt.Printf("Invalid zoom %q: %s\n", os.Args[3], err.Error())
		os.Exit(1)
	}

	shortlink, err := osmshortlink.Create(float32(lat), float32(lon), int(zoom))
	if err != nil {
		fmt.Println("Error generating short link:", err.Error())
		os.Exit(1)
	}
	fmt.Println(shortlink)
}
