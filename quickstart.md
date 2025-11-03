# Quick Start Guide

## 1. Clean Previous Builds (optional)
Removes the `/dist` folder:

```sh
npm run clean
```

## 2. Build the Frontend
Compiles TypeScript to JavaScript in `/dist`:

```sh
npm run build
```

## 3. Serve the Frontend
Start a local server in the `frontend` folder (choose one):

**Using npx:**
```sh
npx serve .
```
**Or with Python:**
```sh
python -m http.server 3000
```

Then open [http://localhost:3000](http://localhost:3000) in your browser.

## 4. Start the Backend
In a separate terminal, from the `backend` folder:

```sh
go run main.go
```

## 5. View the App
- Open your browser to [http://localhost:3000](http://localhost:3000)
- You should see the WMATA stations list.

---

## Troubleshooting

- **No stations?**  
  Make sure the Go backend is running and your frontend is built.
- **Build errors?**  
  Check your TypeScript files and `tsconfig.json`.
- **Script errors?**  
  Open the browser console (F12) for details.

---

## Common Commands

- Clean: `npm run clean`
- Build: `npm run build`
- Serve: `npx serve .`
- Backend: `go run main.go`
- 
#### Quick Commands
Go: `cd backend; cls; go run main.go`

TypeScript: `cd frontend; cls; npm run build; npx serve .`

