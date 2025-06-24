package main

import (
	"fmt"
	"image"
	_ "image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"strings"
	"time"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		displayHelp()
		return
	}

	exe := getExecutableName()
	switch command := args[1]; command {
	case "create":
		if len(args) < 3 {
			fmt.Println("Example usage:", exe, "create 50x20")
			return
		}
		width, height := parseDimensions(args[2])
		options := parseArgs(args[3:])
		template := createImage(width, height, options.PixelsPerMm)
		saveImage(template, "template.png")
		return
	case "replicate":
		if len(args) < 3 {
			fmt.Println("Example usage:", exe, "replicate template.png")
			return
		}
		options := parseArgs(args[3:])
		replicateTemplate(args[2], options)
	case "list", "ls":
		displayTemplates()
		return
	default:
		displayHelp()
		return
	}
}

type Cursor struct {
	x, y int
}

func replicateTemplate(templatePath string, options *Settings) {
	file, err := os.Open(templatePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	template, ext, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	templateBounds := template.Bounds()
	fmt.Println("Extension detected:", ext, "| Template size:", templateBounds.Max)

	canvas := createImage(options.Paper.WidthMm, options.Paper.HeightMm, options.PixelsPerMm)
	canvasBounds := canvas.Bounds()

	canvas_w_px := canvasBounds.Max.X
	canvas_h_px := canvasBounds.Max.Y
	template_w_px := templateBounds.Max.X
	template_h_px := templateBounds.Max.Y
	pixels_per_mm := options.PixelsPerMm

	margin_px := options.MarginMm * pixels_per_mm

	pixelsWritten := 0
	start := time.Now()
	for y := margin_px; y+template_h_px <= canvas_h_px-margin_px; y += template_h_px {
		for x := margin_px; x+template_w_px <= canvas_w_px-margin_px; x += template_w_px {
			canvasCursor := Cursor{x, y}
			templateCursor := Cursor{0, 0}

			for templateCursor.y < template_h_px {
				for templateCursor.x < template_w_px {
					color := template.At(templateCursor.x, templateCursor.y)
					canvas.Set(canvasCursor.x, canvasCursor.y, color)
					pixelsWritten++
					canvasCursor.x++
					templateCursor.x++
				}
				canvasCursor.x = x
				templateCursor.x = 0

				canvasCursor.y++
				templateCursor.y++
			}
		}
	}

	taken := time.Since(start)
	fmt.Println("Replicated", pixelsWritten, "pixels in", taken)

	saveImage(canvas, "generated.png")
}

const GB = 1024 * 1024 * 1024

func createImage(width_mm, height_mm, pixels_per_mm int) *image.NRGBA {
	pixel_width := width_mm * pixels_per_mm
	pixel_height := height_mm * pixels_per_mm
	if pixel_width*pixel_height*4 > GB {
		FailExit("Template image size exceeds 1GB. Reduce dimensions or pixels per mm")
	}
	rectangle := image.Rect(0, 0, pixel_width, pixel_height)
	return image.NewNRGBA(rectangle)
}

func saveImage(img *image.NRGBA, path string) {
	imgFile, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer imgFile.Close()

	fmt.Println("Saving image...")
	if err := png.Encode(imgFile, img); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Saved as", path)
}

func displayTemplates() {
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println(err)
		return
	}

	supportedExtensions := [3]string{"png", "jpg", "jpeg"}
	for _, file := range files {
		for _, ext := range supportedExtensions {
			fileName := file.Name()
			if strings.HasSuffix(fileName, ext) {
				fmt.Println("-", fileName)
			}
		}
	}
}
