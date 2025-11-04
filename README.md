# WMATA Transit Dashboard

A real-time web dashboard for exploring Washington Metro stations with live train predictions. Built with Go, TypeScript, and Leaflet.js as a learning project.

## Features
- **Live train predictions** with auto-refresh (10-second updates)
- **Interactive map** showing all WMATA rail stations and entrances
- **Two-column track display** with platform-grouped, scrollable predictions
- **Station details** with line badges, addresses, and entrance locations

## Stack
- **Backend**: Go (net/http, mutex-protected caching, GeoJSON support)
- **Frontend**: TypeScript (ES modules), Leaflet.js, CSS Grid
- **Data**: WMATA Rail API

## Setup

### Prerequisites
- Go 1.x
- Node.js & npm
- WMATA API key ([get one here](https://developer.wmata.com/))

### Installation

1. Clone the repository
```bash
git clone https://github.com/TDPenguin/transit-dashboard.git
cd transit-dashboard
```

2. Configure backend
```bash
cd backend
# Create .env file
echo "WMATA_API_KEY=your_key_here" > .env
```

3. Install frontend dependencies
```bash
cd ../frontend
npm install
```

### Running

**Terminal 1 - Backend:**
```bash
cd backend
go run .
```

**Terminal 2 - Frontend (optional, for development):**
```bash
cd frontend
npm run watch
```

Open http://localhost:8080

## Project Structure
```
backend/
  ├── main.go           # Entry point, server setup
  ├── types.go          # All structs for API data
  ├── cache.go          # Caching with auto-refresh
  ├── handlers.go       # HTTP handlers & CORS
  └── .env              # API key (gitignored)

frontend/
  ├── index.html        # Page structure
  ├── style.css         # WMATA-themed styling
  ├── types.ts          # TypeScript interfaces
  ├── utils.ts          # Constants and helpers
  ├── api.ts            # Backend API calls
  ├── ui.ts             # HTML rendering
  ├── script.ts         # Main application logic
  ├── tsconfig.json     # TypeScript config
  └── package.json      # Dependencies
```

## License
MIT
