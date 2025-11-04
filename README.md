# WMATA Transit Dashboard

A web dashboard for exploring Washington Metro stations and entrances on an interactive map. Built with Go, TypeScript, and Leaflet.js as a learning project.

## Features
- Interactive map showing all WMATA rail stations
- Station details with line information and addresses
- Station entrance locations displayed as markers
- Server-side caching (24-hour refresh)
- Sequential API fetching (no rate limiting issues)

## Stack
- **Backend**: Go (net/http, server-side caching with mutex)
- **Frontend**: TypeScript, Leaflet.js
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
go run main.go
```

**Terminal 2 - Frontend:**
```bash
cd frontend
npm run build
npx serve .
```

Open http://localhost:3000

## Project Structure
```
backend/
  ├── main.go           # Entry point, server setup
  ├── types.go          # All structs for data
  ├── cache.go          # Caching logic
  ├── handlers.go       # HTTP handlers
  └── .env              # API key (gitignored)

frontend/
  ├── index.html        # Page structure
  ├── script.ts         # TypeScript logic
  ├── style.css         # WMATA-themed styling
  ├── tsconfig.json     # TypeScript config
  └── package.json      # Dependencies
```

## How It Works

### Caching Strategy
- Pre-warms cache on server startup
- Fetches all 91+ stations sequentially (avoids rate limits)
- Cache duration: 24 hours
- Shared across all users (server-side)

### API Endpoints
- `GET /stations` - Returns all station details (cached)
- `GET /entrances?code=XXX` - Returns entrances for a station (cached)

### Frontend Flow
1. Fetch all stations on page load
2. Display markers on map + list in sidebar
3. Click station → fetch entrances → show orange circle markers
4. Pan to station (keeps zoom if already close)

## TODO
- [ ] WMATA brand colors/styling
- [ ] Live train arrivals
- [ ] Bus stop data
- [ ] Mobile responsive design
- [ ] Production deployment

## License
MIT
