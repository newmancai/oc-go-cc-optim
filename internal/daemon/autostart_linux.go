// oc-go-cc — Anthropic-to-OpenAI proxy for OpenCode Go + Claude Code
//
// Copyright (C) 2026  Samuel Tuyizere
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

//go:build linux

package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const desktopFileTemplate = `[Desktop Entry]
Type=Application
Name={{.AppName}}
Comment=Start {{.AppName}} on login
Exec="{{.BinaryPath}}" serve --background{{- if .ConfigPath}} --config "{{.ConfigPath}}"{{- end}}{{- if .Port}} --port {{.Port}}{{- end}}
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
`

type desktopData struct {
	AppName    string
	BinaryPath string
	ConfigPath string
	Port       int
}

func autostartFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "autostart", LaunchAgent+".desktop"), nil
}

// EnableAutostart creates a .desktop file in ~/.config/autostart/.
func EnableAutostart(configPath string, port int) error {
	paths, err := DefaultPaths()
	if err != nil {
		return err
	}
	if err := paths.EnsureConfigDir(); err != nil {
		return err
	}

	desktopPath, err := autostartFilePath()
	if err != nil {
		return err
	}

	// Ensure autostart directory exists
	if err := os.MkdirAll(filepath.Dir(desktopPath), 0755); err != nil {
		return fmt.Errorf("cannot create autostart directory: %w", err)
	}

	data := desktopData{
		AppName:    AppName,
		BinaryPath: paths.BinaryPath,
		ConfigPath: configPath,
		Port:       port,
	}

	tmpl, err := template.New("desktop").Parse(desktopFileTemplate)
	if err != nil {
		return fmt.Errorf("cannot parse desktop template: %w", err)
	}

	f, err := os.Create(desktopPath)
	if err != nil {
		return fmt.Errorf("cannot create desktop file: %w", err)
	}
	defer func() { _ = f.Close() }()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("cannot render desktop file: %w", err)
	}

	fmt.Printf("Autostart enabled. %s will start on login.\n", AppName)
	fmt.Printf("  Desktop file: %s\n", desktopPath)
	return nil
}

// DisableAutostart removes the .desktop file from ~/.config/autostart/.
func DisableAutostart() error {
	desktopPath, err := autostartFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(desktopPath); os.IsNotExist(err) {
		fmt.Println("Autostart is not enabled (no .desktop file found)")
		return nil
	}

	if err := os.Remove(desktopPath); err != nil {
		return fmt.Errorf("cannot remove .desktop file: %w", err)
	}

	fmt.Printf("Autostart disabled. .desktop file removed.\n")
	return nil
}

// AutostartStatus reports whether autostart is enabled.
func AutostartStatus() error {
	desktopPath, err := autostartFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(desktopPath); os.IsNotExist(err) {
		fmt.Println("Autostart: disabled (no .desktop file found)")
		return nil
	}

	fmt.Println("Autostart: enabled (.desktop file installed)")
	fmt.Printf("  Desktop file: %s\n", desktopPath)
	return nil
}
