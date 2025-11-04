// API functions for fetching data from backend

import type { TrainPrediction, StationInfo, StationEntrance } from './types.js';
import { getTimestamp } from './utils.js';

// Fetch real-time train predictions for a specific station (or multiple connected platforms)
export async function fetchTrainPredictions(stationCodes: string[]): Promise<TrainPrediction[]> {
    try {
        const res = await fetch("http://localhost:8080/nexttrains");
        if (!res.ok) throw new Error("Failed to fetch train predictions");
        const allTrains: TrainPrediction[] = await res.json();
        
        // Filter for all provided station codes and sort by arrival time
        const filtered = allTrains
            .filter(train => stationCodes.includes(train.LocationCode))
            .sort((a, b) => {
                // Sort: ARR/BRD first, then by numeric time
                const getOrder = (min: string) => {
                    if (min === "BRD") return 0;
                    if (min === "ARR") return 1;
                    return parseInt(min) || 999;
                };
                return getOrder(a.Min) - getOrder(b.Min);
            });
        
        console.log(`[${getTimestamp()}] [Predictions] ${filtered.length} train(s) for ${stationCodes.join(', ')}`);
        return filtered;
    } catch (err) {
        console.error("Error fetching train predictions:", err);
        return [];
    }
}

// Load station locations from GeoJSON (more accurate than API coordinates, API as fallback)
export async function loadStationLocations(): Promise<Map<string, [number, number]>> {
    const coordsMap = new Map<string, [number, number]>();
    try {
        const res = await fetch("http://localhost:8080/geojson/stations");
        if (!res.ok) throw new Error("Failed to fetch station GeoJSON");
        const geojson = await res.json();
        
        // Extract coordinates from GeoJSON features
        for (const feature of geojson.features) {
            const name = feature.properties?.NAME;
            if (name && feature.geometry?.coordinates) {
                const [lon, lat] = feature.geometry.coordinates;
                coordsMap.set(name, [lat, lon]);
            }
        }
        console.log(`[${getTimestamp()}] [GeoJSON] Loaded ${coordsMap.size} station coordinates`);
    } catch (err) {
        console.error("Error loading station GeoJSON:", err);
    }
    return coordsMap;
}

// Fetch all stations from backend
export async function fetchAllStations(): Promise<StationInfo[]> {
    const res = await fetch("http://localhost:8080/stations");
    if (!res.ok) throw new Error("Failed to fetch stations");
    return await res.json();
}

// Fetch entrances for a specific station
export async function fetchEntrances(stationCode: string): Promise<StationEntrance[]> {
    const res = await fetch(`http://localhost:8080/entrances?code=${stationCode}`);
    if (!res.ok) return [];
    return await res.json();
}

// Fetch rail lines GeoJSON
export async function fetchRailLines(): Promise<any> {
    const res = await fetch("http://localhost:8080/geojson/lines");
    if (!res.ok) throw new Error("Failed to fetch rail lines");
    return await res.json();
}
