# Monitoradlo: Kanshi/Niri Monitor Layout GUI

## Overview

A GUI application for editing Kanshi monitor layout profiles, with live preview via Niri IPC. Built with Wails v2 (Go backend) + Svelte (frontend). Produces a single distributable binary.

## Tech Stack

- **Backend**: Go + Wails v2
- **Frontend**: Svelte + TypeScript + Vite
- **Rendering**: SVG for the monitor canvas (native hit-testing, pointer events)
- **System deps**: webkit2gtk-4.1 (Wails runtime requirement)

## Data Model

### KanshiConfig (Go + mirrored in TypeScript)

```
Config
  └── []Profile
        ├── Name: string
        └── []Output
              ├── Criteria: string    // e.g. "Dell Inc. DELL U3419W 7VK66T2"
              ├── Enabled: bool
              ├── Mode: *string       // e.g. "3440x1440@59.973Hz"
              ├── Scale: *float64
              ├── Position: *Position  // {X: int, Y: int}
              ├── Transform: *string   // "normal", "90", "180", "270", "flipped", etc.
              └── AdaptiveSync: *bool
```

### NiriOutput (live detection)

```
NiriOutput
  ├── Connector: string       // "DP-1"
  ├── Make: string
  ├── Model: string
  ├── Serial: string
  ├── Description: string     // "Make Model Serial" for kanshi matching
  ├── CurrentMode: Mode
  ├── AvailableModes: []Mode
  ├── LogicalPosition: Position
  ├── LogicalSize: Size
  ├── Scale: float64
  ├── Transform: string
  └── PhysicalSizeMM: Size
```

Matching between kanshi outputs and live niri outputs uses the description string ("Make Model Serial").

## Window Layout

Vertical 3-panel, top-to-bottom:

```
+----------------------------------------------+
| Profile: [Home v]    [+ New] [Rename] [Del] [Save] |  <- Profile bar
+----------------------------------------------+
|                                              |
|     +----------+  +------+                   |
|     |DP-1      |  |eDP-1 |                   |  <- Canvas (SVG)
|     |3440x1440 |  |1536x |                   |
|     |          |  |864   |                   |
|     +----------+  +------+                   |
|                                              |
+----------------------------------------------+
| Output: DP-1 - Dell U3419W                  |
| Mode: [3440x1440@60Hz v]  Scale: [1.0]      |  <- Properties panel
| Transform: [Normal v]  [x] Enabled           |
| Position: x [0]  y [0]       [Apply Preview] |
+----------------------------------------------+
```

## Canvas Design

### Coordinate System
- Monitor rectangles sized by logical dimensions (resolution / scale)
- Canvas auto-zooms to fit all monitors with padding
- Zoom factor: `min(canvasW, canvasH) / boundingBox` with margin

### Monitor Rendering
- Each monitor: SVG `<g>` with `<rect>` + `<text>` (connector name + resolution)
- Color-coded per monitor, selected monitor has highlighted border
- Disabled monitors: hatched pattern, reduced opacity

### Drag & Drop
- Mousedown selects monitor, shows properties
- Mousemove updates position in real-time
- Mouseup commits position (snapped)

### Snapping (Edge + Center)
- On drag, check all edges and center lines of dragged monitor against all others
- Snap threshold: 15 canvas pixels
- Draw guide lines showing active snap alignments
- Priority: edge-to-edge > center-to-center > edge-to-center
- Snap targets: top, bottom, left, right edges + horizontal/vertical center

## Properties Panel

Shown when a monitor is selected:
- **Enable/Disable** checkbox
- **Mode** dropdown (populated from niri available_modes when detected)
- **Scale** numeric input (step 0.25)
- **Transform** dropdown (Normal, 90, 180, 270, Flipped, Flipped-90, Flipped-180, Flipped-270)
- **Position** x,y inputs (bidirectional sync with canvas drag)
- **Apply Preview** button: runs `niri msg output <name> <property>` for temporary live preview

## Profile Management

- **Dropdown** selector listing all profiles
- **New**: creates profile pre-populated with currently connected outputs (niri detection)
- **Rename**: inline rename
- **Delete**: with confirmation dialog
- **Save**: serializes all profiles to `~/.config/kanshi/config`, then signals kanshi to reload

## Backend API (Go methods exposed to frontend)

```go
// Config operations
func (a *App) LoadConfig() (*Config, error)
func (a *App) SaveConfig(config *Config) error

// Live output detection
func (a *App) DetectOutputs() ([]NiriOutput, error)

// Live preview
func (a *App) ApplyPreview(connector string, props OutputProps) error

// Reload kanshi
func (a *App) ReloadKanshi() error
```

## Project Structure

```
monitoradlo/
├── main.go
├── app.go                   # Wails-bound App struct
├── kanshi/
│   ├── parser.go
│   ├── parser_test.go
│   ├── serializer.go
│   └── serializer_test.go
├── niri/
│   ├── outputs.go
│   └── outputs_test.go
├── frontend/
│   ├── src/
│   │   ├── App.svelte
│   │   ├── lib/
│   │   │   ├── Canvas.svelte
│   │   │   ├── ProfileBar.svelte
│   │   │   ├── Properties.svelte
│   │   │   ├── stores.ts
│   │   │   └── types.ts
│   │   └── main.ts
│   ├── index.html
│   ├── package.json
│   └── vite.config.ts
├── wails.json
└── go.mod
```

## Build & Distribution

- `wails build` produces a single binary
- Runtime dependency: webkit2gtk-4.1 on Linux
- Config path: `~/.config/kanshi/config` (hardcoded default, could accept flag)
