package kanshi

import (
	"fmt"
	"strings"
)

// Serialize converts a Config struct back into kanshi config format.
func Serialize(config *Config) string {
	var sb strings.Builder

	for i, profile := range config.Profiles {
		if i > 0 {
			sb.WriteString("\n")
		}

		if profile.Name != "" {
			fmt.Fprintf(&sb, "profile \"%s\" {\n", profile.Name)
		} else {
			sb.WriteString("profile {\n")
		}

		for _, output := range profile.Outputs {
			serializeOutput(&sb, &output)
		}

		sb.WriteString("}\n")
	}

	return sb.String()
}

func serializeOutput(sb *strings.Builder, output *Output) {
	fmt.Fprintf(sb, "  output \"%s\" {\n", output.Criteria)

	if output.Enabled != nil {
		if *output.Enabled {
			sb.WriteString("    enable\n")
		} else {
			sb.WriteString("    disable\n")
		}
	}

	if output.Mode != "" {
		fmt.Fprintf(sb, "    mode %s\n", output.Mode)
	}

	if output.Scale != nil {
		// Format scale without trailing zeros, but keep at least one decimal
		s := fmt.Sprintf("%g", *output.Scale)
		if !strings.Contains(s, ".") {
			s += ".0"
		}
		fmt.Fprintf(sb, "    scale %s\n", s)
	}

	if output.Position != nil {
		fmt.Fprintf(sb, "    position %d,%d\n", output.Position.X, output.Position.Y)
	}

	if output.Transform != "" {
		fmt.Fprintf(sb, "    transform %s\n", output.Transform)
	}

	if output.AdaptiveSync != nil {
		if *output.AdaptiveSync {
			sb.WriteString("    adaptive_sync on\n")
		} else {
			sb.WriteString("    adaptive_sync off\n")
		}
	}

	sb.WriteString("  }\n\n")
}
