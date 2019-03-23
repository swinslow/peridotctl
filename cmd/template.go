// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"

	pbc "github.com/swinslow/peridot-core/pkg/controller"
	"github.com/swinslow/peridotctl/internal/config"
)

func init() {
	var cmdTemplate = &cobra.Command{
		Use:   "template",
		Short: "Manage job set template registrations",
		Long: `Manage job set template registrations with the peridot controller;
		list known templates, get information on a particular template, and
		add a new template registration.`,
		//Run: templateList,
	}
	rootCmd.AddCommand(cmdTemplate)

	var cmdTemplateList = &cobra.Command{
		Use:   "list",
		Short: "Get all registered job set templates",
		Long: `Get information about all job set templates registered with the
		peridot controller.`,
		Run: templateList,
	}
	cmdTemplate.AddCommand(cmdTemplateList)
}

func templateList(cmd *cobra.Command, args []string) {
	ctx, cancel := config.GetContext(timeout)
	defer cancel()
	defer conn.Close()

	resp, err := c.GetAllJobSetTemplates(ctx, &pbc.GetAllJobSetTemplatesReq{})
	if err != nil {
		log.Fatalf("could not get job set templates: %v", err)
	}

	fmt.Printf("Registered job set templates:\n\n")

	for _, jst := range resp.Jsts {
		fmt.Printf("  - name: %s\n", jst.Name)
		fmt.Printf("    steps:\n")
		fmt.Println(strings.Join(convertSteps(jst.Steps, 4), "\n"))
	}
	fmt.Printf("\n")
}

func convertSteps(steps []*pbc.StepTemplate, indent int) []string {
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
			subStepLines := convertSteps(x.Concurrent.Steps, indent+4)
			lines = append(lines, line1)
			lines = append(lines, line2)
			lines = append(lines, subStepLines...)
		}
	}

	return lines
}
