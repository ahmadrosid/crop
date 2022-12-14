package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/spf13/cobra"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

var size []int
var padding []int
var outputFolder string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "crop",
		Short: "crop - a simple CLI to crop images",
		Long:  "crop - a simple CLI to crop images",
		Run:   Execute,
	}

	rootCmd.PersistentFlags().IntSliceVar(&size, "size", []int{0, 0}, "Set image size: width x height")
	rootCmd.PersistentFlags().IntSliceVar(&padding, "padding", []int{0, 0}, "Set padding left and top")
	rootCmd.PersistentFlags().StringVar(&outputFolder, "out-folder", "output", "Output folder name.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "There was an error while executing CLI '%s'", err)
		os.Exit(1)
	}
}

func Execute(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		return
	}
	cropImage(args[0])
}

func cropImage(path string) {
	originalImageFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer originalImageFile.Close()

	originalImage, err := png.Decode(originalImageFile)
	if err != nil {
		panic(err)
	}

	bounds := originalImage.Bounds()
	width := bounds.Dx()
	height := width
	if len(size) > 0 && size[0] > 0 {
		width = size[0]
	}
	if len(size) > 1 && size[1] > 0 {
		height = size[1]
	}

	paddingLeft := 0
	paddingTop := width / 3
	if len(padding) > 0 && padding[0] > 0 {
		paddingLeft = padding[0]
	}
	if len(padding) > 1 && padding[1] > 0 {
		paddingTop = padding[1]
	}

	if paddingLeft > bounds.Dx() {
		fmt.Fprintf(os.Stderr, "Error: Left padding: %d, maximum padding: %d\n", paddingTop, bounds.Dx())
		os.Exit(1)
	}

	if paddingTop > bounds.Dy() {
		fmt.Fprintf(os.Stderr, "Error: Top padding: %d, maximum padding: %d\n", paddingTop, bounds.Dy())
		os.Exit(1)
	}

	cropSize := image.Rect(0, 0, width, height)
	paddingSize := image.Point{paddingLeft, paddingTop}
	cropSize = cropSize.Add(paddingSize)
	croppedImage := originalImage.(SubImager).SubImage(cropSize)

	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		os.Mkdir(outputFolder, 0o755)
	}

	croppedImageFile, err := os.Create(fmt.Sprintf("%s/%s", outputFolder, path))
	if err != nil {
		panic(err)
	}

	defer croppedImageFile.Close()
	if err := png.Encode(croppedImageFile, croppedImage); err != nil {
		panic(err)
	}
}
