// This is like a Rust struct definition: describes the fields of a Station

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

// Fetch all stations and display them
async function fetchStations() {
    try {
        const res = await fetch("http://localhost:8080/stations");

        // !res.ok means the HTTP status was not 2xx, like if !res.status().is_success() in Rust
        // throw is like panic! but can be caught (unlike panic! in Rust, which usually aborts)
        if (!res.ok) throw new Error("Failed to fetch stations");
        // If !res.ok is true, throw new Error(...) runs, and the rest of the code in the try block is skipped.

        const stations: Station[] = await res.json(); // Station[] is like Vec<Station> in Rust

        const list = document.getElementById("stations");
        if (!list) return; // If list does not exist, return (like if let Some(list) = ... in Rust, else return)

        list.innerHTML = ""; // Clear any existing items (reset the list)

        // For each station, create a clickable list item
        stations.forEach(station => {
            const li = document.createElement("li");
            li.textContent = `${station.Name} (${station.Code})`;

            // Make the station clickable
            // When clicked it will fetch and display detailed info
            li.style.cursor = "pointer";
            li.style.color = "blue";
            li.style.textDecoration = "underline";

            li.onclick = () => fetchStationInfo(station.Code);

            list.appendChild(li);
        });
    } catch (err) {
        // This is like a match or Result error handling in Rust
        console.error("Error fetching stations:", err);
    }
}

// Fetch detailed info for a specific station
async function fetchStationInfo(stationCode: string) {
    try {
        // Call our backend with the station code as a query parameter
        const res = await fetch(`http://localhost:8080/station-info?code=${stationCode}`);
        if (!res.ok) throw new Error("Failed to fetch station info");

        const info: StationInfo = await res.json();

        // Display the info in the details div
        const detailsDiv = document.getElementById("station-details");
        if (!detailsDiv) return;

        // Build a list of line codes (filter out empty ones)
        const lines = [info.LineCode1, info.LineCode2, info.LineCode3, info.LineCode4]
            .filter(line => line !== ""); // Remove empty strings
        
        detailsDiv.innerHTML = `
            <h2>${info.Name}</h2>
            <p><strong>Station Code:</strong> ${info.Code}</p>
            <p><strong>Lines:</strong> ${lines.join(", ")}</p>
            <p><strong>Address:</strong> ${info.Address.Street}, ${info.Address.City}, ${info.Address.State} ${info.Address.Zip}</p>
            <p><strong>Coordinates:</strong> ${info.Lat}, ${info.Lon}</p>
            ${info.StationTogether1 ? `<p><strong>Connected Platform:</strong> ${info.StationTogether1}</p>` : ""}
        `;
    } catch (err) {
        console.error("Error fetching station info:", err);
    }
}

// Call fetchStations when the page loads
fetchStations();