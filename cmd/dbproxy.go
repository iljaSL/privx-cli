//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"github.com/SSHcom/privx-sdk-go/api/dbproxy"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dbProxyCmd())
}

func dbProxyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db-proxy",
		Short: "Db-proxy service status and config commands",
		Long:  `Show status and configuration for db-proxy service`,
		Example: `
	privx-cli db-proxy status [access flags]
	privx-cli db-proxy config [access flags]
		`,
		SilenceUsage: true,
	}

	cmd.AddCommand(dbproxyStatusCmd())
	cmd.AddCommand(dbproxyConfCmd())

	return cmd
}

//
//
func dbproxyStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Get db-proxy microservice status",
		Long:  `Get db-proxy microservice status. The service status could be one of STOPPED, INITIALIZING, AUTHORIZING, ERROR`,
		Example: `
	privx-cli db-proxy status [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return dbproxyStatus()
		},
	}

	return cmd
}

func dbproxyStatus() error {
	api := dbproxy.New(curl())
	status, err := api.DbProxyStatus()
	if err != nil {
		return err
	}

	return stdout(status)
}

//
//
func dbproxyConfCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Get DB proxy configuration",
		Long:  `Get DB proxy configuration. Includes info about ca_certificate and certificate_chain`,
		Example: `
	privx-cli db-proxy config [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return dbproxyConf()
		},
	}

	return cmd
}

func dbproxyConf() error {
	api := dbproxy.New(curl())
	status, err := api.DbProxyConf()
	if err != nil {
		return err
	}

	return stdout(status)
}
