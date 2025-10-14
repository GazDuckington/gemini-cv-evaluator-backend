package utils

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/unidoc/unipdf/v4/extractor"
	"github.com/unidoc/unipdf/v4/model"
)

// ExtractTextFromPDF loads a document using the klippa-app/go-pdfium library
// and extracts all text content by iterating through pages.
func ExtractTextFromPDF(path string) (string, error) {
	// --- 1. Open the PDF file ---
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF file: %w", err)
	}
	defer file.Close()

	// --- 2. Read the PDF into a UniPDF Reader ---
	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	// --- 3. Get number of pages ---
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", fmt.Errorf("failed to get number of pages: %w", err)
	}

	var allText strings.Builder

	// --- 4. Iterate through pages ---
	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return "", fmt.Errorf("failed to get page %d: %w", i, err)
		}

		// --- 5. Extract text from the page ---
		ex, err := extractor.New(page)
		if err != nil {
			return "", fmt.Errorf("failed to create text extractor for page %d: %w", i, err)
		}

		pageText, err := ex.ExtractText()
		if err != nil {
			return "", fmt.Errorf("failed to extract text from page %d: %w", i, err)
		}

		allText.WriteString(pageText)
		allText.WriteString("\n\n") // separate pages
	}

	return allText.String(), nil
}

// BytesToFloat32sBinary converts a []byte slice to a []float32 slice
// using the encoding/binary package.
// It assumes Little Endian byte order. Use binary.BigEndian for Big Endian.
func BytesToFloat32sBinary(b []byte) ([]float32, error) {
	if len(b)%4 != 0 {
		return nil, fmt.Errorf("byte slice length must be a multiple of 4, got %d", len(b))
	}

	floats := make([]float32, 0, len(b)/4)

	// Iterate over the byte slice 4 bytes at a time
	for i := 0; i < len(b); i += 4 {
		// Read the 4 bytes as a uint32
		bits := binary.LittleEndian.Uint32(b[i : i+4])

		// Convert the uint32 bit pattern to a float32
		f := math.Float32frombits(bits)
		floats = append(floats, f)
	}

	return floats, nil
}
