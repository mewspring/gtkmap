package main

import "os"

import "github.com/mattn/go-gtk/gtk"
import "github.com/mewmew/gtkmap"

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
	m.SetCenter(64.963051, -19.020835)

	// Set zoom level.
	m.SetZoom(6)

	win.ShowAll()
	gtk.Main()
}
