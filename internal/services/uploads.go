package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// UploadService handles file uploads
type UploadService struct {
	uploadDir string
}

// NewUploadService creates a new upload service
func NewUploadService(uploadDir string) *UploadService {
	// Ensure upload directory exists
	receiptsDir := filepath.Join(uploadDir, "receipts")
	os.MkdirAll(receiptsDir, 0755)
	return &UploadService{uploadDir: uploadDir}
}

// AllowedReceiptTypes defines allowed MIME types for receipts
var AllowedReceiptTypes = map[string]bool{
	"image/jpeg":      true,
	"image/jpg":       true,
	"image/png":       true,
	"image/gif":       true,
	"application/pdf": true,
}

// AllowedReceiptExtensions defines allowed file extensions
var AllowedReceiptExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".pdf":  true,
}

// MaxReceiptSize is the maximum file size (10MB)
const MaxReceiptSize = 10 * 1024 * 1024

// ValidateReceipt validates a receipt file
func (s *UploadService) ValidateReceipt(filename string, contentType string, size int64) error {
	// Check file size
	if size > MaxReceiptSize {
		return fmt.Errorf("file too large: maximum size is 10MB")
	}

	// Check content type
	if !AllowedReceiptTypes[contentType] {
		return fmt.Errorf("invalid file type: only JPEG, PNG, GIF, and PDF are allowed")
	}

	// Check extension
	ext := strings.ToLower(filepath.Ext(filename))
	if !AllowedReceiptExtensions[ext] {
		return fmt.Errorf("invalid file extension: only .jpg, .jpeg, .png, .gif, and .pdf are allowed")
	}

	return nil
}

// SaveReceipt saves an uploaded receipt file
func (s *UploadService) SaveReceipt(file io.Reader, filename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = ".jpg"
	}

	// Generate unique filename
	newFilename := uuid.New().String() + ext
	relativePath := filepath.Join("receipts", newFilename)
	fullPath := filepath.Join(s.uploadDir, relativePath)

	// Create file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy content
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(fullPath) // Clean up on error
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return the URL path (not filesystem path)
	return "/uploads/" + relativePath, nil
}

// DeleteReceipt deletes a receipt file
func (s *UploadService) DeleteReceipt(urlPath string) error {
	// Convert URL path to filesystem path
	if !strings.HasPrefix(urlPath, "/uploads/") {
		return fmt.Errorf("invalid path")
	}

	relativePath := strings.TrimPrefix(urlPath, "/uploads/")
	fullPath := filepath.Join(s.uploadDir, relativePath)

	// Security check: ensure path is within upload directory
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return err
	}
	absUploadDir, err := filepath.Abs(s.uploadDir)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(absPath, absUploadDir) {
		return fmt.Errorf("invalid path")
	}

	return os.Remove(fullPath)
}

// GetUploadDir returns the upload directory path
func (s *UploadService) GetUploadDir() string {
	return s.uploadDir
}
