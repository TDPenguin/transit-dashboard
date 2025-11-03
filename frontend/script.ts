declare const L: typeof import('leaflet');

// Station struct: describes basic station info
interface Station {
    Name: string; // like field: String in Rust struct
    Code: string; // like field: String in Rust struct
}

// Address struct: describes address fields
interface Address {
    City: string;
    State: string;
    Street: string;
    Zip: string;
}

// StationInfo struct: detailed information about a station
interface StationInfo {
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
interface StationEntrance {
    Description: string;
    ID: string;
    Lat: number;
    Lon: number;
    Name: string;
    StationCode1: string;
    StationCode2: string;
}

// Global variable to store map and markers
let map: any;
let markers: { [key: string]: any } = {}; // Dictionary mapping station codes to markers
let entranceMarkers: any[] = []; // Array to store entrance markers (cleared when selecting new station)

// Initialize the map centered on Metro Center
function initMap() {
    // Center on Metro Center station
    map = L.map('map').setView([38.898303, -77.028099], 11); // 11 is the zoom.

    // Add OpenStreetMap tiles
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© OpenStreetMap contributors',
        maxZoom: 19
    }).addTo(map);
}

// Fetch all stations and display them
async function fetchStations() {
    try {
        const res = await fetch("http://localhost:8080/stations");

        // !res.ok means the HTTP status was not 2xx, like if !res.status().is_success() in Rust
        // throw is like panic! but can be caught (unlike panic! in Rust, which usually aborts)
        if (!res.ok) throw new Error("Failed to fetch stations");
        // If !res.ok is true, throw new Error(...) runs, and the rest of the code in the try block is skipped.

        const stations: StationInfo[] = await res.json(); // Get full StationInfo

        const list = document.getElementById("stations");
        if (!list) return; // If list does not exist, return (like if let Some(list) = ... in Rust, else return)

        list.innerHTML = ""; // Clear any existing items (reset the list)

        // For each station, add it to both the list and the map
        // No more slow individual API calls!
        for (const station of stations) {
            // Add to list
            const li = document.createElement("li");
            li.textContent = `${station.Name} (${station.Code})`;

            // Make the station clickable
            // When clicked it will display detailed info
            li.style.cursor = "pointer";
            li.style.color = "blue";
            li.style.textDecoration = "underline";

            li.onclick = () => selectStation(station.Code, station);

            list.appendChild(li);

            // Add marker to map - we have coords
            if (station.Lat && station.Lon) {
                const marker = L.marker([station.Lat, station.Lon]).addTo(map);
                marker.bindTooltip(station.Name, { permanent: false, direction: 'top' });
                
                // When you click the marker, show station details
                marker.on('click', () => selectStation(station.Code, station));
                
                // Store marker for later use
                markers[station.Code] = marker;
            }
        }
    } catch (err) {
        console.error("Error fetching stations:", err);
    }
}

// Select a station and display its details
// Always have the info from the initial fetch, so no need for another API call
async function selectStation(stationCode: string, stationInfo: StationInfo) {
    // Clear previous entrance markers from the map
    entranceMarkers.forEach(marker => map.removeLayer(marker));
    entranceMarkers = [];

    // Display the info in the details div
    const detailsDiv = document.getElementById("station-details");
    if (!detailsDiv) return;

    // Build a list of line codes (filter out empty ones)
    const lines = [stationInfo.LineCode1, stationInfo.LineCode2, stationInfo.LineCode3, stationInfo.LineCode4]
        .filter(line => line !== "");
    
    // Fetch entrances for this station
    let entrancesHTML = "";
    try {
        const res = await fetch(`http://localhost:8080/entrances?code=${stationCode}`);
        if (res.ok) {
            const entrances: StationEntrance[] = await res.json();
            if (entrances.length > 0) {
                entrancesHTML = `
                    <p><strong>Entrances (${entrances.length}):</strong></p>
                    <ul style="margin: 0; padding-left: 20px;">
                        ${entrances.map(e => `<li>${e.Name}</li>`).join('')}
                    </ul>
                `;

                // Add entrance markers to the map as small circle markers
                entrances.forEach(entrance => {
                    if (entrance.Lat && entrance.Lon) {
                        const entranceMarker = L.circleMarker([entrance.Lat, entrance.Lon], {
                            radius: 5,           // Small circle
                            color: '#ff7800',    // Orange border
                            fillColor: '#ff7800', // Orange fill
                            fillOpacity: 0.8,
                            weight: 2
                        }).addTo(map);

                        // Add tooltip showing entrance name
                        entranceMarker.bindTooltip(entrance.Name, { 
                            permanent: false, 
                            direction: 'top' 
                        });

                        // Store marker so we can remove it later
                        entranceMarkers.push(entranceMarker);
                    }
                });
            }
        }
    } catch (err) {
        console.error("Error fetching entrances:", err);
    }
    
    detailsDiv.innerHTML = `
        <h2>${stationInfo.Name}</h2>
        <p><strong>Station Code:</strong> ${stationInfo.Code}</p>
        <p><strong>Lines:</strong> ${lines.join(", ")}</p>
        <p><strong>Address:</strong> ${stationInfo.Address.Street}, ${stationInfo.Address.City}, ${stationInfo.Address.State} ${stationInfo.Address.Zip}</p>
        <p><strong>Coordinates:</strong> ${stationInfo.Lat}, ${stationInfo.Lon}</p>
        ${stationInfo.StationTogether1 ? `<p><strong>Connected Platform:</strong> ${stationInfo.StationTogether1}</p>` : ""}
        ${entrancesHTML}
    `;

    // Pan the map to the selected station
    if (stationInfo.Lat && stationInfo.Lon) {
        // Only zoom in if we're currently zoomed out; don't zoom out if already close
        const currentZoom = map.getZoom();
        const targetZoom = Math.max(currentZoom, 16); // Use current zoom if already closer than 15
        
        map.setView([stationInfo.Lat, stationInfo.Lon], targetZoom);
        
        // Highlight the marker (open its popup)
        const marker = markers[stationInfo.Code];
        if (marker) {
            marker.openPopup();
        }
    }
}

// Initialize everything when the page loads
initMap();
fetchStations();