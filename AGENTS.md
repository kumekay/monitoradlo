# Monitoradlo

A GUI application for editing Kanshi monitor layout profiles with live preview via Niri IPC.

## Tech Stack

- **Backend**: Go + Wails v2
- **Frontend**: Svelte + TypeScript + Vite
- **Canvas**: SVG-based monitor layout editor
- **System deps**: webkit2gtk-4.1

## Project Structure

```
monitoradlo/
├── main.go                  # Wails app entry point
├── app.go                   # App struct with methods exposed to frontend
├── kanshi/
│   ├── parser.go            # Parse ~/.config/kanshi/config
│   ├── parser_test.go
│   ├── serializer.go        # Serialize back to kanshi config format
│   └── serializer_test.go
├── niri/
│   ├── outputs.go           # Query niri msg outputs --json
│   └── outputs_test.go
├── frontend/
│   ├── src/
│   │   ├── App.svelte       # Root layout: profile bar + canvas + properties
│   │   ├── lib/
│   │   │   ├── Canvas.svelte      # SVG monitor layout, drag + snap
│   │   │   ├── ProfileBar.svelte  # Top: profile selector + actions
│   │   │   ├── Properties.svelte  # Bottom: output properties editor
│   │   │   ├── stores.ts          # Svelte stores for app state
│   │   │   └── types.ts           # TypeScript types matching Go structs
│   │   └── main.ts
│   ├── index.html
│   ├── package.json
│   └── vite.config.ts
├── wails.json
├── go.mod
└── docs/plans/              # Design documents
```

## Key Concepts

- **Kanshi config**: Profiles with named outputs, each having position, scale, mode, transform, enable/disable
- **Niri IPC**: `niri msg outputs --json` for live detection, `niri msg output <name> <action>` for live preview
- **Canvas coordinates**: Logical pixels (resolution / scale), matching what kanshi position values represent
- **Snapping**: Edge-to-edge and center alignment with visual guide lines

## Development

```bash
# Install system dependency (Arch Linux)
sudo pacman -S webkit2gtk-4.1

# Dev mode with hot reload
wails dev

# Build single binary
wails build
```

## Testing

```bash
# Go tests
go test ./...

# Frontend (from frontend/)
npm run check
```

## Config Paths

- Kanshi config: `~/.config/kanshi/config`
- Niri config: `~/.config/niri/config.kdl` (read-only, for reference)
