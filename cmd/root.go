// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package cmd

import (
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	pbc "github.com/swinslow/peridot-core/pkg/controller"
)

var cfgFile string
var address string
var timeout int

// connection details
var conn *grpc.ClientConn
var c pbc.ControllerClient

var rootCmd = &cobra.Command{
	Use:   "peridotctl",
	Short: "CLI tool for interacting with peridot",
	Long: `peridotctl is a CLI tool that enables interacting
with a peridot controller. It can be used to configure templates,
start new job sets, and get info about running jobs.`,
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute is the root command's execution entry point.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// address on disk for configuration file
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.peridotctl.yaml)")

	// URL (including port) for peridot controller
	rootCmd.PersistentFlags().StringVar(&address, "address", "localhost:8900", "address of peridot controller gRPC server")
	viper.BindPFlag("address", rootCmd.PersistentFlags().Lookup("address"))

	// timeout in seconds to wait for calls; 0 (default) means no timeout
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 0, "timeout in seconds to wait for response to calls")
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))

	fmt.Printf("==> address is %s\n", address)
	dialServer()
	if conn == nil {
		log.Fatalf("could not connect to peridot controller")
	}
}

func initConfig() {
	// check whether config file path is set in flag.
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// search config in home directory with name ".peridotctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".peridotctl")
	}

	// also pull in environment variables, if any detected by viper
	viper.AutomaticEnv()

	// read the config file if we know of one
	err := viper.ReadInConfig()
	if err == nil {
		fmt.Println("Reading from config file: ", viper.ConfigFileUsed())
	}
}

func dialServer() {
	// connect to server
	var err error
	conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("error dialing peridot controller: %v", err)
		return
	}

	// NOTE: each command must close the connection itself when
	// it is done. We cannot defer the Close() here.
	c = pbc.NewControllerClient(conn)
}
