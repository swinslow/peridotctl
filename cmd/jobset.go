// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	pbc "github.com/swinslow/peridot-core/pkg/controller"
	"github.com/swinslow/peridotctl/internal/config"
	"github.com/swinslow/peridotctl/internal/outputfmt"
)

func init() {
	var cmdJobSet = &cobra.Command{
		Use:   "jobset",
		Short: "Manage job sets",
		Long: `Manage job set requests with the peridot controller;
		list existing job sets, get information on a particular job set,
		and request to start a new job set.`,
		//Run: jobSetList,
	}
	rootCmd.AddCommand(cmdJobSet)

	var cmdJobSetList = &cobra.Command{
		Use:   "list",
		Short: "List job sets",
		Long: `Get information about job sets requested for the
		peridot controller.`,
		Run: jobSetList,
	}
	cmdJobSet.AddCommand(cmdJobSetList)

	var cmdJobSetStart = &cobra.Command{
		Use:   "start",
		Short: "Start new job set",
		Long: `Request the peridot controller to start a new job set.

Format: peridotctl jobset start NAME [CFGSTRING]

	NAME:      Name of job set template
	CFGSTRING: Optional: job set configuration values (in format key1:value1;key2:value2;...)`,
		Args: cobra.RangeArgs(1, 2),
		Run:  jobSetStart,
	}
	cmdJobSet.AddCommand(cmdJobSetStart)

	var cmdJobSetGet = &cobra.Command{
		Use:   "get",
		Short: "Get info on job set",
		Long: `Get information about a previously-started job set
		for the peridot controller.`,
		Args: cobra.ExactArgs(1),
		Run:  jobSetGet,
	}
	cmdJobSet.AddCommand(cmdJobSetGet)
}

func jobSetList(cmd *cobra.Command, args []string) {
	ctx, cancel := config.GetContext(timeout)
	defer cancel()
	defer conn.Close()

	resp, err := c.GetAllJobSets(ctx, &pbc.GetAllJobSetsReq{})
	if err != nil {
		log.Fatalf("could not get job sets: %v", err)
	}

	fmt.Printf("Job sets:\n\n")

	for _, jsd := range resp.JobSets {
		fmt.Printf("ID: %d\n", jsd.JobSetID)
		fmt.Printf("template name: %s\n", jsd.TemplateName)
		fmt.Printf("runStatus: %s\n", jsd.St.RunStatus.String())
		fmt.Printf("health: %s\n", jsd.St.HealthStatus.String())
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func jobSetStart(cmd *cobra.Command, args []string) {
	ctx, cancel := config.GetContext(timeout)
	defer cancel()
	defer conn.Close()

	name := args[0]
	var cfgStr string
	if len(args) == 2 {
		cfgStr = args[1]
	}

	// check whether values are okay
	if name == "" {
		log.Fatal("no job set template name specified")
	}

	// extract configuration key-value pairs -- semicolons separating pairs,
	// colons separating key from value within a pair
	cfgs := config.ExtractKVs(cfgStr)

	// build into JobSetConfig list
	kvs := []*pbc.JobSetConfig{}
	for k, v := range cfgs {
		kv := pbc.JobSetConfig{Key: k, Value: v}
		kvs = append(kvs, &kv)
	}

	resp, err := c.StartJobSet(ctx, &pbc.StartJobSetReq{
		JstName: name,
		Cfgs:    kvs,
	})
	if err != nil {
		log.Fatalf("could not start job set for template %s: %v", name, err)
	}

	if resp.Success {
		fmt.Printf("job set started for template %s with ID %d\n", name, resp.JobSetID)
	} else {
		fmt.Printf("error starting job set for template %s: %s\n", name, resp.ErrorMsg)
	}
	fmt.Printf("\n")
}

func jobSetGet(cmd *cobra.Command, args []string) {
	ctx, cancel := config.GetContext(timeout)
	defer cancel()
	defer conn.Close()

	jobSetIDStr := args[0]

	jobSetIDInt, err := strconv.Atoi(jobSetIDStr)
	if err != nil || jobSetIDInt < 0 {
		log.Fatalf("invalid job set ID: %s", jobSetIDStr)
	}

	resp, err := c.GetJobSet(ctx, &pbc.GetJobSetReq{JobSetID: uint64(jobSetIDInt)})
	if err != nil {
		log.Fatalf("could not get job set with ID %d: %v", jobSetIDInt, err)
	}

	if resp.Success == false {
		log.Fatalf("job set with ID %d not found: %s", jobSetIDInt, resp.ErrorMsg)
	}

	jsd := resp.JobSet
	fmt.Printf("job set details:\n\n")
	fmt.Printf("  - id: %d\n", jsd.JobSetID)
	fmt.Printf("    templateName: %s\n", jsd.TemplateName)
	fmt.Printf("    status:\n")
	fmt.Printf("      - runStatus: %s\n", jsd.St.RunStatus.String())
	fmt.Printf("        health: %s\n", jsd.St.HealthStatus.String())
	fmt.Printf("        timeStarted: %s\n", time.Unix(jsd.St.TimeStarted, 0).String())
	fmt.Printf("        timeFinished: %s\n", time.Unix(jsd.St.TimeFinished, 0).String())
	fmt.Printf("        outputMessages: %s\n", jsd.St.OutputMessages)
	fmt.Printf("        errorMessages: %s\n", jsd.St.ErrorMessages)
	fmt.Printf("    steps:\n")
	fmt.Println(strings.Join(outputfmt.ConvertSteps(jsd.Steps, 4), "\n"))

	fmt.Printf("\n")
}
