package utils

import (
	"io"
)

type FileTransform func(io.Reader, io.Writer) error
