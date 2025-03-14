package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// SaveFile saves the uploaded file to a specific location and returns the file path
func SaveFile(src io.Reader, filename string, filetype string) (string, error) {
	// Define the directory where the file will be saved
	var uploadDir string

	if filetype == "profilepicture" {
		uploadDir = "./uploads/profilepictures"
	} else {
		uploadDir = "./uploads/attachments"
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	// Define the full file path
	filePath := filepath.Join(uploadDir, filename)

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy the file content from the source to the destination
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return filePath, nil
}

// GetFileType takes a filename and returns its type (video, picture, document, or unknown)
func GetFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	videoExtensions := map[string]bool{".mp4": true, ".mkv": true, ".avi": true, ".mov": true, ".wmv": true, ".flv": true}
	pictureExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true, ".svg": true, ".webp": true}
	audioExtensions := map[string]bool{".mp3": true, ".wav": true, ".flac": true, ".aac": true, ".ogg": true, ".wma": true, ".m4a": true}

	switch {
	case videoExtensions[ext]:
		return "video"
	case pictureExtensions[ext]:
		return "picture"
	case audioExtensions[ext]:
		return "audio"
	default:
		return "document"
	}
}
