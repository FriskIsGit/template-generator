package main

import (
    "template-generator/utils"
	"fmt"
	"image"
    "os"
	"image/png"
	"time"
	_ "image/color"
	_ "image/jpeg"
)

func main() {
    args := os.Args
    arg_count := len(args)
    if arg_count <= 1 {
        settings.DisplayHelp();
        return
    }

    options := settings.Defaults()
    exe := settings.GetExecutableName()
    switch command := args[1]; command {
    case "create":
        if arg_count < 3 {
            fmt.Println("Example usage:", exe, "create 50x20");
            return
        }
        width, height := settings.ParseDimensions(args[2])
        template := createImage(width, height, options.Pixels_per_mm)
        saveImage(template, "template.png")
        return
    case "replicate":
        if arg_count < 3 {
            fmt.Println("Example usage:", exe, "replicate template.png");
            return
        }
        replicateTemplate(args[2], options)
    default:
        settings.DisplayHelp();
        return
    }
}

type Cursor struct {
    x, y int
}

func replicateTemplate(templatePath string, options settings.Settings) {
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

    canvas := createImage(options.Paper.Width_mm, options.Paper.Height_mm, options.Pixels_per_mm)
    canvasBounds := canvas.Bounds()

    canvas_w_px := canvasBounds.Max.X;
    canvas_h_px := canvasBounds.Max.Y;
    template_w_px := templateBounds.Max.X;
    template_h_px := templateBounds.Max.Y;
    pixels_per_mm := options.Pixels_per_mm;

    margin_px := options.Margin_mm * pixels_per_mm

    pixelsWritten := 0
    start := time.Now()
    for y := margin_px; y + template_h_px <= canvas_h_px - margin_px; y+=template_h_px {
        for x := margin_px; x + template_w_px <= canvas_w_px - margin_px; x+=template_w_px {
            canvasCursor := Cursor{ x, y }
            templateCursor := Cursor{ 0, 0, }

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

    saveImage(canvas, "generated.png");
}

func createImage(width_mm, height_mm, pixels_per_mm int) image.NRGBA {

    pixel_width := width_mm * pixels_per_mm
    pixel_height := height_mm * pixels_per_mm

    rectangle := image.Rect(0, 0, pixel_width, pixel_height)
    img := image.NewNRGBA(rectangle)
    return *img
}

func saveImage(img image.NRGBA, path string) {
    imgFile, err := os.Create(path)
    if err != nil {
        fmt.Println(err)
        return
    }

    defer imgFile.Close()

    fmt.Println("Saving image...")
    if err := png.Encode(imgFile, &img); err != nil {
        fmt.Println(err)
    }
}

