package ui

import (
	"os"

	"github.com/apoorvam/goterminal"
)

//SetupGeneralUI creates and starts new writer
func SetupGeneralUI() *goterminal.Writer {
	writer := goterminal.New(os.Stdout)
	return writer
}

//StopGeneralUI stops writer
func StopGeneralUI(writer *goterminal.Writer) {
	writer.Reset()
}
