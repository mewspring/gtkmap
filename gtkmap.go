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
//
// int getSource(GtkWidget *m) {
//    int source;
//    g_object_get(m, "map-source", &source, NULL);
//    return source;
// }
import "C"
import "fmt"
import "unsafe"

import "github.com/mattn/go-gtk/gtk"

// n returns the Widget as a *C.OsmGpsMap.
func n(m *Map) *C.OsmGpsMap {
	return (*C.OsmGpsMap)(unsafe.Pointer(m.Widget.GWidget))
}

// w returns the Widget as a *C.GtkWidget.
func w(m *Map) *C.GtkWidget {
	return (*C.GtkWidget)(unsafe.Pointer(m.Widget.GWidget))
}

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

// NewMap returns a new Map which uses tile representations from source.
func NewMapWithSource(source Source) (m *Map, err error) {
	if !source.IsValid() {
		return nil, fmt.Errorf("gtkmap.NewMapWithSource: invalid source (%q).", source)
	}
	m = &Map{
		Widget: gtk.WidgetFromNative(unsafe.Pointer(C.newMapWithSource(C.int(source)))),
	}
	return m, nil
}

// SetCenter centers the map on the provided longitude and latitude, which are
// represented in degrees.
func (m *Map) SetCenter(lat, long float64) {
	C.osm_gps_map_set_center(n(m), C.float(lat), C.float(long))
}

// SetCenterAndZoom centers the map on the provided longitude and latitude with
// the provided zoom level. Latitude and longitude are represented in degrees.
func (m *Map) SetCenterAndZoom(lat, long float64, zoom int) {
	C.osm_gps_map_set_center_and_zoom(n(m), C.float(lat), C.float(long), C.int(zoom))
}

// SetZoom sets the zoom level of the map. It returns the new zoom level, which
// may differ if zoom was below the min or above the max zoom level of the
// current source.
func (m *Map) SetZoom(zoom int) int {
	return int(C.osm_gps_map_set_zoom(n(m), C.int(zoom)))
}

// ZoomIn increases the zoom level by one. It returns the new zoom level.
func (m *Map) ZoomIn() int {
	return int(C.osm_gps_map_zoom_in(n(m)))
}

// ZoomOut decreases the zoom level by one. It returns the new zoom level.
func (m *Map) ZoomOut() int {
	return int(C.osm_gps_map_zoom_out(n(m)))
}

// AddGPS adds a GPS marker to the map with the provided latitude, longitude and
// heading. Latitude and longitude are represented in degrees.
func (m *Map) AddGPS(lat, long, heading float64) {
	C.osm_gps_map_gps_add(n(m), C.float(lat), C.float(long), C.float(heading))
}

// Scroll scrolls the map by dx, dy pixels (positive north, east).
func (m *Map) Scroll(dx, dy int) {
	C.osm_gps_map_scroll(n(m), C.gint(dx), C.gint(dy))
}

// Scale returns the scale at the center of the map, in meters/pixel.
func (m *Map) Scale() float64 {
	return float64(C.osm_gps_map_get_scale(n(m)))
}

// ScreenToCoord converts the provided pixel location on the map into the
// corresponding coordinate. Latitude and longitude are represented in degrees.
func (m *Map) ScreenToCoord(x, y int) (lat, long float64) {
	// Convert from pixel location (x, y) to coordinate (lat, long) in radians.
	var pt C.OsmGpsMapPoint
	C.osm_gps_map_convert_screen_to_geographic(n(m), C.gint(x), C.gint(y), &pt)
	// Convert from coordinate (lat, long) in radians to coordinate (lat, long)
	// in degrees.
	var clat C.float
	var clong C.float
	C.osm_gps_map_point_get_degrees(&pt, &clat, &clong)
	return float64(clat), float64(clong)
}

// CoordToScreen converts the provided coordinate on the map into the
// corresponding pixel location. Latitude and longitude are represented in
// degrees.
func (m *Map) CoordToScreen(lat, long float64) (x, y int) {
	// Convert from coordinate (lat, long) in degrees to coordinate (lat, long)
	// in radians.
	var pt *C.OsmGpsMapPoint
	pt = C.osm_gps_map_point_new_degrees(C.float(lat), C.float(long))
	// Convert from coordinate (lat, long) in radians to pixel location (x, y).
	var cx C.gint
	var cy C.gint
	C.osm_gps_map_convert_geographic_to_screen(n(m), pt, &cx, &cy)
	return int(cx), int(cy)
}

// Source returns the current source tile repository.
func (m *Map) Source() (source Source) {
	source = Source(C.getSource(w(m)))
	if source == -1 {
		// default.
		source = SourceOpenStreetMap1
	}
	return source
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

func (source Source) String() string {
	return C.GoString(C.osm_gps_map_source_get_friendly_name(C.OsmGpsMapSource_t(source)))
}

// MinZoom returns the minimum zoom level of source. At zoom level 1 the world
// is 512x512 pixels.
func (source Source) MinZoom() int {
	return int(C.osm_gps_map_source_get_min_zoom(C.OsmGpsMapSource_t(source)))
}

// MaxZoom returns the maximum zoom level of source.
func (source Source) MaxZoom() int {
	return int(C.osm_gps_map_source_get_max_zoom(C.OsmGpsMapSource_t(source)))
}

// IsValid returns true if the source tile repository is valid for use.
func (source Source) IsValid() bool {
	if C.osm_gps_map_source_is_valid(C.OsmGpsMapSource_t(source)) == 1 {
		return true
	}
	return false
}
