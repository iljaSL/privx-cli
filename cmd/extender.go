//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"github.com/SSHcom/privx-sdk-go/api/authorizer"
	"github.com/spf13/cobra"
)

type extenderOptions struct {
	accessGroupID string
	extenderID    string
	fileName      string
}

func init() {
	rootCmd.AddCommand(extenderListCmd())
}

//
//
func extenderListCmd() *cobra.Command {
	options := extenderOptions{}

	cmd := &cobra.Command{
		Use:   "extender",
		Short: "List and download extender certificates",
		Long:  `List and download extender certificates`,
		Example: `
	privx-cli extender [access flags] --group-id <ACCESS-GROUP-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return extenderList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.accessGroupID, "group-id", "", "access group ID filter")

	cmd.AddCommand(extenderShowCmd())
	cmd.AddCommand(extenderRevocationListCmd())

	return cmd
}

func extenderList(options extenderOptions) error {
	api := authorizer.New(curl())

	ca, err := api.ExtenderCACertificates(options.accessGroupID)
	if err != nil {
		return err
	}

	return stdout(ca)
}

//
//
func extenderShowCmd() *cobra.Command {
	options := extenderOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get extender CA certificate",
		Long:  `Get extender CA certificate`,
		Example: `
	privx-cli extender show [access flags] --id <EXTENDER-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return extenderShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.extenderID, "id", "", "extender ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func extenderShow(options extenderOptions) error {
	api := authorizer.New(curl())

	certificate, err := api.ExtenderCACertificate(options.extenderID)
	if err != nil {
		return err
	}

	return stdout(certificate)
}

//
//
func extenderRevocationListCmd() *cobra.Command {
	options := extenderOptions{}

	cmd := &cobra.Command{
		Use:   "revocation-list",
		Short: "Get extender revocation list",
		Long:  `Get extender revocation list`,
		Example: `
	privx-cli extender revocation-list [access flags] --id <EXTENDER-ID> --name <FILE-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return extenderRevocationList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.extenderID, "id", "", "extender ID")
	flags.StringVar(&options.fileName, "name", "", "file name")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("name")

	return cmd
}

func extenderRevocationList(options extenderOptions) error {
	api := authorizer.New(curl())

	err := api.DownloadExtenderCertificateCRL(options.fileName, options.extenderID)
	if err != nil {
		return err
	}

	return nil
}
