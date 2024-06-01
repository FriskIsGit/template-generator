package settings

import (
	"fmt"
	"os"
	"strings"
    "strconv"
    "path/filepath"
)

type Settings struct {
    Paper Paper
    Pixels_per_mm, Margin_mm int
}

func Defaults() Settings {
    return Settings{ A4, 10, 10 }
}

func ApplyArgs() {
    fmt.Println(os.Args)
    // arr := os.Args
}

func ParseDimensions(dimensions string) (width, height int) {
    x := strings.Index(dimensions, "x")
    width, l_err := strconv.Atoi(dimensions[:x])
    if l_err != nil {
        panic(l_err)
    }
    height, r_err := strconv.Atoi(dimensions[x+1:])
    if r_err != nil {
        panic(r_err)
    }
    return width, height
}

func GetExecutableName() string {
    path, err := os.Executable()
    if err != nil {
        return "template"
    }
    path = filepath.Base(path)
    dot := strings.Index(path, ".")
    if dot == -1 {
        return path
    }
    return path[:dot]
}

var VERSION string = "1.0.0";

func DisplayHelp() {
    exe := GetExecutableName()

    fmt.Println("+--------------------------+")
    fmt.Println("|TEMPLATE GENERATOR", "v" + VERSION, "|")
    fmt.Println("+--------------------------+")
    fmt.Println()
    fmt.Println("Usage:")
    fmt.Println("    ", exe, "create <dimensions>        - creates a width x height image [mm]")
    fmt.Println("    ", exe, "replicate <path>           - replicates template across paper")
    fmt.Println()
    fmt.Println("Options:")
    fmt.Println("    -paper A4            Sets paper size (A0-A8) (default: A4)")
    fmt.Println("    -px 10               Pixels per millimeter (default: 10) [px]")
    fmt.Println("    -margin 10           Sets margin to each side during replication (default: 10) [mm]")
    fmt.Println()
    fmt.Println("Example usage:")
    fmt.Println("    ", exe, "create 50x20 -px 16")
    fmt.Println("    ", exe, "replicate label.jpg -margin 20")
}

type Paper struct {
    Width_mm, Height_mm int
}

var A0 = Paper{ 841, 1189 }
var A1 = Paper{ 594, 841 }
var A2 = Paper{ 420, 594 }
var A3 = Paper{ 297, 420 }
var A4 = Paper{ 210, 297 }
var A5 = Paper{ 148, 210 }
var A6 = Paper{ 105, 148 }
var A7 = Paper{ 74, 105 }
var A8 = Paper{ 52, 74 }