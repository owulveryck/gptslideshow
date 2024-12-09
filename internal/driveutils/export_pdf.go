package driveutils

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"google.golang.org/api/drive/v3"
)

// ExtractPDF retrieves the PDF version of a Google Slides document from Google Drive.
// It uses the Drive API to export the file in "application/pdf" format.
//
// Parameters:
// - ctx: The context for managing the API call lifecycle.
// - srvDrive: An authenticated instance of the Google Drive service.
// - ID: The file ID of the Google Slides document to export.
//
// Returns:
// - A byte slice containing the exported PDF file.
// - An error if the operation fails.
func ExtractPDF(ctx context.Context, srvDrive *drive.Service, ID string) ([]byte, error) {
	// Define the MIME type for PDF.
	const pdfMimeType = "application/pdf"

	// Call the Drive API to export the file in the desired format.
	resp, err := srvDrive.Files.Export(ID, pdfMimeType).Context(ctx).Download()
	if err != nil {
		return nil, fmt.Errorf("failed to export file as PDF: %w", err)
	}
	defer resp.Body.Close()

	// Use a buffer to store the response data.
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to read PDF content: %w", err)
	}

	// Return the PDF content as a byte slice.
	return buf.Bytes(), nil
}
