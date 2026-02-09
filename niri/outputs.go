package niri

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Output represents a connected output as reported by niri.
type Output struct {
	Connector      string  `json:"connector"`
	Make           string  `json:"make"`
	Model          string  `json:"model"`
	Serial         string  `json:"serial"`
	Description    string  `json:"description"`
	CurrentMode    Mode    `json:"currentMode"`
	AvailableModes []Mode  `json:"availableModes"`
	LogicalPos     *Pos    `json:"logicalPosition"`
	LogicalSize    *Size   `json:"logicalSize"`
	Scale          float64 `json:"scale"`
	Transform      string  `json:"transform"`
	PhysicalSize   *Size   `json:"physicalSize"`
}

// Mode represents a display mode.
type Mode struct {
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	RefreshRate float64 `json:"refreshRate"` // In Hz (e.g. 59.973)
	IsCurrent   bool    `json:"isCurrent"`
	IsPreferred bool    `json:"isPreferred"`
}

// Pos represents a position.
type Pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Size represents dimensions.
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// niriOutputJSON matches the actual JSON from `niri msg --json outputs`.
type niriOutputJSON struct {
	Name         string           `json:"name"`
	Make         string           `json:"make"`
	Model        string           `json:"model"`
	Serial       *string          `json:"serial"` // nullable
	PhysicalSize json.RawMessage  `json:"physical_size"`
	Modes        []niriModeJSON   `json:"modes"`
	CurrentMode  int              `json:"current_mode"` // index into modes
	Logical      *niriLogicalJSON `json:"logical"`
	VrrSupported bool             `json:"vrr_supported"`
	VrrEnabled   bool             `json:"vrr_enabled"`
}

type niriModeJSON struct {
	Width       int  `json:"width"`
	Height      int  `json:"height"`
	RefreshRate int  `json:"refresh_rate"` // millihertz
	IsPreferred bool `json:"is_preferred"`
}

type niriLogicalJSON struct {
	X         int     `json:"x"`
	Y         int     `json:"y"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	Scale     float64 `json:"scale"`
	Transform string  `json:"transform"`
}

// DetectOutputs queries niri for currently connected outputs.
func DetectOutputs() ([]Output, error) {
	cmd := exec.Command("niri", "msg", "--json", "outputs")
	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running niri msg outputs: %w", err)
	}

	return ParseOutputsJSON(data)
}

// ParseOutputsJSON parses the JSON output from `niri msg --json outputs`.
func ParseOutputsJSON(data []byte) ([]Output, error) {
	// niri returns a map keyed by connector name
	var rawMap map[string]niriOutputJSON
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return nil, fmt.Errorf("parsing niri outputs JSON: %w", err)
	}

	var outputs []Output
	for connector, r := range rawMap {
		o := Output{
			Connector: connector,
			Make:      r.Make,
			Model:     r.Model,
		}

		if r.Serial != nil {
			o.Serial = *r.Serial
		}

		// Build kanshi-style description: "Make Model Serial"
		serial := o.Serial
		if serial == "" {
			serial = "Unknown"
		}
		parts := []string{}
		if r.Make != "" {
			parts = append(parts, r.Make)
		}
		if r.Model != "" {
			parts = append(parts, r.Model)
		}
		parts = append(parts, serial)
		o.Description = strings.Join(parts, " ")

		// Parse modes (refresh_rate is in millihertz)
		for i, m := range r.Modes {
			mode := Mode{
				Width:       m.Width,
				Height:      m.Height,
				RefreshRate: float64(m.RefreshRate) / 1000.0,
				IsCurrent:   i == r.CurrentMode,
				IsPreferred: m.IsPreferred,
			}
			o.AvailableModes = append(o.AvailableModes, mode)
			if i == r.CurrentMode {
				o.CurrentMode = mode
			}
		}

		// Logical info
		if r.Logical != nil {
			o.LogicalPos = &Pos{X: r.Logical.X, Y: r.Logical.Y}
			o.LogicalSize = &Size{Width: r.Logical.Width, Height: r.Logical.Height}
			o.Scale = r.Logical.Scale
			o.Transform = r.Logical.Transform
		}

		// Physical size: array [width, height]
		if len(r.PhysicalSize) > 0 {
			var ps [2]int
			if err := json.Unmarshal(r.PhysicalSize, &ps); err == nil {
				o.PhysicalSize = &Size{Width: ps[0], Height: ps[1]}
			}
		}

		outputs = append(outputs, o)
	}

	return outputs, nil
}
