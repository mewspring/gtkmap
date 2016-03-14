// Package gtkmap provides a GTK map widget with support for GPS coordinates.
// It uses osm-gps-map as a backend.
package gtkmap

// #cgo pkg-config: osmgpsmap-1.0
//
// #include <string.h>
// #include <gtk/gtk.h>
// #include <osm-gps-map.h>
//
// GtkWidget* newMapCustom(int source, char *cache) {
//    GtkWidget *m;
//    char *tile_cache = "friendly://";
//    char *tile_cache_base = NULL;
//
//    if(strlen(cache) == 0) {
//       // default.
//    } else if(strcmp(cache, "none://") == 0) {
//       // cache disabled.
//       tile_cache = cache;
//    } else {
//       // custom cache base dir.
//       tile_cache_base = cache;
//    }
//    m = g_object_new(
//       OSM_TYPE_GPS_MAP,
//       "tile-cache", tile_cache,
//       "tile-cache-base", tile_cache_base,
//       "map-source", source,
//       NULL
//    );
//    return m;
// }
//
// int getDownloadQueueCount(GtkWidget *m) {
//    int n;
//    g_object_get(m, "tiles-queued", &n, NULL);
//    return n;
// }
//
// int getSource(GtkWidget *m) {
//    int source;
//    g_object_get(m, "map-source", &source, NULL);
//    return source;
// }
import "C"

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/mattn/go-gtk/gtk"
)

// A Coordinate is a Lat, Long coordinate pair. The latitude and longitude are
// represented in degrees.
type Coordinate struct {
	// Latitude.
	Lat float64
	// Longitude.
	Long float64
}

// Coord is shorthand for Coordinate{lat, long}.
func Coord(lat, long float64) Coordinate {
	return Coordinate{lat, long}
}

// String returns a string representation of coord like "(3.1,4.2)".
func (coord Coordinate) String() string {
	return fmt.Sprintf("(%v,%v)", coord.Lat, coord.Long)
}

// A Rect contains the coordinates with
// Min.Lat <= Lat < Max.Lat, Min.Long <= Long < Max.Long.
type Rect struct {
	// North west corner.
	Min Coordinate
	// South east corner.
	Max Coordinate
}

// String returns a string representation of rect like "(3.1,4.2)-(6.3,5.4)".
func (rect Rect) String() string {
	return rect.Min.String() + "-" + rect.Max.String()
}

// Add returns the rectangle rect translated by coord.
func (rect Rect) Add(coord Coordinate) Rect {
	return Rect{
		Min: Coord(rect.Min.Lat+coord.Lat, rect.Min.Long+coord.Long),
		Max: Coord(rect.Max.Lat+coord.Lat, rect.Max.Long+coord.Long),
	}
}

// Sub returns the rectangle rect translated by -coord.
func (rect Rect) Sub(coord Coordinate) Rect {
	return Rect{
		Min: Coord(rect.Min.Lat-coord.Lat, rect.Min.Long-coord.Long),
		Max: Coord(rect.Max.Lat-coord.Lat, rect.Max.Long-coord.Long),
	}
}

// Map is a widget for displaying a map, optionally overlaid with tracks of GPS
// coordinates.
//
// Map data is downloaded (and cached for offline use) from a number of
// websites, including http://www.openstreetmap.org.
type Map struct {
	*gtk.Widget
}

// NewMap returns a new Map which uses the default Source and Cache. Use
// NewMapOpt for control over Source and Cache.
func NewMap() (m *Map) {
	m = &Map{
		Widget: gtk.WidgetFromNative(unsafe.Pointer(C.osm_gps_map_new())),
	}
	return m
}

// NewMapOpt returns a new Map which uses the provided options. If no options
// are provided, NewMapOpt() is equivalent to NewMap(). Valid option types are
// Source and Cache.
func NewMapOpt(opts ...interface{}) (m *Map, err error) {
	// Parse options.
	source := SourceDefault
	cache := CacheDefault
	for _, opt := range opts {
		switch v := opt.(type) {
		case Source:
			source = v
			if !source.IsValid() {
				return nil, fmt.Errorf("gtkmap.NewMapOpt: invalid source (%q)", source)
			}
		case Cache:
			cache = v
		default:
			return nil, fmt.Errorf("gtkmap.NewMapOpt: option type (%T) not yet implemented", v)
		}
	}

	m = &Map{
		Widget: gtk.WidgetFromNative(unsafe.Pointer(C.newMapCustom(C.int(source), C.CString(string(cache))))),
	}
	return m, nil
}

// SetCenter centers the map on the provided coordinate.
func (m *Map) SetCenter(coord Coordinate) {
	C.osm_gps_map_set_center(n(m), C.float(coord.Lat), C.float(coord.Long))
}

