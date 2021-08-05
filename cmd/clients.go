//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"fmt"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/spf13/cobra"
)

type clientOptions struct {
	trustedClientID string
}

func init() {
	rootCmd.AddCommand(clientListCmd())
}

//
//
func clientListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clients",
		Short: "Get trusted clients",
		Long:  `Get trusted clients`,
		Example: `
	privx-cli clients [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return trustedClients()
		},
	}

	cmd.AddCommand(clientCreateCmd())
	cmd.AddCommand(clientShowCmd())
	cmd.AddCommand(clientDeleteCmd())
	cmd.AddCommand(clientUpdateCmd())

	return cmd
}

func trustedClients() error {
	api := userstore.New(curl())

	clients, err := api.TrustedClients()
	if err != nil {
		return err
	}

	return stdout(clients)
}

//
//
func clientCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new trusted-client",
		Long:  `Create new trusted client`,
		Example: `
	privx-cli clients create [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return clientCreate(args)
		},
	}

	return cmd
}

func clientCreate(args []string) error {
	var trustedClient userstore.TrustedClient
	api := userstore.New(curl())

	err := decodeJSON(args[0], &trustedClient)
	if err != nil {
		return err
	}

	id, err := api.CreateTrustedClient(trustedClient)
	if err != nil {
		return err
	}

	return stdout(id)
}

//
//
func clientShowCmd() *cobra.Command {
	options := clientOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get trusted client by ID",
		Long:  `Get trusted client by ID. Client ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli clients show [access flags] --id <TRUSTED-CLIENT-ID>,<TRUSTED-CLIENT-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return clientShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.trustedClientID, "id", "", "trusted client ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func clientShow(options clientOptions) error {
	api := userstore.New(curl())
	clients := []userstore.TrustedClient{}

	for _, id := range strings.Split(options.trustedClientID, ",") {
		client, err := api.TrustedClient(id)
		if err != nil {
			return err
		}
		clients = append(clients, *client)
	}

	return stdout(clients)
}

//
//
func clientDeleteCmd() *cobra.Command {
	options := clientOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete trusted client",
		Long:  `Delete trusted client. Client ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli clients delete [access flags] --id <TRUSTED-CLIENT-ID>,<TRUSTED-CLIENT-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return clientDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.trustedClientID, "id", "", "trusted client ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func clientDelete(options clientOptions) error {
	api := userstore.New(curl())

	for _, id := range strings.Split(options.trustedClientID, ",") {
		err := api.DeleteTrustedClient(id)
		if err != nil {
			return err
		} else {
			fmt.Println(id)
		}
	}

	return nil
}

//
//
func clientUpdateCmd() *cobra.Command {
	options := clientOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update trusted client",
		Long:  `Update trusted client`,
		Example: `
	privx-cli clients update [access flags] --id <TRUSTED-CLIENT-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return clientUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.trustedClientID, "id", "", "trusted client ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func clientUpdate(options clientOptions, args []string) error {
	var trustedClient userstore.TrustedClient
	api := userstore.New(curl())

	err := decodeJSON(args[0], &trustedClient)
	if err != nil {
		return err
	}

	err = api.UpdateTrustedClient(options.trustedClientID, &trustedClient)
	if err != nil {
		return err
	}

	return nil
}
