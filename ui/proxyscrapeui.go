package ui

import (
	"fmt"
	"os"

	"github.com/apoorvam/goterminal"
	"github.com/gookit/color"
)

//SetupProxyScrapeUI creates and starts new writer
func SetupProxyScrapeUI() *goterminal.Writer {
	writer := goterminal.New(os.Stdout)
	return writer
}

//UpdateProxyScrapeUI updates parser ui
func UpdateProxyScrapeUI(writer *goterminal.Writer, sourcesTotal, doneSources, scrapedProxiesTotal, validProxiesTotal int) {
	value := fmt.Sprintf("[total sources >> <green>%d</>]\n[done sources >> <green>%d</>]\n[scraped proxies >> <green>%d</>]\n[valid proxies >> <green>%d</>]\n",
		sourcesTotal, doneSources, scrapedProxiesTotal, validProxiesTotal)
	writer.Clear()
	color.Fprint(writer, value)
	writer.Print()
}

//StopProxyScrapeUI stops writer
func StopProxyScrapeUI(writer *goterminal.Writer) {
	writer.Reset()
}
