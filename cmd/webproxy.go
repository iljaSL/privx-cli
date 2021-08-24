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

type webproxyOptions struct {
	accessGroupID   string
	trustedClientID string
	fileName        string
}

func init() {
	rootCmd.AddCommand(webproxyListCmd())
}

//
//
func webproxyListCmd() *cobra.Command {
	options := webproxyOptions{}

	cmd := &cobra.Command{
		Use:   "web-proxy",
		Short: "List and download webproxy certificates/configs",
		Long:  `List and download webproxy certificates/configs`,
		Example: `
	privx-cli web-proxy [access flags] --group-id <ACCESS-GROUP-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return webproxyList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.accessGroupID, "group-id", "", "access group ID filter")

	cmd.AddCommand(webproxyShowCmd())
	cmd.AddCommand(webproxyRevocationListCmd())

	return cmd
}

func webproxyList(options webproxyOptions) error {
	api := authorizer.New(curl())

	ca, err := api.WebProxyCACertificates(options.accessGroupID)
	if err != nil {
		return err
	}

	return stdout(ca)
}

//
//
func webproxyShowCmd() *cobra.Command {
	options := webproxyOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get web-proxy CA certificate",
		Long:  `Get web-proxy CA certificate`,
		Example: `
	privx-cli web-proxy show [access flags] --id <TRUSTED-CLIENT-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return webproxyShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.trustedClientID, "id", "", "trusted client ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func webproxyShow(options webproxyOptions) error {
	api := authorizer.New(curl())

	certificate, err := api.WebProxyCACertificate(options.trustedClientID)
	if err != nil {
		return err
	}

	return stdout(certificate)
}

//
//
func webproxyRevocationListCmd() *cobra.Command {
	options := webproxyOptions{}

	cmd := &cobra.Command{
		Use:   "revocation-list",
		Short: "Get web-proxy revocation list",
		Long:  `Get web-proxy revocation list`,
		Example: `
	privx-cli web-proxy revocation-list [access flags] --id <TRUSTED-CLIENT-ID> --name <FILE-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return webproxyRevocationList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.trustedClientID, "id", "", "trusted client ID")
	flags.StringVar(&options.fileName, "name", "", "file name")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("name")

	return cmd
}

func webproxyRevocationList(options webproxyOptions) error {
	api := authorizer.New(curl())

	err := api.DownloadWebProxyCertificateCRL(options.fileName, options.trustedClientID)
	if err != nil {
		return err
	}

	return nil
}
