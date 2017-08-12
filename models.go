package main

//
// import "math"
//
// //===================AGENCY===================
// type agency struct {
// 	ID          string   `json:"id"`           // ID of the region (Country.City.Agency)
// 	Name        string   `json:"name"`         // Displayed name of the Agency
// 	URL         string   `json:"url"`          // The URL to the agency's website/app...
// 	Git         string   `json:"git"`          // The URL to the git repo
// 	Center      position `json:"center"`       // Center of the Agency
// 	Radius      float64  `json:"radius"`       // Radius of the Agency in meters
// 	Types       []int    `json:"types"`        // The type of transports handled by the agency
// 	TypesString []string `json:"types_string"` // Name for the type of transports
// }
//
// //===================POSITION===================
// type position struct {
// 	Latitude  float64 `json:"latitude"`  // Nord-South position
// 	Longitude float64 `json:"longitude"` // East-West position
// }
//
// // DistanceFrom - Compute the distance between two positions in meters
// // See http://www.movable-type.co.uk/scripts/latlong.html
// // @param pos2: the position
// // @return the distance in meters
// func (pos *position) DistanceFrom(pos2 *position) float64 {
// 	// Get radian diff between latitude and longitude the position
// 	dLat := toRadians(pos2.Latitude - pos.Latitude)
// 	dLon := toRadians(pos2.Longitude - pos.Longitude)
// 	// Some complexe computations with sin and cos
// 	a := math.Sin(dLat / 2)
// 	b := math.Sin(dLon / 2)
// 	c := math.Cos(toRadians(pos.Latitude)) * math.Cos(toRadians(pos2.Latitude))
// 	d := (a * a) + c*(b*b)
// 	// 6371000 => Average earth radius in meters
// 	return 6371000 * 2 * math.Atan2(math.Sqrt(d), math.Sqrt(1-d))
// }
//
// // Convert a degres angle to radian angle
// // Exemple: 180° ==> π rad
// func toRadians(degre float64) float64 {
// 	return degre * math.Pi / 180
// }
//
// //===================TRANSPORT===================
// type transport struct {
// 	ID       string     `json:"ID"`       // ID of the Transport, should be specific to the Agency
// 	AgencyID string     `json:"agencyID"` // ID of the associated agency
// 	Type     int        `json:"type"`     // String identifing the kind of transport
// 	Image    string     `json:"image"`    // The image to display that represent the transport
// 	Name     string     `json:"name"`     // The name of the transport, doesn't have to be unique
// 	Line     string     `json:"line"`     // The group of the transport
// 	Position position   `json:"position"` // Position of the transport
// 	Passages []*passage `json:"passages"` // Next passage for public transports
// }
//
// type passage struct {
// 	Direction string   `json:"direction"` // Direction of the passage
// 	Times     []string `json:"times"`     // Time, is array of string to support non numeric values
// }
//
// // DistanceFrom - Compute the distance between the transport position and a given position
// // @param pos: the position
// // @return the distance in meters
// func (t *transport) DistanceFrom(pos *position) float64 {
// 	return t.Position.DistanceFrom(pos)
// }
//
// //===================TRANSPORT TYPES===================
// const (
// 	tram = iota
// 	metro
// 	rail
// 	bus
// 	ferry
// 	unknown
// )
//
// const (
// 	tramString    = "Tram"
// 	metroString   = "Metro"
// 	railString    = "Train"
// 	busString     = "Bus"
// 	ferryString   = "Ferry"
// 	unknownString = "unknown"
// )
