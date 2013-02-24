// Package gtkmap provides a GTK map widget with support for GPS coordinates.
// It uses osm-gps-map as a backend.
package gtkmap

// #cgo pkg-config: osmgpsmap
//
// #include <gtk/gtk.h>
// #include <osmgpsmap/osm-gps-map.h>
//
// GtkWidget* newMapWithSource(int source) {
//    GtkWidget *m = g_object_new(
//       OSM_TYPE_GPS_MAP,
//       "map-source", source,
//       NULL
//    );
//    return m;
// }
import "C"
import "unsafe"

import "github.com/mattn/go-gtk/gtk"

// Map is a widget for displaying a map, optionally overlaid with tracks of GPS
// coordinates, images, points of interest or on screen display controls.
//
// Map data is downloaded (and cached for offline use) from a number of
// websites, including http://www.openstreetmap.org.
type Map struct {
	*gtk.Widget
}

// NewMap returns a new Map which uses tile representations from OpenStreetMap
// by default. Use NewMapWithSource to use a custom source.
func NewMap() (m *Map) {
	m = &Map{
		Widget: gtk.WidgetFromNative(unsafe.Pointer(C.osm_gps_map_new())),
	}
	return m
}

// Source represents the tile repository to use.
type Source int

// Source tile repositories.
const (
	SourceNone Source = C.OSM_GPS_MAP_SOURCE_NULL
	// default.
	SourceOpenStreetMap1        Source = C.OSM_GPS_MAP_SOURCE_OPENSTREETMAP
	SourceOpenStreetMap2        Source = C.OSM_GPS_MAP_SOURCE_OPENSTREETMAP_RENDERER
	SourceOpenAerialMap         Source = C.OSM_GPS_MAP_SOURCE_OPENAERIALMAP
	SourceMapsForFree           Source = C.OSM_GPS_MAP_SOURCE_MAPS_FOR_FREE
	SourceOpenCycleMap          Source = C.OSM_GPS_MAP_SOURCE_OPENCYCLEMAP
	SourcePublicTransport       Source = C.OSM_GPS_MAP_SOURCE_OSM_PUBLIC_TRANSPORT
	SourceGoogleMaps            Source = C.OSM_GPS_MAP_SOURCE_GOOGLE_STREET
	SourceGoogleSatellite       Source = C.OSM_GPS_MAP_SOURCE_GOOGLE_SATELLITE
	SourceGoogleHybrid          Source = C.OSM_GPS_MAP_SOURCE_GOOGLE_HYBRID
	SourceVirtualEarth          Source = C.OSM_GPS_MAP_SOURCE_VIRTUAL_EARTH_STREET
	SourceVirtualEarthSatellite Source = C.OSM_GPS_MAP_SOURCE_VIRTUAL_EARTH_SATELLITE
	SourceVirtualEarthHybrid    Source = C.OSM_GPS_MAP_SOURCE_VIRTUAL_EARTH_HYBRID
	SourceYahooMaps             Source = C.OSM_GPS_MAP_SOURCE_YAHOO_STREET
	SourceYahooSatellite        Source = C.OSM_GPS_MAP_SOURCE_YAHOO_SATELLITE
	SourceYahooHybrid           Source = C.OSM_GPS_MAP_SOURCE_YAHOO_HYBRID
	SourceOSMCTrails            Source = C.OSM_GPS_MAP_SOURCE_OSMC_TRAILS
)

// NewMap returns a new Map which uses tile representations from source.
func NewMapWithSource(source Source) (m *Map) {
	m = &Map{
		Widget: gtk.WidgetFromNative(unsafe.Pointer(C.newMapWithSource(C.int(source)))),
	}
	return m
}

// n returns the native type of the map.
func (m *Map) n() *C.OsmGpsMap {
	return (*C.OsmGpsMap)(unsafe.Pointer(m.Widget.GWidget))
}

// SetCenter centers the map around the provided longitude and latitude.
func (m *Map) SetCenter(lat, long float64) {
	C.osm_gps_map_set_center(m.n(), C.float(lat), C.float(long))
}

// Min and max zoom levels for OpenStreetMap. At zoom level 1 the world is
// 512x512 pixels.
const (
	MinZoomOSM = 1
	MaxZoomOSM = 18
)

// SetZoom sets the zoom level of the map. It returns the new zoom level, which
// may differ if zoom was below the min or above the max zoom level of the
// current source.
func (m *Map) SetZoom(zoom int) int {
	return int(C.osm_gps_map_set_zoom(m.n(), C.int(zoom)))
}

// ZoomIn increases the zoom level by one. It returns the new zoom level.
func (m *Map) ZoomIn() int {
	return int(C.osm_gps_map_zoom_in(m.n()))
}

// ZoomOut decreases the zoom level by one. It returns the new zoom level.
func (m *Map) ZoomOut() int {
	return int(C.osm_gps_map_zoom_out(m.n()))
}
