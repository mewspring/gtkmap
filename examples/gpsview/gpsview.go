package main

import "flag"
import "fmt"
import "log"
import "math"
import "os"
import "unsafe"

import "github.com/mattn/go-gtk/gdk"
import "github.com/mattn/go-gtk/glib"
import "github.com/mattn/go-gtk/gtk"
import "github.com/mewmew/gtkmap"
import "github.com/rwcarlsen/goexif/exif"
import gps "github.com/kurtcc/goexifgps"

// flagLat specifies the latitude on which to center the map.
var flagLat float64

// flagLong specifies the longitude on which to center the map.
var flagLong float64

// flagSource specifies the tile representation source.
var flagSource int

// flagZoom specifies the zoom level of the map.
var flagZoom int

func init() {
	// Cat Ba.
	lat := 20.793415
	long := 106.99894
	flag.Float64Var(&flagLat, "lat", lat, "Latitude.")
	flag.Float64Var(&flagLong, "long", long, "Longitude.")
	flag.IntVar(&flagSource, "s", int(gtkmap.SourceVirtualEarthSatellite), "Tile representation source (1-16).")
	flag.IntVar(&flagZoom, "z", 11, "Zoom level (1-20).")
}

func main() {
	// Parse image GPS coordinates.
	flag.Parse()
	var coords []*gps.GeoFields
	for _, imgPath := range flag.Args() {
		coord, err := getCoordinate(imgPath)
		if err != nil {
			log.Println(err)
			continue
		}
		coords = append(coords, coord)
	}

	gtk.Init(&os.Args)

	// Create a window.
	win := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	win.Connect("destroy", gtk.MainQuit)

	// Create a map widget which uses flagSource as source for the map tiles.
	m, err := gtkmap.NewMapWithSource(gtkmap.Source(flagSource))
	if err != nil {
		// Fall back to using OpenStreetMap if Google Maps could not be used as
		// source.
		m = gtkmap.NewMap()
	}

	fmt.Println("Map tile representations from:", m.Source())
	m.SetSizeRequest(640, 480)
	win.Add(m)

	if len(coords) == 0 {
		// Center on the provided latitude and longitude if no image GPS
		// coordinates were found .
		fmt.Println("Center on latitude and longitude from flags.")
		m.SetCenter(flagLat, flagLong)
	} else {
		for _, coord := range coords {
			// Place GPS coordinates on the map.
			lat := float64(coord.Lat)
			long := float64(coord.Long)
			m.AddGPS(lat, long, 0)

			// Center on the last GPS coordinate.
			m.SetCenter(lat, long)
		}
	}

	// Set zoom level.
	m.SetZoom(flagZoom)

	// Connect mouse button events.
	onButtonPress := func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		ev := (*gdk.EventButton)(unsafe.Pointer(arg))
		// Double click.
		if ev.Type == int(gdk.BUTTON2_PRESS) {
			lat, long := m.ScreenToCoord(int(ev.X), int(ev.Y))
			switch ev.Button {
			case 1:
				// Left mouse button.
				m.SetCenter(lat, long)
				m.ZoomIn()
			case 2:
				// Middle mouse button.
				m.ClearGPS()
			case 3:
				// Right mouse button.
				m.SetCenter(lat, long)
				m.ZoomOut()
			}
		}
	}
	m.Widget.Connect("button-press-event", onButtonPress)

	win.ShowAll()
	gtk.Main()
}

// getCoordinate returns the GPS coordinate of an image. The information is
// stored in the EXIF data of the image.
func getCoordinate(imgPath string) (coord *gps.GeoFields, err error) {
	fr, err := os.Open(imgPath)
	if err != nil {
		return nil, err
	}
	defer fr.Close()
	x, err := exif.Decode(fr)
	if err != nil {
		return nil, err
	}
	coord, err = gps.GetGPS(x)
	if err != nil {
		return nil, err
	}
	if math.IsNaN(float64(coord.Lat)) || math.IsNaN(float64(coord.Long)) {
		return nil, fmt.Errorf("getCoordinate: unable to locate lat and long in EXIF data for %q.", imgPath)
	}
	return coord, nil
}
