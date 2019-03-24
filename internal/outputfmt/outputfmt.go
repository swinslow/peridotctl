// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package outputfmt

import (
	"fmt"
	"strings"

	pbc "github.com/swinslow/peridot-core/pkg/controller"
)

// ConvertSteps converts a slice of pointers to Steps into a formatted,
// YAML-style slice of string lines.
func ConvertSteps(steps []*pbc.Step, indent int) []string {
	lines := []string{}

	for _, step := range steps {
		commonLines := []string{
			fmt.Sprintf("%s    stepID: %d", strings.Repeat(" ", indent), step.StepID),
			fmt.Sprintf("%s    stepOrder: %d", strings.Repeat(" ", indent), step.StepOrder),
			fmt.Sprintf("%s    runStatus: %s", strings.Repeat(" ", indent), step.RunStatus),
			fmt.Sprintf("%s    health: %s", strings.Repeat(" ", indent), step.HealthStatus),
		}
		switch x := step.S.(type) {
		case *pbc.Step_Agent:
			newLines := []string{
				fmt.Sprintf("%s  - type: agent", strings.Repeat(" ", indent)),
				fmt.Sprintf("%s    name: %s", strings.Repeat(" ", indent), x.Agent.AgentName),
				fmt.Sprintf("%s    jobID: %d", strings.Repeat(" ", indent), x.Agent.JobID),
			}
			lines = append(lines, newLines...)
			lines = append(lines, commonLines...)
		case *pbc.Step_Jobset:
			newLines := []string{
				fmt.Sprintf("%s  - type: jobset", strings.Repeat(" ", indent)),
				fmt.Sprintf("%s    templateName: %s", strings.Repeat(" ", indent), x.Jobset.TemplateName),
				fmt.Sprintf("%s    jobSetID: %d", strings.Repeat(" ", indent), x.Jobset.JobSetID),
			}
			lines = append(lines, newLines...)
			lines = append(lines, commonLines...)
		case *pbc.Step_Concurrent:
			line1 := fmt.Sprintf("%s  - type: concurrent", strings.Repeat(" ", indent))
			line2 := fmt.Sprintf("%s    steps:", strings.Repeat(" ", indent))
			subStepLines := ConvertSteps(x.Concurrent.Steps, indent+4)
			lines = append(lines, line1)
			lines = append(lines, commonLines...)
			lines = append(lines, line2)
			lines = append(lines, subStepLines...)
		}
	}

	return lines
}

// ConvertStepTemplates converts a slice of pointers to StepTemplates into
// a formatted, YAML-style slice of string lines.
func ConvertStepTemplates(steps []*pbc.StepTemplate, indent int) []string {
	lines := []string{}

	for _, step := range steps {
		switch x := step.S.(type) {
		case *pbc.StepTemplate_Agent:
			line1 := fmt.Sprintf("%s  - type: agent", strings.Repeat(" ", indent))
			line2 := fmt.Sprintf("%s    name: %s", strings.Repeat(" ", indent), x.Agent.Name)
			lines = append(lines, line1)
			lines = append(lines, line2)
		case *pbc.StepTemplate_Jobset:
			line1 := fmt.Sprintf("%s  - type: jobset", strings.Repeat(" ", indent))
			line2 := fmt.Sprintf("%s    name: %s", strings.Repeat(" ", indent), x.Jobset.Name)
			lines = append(lines, line1)
			lines = append(lines, line2)
		case *pbc.StepTemplate_Concurrent:
			line1 := fmt.Sprintf("%s  - type: concurrent", strings.Repeat(" ", indent))
			line2 := fmt.Sprintf("%s    steps:", strings.Repeat(" ", indent))
			subStepLines := ConvertStepTemplates(x.Concurrent.Steps, indent+4)
			lines = append(lines, line1)
			lines = append(lines, line2)
			lines = append(lines, subStepLines...)
		}
	}

	return lines
}
