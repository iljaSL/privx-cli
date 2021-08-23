//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"os"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/spf13/cobra"
)

type principalkeyOptions struct {
	roleID string
	keyID  string
}

func init() {
	rootCmd.AddCommand(principalkeyListCmd())
}

//
//
func principalkeyListCmd() *cobra.Command {
	options := principalkeyOptions{}

	cmd := &cobra.Command{
		Use:   "principal-keys",
		Short: "List and manage role's principal key's",
		Long: `List and manage role's principal key's.
Role ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli principal-keys [access flags] --role-id <ROLE-ID>,<ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return principalkeyList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.roleID, "role-id", "", "role ID")
	cmd.MarkFlagRequired("role-id")

	cmd.AddCommand(principalkeyGenerateCmd())
	cmd.AddCommand(principalkeyImportCmd())
	cmd.AddCommand(principalkeyShowCmd())
	cmd.AddCommand(principalkeyDeleteCmd())

	return cmd
}

func principalkeyList(options principalkeyOptions) error {
	api := rolestore.New(curl())
	keys := []rolestore.PrincipalKey{}

	for _, id := range strings.Split(options.roleID, ",") {
		key, err := api.PrincipalKeys(id)
		if err != nil {
			return err
		}
		keys = append(keys, key...)
	}

	return stdout(keys)
}

//
//
func principalkeyGenerateCmd() *cobra.Command {
	options := principalkeyOptions{}

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate new principal key for role",
		Long:  `Generate new principal key for role`,
		Example: `
	privx-cli principal-keys generate [access flags] --role-id <ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return principalkeyGenerate(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.roleID, "role-id", "", "role ID")
	cmd.MarkFlagRequired("role-id")

	return cmd
}

func principalkeyGenerate(options principalkeyOptions) error {
	api := rolestore.New(curl())

	key, err := api.GeneratePrincipalKey(options.roleID)
	if err != nil {
		return err
	}

	return stdout(key)
}

//
//
func principalkeyImportCmd() *cobra.Command {
	options := principalkeyOptions{}

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import new principal key for role",
		Long: `Import new principal key for role.
PEM encoded private key, pkcs#8, RSA, ECDSA and Ed25519 private keys are supported`,
		Example: `
	privx-cli principal-keys import [access flags] --role-id <ROLE-ID> KEY-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return principalkeyImport(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.roleID, "role-id", "", "role ID")
	cmd.MarkFlagRequired("role-id")

	return cmd
}

func principalkeyImport(options principalkeyOptions, args []string) error {
	key, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	newKey := rolestore.PrivateKey{
		PrivateKey: string(key),
	}
	api := rolestore.New(curl())

	id, err := api.ImportPrincipalKey(newKey, options.roleID)
	if err != nil {
		return err
	}

	return stdout(id)
}

//
//
func principalkeyShowCmd() *cobra.Command {
	options := principalkeyOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get role's principal key",
		Long:  `Get role's principal key`,
		Example: `
	privx-cli principal-keys show [access flags] --id <KEY-ID> --role-id <ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return principalkeyShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.roleID, "role-id", "", "role ID")
	flags.StringVar(&options.keyID, "id", "", "key ID")
	cmd.MarkFlagRequired("role-id")
	cmd.MarkFlagRequired("id")

	return cmd
}

func principalkeyShow(options principalkeyOptions) error {
	api := rolestore.New(curl())

	keys, err := api.PrincipalKey(options.roleID, options.keyID)
	if err != nil {
		return err
	}

	return stdout(keys)
}

//
//
func principalkeyDeleteCmd() *cobra.Command {
	options := principalkeyOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a role's principal key",
		Long:  `Delete a role's principal key.`,
		Example: `
	privx-cli principal-keys delete [access flags] --id <KEY_ID> --role-id <ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return principalkeyDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.roleID, "role-id", "", "role ID")
	flags.StringVar(&options.keyID, "id", "", "key ID")
	cmd.MarkFlagRequired("role-id")
	cmd.MarkFlagRequired("id")

	return cmd
}

func principalkeyDelete(options principalkeyOptions) error {
	api := rolestore.New(curl())

	err := api.DeletePrincipalKey(options.roleID, options.keyID)
	if err != nil {
		return err
	}

	return nil
}
