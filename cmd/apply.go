// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	pbc "github.com/swinslow/peridot-core/pkg/controller"
	"github.com/swinslow/peridotctl/internal/config"
	"github.com/swinslow/peridotctl/internal/parser"
)

func init() {
	var cmdApply = &cobra.Command{
		Use:   "apply",
		Short: "Apply YAML file",
		Long: `Apply configurations and actions from a YAML file to
the peridot controller.

Format: peridotctl apply YAMLFILE

YAMLFILE: path to YAML file to apply`,
		Args: cobra.ExactArgs(1),
		Run:  apply,
	}
	rootCmd.AddCommand(cmdApply)
}

func apply(cmd *cobra.Command, args []string) {
	ctx, cancel := config.GetContext(timeout)
	defer cancel()
	defer conn.Close()

	// load and parse YAML file, and confirm it is valid
	req, err := parser.ParseYAML(args[0])
	if err != nil {
		log.Fatalf("error parsing %s: %v", args[0], err)
	}

	// build any Agents
	err = applyAgents(ctx, req.Agents)
	if err != nil {
		log.Fatalf("error requesting agents: %v", err)
	}

	// build any JobSetTemplates
	err = applyTemplates(ctx, req.Templates)
	if err != nil {
		log.Fatalf("error requesting JobSetTemplates: %v", err)
	}

	// we're done! will cancel and close connection
}

func applyAgents(ctx context.Context, agents []parser.PeridotAgent) error {
	for _, agent := range agents {
		// build configs into AgentKV list
		kvs := []*pbc.AgentConfig_AgentKV{}
		for k, v := range agent.Configs {
			kv := pbc.AgentConfig_AgentKV{Key: k, Value: v}
			kvs = append(kvs, &kv)
		}

		// build AgentConfig object
		ac := &pbc.AgentConfig{
			Name: agent.Name,
			Url:  agent.URL,
			Port: agent.Port,
			Type: agent.TypeStr,
			Kvs:  kvs,
		}

		resp, err := c.AddAgent(ctx, &pbc.AddAgentReq{Cfg: ac})
		if err != nil {
			return fmt.Errorf("could not add agent %s: %v", agent.Name, err)
		}

		if resp.Success {
			fmt.Printf("agent %s successfully registered\n", agent.Name)
		} else {
			fmt.Printf("error registering agent %s: %s\n", agent.Name, resp.ErrorMsg)
		}
	}

	return nil
}

func applyTemplates(ctx context.Context, templates []parser.PeridotJobSetTemplate) error {
	for _, template := range templates {

		// translate template object into protobuf version of StepTemplates
		steps, err := buildStepTemplates(template.Steps)
		if err != nil {
			return err
		}

		// build JobSetTemplate object
		jst := &pbc.JobSetTemplate{
			Name:  template.Name,
			Steps: steps,
		}

		resp, err := c.AddJobSetTemplate(ctx, &pbc.AddJobSetTemplateReq{Jst: jst})
		if err != nil {
			return fmt.Errorf("could not add job set template %s: %v", template.Name, err)
		}

		if resp.Success {
			fmt.Printf("job set template %s successfully registered\n", template.Name)
		} else {
			fmt.Printf("error registering job set template %s: %s\n", template.Name, resp.ErrorMsg)
		}
	}

	return nil
}

func buildStepTemplates(jstSteps []parser.PeridotJSTStep) ([]*pbc.StepTemplate, error) {
	stepTemplates := []*pbc.StepTemplate{}

	for _, jstStep := range jstSteps {
		newStep := &pbc.StepTemplate{}
		switch jstStep.TypeStr {
		case "agent":
			newStep.S = &pbc.StepTemplate_Agent{
				Agent: &pbc.StepAgentTemplate{Name: jstStep.Name},
			}
			stepTemplates = append(stepTemplates, newStep)
		case "jobset":
			newStep.S = &pbc.StepTemplate_Jobset{
				Jobset: &pbc.StepJobSetTemplate{Name: jstStep.Name},
			}
			stepTemplates = append(stepTemplates, newStep)
		case "concurrent":
			subSteps, err := buildStepTemplates(jstStep.Steps)
			if err != nil {
				return nil, err
			}
			newStep.S = &pbc.StepTemplate_Concurrent{
				Concurrent: &pbc.StepConcurrentTemplate{Steps: subSteps},
			}
			stepTemplates = append(stepTemplates, newStep)
		default:
			return nil, fmt.Errorf("invalid step type %s in request", jstStep.TypeStr)
		}
	}

	return stepTemplates, nil
}
