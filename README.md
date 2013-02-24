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

Examples
--------

Simple example which creates a new GTK window with a map widget and center the
map around Iceland.

    go get github.com/mewmew/gtkmap/examples/mapview

![Screenshot - OpenStreetMap](https://github.com/mewmew/gtkmap/blob/master/examples/mapview/mapview.png?raw=true)

This example uses Google Maps as source for the map tiles (default is
OpenStreetMap).

    go get github.com/mewmew/gtkmap/examples/gmapview

![Screenshot - Google Maps](https://github.com/mewmew/gtkmap/blob/master/examples/gmapview/gmapview.png?raw=true)

public domain
-------------

This code is hereby released into the *[public domain][]*.

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
