package ui

import (
	"fmt"

	"github.com/apoorvam/goterminal"
	"github.com/gookit/color"
)

//UpdateParseUI updates parser ui
func UpdateParseUI(writer *goterminal.Writer, engine string, dorksTotal, parsedURLTotal, validTotal, parsedDorks int) {
	value := fmt.Sprintf("[search engine >> <green>%s</>]\n[total dorks >> <green>%d</>]\n[parsed urls >> <green>%d</>]\n[valid urls >> <green>%d</>]\n[parsed dorks >> <green>%d</>]\n",
		engine, dorksTotal, parsedURLTotal, validTotal, parsedDorks)
	writer.Clear()
	color.Fprint(writer, value)
	writer.Print()
}
