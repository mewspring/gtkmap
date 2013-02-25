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

Run the following command after installing the [osm-gps-map] dependency:

	$ go get github.com/mewmew/gtkmap

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

The gpsview example can parse image GPS coordinates and plot them on a map.
Command line flags control the tile representation source and map zoom level.

    go get github.com/mewmew/gtkmap/examples/gpsview

![Screenshot - Ha Long Bay](https://github.com/mewmew/gtkmap/blob/master/examples/gpsview/gpsview1.png?raw=true)

![Screenshot - Angkor Wat](https://github.com/mewmew/gtkmap/blob/master/examples/gpsview/gpsview2.png?raw=true)

public domain
-------------

This code is hereby released into the *[public domain][]*.

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
