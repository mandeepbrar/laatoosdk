package utils

import (
	"io"
	"path"
	"path/filepath"

	"laatoo.io/sdk/constants"
	"laatoo.io/sdk/ctx"
)

type FileTransform func(io.Reader, io.Writer) error

func GetAbsFilePath(ctx ctx.Context, fpath string) string {
	if filepath.IsAbs(fpath) {
		return fpath
	} else {
		basedir, _ := ctx.GetString(constants.BASEDIR)
		return path.Join(basedir, fpath)
	}
}
