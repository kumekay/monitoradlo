package kanshi

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Config represents a complete kanshi configuration file.
type Config struct {
	Profiles []Profile `json:"profiles"`
}

// Profile represents a kanshi profile with a name and list of outputs.
type Profile struct {
	Name    string   `json:"name"`
	Outputs []Output `json:"outputs"`
}

// Output represents a single output entry within a kanshi profile.
type Output struct {
	Criteria     string    `json:"criteria"`
	Enabled      *bool     `json:"enabled,omitempty"`
	Mode         string    `json:"mode,omitempty"`
	Scale        *float64  `json:"scale,omitempty"`
	Position     *Position `json:"position,omitempty"`
	Transform    string    `json:"transform,omitempty"`
	AdaptiveSync *bool     `json:"adaptiveSync,omitempty"`
}

// Position represents an x,y coordinate pair.
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Parse parses a kanshi config string into a Config struct.
func Parse(input string) (*Config, error) {
	p := &parser{input: input, pos: 0}
	config := &Config{}

	for p.pos < len(p.input) {
		p.skipWhitespaceAndComments()
		if p.pos >= len(p.input) {
			break
		}

		word := p.readWord()
		switch word {
		case "profile":
			profile, err := p.parseProfile()
			if err != nil {
				return nil, err
			}
			config.Profiles = append(config.Profiles, profile)
		case "":
			// skip
		default:
			// Skip unknown top-level directives (include, output defaults, etc.)
			p.skipUntilNewline()
		}
	}

	return config, nil
}

type parser struct {
	input string
	pos   int
}

func (p *parser) parseProfile() (Profile, error) {
	profile := Profile{}
	p.skipWhitespace()

	// Check for optional profile name (quoted or unquoted)
	if p.pos < len(p.input) && p.input[p.pos] != '{' {
		profile.Name = p.readStringOrWord()
		p.skipWhitespace()
	}

	// Expect opening brace
	if p.pos >= len(p.input) || p.input[p.pos] != '{' {
		return profile, fmt.Errorf("expected '{' at position %d", p.pos)
	}
	p.pos++ // skip '{'

	for p.pos < len(p.input) {
		p.skipWhitespaceAndComments()
		if p.pos >= len(p.input) {
			return profile, fmt.Errorf("unexpected end of input in profile")
		}
		if p.input[p.pos] == '}' {
			p.pos++
			break
		}

		word := p.readWord()
		switch word {
		case "output":
			output, err := p.parseOutput()
			if err != nil {
				return profile, err
			}
			profile.Outputs = append(profile.Outputs, output)
		case "exec":
			p.skipUntilNewline()
		default:
			p.skipUntilNewline()
		}
	}

	return profile, nil
}

func (p *parser) parseOutput() (Output, error) {
	output := Output{}
	p.skipWhitespace()

	// Read criteria (quoted string or unquoted word)
	output.Criteria = p.readStringOrWord()
	p.skipWhitespace()

	// Check if directives are in a block or inline
	if p.pos < len(p.input) && p.input[p.pos] == '{' {
		p.pos++ // skip '{'
		if err := p.parseOutputDirectives(&output, '}'); err != nil {
			return output, err
		}
		p.pos++ // skip '}'
	} else {
		if err := p.parseOutputDirectives(&output, '\n'); err != nil {
			return output, err
		}
	}

	return output, nil
}

func (p *parser) parseOutputDirectives(output *Output, terminator byte) error {
	for p.pos < len(p.input) {
		p.skipWhitespaceAndComments()
		if p.pos >= len(p.input) {
			break
		}
		if p.input[p.pos] == terminator {
			break
		}
		if terminator == '\n' && p.pos > 0 && p.input[p.pos-1] == '\n' {
			// For inline directives, stop at newline
			break
		}

		word := p.readWord()
		switch word {
		case "enable":
			b := true
			output.Enabled = &b
		case "disable":
			b := false
			output.Enabled = &b
		case "mode":
			p.skipWhitespace()
			output.Mode = p.readStringOrWord()
		case "scale":
			p.skipWhitespace()
			s := p.readWord()
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return fmt.Errorf("invalid scale value %q: %w", s, err)
			}
			output.Scale = &f
		case "position":
			p.skipWhitespace()
			posStr := p.readWord()
			parts := strings.SplitN(posStr, ",", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid position %q, expected x,y", posStr)
			}
			x, err := strconv.Atoi(parts[0])
			if err != nil {
				return fmt.Errorf("invalid position x %q: %w", parts[0], err)
			}
			y, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("invalid position y %q: %w", parts[1], err)
			}
			output.Position = &Position{X: x, Y: y}
		case "transform":
			p.skipWhitespace()
			output.Transform = p.readStringOrWord()
		case "adaptive_sync":
			p.skipWhitespace()
			val := p.readWord()
			b := val == "on"
			output.AdaptiveSync = &b
		case "":
			// skip
		default:
			// Unknown directive, skip its value if any
			p.skipWhitespace()
			if p.pos < len(p.input) && p.input[p.pos] != '\n' && p.input[p.pos] != '}' {
				p.readStringOrWord()
			}
		}
	}
	return nil
}

// readWord reads a non-whitespace, non-brace word.
func (p *parser) readWord() string {
	start := p.pos
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if unicode.IsSpace(rune(ch)) || ch == '{' || ch == '}' {
			break
		}
		p.pos++
	}
	return p.input[start:p.pos]
}

// readStringOrWord reads a quoted string or an unquoted word.
func (p *parser) readStringOrWord() string {
	if p.pos < len(p.input) && p.input[p.pos] == '"' {
		return p.readQuotedString()
	}
	return p.readWord()
}

// readQuotedString reads a double-quoted string.
func (p *parser) readQuotedString() string {
	if p.pos >= len(p.input) || p.input[p.pos] != '"' {
		return ""
	}
	p.pos++ // skip opening quote
	start := p.pos
	for p.pos < len(p.input) && p.input[p.pos] != '"' {
		p.pos++
	}
	s := p.input[start:p.pos]
	if p.pos < len(p.input) {
		p.pos++ // skip closing quote
	}
	return s
}

func (p *parser) skipWhitespace() {
	for p.pos < len(p.input) && unicode.IsSpace(rune(p.input[p.pos])) {
		p.pos++
	}
}

func (p *parser) skipWhitespaceAndComments() {
	for p.pos < len(p.input) {
		if unicode.IsSpace(rune(p.input[p.pos])) {
			p.pos++
			continue
		}
		if p.input[p.pos] == '#' {
			p.skipUntilNewline()
			continue
		}
		break
	}
}

func (p *parser) skipUntilNewline() {
	for p.pos < len(p.input) && p.input[p.pos] != '\n' {
		p.pos++
	}
}
