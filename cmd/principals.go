//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"fmt"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/authorizer"
	"github.com/spf13/cobra"
)

type principalOptions struct {
	groupID string
	keyID   string
	filter  string
}

func init() {
	rootCmd.AddCommand(principalListCmd())
}

//
//
func principalListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "principals",
		Short: "List and manage defined principals",
		Long:  `List and manage defined principals`,
		Example: `
	privx-cli principals [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return principalList()
		},
	}

	cmd.AddCommand(principalShowCmd())
	cmd.AddCommand(principalDeleteCmd())
	cmd.AddCommand(principalCreateCmd())
	cmd.AddCommand(principalImportCmd())
	cmd.AddCommand(principalSignCmd())

	return cmd
}

func principalList() error {
	api := authorizer.New(curl())

	principals, err := api.Principals()
	if err != nil {
		return err
	}

	return stdout(principals)
}

//
//
func principalShowCmd() *cobra.Command {
	options := principalOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get the principal key by its group ID",
		Long:  `Get the principal key by its group ID`,
		Example: `
	privx-cli principals show [access flags] --id <GROUP-ID>
	privx-cli principals show [access flags] --id <GROUP-ID> --key-id <KEY-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return principalShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.groupID, "id", "", "principal group ID")
	flags.StringVar(&options.keyID, "key-id", "", "request specific principal key")
	flags.StringVar(&options.filter, "filter", "", "if filter=all then all principal keys are returned")
	cmd.MarkFlagRequired("id")

	return cmd
}

func principalShow(options principalOptions) error {
	api := authorizer.New(curl())

	key, err := api.Principal(options.groupID, options.keyID, options.filter)
	if err != nil {
		return err
	}

	return stdout(key)
}

//
//
func principalDeleteCmd() *cobra.Command {
	options := principalOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes the principal key by its group ID",
		Long:  `Deletes the principal key by its group ID. Group ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli principals delete [access flags] --id <GROUP-ID>,<GROUP-ID>
	privx-cli principals delete [access flags] --id <GROUP-ID> --key-id <KEY-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return principalDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.groupID, "id", "", "group ID")
	flags.StringVar(&options.keyID, "key-id", "", "request specific principal key")
	cmd.MarkFlagRequired("id")

	return cmd
}

func principalDelete(options principalOptions) error {
	api := authorizer.New(curl())

	for _, id := range strings.Split(options.groupID, ",") {
		err := api.DeletePrincipalKey(id, options.keyID)
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
func principalCreateCmd() *cobra.Command {
	options := principalOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create principal key pair",
		Long:  `Create principal key pair`,
		Example: `
	privx-cli principals create [access flags] --id <GROUP-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return principalCreate(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.groupID, "id", "", "group ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func principalCreate(options principalOptions) error {
	api := authorizer.New(curl())

	signature, err := api.CreatePrincipalKey(options.groupID)
	if err != nil {
		return err
	}

	return stdout(signature)
}

//
//
func principalImportCmd() *cobra.Command {
	options := principalOptions{}

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import a principal key pair",
		Long:  `Import a principal key pair`,
		Example: `
	privx-cli principals import [access flags] --id <GROUP-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return principalImport(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.groupID, "id", "", "group ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func principalImport(options principalOptions, args []string) error {
	var importRequest authorizer.PrincipalKeyImportRequest
	api := authorizer.New(curl())

	err := decodeJSON(args[0], &importRequest)
	if err != nil {
		return err
	}

	signature, err := api.ImportPrincipalKey(options.groupID, &importRequest)
	if err != nil {
		return err
	}

	return stdout(signature)
}

//
//
func principalSignCmd() *cobra.Command {
	options := principalOptions{}

	cmd := &cobra.Command{
		Use:   "sign",
		Short: "Get a signature",
		Long:  `Get a signature`,
		Example: `
	privx-cli principals sign [access flags] --id <GROUP-ID> JSON-FILE
	privx-cli principals sign [access flags] --id <GROUP-ID> --key-id <KEY-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return principalSign(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.groupID, "id", "", "group ID")
	flags.StringVar(&options.keyID, "key-id", "", "key ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func principalSign(options principalOptions, args []string) error {
	var signRequest authorizer.Credential
	api := authorizer.New(curl())

	err := decodeJSON(args[0], &signRequest)
	if err != nil {
		return err
	}

	signature, err := api.SignPrincipalKey(options.groupID, options.keyID, &signRequest)
	if err != nil {
		return err
	}

	return stdout(signature)
}
