// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	pbc "github.com/swinslow/peridot-core/pkg/controller"
	"github.com/swinslow/peridotctl/internal/config"
)

// vars for AddAgent
var vAddAgentName string
var vAddAgentURL string
var vAddAgentPort uint32
var vAddAgentType string
var vAddAgentKvs string

func init() {
	var cmdAgent = &cobra.Command{
		Use:   "agent",
		Short: "Manage agent registrations",
		Long: `Manage agent registrations with the peridot controller;
		list known agents, get information on a particular agent, and
		add a new agent registration.`,
		//Run: agentList,
	}
	rootCmd.AddCommand(cmdAgent)

	var cmdAgentList = &cobra.Command{
		Use:   "list",
		Short: "Get all registered agents",
		Long: `Get information about all agents registered with the
		peridot controller.`,
		Run: agentList,
	}
	cmdAgent.AddCommand(cmdAgentList)

	var cmdAgentAdd = &cobra.Command{
		Use:   "add",
		Short: "Add new registered agent",
		Long:  `Register an agent with the peridot controller.`,
		Run:   agentAdd,
	}
	cmdAgent.AddCommand(cmdAgentAdd)
	cmdAgentAdd.Flags().StringVarP(&vAddAgentName, "name", "n", "", "Unique name for agent instance")
	cmdAgentAdd.Flags().StringVarP(&vAddAgentURL, "url", "u", "localhost", "Agent instance hostname (omitting port)")
	cmdAgentAdd.Flags().Uint32VarP(&vAddAgentPort, "port", "p", 0, "Agent instance port")
	cmdAgentAdd.Flags().StringVarP(&vAddAgentType, "type", "t", "", "Agent instance type (may be repeated with other agents)")
	cmdAgentAdd.Flags().StringVarP(&vAddAgentKvs, "cfg", "c", "", "Agent configuration values (in format key1:value1;key2:value2;...)")

}

func agentList(cmd *cobra.Command, args []string) {
	var ctx context.Context
	var cancel context.CancelFunc

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()
	defer conn.Close()

	resp, err := c.GetAllAgents(ctx, &pbc.GetAllAgentsReq{})
	if err != nil {
		log.Fatalf("could not get status: %v", err)
	}

	fmt.Printf("Registered agents:\n\n")

	for _, agentConfig := range resp.Cfgs {
		fmt.Printf("name: %s\n", agentConfig.Name)
		fmt.Printf("url: %s\n", agentConfig.Url)
		fmt.Printf("port: %d\n", agentConfig.Port)
		fmt.Printf("type: %s\n", agentConfig.Type)
		fmt.Printf("Key-value configs:\n")
		for _, kv := range agentConfig.Kvs {
			fmt.Printf("  %s: %s\n", kv.Key, kv.Value)
		}
		fmt.Printf("\n")
	}
}

func agentAdd(cmd *cobra.Command, args []string) {
	var ctx context.Context
	var cancel context.CancelFunc

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()
	defer conn.Close()

	// check whether values are okay
	if vAddAgentName == "" {
		log.Fatal("no agent name specified")
	}
	// URL defaulting to "localhost" is acceptable, but empty string isn't
	if vAddAgentURL == "" {
		log.Fatal("agent URL cannot be empty string")
	}
	if vAddAgentPort == 0 {
		log.Fatal("no agent port specified")
	}
	if vAddAgentType == "" {
		log.Fatal("no agent type specified")
	}

	// extract configuration key-value pairs -- semicolons separating pairs,
	// colons separating key from value within a pair
	cfgs := config.ExtractKVs(vAddAgentKvs)

	// build into AgentKV list
	kvs := []*pbc.AgentConfig_AgentKV{}
	for k, v := range cfgs {
		kv := pbc.AgentConfig_AgentKV{Key: k, Value: v}
		kvs = append(kvs, &kv)
	}

	// build AgentConfig object
	ac := &pbc.AgentConfig{
		Name: vAddAgentName,
		Url:  vAddAgentURL,
		Port: vAddAgentPort,
		Type: vAddAgentType,
		Kvs:  kvs,
	}

	resp, err := c.AddAgent(ctx, &pbc.AddAgentReq{Cfg: ac})
	if err != nil {
		log.Fatalf("could not add agent: %v", err)
	}

	if resp.Success {
		fmt.Printf("agent %s successfully registered\n", vAddAgentName)
	} else {
		fmt.Printf("error registering agent %s: %s\n", vAddAgentName, resp.ErrorMsg)
	}
	fmt.Printf("\n")
}
