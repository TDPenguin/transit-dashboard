// Main entry point for transit dashboard

import type { StationInfo } from './types.js';
import { getTimestamp, getStationCodes, LINE_COLORS } from './utils.js';
import { fetchTrainPredictions, loadStationLocations, fetchAllStations, fetchEntrances, fetchRailLines } from './api.js';
import { buildTrainPredictionsHTML } from './ui.js';

declare const L: typeof import('leaflet');

// Global state
let map: any;
let markers: { [key: string]: any } = {};
let entranceMarkers: any[] = [];
let allStations: { [key: string]: StationInfo } = {};
let currentStationCode: string | null = null;
let currentStationInfo: StationInfo | null = null;
let predictionRefreshInterval: number | null = null;

// Helper to create circle markers (reduces duplication)
function createCircleMarker(coords: [number, number], color: string, name: string, onClick?: () => void) {
    const marker = L.circleMarker(coords, {
        radius: 5,
        fillColor: color,
        color: color === '#0066cc' ? '#ffffff' : color,
        weight: color === '#0066cc' ? 1 : 2,
        opacity: 1,
        fillOpacity: 0.8
    }).addTo(map);
    marker.bindTooltip(name, { permanent: false, direction: 'top' });
    if (onClick) marker.on('click', onClick);
    return marker;
}

// Helper to get all line codes from station and connected platforms
function getAllLineCodes(stationInfo: StationInfo): string[] {
    const lines: string[] = [stationInfo.LineCode1, stationInfo.LineCode2, stationInfo.LineCode3, stationInfo.LineCode4];
    [stationInfo.StationTogether1, stationInfo.StationTogether2].forEach(code => {
        if (code && allStations[code]) {
            const together = allStations[code];
            lines.push(together.LineCode1, together.LineCode2, together.LineCode3, together.LineCode4);
        }
    });
    return [...new Set(lines.filter(line => line && line.trim() !== ""))];
}

// Start prediction refresh timer for current station
function startPredictionRefresh() {
    // Clear any existing interval first
    if (predictionRefreshInterval !== null) {
        clearInterval(predictionRefreshInterval);
        predictionRefreshInterval = null;
    }
    
    // Only start if we have a station selected
    if (!currentStationCode || !currentStationInfo) return;
    
    predictionRefreshInterval = window.setInterval(async () => {
        if (currentStationCode && currentStationInfo) {
            console.log(`[${getTimestamp()}] [Auto-Refresh] Updating predictions for ${currentStationInfo.Name}`);
            await updatePredictions(currentStationCode, currentStationInfo);
        }
    }, 10000); // 10 seconds
}

// Stop prediction refresh timer
function stopPredictionRefresh() {
    if (predictionRefreshInterval !== null) {
        clearInterval(predictionRefreshInterval);
        predictionRefreshInterval = null;
    }
}

// Initialize the map centered on Metro Center
function initMap() {
    map = L.map('map').setView([38.898303, -77.028099], 11); // 11 is the zoom level
    L.tileLayer('https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>',
        subdomains: 'abcd',
        maxZoom: 20
    }).addTo(map);
}

// Fetch all stations and display them
async function fetchStations() {
    try {
        console.log(`[${getTimestamp()}] [Stations] Loading...`);
        const stationCoords = await loadStationLocations();
        const stations = await fetchAllStations();
        console.log(`[${getTimestamp()}] [Stations] Loaded ${stations.length} stations`);

        const list = document.getElementById("stations");
        if (!list) return;
        list.innerHTML = ""; // Clear any existing items (reset the list)

        // Store all stations in global dictionary for lookup
        for (const station of stations) allStations[station.Code] = station;

        // Deduplicate stations by name (some stations appear twice as different platforms)
        const uniqueStations = new Map<string, StationInfo>();
        for (const station of stations) {
            if (!uniqueStations.has(station.Name)) uniqueStations.set(station.Name, station);
        }

        // Sort stations alphabetically by name
        const sortedStations = Array.from(uniqueStations.values()).sort((a, b) => a.Name.localeCompare(b.Name));

        // For each unique station, add it to both the list and the map
        for (const station of sortedStations) {
            // Add to list
            const li = document.createElement("li");
            li.textContent = station.Name;
            
            // Make the station clickable - When clicked it will display detailed info
            li.style.cursor = "pointer";
            li.style.color = "blue";
            li.style.textDecoration = "underline";
            li.onclick = () => selectStation(station.Code, station);
            list.appendChild(li);

            // Use GeoJSON coordinates if available, fallback to API coordinates
            const coords = stationCoords.get(station.Name) || [station.Lat, station.Lon];
            if (coords[0] && coords[1]) {
                markers[station.Code] = createCircleMarker(coords, '#0066cc', station.Name, () => selectStation(station.Code, station));
            }
        }
    } catch (err) {
        console.error("Error fetching stations:", err);
    }
}

