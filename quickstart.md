# Quick Start Guide

## Development (Two Terminals)

**Terminal 1:** TypeScript watch (auto-compile on save)
```sh
cd frontend
npm run watch
```

**Terminal 2:** Go server (serves both API and frontend)
```sh
cd backend
go run .
```

Then open **http://localhost:8080** in your browser.

Edit `script.ts` → Save → Refresh browser (F5) → See changes.

---

## Quick Run (One Terminal)

Build TypeScript once, then run Go:
```sh
cd frontend; npm run build; cd ../backend; go run .
```

Open **http://localhost:8080** in your browser.

When you edit `script.ts`, rebuild with:
```sh
cd frontend; npm run build
```
Then refresh browser.

---

## Troubleshooting

- **No stations showing?**  
  Backend still loading. Check terminal for "Cache pre-warmed successfully!"
  
- **TypeScript errors?**  
  Check `tsconfig.json` and run `npm run build` to see errors.
  
- **Map not showing?**  
  Open browser console (F12) for errors.
  
- **Changes not showing?**  
  Make sure TypeScript compiled (check terminal 1) and you refreshed browser (F5).

---

## Useful Commands

- `npm run build` - Compile TypeScript once
- `npm run watch` - Auto-compile on save
- `npm run clean` - Delete compiled files
- `go run .` - Start Go server (from backend folder)

