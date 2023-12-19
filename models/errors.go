package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

var (
	ErrEmailExists  = errors.New("models: users with such email already exists")
	ErrNotFound     = errors.New("models: couldn't find an entity")
	ErrUnauthorized = errors.New("models: insufficient level of access")
)

type FileError struct {
	Issue string
}

func (err FileError) Error() string {
	return fmt.Sprintf("file error: %v", err.Issue)
}

func checkContentType(content io.Reader, allowedContentTypes []string) ([]byte, error) {
	testBytes := make([]byte, 512)
	n, err := content.Read(testBytes)
	if err != nil {
		return testBytes, FileError{
			Issue: fmt.Sprintf("checking content type: %v", err),
		}
	}

	contentType := http.DetectContentType(testBytes)
	for _, allowedContentType := range allowedContentTypes {
		if contentType == allowedContentType {
			return testBytes[:n], nil
		}
	}

	return testBytes, FileError{
		Issue: fmt.Sprintf("invalid content type %s", contentType),
	}
}

func checkExtension(filename string, allowedExtensions []string) error {
	if hasExtension(filename, allowedExtensions) {
		return nil
	}
	return FileError{
		Issue: fmt.Sprintf("invalid extension %s", filepath.Ext(filename)),
	}
}
