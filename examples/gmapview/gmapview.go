// gmapview is a simple example tool which uses Google Maps as source for the
// map tiles (the default is OpenStreetMap).
package main

import (
	"fmt"
	"os"

	"github.com/mewspring/gtkmap"
	gtk "github.com/zurek87/go-gtk3/gtk3"
)

func main() {
	gtk.Init(&os.Args)

	// Create a window.
	win := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	win.Connect("destroy", gtk.MainQuit)

	// Create a map widget which uses Google Maps as source for the map tiles.
	source := gtkmap.SourceGoogleMaps
	m, err := gtkmap.NewMapOpt(source)
	if err != nil {
		// Fall back to using OpenStreetMap if Google Maps could not be used as
		// source.
		m = gtkmap.NewMap()
	}
	fmt.Println("Tile source repository:", m.Source().FriendlyName())
	m.SetSizeRequest(640, 480)
	win.Add(m)

	// Center the map on Guangzhou and set zoom level.
	coord := gtkmap.Coord(45.963051, -76.020835)
	zoom := 14
	m.SetCenterAndZoom(coord, zoom)

	win.ShowAll()
	gtk.Main()
}
