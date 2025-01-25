package files

import (
	"errors"
	"goStreamer/modules/config"
	"goStreamer/modules/db"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Output struct {
	source file
	target file
	output file
}
type file struct {
	file string
}

func (o *Output) Update(source, target, output string) {
	o.source.file = source
	o.target.file = target
	o.output.file = output
	config.Config.InputSource = o.source.file
	config.Config.InputTarget = o.target.file
	config.Config.OutputFile = o.output.file
	db.Write("config.json", config.Config)
}

func (o *Output) UpdateSingle(source, webcam string) {
	o.source.file = source
	config.Config.InputSource = o.source.file
	config.Config.InputTarget = webcam
	db.Write("config.json", config.Config)
}

// IsFile checks if the file is a video based on its MIME type.
func IsFileAndExist(filePath, fileType string) bool {
	ok, err := isFileOfType(filePath, fileType)
	if err != nil {
		log.Println("Error checking filetype:", err)
	}
	return ok
}

// isFileOfType is a helper function to determine the file type based on MIME type prefix.
func isFileOfType(filePath string, fileType string) (bool, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read the file header to determine MIME type
	buffer := make([]byte, 512) // 512 bytes for detecting MIME type
	_, err = file.Read(buffer)
	if err != nil {
		return false, err
	}

	// Detect the content type
	contentType := http.DetectContentType(buffer)

	// Parse the MIME type
	mimeType := strings.Split(contentType, "/")
	if len(mimeType) < 2 {
		return false, errors.New("invalid MIME type")
	}

	// Check if the MIME type matches the expected file type
	return mimeType[0] == fileType, nil
}

func (o *Output) Source() string {
	return o.source.file
}

func (o *Output) Target() string {
	return o.target.file
}

func (o *Output) Output() string {
	return o.output.file
}

func IsVideoOrImageFileName(fileName string) bool {
	// Define supported video and image extensions
	videoExtensions := map[string]bool{
		".mp4":  true,
		".mkv":  true,
		".avi":  true,
		".mov":  true,
		".wmv":  true,
		".flv":  true,
		".webm": true,
		".mpeg": true,
		".3gp":  true,
	}

	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".tiff": true,
		".svg":  true,
		".webp": true,
	}

	// Extract the file extension
	ext := strings.ToLower(filepath.Ext(fileName))

	// Check against the extensions maps
	isVideo := videoExtensions[ext]
	isImage := imageExtensions[ext]

	return (isVideo || isImage)
}
