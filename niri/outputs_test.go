package niri

import (
	"testing"
)

const testJSON = `{
  "DP-1": {
    "name": "DP-1",
    "make": "Dell Inc.",
    "model": "DELL U3419W",
    "serial": "7VK66T2",
    "physical_size": [800, 330],
    "modes": [
      {"width": 3440, "height": 1440, "refresh_rate": 59973, "is_preferred": true},
      {"width": 1920, "height": 1080, "refresh_rate": 60000, "is_preferred": false}
    ],
    "current_mode": 0,
    "is_custom_mode": false,
    "vrr_supported": false,
    "vrr_enabled": false,
    "logical": {
      "x": 0, "y": 0,
      "width": 3440, "height": 1440,
      "scale": 1.0,
      "transform": "Normal"
    }
  },
  "eDP-1": {
    "name": "eDP-1",
    "make": "Lenovo Group Limited",
    "model": "0x40A9",
    "serial": null,
    "physical_size": [310, 170],
    "modes": [
      {"width": 1920, "height": 1080, "refresh_rate": 60033, "is_preferred": true}
    ],
    "current_mode": 0,
    "is_custom_mode": false,
    "vrr_supported": false,
    "vrr_enabled": false,
    "logical": {
      "x": 3440, "y": 288,
      "width": 1536, "height": 864,
      "scale": 1.25,
      "transform": "Normal"
    }
  }
}`

func TestParseOutputsJSON(t *testing.T) {
	outputs, err := ParseOutputsJSON([]byte(testJSON))
	if err != nil {
		t.Fatalf("ParseOutputsJSON failed: %v", err)
	}

	if len(outputs) != 2 {
		t.Fatalf("expected 2 outputs, got %d", len(outputs))
	}

	// Find each by connector (map iteration order is random)
	var dp1, edp1 *Output
	for i := range outputs {
		switch outputs[i].Connector {
		case "DP-1":
			dp1 = &outputs[i]
		case "eDP-1":
			edp1 = &outputs[i]
		}
	}

	if dp1 == nil || edp1 == nil {
		t.Fatal("missing expected connectors")
	}

	// DP-1 checks
	if dp1.Make != "Dell Inc." {
		t.Errorf("DP-1 make: got %q", dp1.Make)
	}
	if dp1.Description != "Dell Inc. DELL U3419W 7VK66T2" {
		t.Errorf("DP-1 description: got %q", dp1.Description)
	}
	if len(dp1.AvailableModes) != 2 {
		t.Errorf("DP-1 modes: expected 2, got %d", len(dp1.AvailableModes))
	}
	if dp1.CurrentMode.Width != 3440 || dp1.CurrentMode.Height != 1440 {
		t.Errorf("DP-1 current mode: got %dx%d", dp1.CurrentMode.Width, dp1.CurrentMode.Height)
	}
	// refresh_rate 59973 millihertz = 59.973 Hz
	if dp1.CurrentMode.RefreshRate < 59.97 || dp1.CurrentMode.RefreshRate > 59.98 {
		t.Errorf("DP-1 refresh rate: expected ~59.973, got %f", dp1.CurrentMode.RefreshRate)
	}
	if dp1.LogicalPos == nil || dp1.LogicalPos.X != 0 || dp1.LogicalPos.Y != 0 {
		t.Errorf("DP-1 logical pos: got %v", dp1.LogicalPos)
	}
	if dp1.Scale != 1.0 {
		t.Errorf("DP-1 scale: got %f", dp1.Scale)
	}
	if dp1.PhysicalSize == nil || dp1.PhysicalSize.Width != 800 {
		t.Errorf("DP-1 physical size: got %v", dp1.PhysicalSize)
	}

	// eDP-1 checks
	if edp1.Serial != "" {
		t.Errorf("eDP-1 serial: expected empty, got %q", edp1.Serial)
	}
	if edp1.Description != "Lenovo Group Limited 0x40A9 Unknown" {
		t.Errorf("eDP-1 description: got %q", edp1.Description)
	}
	if edp1.Scale != 1.25 {
		t.Errorf("eDP-1 scale: got %f", edp1.Scale)
	}
	if edp1.LogicalSize == nil || edp1.LogicalSize.Width != 1536 || edp1.LogicalSize.Height != 864 {
		t.Errorf("eDP-1 logical size: got %v", edp1.LogicalSize)
	}
}
