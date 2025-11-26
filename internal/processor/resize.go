package processor

import (
	"fmt"
	"image"
	"image/png"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

// ResizeImage takes an input image and generates a resized version at the target width and height.
// It uses the Lanczos filter for high-quality downscaling and saves the result to the output directory.
func ResizeImage(img image.Image, width, height int, outputDir, format string, quality int) error {
	// Use Lanczos filter for high-quality resizing (supersampling effect)
	// imaging.Fill preserves the aspect ratio and crops the center if necessary.
	resizedImg := imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)

	// Construct output filename
	filename := fmt.Sprintf("wallpaper_%dx%d.%s", width, height, strings.ToLower(format))
	outputPath := filepath.Join(outputDir, filename)

	// Save the image
	var err error
	switch strings.ToLower(format) {
	case "jpg", "jpeg":
		err = imaging.Save(resizedImg, outputPath, imaging.JPEGQuality(quality))
	case "png":
		err = imaging.Save(resizedImg, outputPath, pngCompression(quality))
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to save image to %s: %w", outputPath, err)
	}

	return nil
}

// pngCompression maps a quality score (0-100) to PNG compression level.
// Higher quality input means lower compression level (faster, larger file).
// Lower quality input means higher compression level (slower, smaller file).
func pngCompression(quality int) imaging.EncodeOption {
	var level png.CompressionLevel
	if quality >= 90 {
		level = png.NoCompression
	} else if quality >= 70 {
		level = png.BestSpeed
	} else if quality >= 50 {
		level = png.DefaultCompression
	} else {
		level = png.BestCompression
	}
	return imaging.PNGCompressionLevel(level)
}
