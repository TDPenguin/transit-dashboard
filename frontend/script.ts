// This is like a Rust struct definition: describes the fields of a Station
interface Station {
    Name: string; // like field: String in Rust struct
    Code: string; // like field: String in Rust struct
}

// This is an async function, like async fn in Rust
async function fetchStations() {
    try {
        // This is like: let res = reqwest::get("http://localhost:8080/stations").await;
        const res = await fetch("http://localhost:8080/stations");

        // !res.ok means the HTTP status was not 2xx, like if !res.status().is_success() in Rust
        // throw is like panic! but can be caught (unlike panic! in Rust, which usually aborts)
        if (!res.ok) throw new Error("Failed to fetch stations");
        // If !res.ok is true, throw new Error(...) runs, and the rest of the code in the try block is skipped.

        // This is like: let stations: Vec<Station> = res.json().await?;
        const stations: Station[] = await res.json(); // Station[] is like Vec<Station> in Rust

        // This is like getting a reference to a DOM node; not directly in Rust, but similar to manipulating a UI tree
        const list = document.getElementById("stations");
        if (!list) return; // If list does not exist, return (like if let Some(list) = ... in Rust, else return)

        list.innerHTML = ""; // Clear any existing items (reset the list)

        // For each station, create a new list item and add it to the DOM
        // This is like for station in stations { ... } in Rust
        stations.forEach(station => {
            const li = document.createElement("li");
            li.textContent = `${station.Name} (${station.Code})`;
            list.appendChild(li);
        });
    } catch (err) {
        // This is like a match or Result error handling in Rust
        console.error("Error fetching stations:", err);
    }
}

// Call it once the page loads
fetchStations(); // Like calling the async fn in main() in Rust with .await