// Select a station and display its details
async function selectStation(stationCode: string, stationInfo: StationInfo) {
    console.log(`[${getTimestamp()}] [Station] Selected: ${stationInfo.Name} (${stationCode})`);
    
    // Store current station for auto-refresh
    currentStationCode = stationCode;
    currentStationInfo = stationInfo;
    
    // Clear previous entrance markers from the map
    entranceMarkers.forEach(marker => map.removeLayer(marker));
    entranceMarkers = [];

    const detailsDiv = document.getElementById("station-details-content");
    if (!detailsDiv) return;

    const closeBtn = document.getElementById("close-details");
    if (closeBtn) closeBtn.style.display = "block";

    // Get all unique line codes from this station and connected platforms
    const uniqueLines = getAllLineCodes(stationInfo);
    
    // Fetch entrances for this station
    let entrancesHTML = "";
    try {
        const entrances = await fetchEntrances(stationCode);
        console.log(`[${getTimestamp()}] [Entrances] Found ${entrances.length} entrance(s)`);
        if (entrances.length > 0) {
            entrancesHTML = `
                <p><strong>Entrances (${entrances.length}):</strong></p>
                <ul style="margin: 0; padding-left: 20px;">
                    ${entrances.map(e => `<li>${e.Name}</li>`).join('')}
                </ul>
            `;

            // Add entrance markers to the map
            entrances.forEach(entrance => {
                if (entrance.Lat && entrance.Lon) {
                    entranceMarkers.push(createCircleMarker([entrance.Lat, entrance.Lon], '#ff7800', entrance.Name));
                }
            });
        }
    } catch (err) {
        console.error("Error fetching entrances:", err);
    }
    
    // Fetch train predictions
    const predictions = await fetchTrainPredictions(getStationCodes(stationInfo));
    console.log(`[${getTimestamp()}] [Predictions] Fetch completed`);
    const predictionsHTML = buildTrainPredictionsHTML(predictions);

    // Build line badges HTML
    const lineBadgesHTML = uniqueLines.map(line => 
        `<img src="assets/${line}.svg" alt="${line}" style="height: 20px; width: 20px; margin-right: 4px;">`
    ).join('');

    detailsDiv.innerHTML = `
        <h2 style="display: flex; align-items: center; gap: 8px;">
            ${stationInfo.Name}
            <span style="display: flex; align-items: center;">${lineBadgesHTML}</span>
        </h2>
        ${predictionsHTML}
        <p><strong>Address:</strong> ${stationInfo.Address.Street}, ${stationInfo.Address.City}, ${stationInfo.Address.State} ${stationInfo.Address.Zip}</p>
        ${entrancesHTML}
    `;

    // Pan the map to the selected station
    if (stationInfo.Lat && stationInfo.Lon) {
        // Only zoom in if we're currently zoomed out, don't zoom out if already close
        const currentZoom = map.getZoom();
        const targetZoom = Math.max(currentZoom, 16);
        map.setView([stationInfo.Lat, stationInfo.Lon], targetZoom, { animate: false });
        
        // Highlight the marker (open its popup)
        const marker = markers[stationInfo.Code];
        if (marker) marker.openPopup();
    }
    console.log(`[${getTimestamp()}] [Station] Selection complete`);
    
    // Start auto-refresh timer for this station
    startPredictionRefresh();
}

// Update just the predictions for the current station (used for auto-refresh)
async function updatePredictions(stationCode: string, stationInfo: StationInfo) {
    const predictions = await fetchTrainPredictions(getStationCodes(stationInfo));
    const predictionsHTML = buildTrainPredictionsHTML(predictions);
    
    const detailsDiv = document.getElementById("station-details-content");
    if (!detailsDiv) return;
    
    // Find the "Next Trains:" section and replace it
    const parser = new DOMParser();
    const newDoc = parser.parseFromString(detailsDiv.innerHTML, 'text/html');
    const predSection = newDoc.querySelector('div[style*="margin-top: 16px"]');
    
    if (predSection) {
        const tempDiv = document.createElement('div');
        tempDiv.innerHTML = predictionsHTML;
        const newPredSection = tempDiv.firstElementChild;
        if (newPredSection) {
            predSection.replaceWith(newPredSection);
            detailsDiv.innerHTML = newDoc.body.innerHTML;
        }
    }
}

// Fetch and display rail lines on the map
async function loadRailLines() {
    try {
        console.log(`[${getTimestamp()}] [Lines] Loading rail lines...`);
        const geojson = await fetchRailLines();
        console.log(`[${getTimestamp()}] [Lines] Loaded ${geojson.features.length} rail line(s)`);
        
        // Adds GeoJSON to map with custom styling
        L.geoJSON(geojson, {
            style: (feature) => {
                // Get the line name from the GeoJSON properties
                const lineName = (feature?.properties?.LINE || "").toLowerCase();
                // Pick a color based on the line name
                const color = Object.entries(LINE_COLORS).find(([key]) => lineName.includes(key))?.[1] || '#888';
                return { color, weight: 3, opacity: 0.8 };
            }
        }).addTo(map);
    } catch (err) {
        console.error("Error loading rail lines:", err);
    }
}

// Close the station details panel
function closeStationDetails() {
    // Stop auto-refresh timer
    stopPredictionRefresh();
    
    // Clear current station tracking
    currentStationCode = null;
    currentStationInfo = null;
    
    const detailsDiv = document.getElementById("station-details-content");
    if (detailsDiv) detailsDiv.innerHTML = '<p>Click on a station in the list or on the map to view details</p>';
    
    const closeBtn = document.getElementById("close-details");
    if (closeBtn) closeBtn.style.display = "none";
    
    // Clear entrance markers from the map
    entranceMarkers.forEach(marker => map.removeLayer(marker));
    entranceMarkers = [];
    map.closePopup();
}

// Initialize everything when the page loads
initMap();
loadRailLines();
fetchStations();

// Wire up close button
const closeBtn = document.getElementById("close-details");
if (closeBtn) closeBtn.addEventListener("click", closeStationDetails);
