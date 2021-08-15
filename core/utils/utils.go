package utils

import (
	"github.com/pkg/errors"
	"os"
	"os/exec"
)

// vimで対象noteを開く
func OpenVim(filePath string) error {
	c := exec.Command("vim", filePath)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		return errors.Wrap(err, "vim起動エラー")
	}
	return nil
}
