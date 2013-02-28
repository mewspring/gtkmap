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

// flagCache specifies the cache directory for the map tiles.
var flagCache string

// flagLat specifies the latitude on which to center the map.
var flagLat float64

// flagLong specifies the longitude on which to center the map.
var flagLong float64

// flagSource specifies the tile source repository.
var flagSource = gtkmap.SourceVirtualEarthSatellite

// When flagVerbose is true, enable verbose output.
var flagVerbose bool

// flagZoom specifies the zoom level of the map.
var flagZoom int

func init() {
	// Cat Ba.
	coord := gtkmap.Coord(20.793415, 106.99894)
	flag.StringVar(&flagCache, "cache", string(gtkmap.CacheDefault), `Cache directory ("" represent "$HOME/.cache", "none://" disables cache.).`)
	flag.Float64Var(&flagLat, "lat", coord.Lat, "Latitude.")
	flag.Float64Var(&flagLong, "long", coord.Long, "Longitude.")
	flag.Var(&flagSource, "s", "Tile source repository (1-16).")
	flag.BoolVar(&flagVerbose, "v", false, "Verbose.")
	flag.IntVar(&flagZoom, "z", 11, "Zoom level (1-18).")
	flag.Usage = usage
}

/// ### [ todo ] ###
///    - extract date and sort images before plotting them.
/// ### [/ todo ] ###

func usage() {
	fmt.Fprintln(os.Stderr, "gpsview [OPTION]... [IMAGE]...")
	fmt.Fprintln(os.Stderr, "Parses image GPS coordinates and plots them on a map.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, `  Plot all images in the "images/" directory.`)
	fmt.Fprintln(os.Stderr, "    gpsmap images/*")
	fmt.Fprintln(os.Stderr, "  Disable cache, use Google Maps as source and set zoom level to 16.")
	fmt.Fprintln(os.Stderr, `    gpsmap -cache="none://" -s=6 -z=16 *`)
}

func main() {
	// Parse image GPS coordinates.
	flag.Parse()
	var geoCoords []*gps.GeoFields
	for _, imgPath := range flag.Args() {
		geoCoord, err := getCoordinate(imgPath)
		if err != nil {
			if flagVerbose {
				log.Println(err)
			}
			continue
		}
		geoCoords = append(geoCoords, geoCoord)
	}

	gtk.Init(&os.Args)

	// Create a window.
	win := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	win.Connect("destroy", gtk.MainQuit)

	// Create a map widget which uses flagSource as tile source repository and
	// flagCache as cache directory.
	m, err := gtkmap.NewMapOpt(flagSource, gtkmap.Cache(flagCache))
	if err != nil {
		// Fall back to using OpenStreetMap if flagSource could not be used as
		// source.
		m = gtkmap.NewMap()
	}

	if flagVerbose {
		fmt.Println("Tile source repository:", m.Source().FriendlyName())
	}
	m.SetSizeRequest(640, 480)
	win.Add(m)

	if len(geoCoords) == 0 {
		// Center on the provided latitude and longitude if no image GPS
		// coordinates were found .
		coord := gtkmap.Coord(flagLat, flagLong)
		if flagVerbose {
			fmt.Printf("Center on coordinate %v provided from flags.\n", coord)
		}
		m.SetCenter(coord)
	} else {
		for i, geoCoord := range geoCoords {
			// Place GPS coordinates on the map.
			lat := float64(geoCoord.Lat)
			long := float64(geoCoord.Long)
			coord := gtkmap.Coord(lat, long)
			m.AddGPS(coord, 0)
			if flagVerbose {
				fmt.Printf("Add GPS marker at coordinate %v.\n", coord)
			}

			if i == len(geoCoords)-1 {
				// Center on the last GPS coordinate.
				m.SetCenter(coord)
				if flagVerbose {
					fmt.Printf("Center on coordinate %v provided from image.\n", coord)
				}
			}
		}
	}

	// Set zoom level.
	m.SetZoom(flagZoom)

	// Connect mouse button events.
	onButtonPress := func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		ev := (*gdk.EventButton)(unsafe.Pointer(arg))
		coord := m.ScreenToCoord(int(ev.X), int(ev.Y))

		// Single-click.
		if ev.Type == int(gdk.BUTTON_PRESS) {
			switch ev.Button {
			case 1:
				// left click
				switch {
				case ev.State&uint32(gdk.SHIFT_MASK) != 0:
					// [shift] + left click
					m.AddGPS(coord, 0)
				case ev.State&uint32(gdk.CONTROL_MASK) != 0:
					// [ctrl] + left click
					m.ClearGPS()
				}
			case 2:
				// middle click
				fmt.Println("coordinate:", coord)
			}
		}

		// Double-click.
		if ev.Type == int(gdk.BUTTON2_PRESS) {
			switch ev.Button {
			case 1:
				// left double-click
				m.SetCenter(coord)
				m.ZoomIn()
			case 3:
				// right double-click
				m.SetCenter(coord)
				m.ZoomOut()
			}
		}
	}
	m.Widget.Connect("button-press-event", onButtonPress)

	win.ShowAll()
	gtk.Main()
}

// getCoordinate returns the GPS coordinate of an image. The information is
// stored in the image's EXIF data.
func getCoordinate(imgPath string) (coord *gps.GeoFields, err error) {
	fr, err := os.Open(imgPath)
	if err != nil {
		return nil, err
	}
	defer fr.Close()
	x, err := exif.Decode(fr)
	if err != nil {
		return nil, fmt.Errorf("exif.Decode: failed for %q; %s.", imgPath, err.Error())
	}
	coord, err = gps.GetGPS(x)
	if err != nil {
		return nil, fmt.Errorf("gps.GetGPS: failed for %q; %s.", imgPath, err.Error())
	}
	if math.IsNaN(float64(coord.Lat)) || math.IsNaN(float64(coord.Long)) {
		return nil, fmt.Errorf("getCoordinate: failed for %q; unable to locate lat and long in EXIF data.", imgPath)
	}
	return coord, nil
}
