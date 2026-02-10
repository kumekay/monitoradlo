package main

import (
	"context"
	"fmt"
	"monitoradlo/kanshi"
	"monitoradlo/niri"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// App struct holds application state and is bound to the frontend.
type App struct {
	ctx context.Context
}

// NewApp creates a new App instance.
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// configPath returns the path to the kanshi config file.
func configPath() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "kanshi", "config")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}
	return filepath.Join(home, ".config", "kanshi", "config")
}

// LoadConfig reads and parses the kanshi config file.
func (a *App) LoadConfig() (*kanshi.Config, error) {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return nil, fmt.Errorf("reading kanshi config: %w", err)
	}
	config, err := kanshi.Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parsing kanshi config: %w", err)
	}
	return config, nil
}

// SaveConfig serializes and writes the kanshi config file.
// Creates a .bak backup of the existing file before overwriting.
func (a *App) SaveConfig(config *kanshi.Config) error {
	path := configPath()
	data := kanshi.Serialize(config)

	// Backup existing file before overwriting
	if existing, err := os.ReadFile(path); err == nil {
		bakPath := path + ".bak"
		if err := os.WriteFile(bakPath, existing, 0644); err != nil {
			return fmt.Errorf("creating backup: %w", err)
		}
	}

	err := os.WriteFile(path, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("writing kanshi config: %w", err)
	}
	return nil
}

// DetectOutputs queries niri for currently connected outputs.
func (a *App) DetectOutputs() ([]niri.Output, error) {
	return niri.DetectOutputs()
}

// ApplyPreview applies temporary output settings via niri msg.
func (a *App) ApplyPreview(connector string, props map[string]string) error {
	// Apply in a deterministic order to avoid transform/position races.
	// If turning off, just do that and return.
	if _, ok := props["off"]; ok {
		cmd := exec.Command("niri", "msg", "output", connector, "off")
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("niri msg output %s off: %s: %w", connector, string(out), err)
		}
		return nil
	}

	order := []string{"on", "mode", "scale", "transform", "position", "vrr"}
	for _, action := range order {
		value, ok := props[action]
		if !ok {
			continue
		}
		var args []string
		switch action {
		case "position":
			// niri msg output <NAME> position set <X> <Y>
			args = []string{"msg", "output", connector, "position", "set", "--"}
			args = append(args, strings.Fields(value)...)
		case "on":
			// niri msg output <NAME> on
			args = []string{"msg", "output", connector, "on"}
		default:
			// mode, scale, transform: niri msg output <NAME> <action> <value>
			args = []string{"msg", "output", connector, action}
			if value != "" {
				args = append(args, strings.Fields(value)...)
			}
		}
		cmd := exec.Command("niri", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("niri msg output %s %s: %s: %w", connector, action, string(out), err)
		}
	}
	return nil
}

// ReloadKanshi signals kanshi to reload its config.
func (a *App) ReloadKanshi() error {
	// Try kanshictl first, fall back to pkill
	if err := exec.Command("kanshictl", "reload").Run(); err != nil {
		if err := exec.Command("pkill", "-HUP", "kanshi").Run(); err != nil {
			return fmt.Errorf("reloading kanshi: %w", err)
		}
	}
	return nil
}
