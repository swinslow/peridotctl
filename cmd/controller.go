// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	pbc "github.com/swinslow/peridot-core/pkg/controller"
)

func init() {
	var cmdController = &cobra.Command{
		Use:   "controller",
		Short: "Manage peridot controller",
		Long: `Manage the overall functionality of the
	peridot controller, such as starting and stopping
	it, and getting its current status.`,
		Run: controllerStatus,
	}

	var cmdControllerStatus = &cobra.Command{
		Use:   "status",
		Short: "Get peridot controller status",
		Long: `Get the current status, health, output and
	error messages for the peridot controller.`,
		Run: controllerStatus,
	}

	rootCmd.AddCommand(cmdController)
	cmdController.AddCommand(cmdControllerStatus)
}

func controllerStatus(cmd *cobra.Command, args []string) {
	var ctx context.Context
	var cancel context.CancelFunc

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	resp, err := c.GetStatus(ctx, &pbc.GetStatusReq{})
	if err != nil {
		log.Fatalf("could not get status: %v", err)
	}

	fmt.Printf("status: %s\n", resp.RunStatus.String())
	fmt.Printf("health: %s\n", resp.HealthStatus.String())
	fmt.Printf("output: %s\n", resp.OutputMsg)
	fmt.Printf("errors: %s\n", resp.ErrorMsg)

	conn.Close()
}
