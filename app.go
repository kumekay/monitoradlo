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
	home, _ := os.UserHomeDir()
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
func (a *App) SaveConfig(config *kanshi.Config) error {
	data := kanshi.Serialize(config)
	err := os.WriteFile(configPath(), []byte(data), 0644)
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
	for action, value := range props {
		var args []string
		switch action {
		case "position":
			// niri msg output <NAME> position set <X> <Y>
			args = []string{"msg", "output", connector, "position", "set"}
			args = append(args, strings.Fields(value)...)
		case "on", "off":
			// niri msg output <NAME> on/off
			args = []string{"msg", "output", connector, action}
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
