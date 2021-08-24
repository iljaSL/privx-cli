//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/authorizer"
	"github.com/spf13/cobra"
)

type preconfigurationOptions struct {
	caType          string
	trustedClientID string
	fileName        string
}

func init() {
	rootCmd.AddCommand(preconfigurationCmd())
}

//
//
func preconfigurationCmd() *cobra.Command {
	options := preconfigurationOptions{}

	cmd := &cobra.Command{
		Use:   "pre-configurations",
		Short: "Download a pre-configured config file for extender, web proxy or carrier",
		Long:  `Download a pre-configured config file for extender, web proxy or carrier`,
		Example: `
	privx-cli pre-configurations [access flags] --group-id <ACCESS-GROUP-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return preconfiguration(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.trustedClientID, "id", "", "trusted client ID")
	flags.StringVar(&options.caType, "type", "", "ca type: extender, webproxy or carrier")
	flags.StringVar(&options.fileName, "name", "", "file name")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("name")

	return cmd
}

func preconfiguration(options preconfigurationOptions) error {
	switch options.caType {
	case "extender":
		downloadExtenderPreConf(options)
	case "webproxy":
		downloadWebProxyPreConf(options)
	case "carrier":
		downloadCarrierPreConf(options)
	default:
		return fmt.Errorf("ca type does not exist: %s. You have the following ca type options: extender, webproxy or carrier", options.caType)
	}

	return nil
}

func downloadExtenderPreConf(options preconfigurationOptions) error {
	api := authorizer.New(curl())

	handler, err := api.ExtenderConfigDownloadHandle(options.trustedClientID)
	if err != nil {
		return err
	}

	err = api.DownloadExtenderConfig(options.trustedClientID, handler.SessionID, options.fileName)
	if err != nil {
		return err
	}

	return nil
}

func downloadWebProxyPreConf(options preconfigurationOptions) error {
	api := authorizer.New(curl())

	handler, err := api.WebProxySessionDownloadHandle(options.trustedClientID)
	if err != nil {
		return err
	}

	err = api.DownloadWebProxyConfig(options.trustedClientID, handler.SessionID, options.fileName)
	if err != nil {
		return err
	}

	return nil
}

func downloadCarrierPreConf(options preconfigurationOptions) error {
	api := authorizer.New(curl())

	handler, err := api.CarrierConfigDownloadHandle(options.trustedClientID)
	if err != nil {
		return err
	}

	err = api.DownloadCarrierConfig(options.trustedClientID, handler.SessionID, options.fileName)
	if err != nil {
		return err
	}

	return nil
}
