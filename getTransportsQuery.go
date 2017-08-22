package main

// Template of the query made the the HSL server
// lat, lon and maxDistance need to be change accordingly
type queryStruct struct {
	Nearest struct {
		Edges []struct {
			Node struct {
				Stop struct {
					GtfsID                   string
					Name                     string
					Lat                      float64
					Lon                      float64
					StoptimesWithoutPatterns []struct {
						RealtimeDeparture int
						Trip              struct {
							TripHeadsign string
							Route        struct {
								LongName  string
								ShortName string
								Type      int
							}
						}
					}
				}
			}
		}
	} `graphql:"stopsByRadius(lat: $lat, lon: $lon, radius: $radius)"`
}
