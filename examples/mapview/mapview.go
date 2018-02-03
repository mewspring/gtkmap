// mapview is a simple example tool which creates a new GTK window with a map
// widget and center the map on Iceland.
package main

import (
	"os"

	"github.com/mewspring/gtkmap"
	gtk "github.com/zurek87/go-gtk3/gtk3"
)

func main() {
	gtk.Init(&os.Args)

	// Create a window.
	win := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	win.Connect("destroy", gtk.MainQuit)

	// Create a map widget.
	m := gtkmap.NewMap()
	m.SetSizeRequest(640, 480)
	win.Add(m)

	// Center the map on Iceland.
	m.SetCenter(gtkmap.Coord(45.963051, -76.020835))

	// Set zoom level.
	m.SetZoom(6)

	win.ShowAll()
	gtk.Main()
}
