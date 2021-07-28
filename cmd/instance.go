//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"github.com/SSHcom/privx-sdk-go/api/monitor"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(instanceListCmd())
}

//
//
func instanceListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "List and manage the whole instance",
		Long:  `List and manage the whole instance`,
		Example: `
	privx-cli instance [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return instanceList()
		},
	}

	cmd.AddCommand(instanceTerminateCmd())

	return cmd
}

func instanceList() error {
	api := monitor.New(curl())

	status, err := api.InstanceStatus()
	if err != nil {
		return err
	}

	return stdout(status)
}

//
//
func instanceTerminateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "terminate",
		Short: "Terminate PrivX instances",
		Long:  `Terminate PrivX instances`,
		Example: `
	privx-cli instance terminate [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return instanceTerminate()
		},
	}

	return cmd
}

func instanceTerminate() error {
	api := monitor.New(curl())

	err := api.TerminateInstances()
	if err != nil {
		return err
	}

	return nil
}
