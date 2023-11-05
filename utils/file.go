package utils

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// CreateImageFile creates a new image file in the assets/image folder
// name.file_extension
// Example: image.png
func CreateImageFile(file []byte, name string) error {
	checkFolderExist("./assets/image")
	filePath := "./assets/image/" + name
	if err := os.WriteFile(filePath, file, 0644); err != nil {
		return err
	}
	return nil
}

// Delete Image file in the assets/image folder
// name.file_extension
// Example: image.png

func DeleteImageFile(fileName string) error {
	err := os.Remove("./assets/image/" + fileName)
	if err != nil {
		return err
	}

	return nil
}

func ConvertImageFileToBytes(file *multipart.FileHeader) ([]byte, error) {
	fileContent, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileContent.Close()

	fileBytes, err := io.ReadAll(fileContent)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

func ConvertImageToPNG(imageBytes []byte) ([]byte, error) {
	contentType := http.DetectContentType(imageBytes)

	switch contentType {
	case "image/png":
		return imageBytes, nil // Already in PNG format
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, fmt.Errorf("unable to decode jpeg: %v", err.Error())
		}

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return nil, fmt.Errorf("unable to encode png: %v", err.Error())
		}

		return buf.Bytes(), nil
	default:
		return nil, fmt.Errorf("unsupported image format: %v", contentType)
	}
}

func CheckImageContentType(contentType string) error {
	if contentType == "image/jpeg" || contentType == "image/png" {
		return nil
	}
	return fmt.Errorf("unsupported content type: only jpeg and png formats are supported")
}

func checkFolderExist(folderPath string) error {
	// Check if the folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// Folder doesn't exist, so create it
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create folder: %v\n", err)
			return err
		}
		fmt.Printf("Folder created: %s\n", folderPath)
	}
	return nil
}
