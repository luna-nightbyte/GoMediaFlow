package local

import (
	"errors"
	"goStreamer/modules/config"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var Files Output

type Output struct {
	source folder
	target folder
	output folder
}
type folder struct {
	folder string
}

func (o *Output) Update(sourceFolder, targetFolder, outputFolder string) {
	o.source.set(sourceFolder)
	o.target.set(targetFolder)
	o.output.set(outputFolder)
	config.Config.Local.SourceFolder = sourceFolder
	config.Config.Local.SourceFolder = targetFolder
	config.Config.Local.SourceFolder = outputFolder
	config.Config.Update()
}
func (f *folder) set(input string) {
	f.folder = input
}
func (o *Output) UpdateSingle(sourceFolder, webcamTarget string) {
	o.source.folder = sourceFolder
	config.Config.Local.SourceFolder = sourceFolder
	config.Config.Local.Webcam.Target = webcamTarget
	config.Config.Update()
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

func (o *Output) SourceFolder() string {
	return o.source.folder
}
func (o *Output) TargetFolder() string {
	return o.target.folder
}
func (o *Output) OutputFolder() string {
	return o.output.folder
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
