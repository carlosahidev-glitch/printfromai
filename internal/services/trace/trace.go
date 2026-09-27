package main

import (
	"bufio"
	"fmt"
	"image"
	_ "image/png" // Register PNG decoder
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path-to-png>")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPDFPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".pdf"

	// 1. Convert the PNG to a temporary PBM (Portable Bitmap) file
	pbmPath, err := convertPNGToPBM(inputPath)
	if err != nil {
		fmt.Printf("Error converting image: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(pbmPath) // Clean up the temp file when done

	// 2. Call Potrace to vectorize the PBM and output a PDF
	err = vectorizeToPDF(pbmPath, outputPDFPath)
	if err != nil {
		fmt.Printf("Error vectorizing image: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created vector PDF: %s\n", outputPDFPath)
}

// convertPNGToPBM reads a PNG and outputs a 1-bit PBM file.
func convertPNGToPBM(inputPath string) (string, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Create a temporary file for the PBM output
	tempFile, err := os.CreateTemp("", "trace-*.pbm")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	writer := bufio.NewWriter(tempFile)

	// Write PBM Header (P1 = Plain PBM format)
	fmt.Fprintf(writer, "P1\n%d %d\n", width, height)

	// Write pixel data: Potrace expects 1 for black (foreground) and 0 for white (background)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Handle transparency (treat transparent as white)
			if a == 0 {
				writer.WriteString("0 ")
				continue
			}

			// Simple luminance calculation to determine if pixel is black or white
			// (r, g, b are 16-bit, max 65535)
			luminance := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 65535.0

			// Threshold at 0.5 (closer to 0 is black, closer to 1 is white)
			if luminance < 0.5 {
				writer.WriteString("1 ") // Black pixel
			} else {
				writer.WriteString("0 ") // White pixel
			}
		}
		writer.WriteString("\n")
	}

	err = writer.Flush()
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

// vectorizeToPDF runs the potrace CLI command to generate a PDF.
func vectorizeToPDF(inputPBM, outputPDF string) error {
	// potrace <input> -b pdf -o <output>
	cmd := exec.Command("potrace", inputPBM, "-b", "pdf", "-o", outputPDF)

	// Capture any standard error output from Potrace for debugging
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
