# Monitoradlo

A GUI application for editing [Kanshi](https://sr.ht/~emersion/kanshi/) monitor layout profiles with live preview via [Niri](https://github.com/YaLTeR/niri) IPC.

Drag monitors on an SVG canvas, edit properties, save to your kanshi config, and preview changes live on your displays — all from a single binary.

## Install

### From source (Arch Linux)

```bash
# Install dependencies
sudo pacman -S webkit2gtk-4.1 go
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Build and install
git clone https://github.com/kumekay/monitoradlo.git
cd monitoradlo
wails build -tags webkit2_41
cp build/bin/monitoradlo ~/.local/bin/
```

### Runtime requirements

- [Niri](https://github.com/YaLTeR/niri) compositor (for output detection and live preview)
- [Kanshi](https://sr.ht/~emersion/kanshi/) (the config file this app edits)
- `webkit2gtk-4.1` (runtime dependency)

## Usage

1. Run `monitoradlo`.
2. It reads your kanshi config (`$XDG_CONFIG_HOME/kanshi/config` or `~/.config/kanshi/config`) and detects connected outputs via niri IPC.
3. Select a profile from the dropdown.
4. Drag monitors on the canvas to reposition them. Edges and centers snap to each other.
5. Click a monitor to edit its properties (mode, scale, transform, position, enable/disable).
6. Click **Apply Preview** to temporarily apply changes to your live display via niri.
7. Press **Ctrl+S** or click **Save** to write the config and reload kanshi.

A `.bak` backup is created before each save.

## Development

```bash
# Dev mode with hot reload
wails dev -tags webkit2_41

# Go tests
go test ./...

# Frontend type checking
cd frontend && npm run check
```

### Build requirements

- Go 1.21+
- Node.js 18+
- [Wails v2](https://wails.io/) CLI
- `webkit2gtk-4.1` (system package)

## License

[MIT](LICENSE)
