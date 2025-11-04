package main

// Station struct: like a struct in Rust, defines fields and their types
type Station struct {
	Name string `json:"Name"` // field maps to "Name" in JSON
	Code string `json:"Code"` // field maps to "Code" in JSON
}

// StationsResponse struct: holds a slice (dynamic array, like Vec in Rust) of Station
type StationsResponse struct {
	Stations []Station `json:"Stations"` // maps to "Stations" in JSON
}

// Address struct: holds address information for a station
type Address struct {
	City   string `json:"City"`
	State  string `json:"State"`
	Street string `json:"Street"`
	Zip    string `json:"Zip"`
}

// StationInfo struct: Detailed info about a single station
type StationInfo struct {
	Address          Address `json:"Address"`
	Code             string  `json:"Code"`
	Lat              float64 `json:"Lat"`
	LineCode1        string  `json:"LineCode1"`
	LineCode2        string  `json:"LineCode2"`
	LineCode3        string  `json:"LineCode3"`
	LineCode4        string  `json:"LineCode4"`
	Lon              float64 `json:"Lon"`
	Name             string  `json:"Name"`
	StationTogether1 string  `json:"StationTogether1"`
	StationTogether2 string  `json:"StationTogether2"`
}

// StationEntrance struct: Info about a station entrance
type StationEntrance struct {
	Description  string  `json:"Description"`
	ID           string  `json:"ID"`
	Lat          float64 `json:"Lat"`
	Lon          float64 `json:"Lon"`
	Name         string  `json:"Name"`
	StationCode1 string  `json:"StationCode1"`
	StationCode2 string  `json:"StationCode2"`
}

// EntrancesResponse struct: Holds all station entrances
type EntrancesResponse struct {
	Entrances []StationEntrance `json:"Entrances"`
}

// TrainPrediction struct: Info about train predictions
type TrainPrediction struct {
	Car             string `json:"Car"`
	Destination     string `json:"Destination"`
	DestinationCode string `json:"DestinationCode"`
	DestinationName string `json:"DestinationName"`
	Group           string `json:"Group"`
	Line            string `json:"Line"`
	LocationCode    string `json:"LocationCode"`
	LocationName    string `json:"LocationName"`
	Min             string `json:"Min"`
}

// TrainPredictionsResponse struct: Holds all train prediction responses
type TrainPredictionsResponse struct {
	Trains []TrainPrediction `json:"Trains"`
}

// Lines struct: Info about lines
type Lines struct {
	DisplayName          string `json:"DisplayName"`
	EndStationCode       string `json:"EndStationCode"`
	InternalDestination1 string `json:"InternalDestination1"`
	InternalDestination2 string `json:"InternalDestination2"`
	LineCode             string `json:"LineCode"`
	StartStationCode     string `json:"StartStationCode"`
}

// LinesResponse struct: Holds all lines responses
type LinesResponse struct {
	Lines []Lines `json:"Lines"`
}

/*
Parking information will not be used for now as this project focuses on accessibility and walkability.
It might be used in the future for visualisations on the site, or to see how "car-dependant" a station is.

Kept here for futureproofing.

Bike parking would have been nice but this is manually obtained, not via API.
*/

// AllDayParking struct: Info about all-day parking options
type AllDayParking struct {
	TotalCount   int      `json:"TotalCount"`
	RiderCost    *float64 `json:"RiderCost"`    // Nullable, a pointer type like *float64 or *string means the variable can either:
	NonRiderCost *float64 `json:"NonRiderCost"` // point to a value like a real float64 or string, or be "nil", or null.
}

// ShortTermParking struct: Info about short-term parking options
type ShortTermParking struct {
	SaturdayRiderCost    *float64 `json:"SaturdayRiderCost"`
	SaturdayNonRiderCost *float64 `json:"SaturdayNonRiderCost"`
	TotalCount           int      `json:"TotalCount"`
	Notes                *string  `json:"Notes"`
}

// StationParking struct: Info about parking at a station
type StationParking struct {
	Code             string           `json:"Code"`
	Notes            *string          `json:"Notes"`
	AllDayParking    AllDayParking    `json:"AllDayParking"`
	ShortTermParking ShortTermParking `json:"ShortTermParking"`
}

// StationsParkingResponse struct: Holds all station parking responses
type StationsParkingResponse struct {
	StationsParking []StationParking `json:"StationsParking"`
}
