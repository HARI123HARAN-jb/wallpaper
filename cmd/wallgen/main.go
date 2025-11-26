package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"wallgen/internal/processor"

	"github.com/disintegration/imaging"
)

func main() {
	// 1. Parse Flags
	inputPath := flag.String("input", "", "Path to the 4K source image")
	outputDir := flag.String("out", "output", "Directory to save generated wallpapers")
	resStr := flag.String("res", "1366x768,1920x1080,2560x1440,3840x2160,1080x2400,1440x3200", "Comma-separated list of target resolutions")
	format := flag.String("format", "jpg", "Output format (jpg or png)")
	quality := flag.Int("quality", 95, "Output quality (0-100)")
	flag.Parse()

	if *inputPath == "" {
		fmt.Println("Error: -input flag is required")
		flag.Usage()
		os.Exit(1)
	}

	start := time.Now()
	fmt.Println("========================================")
	fmt.Println("   WallGen - High-Quality Wallpaper Generator")
	fmt.Println("========================================")
	fmt.Printf("Input: %s\n", *inputPath)
	fmt.Printf("Output Dir: %s\n", *outputDir)
	fmt.Printf("Format: %s, Quality: %d\n", *format, *quality)
	fmt.Println("----------------------------------------")

	// 2. Validate Input & Create Output Directory
	if _, err := os.Stat(*inputPath); os.IsNotExist(err) {
		log.Fatalf("Error: Input file does not exist: %s", *inputPath)
	}

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Error: Failed to create output directory: %v", err)
	}

	// 3. Parse Resolutions
	resolutions, err := processor.ParseResolutions(*resStr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// 4. AI Upscale Hook (Optional)
	processedInput := processor.AiUpscaleIfNeeded(*inputPath)

	// 5. Load Image
	fmt.Println("Loading source image...")
	srcImg, err := imaging.Open(processedInput)
	if err != nil {
		log.Fatalf("Error: Failed to open image: %v", err)
	}

	// 6. Process Each Resolution
	fmt.Println("Generating wallpapers...")
	for _, res := range resolutions {
		fmt.Printf("  -> Processing %dx%d... ", res.Width, res.Height)
		err := processor.ResizeImage(srcImg, res.Width, res.Height, *outputDir, *format, *quality)
		if err != nil {
			fmt.Printf("FAILED: %v\n", err)
		} else {
			fmt.Println("DONE")
		}
	}

	elapsed := time.Since(start)
	fmt.Println("----------------------------------------")
	fmt.Printf("All tasks completed in %s\n", elapsed)
	fmt.Println("========================================")
}
