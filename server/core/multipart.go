package core

import "io"

type MultipartFile struct {
	File     io.ReadCloser
	FileName string
	MimeType string
}
