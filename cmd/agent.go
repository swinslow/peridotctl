// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"

	pbc "github.com/swinslow/peridot-core/pkg/controller"
	"github.com/swinslow/peridotctl/internal/config"
)

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
		Long: `Register an agent with the peridot controller.

Format: peridotctl agent add NAME URL PORT TYPE [CFGSTRING]

	NAME:      Unique name for agent instance
	URL:       Agent instance hostname (omitting port)
	PORT:      Agent instance port
	TYPE:      Agent instance type (may be repeated with other agents)
	CFGSTRING: Optional: agent configuration values (in format key1:value1;key2:value2;...)`,
		Args: cobra.RangeArgs(4, 5),
		Run:  agentAdd,
	}
	cmdAgent.AddCommand(cmdAgentAdd)

	var cmdAgentGet = &cobra.Command{
		Use:   "get",
		Short: "Get info on registered agent",
		Long: `Get information about an agent that has already been
		registered with the peridot controller.`,
		Args: cobra.ExactArgs(1),
		Run:  agentGet,
	}
	cmdAgent.AddCommand(cmdAgentGet)
}

func agentList(cmd *cobra.Command, args []string) {
	ctx, cancel := config.GetContext(timeout)
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
	ctx, cancel := config.GetContext(timeout)
	defer cancel()
	defer conn.Close()

	name := args[0]
	url := args[1]
	portStr := args[2]
	typeStr := args[3]
	var cfgStr string
	if len(args) == 5 {
		cfgStr = args[4]
	}

	portInt, err := strconv.Atoi(portStr)
	if err != nil || portInt <= 0 {
		log.Fatalf("invalid agent port: %s", portStr)
	}

	// check whether values are okay
	if name == "" {
		log.Fatal("no agent name specified")
	}
	// URL defaulting to "localhost" is acceptable, but empty string isn't
	if url == "" {
		log.Fatal("agent URL cannot be empty string")
	}
	if typeStr == "" {
		log.Fatal("no agent type specified")
	}

	// extract configuration key-value pairs -- semicolons separating pairs,
	// colons separating key from value within a pair
	cfgs := config.ExtractKVs(cfgStr)

	// build into AgentKV list
	kvs := []*pbc.AgentConfig_AgentKV{}
	for k, v := range cfgs {
		kv := pbc.AgentConfig_AgentKV{Key: k, Value: v}
		kvs = append(kvs, &kv)
	}

	// build AgentConfig object
	ac := &pbc.AgentConfig{
		Name: name,
		Url:  url,
		Port: uint32(portInt),
		Type: typeStr,
		Kvs:  kvs,
	}

	resp, err := c.AddAgent(ctx, &pbc.AddAgentReq{Cfg: ac})
	if err != nil {
		log.Fatalf("could not add agent: %v", err)
	}

	if resp.Success {
		fmt.Printf("agent %s successfully registered\n", name)
	} else {
		fmt.Printf("error registering agent %s: %s\n", name, resp.ErrorMsg)
	}
	fmt.Printf("\n")
}

func agentGet(cmd *cobra.Command, args []string) {
	ctx, cancel := config.GetContext(timeout)
	defer cancel()
	defer conn.Close()

	name := args[0]

	resp, err := c.GetAgent(ctx, &pbc.GetAgentReq{Name: name})
	if err != nil {
		log.Fatalf("could not get status for %s: %v", name, err)
	}

	if resp.Success == false {
		log.Fatalf("agent %s not found: %s", name, resp.ErrorMsg)
	}

	agentConfig := resp.Cfg
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
