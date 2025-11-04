package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	utils_v1 "github.com/FDSAP-Git-Org/hephaestus/utils/v1"
	"github.com/gofiber/fiber/v3"
)

// FileUploadConfig holds configuration for file uploads
type FileUploadConfig struct {
	MaxSize      int64
	AllowedTypes []string
	UploadPath   string
}

// DefaultFileUploadConfig returns default configuration
func DefaultFileUploadConfig() FileUploadConfig {
	return FileUploadConfig{
		MaxSize:      5 * 1024 * 1024, // 5MB
		AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
		UploadPath:   "./assets/images/uploads/expenses",
	}
}

// UploadFile handles file upload and returns the file path/URL
func UploadFile(c fiber.Ctx, fileHeader *multipart.FileHeader, config FileUploadConfig) (string, error) {
	// Validate file size
	if fileHeader.Size > config.MaxSize {
		return "", fmt.Errorf("file size exceeds limit of %d bytes", config.MaxSize)
	}

	// Validate file type
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	contentType := http.DetectContentType(buffer)
	allowed := false
	for _, allowedType := range config.AllowedTypes {
		if contentType == allowedType {
			allowed = true
			break
		}
	}

	if !allowed {
		return "", fmt.Errorf("file type %s is not allowed", contentType)
	}

	// Reset file pointer
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(config.UploadPath, 0755); err != nil {
		return "", err
	}

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	randomString := utils_v1.GenerateRandomStrings(8, []string{utils_v1.UpperString, utils_v1.LowerString, utils_v1.NumericString})
	filename := fmt.Sprintf("%d%s%s", time.Now().UnixNano(), randomString, ext)
	filePath := filepath.Join(config.UploadPath, filename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", err
	}

	// Return relative path or full URL
	return fmt.Sprintf("/assets/images/uploads/expenses/%s", filename), nil
}

// DeleteUploadedFile deletes an uploaded file from the filesystem
func DeleteUploadedFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	// Extract filename from URL path
	filename := filepath.Base(filePath)
	fullPath := filepath.Join("./assets/images/uploads/expenses", filename)

	// Check if file exists before attempting to delete
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}

	err := os.Remove(fullPath)
	if err != nil {
		return fmt.Errorf("failed to delete file %s: %w", filename, err)
	}

	return nil
}

// ExtractFilenameFromURL extracts the filename from a URL
func ExtractFilenameFromURL(url string) string {
	if url == "" {
		return ""
	}
	return filepath.Base(url)
}
