package processor

import (
	"fmt"
	"strconv"
	"strings"
)

// Resolution represents a width x height pair
type Resolution struct {
	Width  int
	Height int
}

// ParseResolutions parses a comma-separated string of resolutions (e.g., "1920x1080,1366x768").
func ParseResolutions(resStr string) ([]Resolution, error) {
	var resolutions []Resolution
	parts := strings.Split(resStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		dims := strings.Split(part, "x")
		if len(dims) != 2 {
			return nil, fmt.Errorf("invalid resolution format: %s", part)
		}

		width, err := strconv.Atoi(strings.TrimSpace(dims[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid width in resolution %s: %w", part, err)
		}

		height, err := strconv.Atoi(strings.TrimSpace(dims[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid height in resolution %s: %w", part, err)
		}

		resolutions = append(resolutions, Resolution{Width: width, Height: height})
	}

	return resolutions, nil
}

// AiUpscaleIfNeeded is a placeholder for future AI upscaling integration.
// Currently, it just returns the input path as is.
func AiUpscaleIfNeeded(inputPath string) string {
	// TODO: Implement Real-ESRGAN-ncnn-vulkan call here
	// Example:
	// cmd := exec.Command("realesrgan-ncnn-vulkan", "-i", inputPath, "-o", "upscaled.png")
	// cmd.Run()
	// return "upscaled.png"
	return inputPath
}
