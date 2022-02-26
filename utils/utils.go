package utils

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func GetRshinMamoBaseDirPath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return filepath.Join(homedir, "rshinmemo"), nil
}
