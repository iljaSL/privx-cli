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
	rootCmd.AddCommand(instanceShowCmd())
}

//
//
func instanceShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Show instance status and restart instance",
		Long:  `Show instance status and restart instance`,
		Example: `
	privx-cli instance [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return instanceShow()
		},
	}

	cmd.AddCommand(instanceTerminateCmd())

	return cmd
}

func instanceShow() error {
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
		Use:   "reset",
		Short: "Reset PrivX instances",
		Long:  `Reset PrivX instances`,
		Example: `
	privx-cli instance reset [access flags]
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
