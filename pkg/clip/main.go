package clip

import (
	"github.com/charmbracelet/log"
	"golang.design/x/clipboard"
)

func CopyToClipboard(text string, verbose int) {
	err := clipboard.Write(clipboard.FmtText, []byte(text))
	if err != nil {
		// Do not exit on error
		if verbose > 0 {
			log.Info("failed to copy to clipboard", "error", err)
		}

	}
}