// SetCenterAndZoom centers the map on the provided coordinate with the provided
// zoom level.
func (m *Map) SetCenterAndZoom(coord Coordinate, zoom int) {
	C.osm_gps_map_set_center_and_zoom(n(m), C.float(coord.Lat), C.float(coord.Long), C.int(zoom))
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

// AddGPS adds a GPS marker to the map with the provided coordinate and heading.
func (m *Map) AddGPS(coord Coordinate, heading float64) {
	C.osm_gps_map_gps_add(n(m), C.float(coord.Lat), C.float(coord.Long), C.float(heading))
}

// ClearGPS clears the lines between all GPS markers.
func (m *Map) ClearGPS() {
	C.osm_gps_map_gps_clear(n(m))
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
// corresponding coordinate.
func (m *Map) ScreenToCoord(x, y int) (coord Coordinate) {
	// Convert from pixel location (x, y) to coordinate (lat, long) in radians.
	var cpt C.OsmGpsMapPoint
	C.osm_gps_map_convert_screen_to_geographic(n(m), C.gint(x), C.gint(y), &cpt)
	// Convert from coordinate (lat, long) in radians to coordinate (lat, long)
	// in degrees.
	var clat C.float
	var clong C.float
	C.osm_gps_map_point_get_degrees(&cpt, &clat, &clong)
	return Coord(float64(clat), float64(clong))
}

// CoordToScreen converts the provided coordinate on the map into the
// corresponding pixel location.
func (m *Map) CoordToScreen(coord Coordinate) (x, y int) {
	// Convert from coordinate (lat, long) in degrees to coordinate (lat, long)
	// in radians.
	var cpt *C.OsmGpsMapPoint
	cpt = C.osm_gps_map_point_new_degrees(C.float(coord.Lat), C.float(coord.Long))
	// Convert from coordinate (lat, long) in radians to pixel location (x, y).
	var cx C.gint
	var cy C.gint
	C.osm_gps_map_convert_geographic_to_screen(n(m), cpt, &cx, &cy)
	return int(cx), int(cy)
}

// DownloadTiles downloads all tiles over the supplied zoom range in the
// rectangular region.
func (m *Map) DownloadTiles(rect Rect, zoomStart, zoomEnd int) {
	// Convert from coordinate (lat, long) in degrees to coordinate (lat, long)
	// in radians.
	var cpt1, cpt2 *C.OsmGpsMapPoint
	cpt1 = C.osm_gps_map_point_new_degrees(C.float(rect.Min.Lat), C.float(rect.Min.Long))
	cpt2 = C.osm_gps_map_point_new_degrees(C.float(rect.Max.Lat), C.float(rect.Max.Long))
	C.osm_gps_map_download_maps(n(m), cpt1, cpt2, C.int(zoomStart), C.int(zoomEnd))
}

// CancelDownloads cancels all tiles currently being downloaded. Typically used
// if you wish to cacel a large number of tiles queued using Map.DownloadTiles.
func (m *Map) CancelDownloads() {
	C.osm_gps_map_download_cancel_all(n(m))
}

// DownloadQueueCount returns the number of tiles currently waiting to download.
func (m *Map) DownloadQueueCount() int {
	return int(C.getDownloadQueueCount(w(m)))
}

// Source returns the current tile source repository.
func (m *Map) Source() (source Source) {
	source = Source(C.getSource(w(m)))
	if source == -1 {
		// default.
		source = SourceOpenStreetMap1
	}
	return source
}

// n returns the Widget as a *C.OsmGpsMap.
func n(m *Map) *C.OsmGpsMap {
	return (*C.OsmGpsMap)(unsafe.Pointer(m.Widget.GWidget))
}

// w returns the Widget as a *C.GtkWidget.
func w(m *Map) *C.GtkWidget {
	return (*C.GtkWidget)(unsafe.Pointer(m.Widget.GWidget))
}

// Source represents the tile source repository to use.
type Source int

// Tile source repositories.
const (
	SourceNone    Source = C.OSM_GPS_MAP_SOURCE_NULL
	SourceDefault Source = SourceOpenStreetMap1
	// Sources.
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
	SourceOSMCTrails            Source = C.OSM_GPS_MAP_SOURCE_OSMC_TRAILS
)

// FriendlyName returns the friendly name of the tile source repository.
func (source Source) FriendlyName() string {
	return C.GoString(C.osm_gps_map_source_get_friendly_name(C.OsmGpsMapSource_t(source)))
}

func (source Source) String() string {
	return strconv.Itoa(int(source))
}

// Set sets the source based on the provided flag value. Source satisfies the
// flag.Value interface.
func (source *Source) Set(s string) (err error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	if n < 1 || n > 16 {
		return fmt.Errorf("Source.Set: invalid source (%d); valid range is 1-16.", n)
	}
	*source = Source(n)
	return nil
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

// IsValid returns true if the tile source repository is valid for use.
func (source Source) IsValid() bool {
	if C.osm_gps_map_source_is_valid(C.OsmGpsMapSource_t(source)) == 1 {
		return true
	}
	return false
}

// Cache specifies the base directory of the tile cache. See CacheDefault and
// CacheDisabled for predefined values.
type Cache string

const (
	// CacheDefault caches map tiles in the users cache directory (as outlined in
	// the XDG Base Directory Specification).
	//
	// ref: http://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html
	CacheDefault Cache = ""
	// CacheDisabled disables the on disk tile cache (so all tiles are fetched
	// from the network).
	//
	// Note: not all sources work with cache disabled.
	CacheDisabled Cache = "none://"
)
