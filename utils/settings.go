package settings

import (
    "fmt"
    "os"
    "strings"
    "strconv"
    "path/filepath"
)

type Paper struct {
    WidthMm, HeightMm int
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

var A_PAPERS = [9]Paper{A0, A1, A2, A3, A4, A5, A6, A7, A8}

type Settings struct {
    Paper Paper
    PixelsPerMm, MarginMm int
}

func Defaults() Settings {
    return Settings{ A4, 10, 10 }
}

func FromArgs() Settings {
    args := os.Args
    arg_count := len(args)
    settings := Defaults()
    for i := 2; i < arg_count ; i++ {
        if !strings.HasPrefix(args[i], "-") && i + 1 >= arg_count {
            continue
        }
        flag := args[i][1:]
        switch value := args[i+1]; flag {
        case "paper":
            i++
            if !strings.HasPrefix(args[i], "A") && i + 1 >= arg_count {
                errPrintln("WARNING: Invalid paper kind", nil)
                continue
            }
            a, err := strconv.Atoi(value[1:])
            if err != nil {
                errPrintln("WARNING: Failed to parse paper size", err)
                continue
            }
            settings.Paper = A_PAPERS[a]
        case "margin":
            i++
            margin, err := strconv.Atoi(value)
            if err != nil {
                errPrintln("WARNING: Failed to parse margin", err)
                continue
            }
            settings.MarginMm = margin
        case "px":
            i++
            pixels, err := strconv.Atoi(value)
            if err != nil {
                errPrintln("WARNING: Failed to parse pixels per mm", err)
                continue
            }
            settings.PixelsPerMm = pixels
        }
    }
    return settings
}

func ParseDimensions(dimensions string) (width, height int) {
    x := strings.Index(dimensions, "x")
    if x == -1 {
        errPrintln("WARNING: No `x` delimiter found", nil)
        os.Exit(1)
    }
    width, l_err := strconv.Atoi(dimensions[:x])
    if l_err != nil {
        errPrintln("WARNING: Failed to parse `width`", l_err)
        os.Exit(1)
    }
    height, r_err := strconv.Atoi(dimensions[x+1:])
    if r_err != nil {
        errPrintln("WARNING: Failed to parse `height`", r_err)
        os.Exit(1)
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

var VERSION string = "1.0.1";

func DisplayHelp() {
    exe := GetExecutableName()

    fmt.Println("+--------------------------+")
    fmt.Println("|TEMPLATE GENERATOR", "v" + VERSION, "|")
    fmt.Println("+--------------------------+")
    fmt.Println()
    fmt.Println("Usage:")
    fmt.Println("    ", exe, "create <dimensions>        - creates a width x height image [mm]")
    fmt.Println("    ", exe, "replicate <path>           - replicates template across paper")
    fmt.Println("    ", exe, "ls                         - lists images in current folder")
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

func errPrintln(format string, err error) {
    if err == nil {
        fmt.Fprintf(os.Stderr, format + "\n")
        return
    }
    fmt.Fprintf(os.Stderr, format + " [%s]\n", err.Error())
}