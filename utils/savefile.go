package utils

import (
	"io"
	"os"
	"path/filepath"
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
