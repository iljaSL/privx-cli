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

type componentsOptions struct {
	hostName string
}

func init() {
	rootCmd.AddCommand(componentsListCmd())
}

//
//
func componentsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "components",
		Short: "List and manage privx components",
		Long:  `List and manage privx components`,
		Example: `
	privx-cli components [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return componentsList()
		},
	}

	cmd.AddCommand(componentShowCmd())

	return cmd
}

func componentsList() error {
	api := monitor.New(curl())

	status, err := api.ComponentsStatus()
	if err != nil {
		return err
	}

	return stdout(status)
}

//
//
func componentShowCmd() *cobra.Command {
	options := componentsOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get component status by hostname",
		Long:  `Get component status by hostname`,
		Example: `
	privx-cli components show [access flags] --name <HOST-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return componentShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.hostName, "name", "", "host name")
	cmd.MarkFlagRequired("name")

	return cmd
}

func componentShow(options componentsOptions) error {
	api := monitor.New(curl())

	status, err := api.ComponentStatus(options.hostName)
	if err != nil {
		return err
	}

	return stdout(status)
}
