package kanshi

import (
	"os"
	"testing"
)

func TestParseRealConfig(t *testing.T) {
	input := `# ThinkPad T14

profile "Home" {
  output "Lenovo Group Limited 0x40A9 Unknown" {
    enable
    scale 1.25
    position 384,1200
  }

  output "Samsung Electric Company SMS24A850 HTRCC00024" {
    enable
    position 0,0
  }
}

profile "Office" {
  output "Lenovo Group Limited 0x40A9 Unknown" {
    enable
    scale 1.25
    position 3440,288
  }

  output "Dell Inc. DELL U3419W 7VK66T2" {
    enable
    position 0,0
  }
}


# Legion Go

profile "Legion Go Kitchen" {
  output "HP Inc. HP U28 4K HDR 1CR1411RL6" {
    scale 1.5
    position 0,0
  }

  output "Lenovo Group Limited Go Display 0x00888888" {
    scale 2.25
    position 0,1440
  }
}
`
	config, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(config.Profiles) != 3 {
		t.Fatalf("expected 3 profiles, got %d", len(config.Profiles))
	}

	// Profile "Home"
	home := config.Profiles[0]
	if home.Name != "Home" {
		t.Errorf("expected profile name 'Home', got %q", home.Name)
	}
	if len(home.Outputs) != 2 {
		t.Fatalf("Home: expected 2 outputs, got %d", len(home.Outputs))
	}
	if home.Outputs[0].Criteria != "Lenovo Group Limited 0x40A9 Unknown" {
		t.Errorf("Home output 0 criteria: got %q", home.Outputs[0].Criteria)
	}
	if home.Outputs[0].Scale == nil || *home.Outputs[0].Scale != 1.25 {
		t.Errorf("Home output 0 scale: expected 1.25, got %v", home.Outputs[0].Scale)
	}
	if home.Outputs[0].Position == nil || home.Outputs[0].Position.X != 384 || home.Outputs[0].Position.Y != 1200 {
		t.Errorf("Home output 0 position: expected 384,1200, got %v", home.Outputs[0].Position)
	}
	if home.Outputs[0].Enabled == nil || *home.Outputs[0].Enabled != true {
		t.Errorf("Home output 0 enabled: expected true, got %v", home.Outputs[0].Enabled)
	}

	// Profile "Office"
	office := config.Profiles[1]
	if office.Name != "Office" {
		t.Errorf("expected profile name 'Office', got %q", office.Name)
	}
	if len(office.Outputs) != 2 {
		t.Fatalf("Office: expected 2 outputs, got %d", len(office.Outputs))
	}

	// Profile "Legion Go Kitchen"
	legion := config.Profiles[2]
	if legion.Name != "Legion Go Kitchen" {
		t.Errorf("expected profile name 'Legion Go Kitchen', got %q", legion.Name)
	}
	if len(legion.Outputs) != 2 {
		t.Fatalf("Legion Go Kitchen: expected 2 outputs, got %d", len(legion.Outputs))
	}
	if legion.Outputs[0].Scale == nil || *legion.Outputs[0].Scale != 1.5 {
		t.Errorf("Legion output 0 scale: expected 1.5, got %v", legion.Outputs[0].Scale)
	}
	// No explicit enable/disable in this profile
	if legion.Outputs[0].Enabled != nil {
		t.Errorf("Legion output 0 enabled: expected nil, got %v", *legion.Outputs[0].Enabled)
	}
}

func TestParseAndSerializeRoundTrip(t *testing.T) {
	input := `profile "Test" {
  output "Some Monitor 1234" {
    enable
    mode 1920x1080@60Hz
    scale 2.0
    position 100,200
    transform 90
  }

  output "Other Monitor 5678" {
    disable
  }
}
`
	config, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	output := Serialize(config)

	// Re-parse the serialized output
	config2, err := Parse(output)
	if err != nil {
		t.Fatalf("Re-parse failed: %v", err)
	}

	if len(config2.Profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(config2.Profiles))
	}

	p := config2.Profiles[0]
	if p.Name != "Test" {
		t.Errorf("expected name 'Test', got %q", p.Name)
	}
	if len(p.Outputs) != 2 {
		t.Fatalf("expected 2 outputs, got %d", len(p.Outputs))
	}

	o := p.Outputs[0]
	if o.Criteria != "Some Monitor 1234" {
		t.Errorf("criteria: got %q", o.Criteria)
	}
	if o.Enabled == nil || *o.Enabled != true {
		t.Errorf("enabled: expected true")
	}
	if o.Mode != "1920x1080@60Hz" {
		t.Errorf("mode: got %q", o.Mode)
	}
	if o.Scale == nil || *o.Scale != 2.0 {
		t.Errorf("scale: expected 2.0")
	}
	if o.Position == nil || o.Position.X != 100 || o.Position.Y != 200 {
		t.Errorf("position: expected 100,200")
	}
	if o.Transform != "90" {
		t.Errorf("transform: got %q", o.Transform)
	}

	o2 := p.Outputs[1]
	if o2.Enabled == nil || *o2.Enabled != false {
		t.Errorf("output 2 enabled: expected false")
	}
}

// TestParseActualFile attempts to parse the real kanshi config if it exists.
func TestParseActualFile(t *testing.T) {
	home, _ := os.UserHomeDir()
	path := home + "/.config/kanshi/config"
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skip("kanshi config not found, skipping")
	}

	config, err := Parse(string(data))
	if err != nil {
		t.Fatalf("Parse failed on real config: %v", err)
	}

	if len(config.Profiles) == 0 {
		t.Error("expected at least one profile in real config")
	}

	t.Logf("Parsed %d profiles from %s", len(config.Profiles), path)
	for _, p := range config.Profiles {
		t.Logf("  Profile %q: %d outputs", p.Name, len(p.Outputs))
	}
}
