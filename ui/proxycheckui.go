package ui

import (
	"fmt"

	"github.com/apoorvam/goterminal"
	"github.com/gookit/color"
)

//UpdateProxyCheckUI updates checker ui
func UpdateProxyCheckUI(writer *goterminal.Writer, proxiesTotal, checkTotal, workingTotal, brokenTotal int) {
	value := fmt.Sprintf("[total proxies >> <green>%d</>]\n[checked proxies >> <green>%d</>]\n[working proxies >> <green>%d</>]\n[broken proxies >> <red>%d</>]\n",
		proxiesTotal, checkTotal, workingTotal, brokenTotal)
	writer.Clear()
	color.Fprint(writer, value)
	writer.Print()
}
