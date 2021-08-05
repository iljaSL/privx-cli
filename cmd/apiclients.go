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

type apiClientOptions struct {
	clientID       string
	apiClientRoles string
	name           string
}

func init() {
	rootCmd.AddCommand(apiClientListCmd())
}

//
//
func apiClientListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api-clients",
		Short: "List and manage API clients",
		Long:  `List and manage API clients`,
		Example: `
	privx-cli api-clients [access flags]
		`,
		SilenceUsage: true,
		RunE:         apiClientList,
	}

	cmd.AddCommand(apiClientCreateCmd())
	cmd.AddCommand(apiClientShowCmd())
	cmd.AddCommand(apiClientDeleteCmd())
	cmd.AddCommand(apiClientUpdateCmd())

	return cmd
}

func apiClientList(cmd *cobra.Command, args []string) error {
	api := userstore.New(curl())

	clients, err := api.APIClients()
	if err != nil {
		return err
	}

	return stdout(clients)
}

//
//
func apiClientCreateCmd() *cobra.Command {
	options := apiClientOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new API client",
		Long:  `Create new API client. Role ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli api-clients create [access flags] --name <CLIENT-NAME> --roles <ROLE-ID>,<ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return apiClientCreate(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.name, "name", "", "API client name")
	flags.StringVar(&options.apiClientRoles, "roles", "", "list of roles possessed by the API client")
	cmd.MarkFlagRequired("name")

	return cmd
}

func apiClientCreate(options apiClientOptions) error {
	api := userstore.New(curl())

	id, err := api.CreateAPIClient(options.name, strings.Split(options.apiClientRoles, ","))
	if err != nil {
		return err
	}

	return stdout(id)
}

//
//
func apiClientShowCmd() *cobra.Command {
	options := apiClientOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get API client by ID",
		Long:  `Get API client by ID. API Client ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli api-clients show [access flags] --id <API-CLIENT-ID>,<API-CLIENT-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return apiClientShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.clientID, "id", "", "API client ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func apiClientShow(options apiClientOptions) error {
	api := userstore.New(curl())
	clients := []userstore.APIClient{}

	for _, id := range strings.Split(options.clientID, ",") {
		client, err := api.APIClient(id)
		if err != nil {
			return err
		}
		clients = append(clients, *client)
	}

	return stdout(clients)
}

//
//
func apiClientDeleteCmd() *cobra.Command {
	options := apiClientOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete API client",
		Long:  `Delete a API client`,
		Example: `
	privx-cli api-clients delete [access flags] --id API-CLIENT-ID
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return apiClientDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.clientID, "id", "", "API client ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func apiClientDelete(options apiClientOptions) error {
	api := userstore.New(curl())

	for _, id := range strings.Split(options.clientID, ",") {
		err := api.DeleteAPIClient(id)
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
func apiClientUpdateCmd() *cobra.Command {
	options := apiClientOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update API client",
		Long:  `Update an existing API client`,
		Example: `
	privx-cli api-clients update [access flags] --id <API-CLIENT-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return apiClientUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.clientID, "id", "", "API client ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func apiClientUpdate(options apiClientOptions, args []string) error {
	var apiClient userstore.APIClient
	api := userstore.New(curl())

	err := decodeJSON(args[0], &apiClient)
	if err != nil {
		return err
	}

	err = api.UpdateAPIClient(options.clientID, &apiClient)
	if err != nil {
		return err
	}

	return nil
}
