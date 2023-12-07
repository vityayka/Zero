package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

var (
	ErrEmailExists  error = errors.New("models: users with such email already exists")
	ErrNotFound     error = errors.New("models: couldn't find an entity")
	ErrUnauthorized error = errors.New("models: insufficient level of access")
)

type FileError struct {
	Issue string
}

func (err FileError) Error() string {
	return fmt.Sprintf("file error: %v", err.Issue)
}

func checkContentType(content io.ReadSeeker, allowedContentTypes []string) error {
	testBytes := make([]byte, 512)
	_, err := content.Read(testBytes)
	if err != nil {
		return FileError{
			Issue: fmt.Sprintf("checking content type: %v", err),
		}
	}
	_, err = content.Seek(0, 0)
	if err != nil {
		return FileError{
			Issue: fmt.Sprintf("checking content type: %v", err),
		}
	}

	contentType := http.DetectContentType(testBytes)
	for _, allowedContentType := range allowedContentTypes {
		if contentType == allowedContentType {
			return nil
		}
	}

	return FileError{
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
