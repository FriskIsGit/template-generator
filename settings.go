package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Paper struct {
	WidthMm, HeightMm int
}

var A0 = Paper{841, 1189}
var A1 = Paper{594, 841}
var A2 = Paper{420, 594}
var A3 = Paper{297, 420}
var A4 = Paper{210, 297}
var A5 = Paper{148, 210}
var A6 = Paper{105, 148}
var A7 = Paper{74, 105}
var A8 = Paper{52, 74}

var A_PAPERS = [9]Paper{A0, A1, A2, A3, A4, A5, A6, A7, A8}

type Settings struct {
	Paper                 Paper
	PixelsPerMm, MarginMm int
}

func defaultSettings() Settings {
	return Settings{A4, 10, 10}
}

func parseArgs(args []string) *Settings {
	settings := defaultSettings()
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			continue
		}
		switch arg {
		case "-h", "-help", "--help":
			displayHelp()
			os.Exit(0)
		}
		if i+1 == len(args) {
			FailExit(fmt.Sprintf("ERROR: Expected value after flag %v\n", arg))
			break
		}
		switch value := args[i+1]; arg {
		case "-paper":
			if !strings.HasPrefix(value, "A") && !strings.HasPrefix(value, "a") {
				FailExit(fmt.Sprintf("ERROR: Paper's size should start with letter 'A', given: %v\n", value))
				break
			}
			a, err := strconv.Atoi(value[1:])
			if err != nil {
				FailExit(fmt.Sprintf("ERROR: Failed to parse paper's size from: %v\n", value[1:]))
				break
			}
			if a < 0 || a >= len(A_PAPERS) {
				FailExit(fmt.Sprintf("ERROR: Paper's A size is out of range, given: %v\n", a))
				break
			}
			i++
			settings.Paper = A_PAPERS[a]
		case "-margin":
			margin, err := strconv.Atoi(value)
			if err != nil {
				FailExit(fmt.Sprintf("ERROR: Failed to parse margin, given: %v\n", value))
				break
			}
			i++
			settings.MarginMm = margin
		case "-px":
			pixels, err := strconv.Atoi(value)
			if err != nil {
				FailExit(fmt.Sprintf("ERROR: Failed to parse pixels per mm, given: %v\n", value))
				break
			}
			i++
			settings.PixelsPerMm = pixels
		}
	}
	return &settings
}

func parseDimensions(dimensions string) (width, height int) {
	x := strings.Index(dimensions, "x")
	if x == -1 {
		FailExit("ERROR: No `x` delimiter found in dimensions")
	}
	width, err := strconv.Atoi(dimensions[:x])
	if err != nil {
		FailExit("ERROR: Failed to parse `width`")
	}
	height, err = strconv.Atoi(dimensions[x+1:])
	if err != nil {
		FailExit("ERROR: Failed to parse `height`")
	}
	return width, height
}

func getExecutableName() string {
	exec, err := os.Executable()
	if err != nil {
		exec = os.Args[0]
	}
	return filepath.Base(exec)
}

const VERSION = "1.0.3"

func displayHelp() {
	exe := getExecutableName()

	fmt.Println("+--------------------------+")
	fmt.Println("|TEMPLATE GENERATOR", "v"+VERSION, "|")
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
	fmt.Println("    -h, --help           Displays this help message (any order)")
	fmt.Println()
	fmt.Println("Example usage:")
	fmt.Println("    ", exe, "create 50x20 -px 16")
	fmt.Println("    ", exe, "replicate label.jpg -margin 20")
}

func FailExit(message string) {
	fmt.Println(message)
	os.Exit(1)
}
