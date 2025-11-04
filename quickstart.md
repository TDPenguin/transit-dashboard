# Quick Start Guide

## Development (Two Terminals)

**Terminal 1:** TypeScript watch (auto-compile on save)
```sh
cd frontend
npm run watch
```

**Terminal 2:** Go server (serves API + frontend + predictions)
```sh
cd backend
go run .
```

Then open **http://localhost:8080** in your browser.

Edit any `.ts` file → Save → Refresh browser (F5) → See changes immediately.

---

## Quick Run (One Terminal)

Build TypeScript once, then run Go:
```sh
cd frontend; npm run build; cd ../backend; go run .
```

Open **http://localhost:8080** in your browser.

When you edit any `.ts` file, rebuild with:
```sh
cd frontend; npm run build
```
Then refresh browser.

---

## File Guide

**Backend (Go):**
- `main.go` - Server setup, routes, CORS
- `types.go` - Data structures for WMATA API
- `cache.go` - Caching with auto-refresh timers
- `handlers.go` - HTTP endpoint handlers

**Frontend (TypeScript):**
- `types.ts` - TypeScript interfaces (matches Go types)
- `utils.ts` - Constants (colors, expansions) and helpers
- `api.ts` - All backend API calls
- `ui.ts` - HTML rendering for predictions
- `script.ts` - Main app logic, map, event handlers

---

## Troubleshooting

- **No stations showing?**  
  Backend still loading. Check terminal for "Cache pre-warmed successfully!"
  
- **No predictions showing?**  
  Backend refreshing data. Wait ~20 seconds and click station again.
  
- **TypeScript errors?**  
  Check `tsconfig.json` and run `npm run build` to see errors.
  
- **Map not showing?**  
  Open browser console (F12) for errors. Check Leaflet.js loaded.
  
- **Changes not showing?**  
  Make sure TypeScript compiled (check terminal 1) and you refreshed browser (F5).

- **Predictions not auto-updating?**  
  Timer only runs when station is selected. Close and reopen station panel.

---

## Useful Commands

- `npm run build` - Compile all TypeScript files once
- `npm run watch` - Auto-compile on save (dev mode)
- `npm run clean` - Delete all compiled `.js` files
- `go run .` - Start Go server (from backend folder)

