// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser

import "fmt"

// ValidateReq checks the request object to confirm it is valid.
// It returns the first found or nil if request is okay.
func ValidateReq(req *PeridotReq) error {
	// check the API version is valid
	if req.APIVersion != "v0-alpha1" {
		return fmt.Errorf("unknown apiVersion: %s", req.APIVersion)
	}

	// check the agents list
	if err := ValidateAgents(req.Agents); err != nil {
		return err
	}

	// check the templates list
	if err := ValidateJobSetTemplates(req.Templates); err != nil {
		return err
	}

	// looks good!
	return nil
}

// ValidateAgents checks the requested Agents object to confirm it is valid.
// It returns the first found or nil if request is okay.
func ValidateAgents(agents []PeridotAgent) error {
	for _, agent := range agents {
		// check that it has a name
		if agent.Name == "" {
			return fmt.Errorf("got agent with no name, expected name")
		}
		// check that it has a URL
		if agent.URL == "" {
			return fmt.Errorf("got agent name %s with no URL, expected URL", agent.Name)
		}
		// check that it has a type
		if agent.TypeStr == "" {
			return fmt.Errorf("got agent name %s with no type, expected type", agent.Name)
		}
		// port needs to be non-zero
		if agent.Port == 0 {
			return fmt.Errorf("invalid Port for agent %s: got 0, must be non-zero", agent.Name)
		}
	}

	// looks good!
	return nil
}

// ValidateJobSetTemplates checks the requested JobSetTemplates object to
// confirm it is valid. It returns the first found or nil if request is okay.
func ValidateJobSetTemplates(templates []PeridotJobSetTemplate) error {
	for _, template := range templates {
		// check that it has a name
		if template.Name == "" {
			return fmt.Errorf("got template with no name, expected name")
		}
		// check that steps are valid
		if err := ValidateJSTSteps(template.Steps); err != nil {
			return err
		}

	}

	// looks good!
	return nil
}

// ValidateJSTSteps checks the requested JobSetTemplate Steps object to
// confirm it is valid, including recursively check concurrent sub-steps.
// It returns the first found or nil if request is okay.
func ValidateJSTSteps(steps []PeridotJSTStep) error {
	if len(steps) == 0 {
		return fmt.Errorf("got zero steps for template, expected greater than zero steps")
	}

	for _, step := range steps {
		switch step.TypeStr {
		case "agent":
			if step.Name == "" {
				return fmt.Errorf("got template step type agent with no name, expected name")
			}
			if len(step.Steps) != 0 {
				return fmt.Errorf("got steps for template type agent, name %s, expected zero steps", step.Name)
			}

		case "jobset":
			if step.Name == "" {
				return fmt.Errorf("got template step type jobset with no name, expected name")
			}
			if len(step.Steps) != 0 {
				return fmt.Errorf("got steps for template type jobset, name %s, expected zero steps", step.Name)
			}

		case "concurrent":
			if step.Name != "" {
				return fmt.Errorf("got template step type concurrent with name, expected no name")
			}
			if len(step.Steps) == 0 {
				return fmt.Errorf("got zero steps for template type concurrent, expected greater than zero steps")
			}
			err := ValidateJSTSteps(step.Steps)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("got template step type %s, invalid type", step.TypeStr)
		}
	}

	return nil
}
