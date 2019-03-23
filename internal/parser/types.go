// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser

// PeridotReq represents the parsed YAML data for a collection of
// peridotctl objects.
type PeridotReq struct {
	APIVersion string                  `yaml:"apiVersion"`
	Agents     []PeridotAgent          `yaml:",omitempty"`
	Templates  []PeridotJobSetTemplate `yaml:"jobSetTemplates,omitempty"`
}

// PeridotAgent represents the parsed YAML data for a peridotctl
// agent object.
type PeridotAgent struct {
	Name    string
	URL     string `yaml:"url"`
	Port    uint32
	TypeStr string            `yaml:"type"`
	Configs map[string]string `yaml:",omitempty"`
}

// PeridotJobSetTemplate represents the parsed YAML data for a
// peridotctl JobSetTemplate object.
type PeridotJobSetTemplate struct {
	Name  string
	Steps []PeridotJSTStep
}

// PeridotJSTStep represents the parsed YAML data for a single step in
// a peridotctl JobSetTemplate object.
type PeridotJSTStep struct {
	TypeStr string `yaml:"type"`
	Name    string
	Steps   []PeridotJSTStep
}
