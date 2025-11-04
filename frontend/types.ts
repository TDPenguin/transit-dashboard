// Type definitions for WMATA API responses

// Address struct: describes address fields
export interface Address {
    City: string;
    State: string;
    Street: string;
    Zip: string;
}

// StationInfo struct: detailed information about a station
export interface StationInfo {
    Address: Address;
    Code: string;
    Lat: number;
    Lon: number;
    LineCode1: string;
    LineCode2: string;
    LineCode3: string;
    LineCode4: string;
    Name: string;
    StationTogether1: string;
    StationTogether2: string;
}

// StationEntrance struct: information about a station entrance
export interface StationEntrance {
    Description: string;
    ID: string;
    Lat: number;
    Lon: number;
    Name: string;
    StationCode1: string;
    StationCode2: string;
}

// TrainPrediction struct: rail-time train arrival info
export interface TrainPrediction {
    Car: string;
    Destination: string;
    DestinationCode: string;
    DestinationName: string;
    Group: string;
    Line: string;
    LocationCode: string;
    LocationName: string;
    Min: string; // String because "2","BRD" (Boarding), "ARR" (Arriving).
}
