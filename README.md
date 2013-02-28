gtkmap
======

This package provides a GTK map widget with support for GPS coordinates. It uses
[osm-gps-map][] as a backend.

[osm-gps-map]: http://nzjrs.github.com/osm-gps-map/

Documentation
-------------

Documentation provided by GoDoc.

   - [gtkmap][]

[gtkmap]: http://godoc.org/github.com/mewmew/gtkmap

Installation
------------

Install the [osm-gps-map] dependency and run:

	go get github.com/mewmew/gtkmap

Examples
--------

mapview is a simple example which creates a new GTK window with a map widget and
center the map on Iceland.

	go get github.com/mewmew/gtkmap/examples/mapview

![Screenshot - OpenStreetMap](https://github.com/mewmew/gtkmap/blob/master/examples/mapview/mapview.png?raw=true)

The gmapview example uses Google Maps as source for the map tiles (the default
is OpenStreetMap).

	go get github.com/mewmew/gtkmap/examples/gmapview

![Screenshot - Google Maps](https://github.com/mewmew/gtkmap/blob/master/examples/gmapview/gmapview.png?raw=true)

gpsview
=======

gpsview parses image GPS coordinates and plots them on a map. The tile source
repository and cache settings are customizeable.

Installation
------------

	go get github.com/mewmew/gtkmap/cmd/gpsview

Usage
-----

	gpsview [OPTION]... [IMAGE]...

Flags:

	-cache (default="")
		Cache directory ("" represent "$HOME/.cache", "none://" disables cache.).
	-lat (default=20.793415)
		Latitude.
	-long (default=106.99894
		Longitude.
	-s (default=11)
		Tile source repository (1-16).
	-v (default=false
		Verbose.
	-z (default=11)
		Zoom level (1-18).

Mouse button events:

	* left double-click
		Center on mouse cursor and zoom in.
	* right double-click
		Center on mouse cursor and zoom out.

	* middle click
		Print coordinate at mouse cursor.

	* [shift] + left click
		Add GPS marker at mouse cursor.
	* [ctrl] + left click
		Clear lines between GPS markers.


Examples
--------

1. Plot all images in the "images/" directory.

		gpsview images/*

![Screenshot - Ha Long Bay](https://github.com/mewmew/gtkmap/blob/master/cmd/gpsview/gpsview1.png?raw=true)

2. Disable cache, use Google Maps as source and set zoom level to 16.

		gpsview -cache="none://" -s=6 -z=16 *

![Screenshot - Angkor Wat](https://github.com/mewmew/gtkmap/blob/master/cmd/gpsview/gpsview2.png?raw=true)

public domain
-------------

This code is hereby released into the *[public domain][]*.

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
