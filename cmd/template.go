// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"

	pbc "github.com/swinslow/peridot-core/pkg/controller"
	"github.com/swinslow/peridotctl/internal/config"
	"github.com/swinslow/peridotctl/internal/outputfmt"
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
		fmt.Println(strings.Join(outputfmt.ConvertStepTemplates(jst.Steps, 4), "\n"))
	}
	fmt.Printf("\n")
}
