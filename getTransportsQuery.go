package main

// Structure receiving the result of the query
type answerType struct {
	Data struct {
		Nearest struct {
			Edges []struct {
				Node struct {
					Distance int `json:"distance"`
					Place    struct {
						ID        string  `json:"ID"`
						Name      string  `json:"name"`
						Latitude  float64 `json:"lat"`
						Longitude float64 `json:"lon"`
						Times     []struct {
							Date int `json:"realtimeDeparture"`
							Trip struct {
								Direction string `json:"tripHeadsign"`
								Route     struct {
									Name string `json:"shortname"`
									Mode string `json:"mode"`
								}
							}
						} `json:"stoptimesWithoutPatterns"`
					} `json:"place"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"nearest"`
	} `json:"data"`
}

// Template of the query made the the HSL server
// lat, lon and maxDistance need to be change accordingly
const queryTemplate = `{
	nearest(
		lat: %v,
		lon: %v,
		maxDistance: %v,
		filterByPlaceTypes: STOP
	) {
		edges {
			node {
				distance
					place {
						lat
						lon
						__typename
						...on Stop {
							name
							stoptimesWithoutPatterns {
								realtimeDeparture
								trip {
									tripHeadsign
									route {
										longName
										shortName
										mode
									}
								}
							}
						}
					}
				}
			}
		}
	}
}`
