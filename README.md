# Monitoradlo

A GUI application for editing [Kanshi](https://sr.ht/~emersion/kanshi/) monitor layout profiles with live preview via [Niri](https://github.com/YaLTeR/niri) IPC.

Drag monitors on an SVG canvas, edit properties, save to your kanshi config, and preview changes live on your displays — all from a single binary.

## Requirements

- [Niri](https://github.com/YaLTeR/niri) compositor (for output detection and live preview)
- [Kanshi](https://sr.ht/~emersion/kanshi/) (the config file this app edits)
- Go 1.21+
- Node.js 18+
- [Wails v2](https://wails.io/) CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- `webkit2gtk-4.1` (system package)

### Arch Linux

```bash
sudo pacman -S webkit2gtk-4.1
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

## Build

```bash
wails build -tags webkit2_41
```

The binary is at `build/bin/monitoradlo`.

## Development

```bash
wails dev -tags webkit2_41
```

This starts the app with hot reload — frontend changes are applied instantly.

## Testing

```bash
# Go tests
go test ./...

# Frontend type checking
cd frontend && npm run check
```

## Usage

1. Launch the app. It reads your kanshi config from `~/.config/kanshi/config` and detects connected outputs via `niri msg --json outputs`.
2. Select a profile from the dropdown.
3. Drag monitors on the canvas to reposition them. Edges and centers snap to each other.
4. Click a monitor to edit its properties (mode, scale, transform, position, enable/disable).
5. Click **Apply Preview** to temporarily apply changes to your live display via niri.
6. Click **Save** to write the config back and reload kanshi.
