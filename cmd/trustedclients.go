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
	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/spf13/cobra"
)

type extenderOptions struct {
	accessGroupID   string
	extenderID      string
	fileName        string
	clientType      string
	trustedClientID string
}

func (m extenderOptions) normalizeClientType() string {
	return strings.ToUpper(m.clientType)
}

func init() {
	rootCmd.AddCommand(extenderCmd())
}

//
//
func extenderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "trusted-clients",
		Short:        "List trusted clients and donwload pre configs",
		Long:         `List trusted clients and donwload pre configs`,
		SilenceUsage: true,
	}

	cmd.AddCommand(caListCmd())
	cmd.AddCommand(caShowCmd())
	cmd.AddCommand(revocationListCmd())
	cmd.AddCommand(extenderListCmd())
	cmd.AddCommand(extenderShowCmd())
	cmd.AddCommand(preconfigurationDownloadCmd())

	return cmd
}

//
//
func extenderListCmd() *cobra.Command {
	options := extenderOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List trusted clients",
		Long:  `List trusted clients (extender | web-proxy | carrier)`,
		Example: `
	privx-cli extender [access flags] --type extender | webproxy | carrier
	privx-cli extender [access flags] --group-id <ACCESS-GROUP-ID> --type extender | webproxy | carrier
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return extenderList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.clientType, "type", "", " trusted client type")
	cmd.MarkFlagRequired("type")

	return cmd
}

func extenderList(options extenderOptions) error {
	var clients []userstore.TrustedClient
	api := userstore.New(curl())

	res, err := api.TrustedClients()
	if err != nil {
		return err
	}

	switch options.clientType {
	case "extender":
		clients = extenderListHelper(res, options.normalizeClientType())
	case "webproxy":
		clients = extenderListHelper(res, "ICAP")
	case "carrier":
		clients = extenderListHelper(res, options.normalizeClientType())
	default:
		return fmt.Errorf("client type does not exist: %s", options.clientType)
	}

	return stdout(clients)
}

func extenderListHelper(trustedClients []userstore.TrustedClient, clientType string) []userstore.TrustedClient {
	clients := []userstore.TrustedClient{}

	for _, client := range trustedClients {
		if client.Type == userstore.ClientType(clientType) {
			clients = append(clients, client)
		}
	}

	return clients
}

//
//
func extenderShowCmd() *cobra.Command {
	options := extenderOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get trusted client by ID",
		Long:  `Get trusted client by ID`,
		Example: `
	privx-cli extender [access flags] --client-id <TRUSTED-CLIENT-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return extenderShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.trustedClientID, "client-id", "", "trusted client ID")
	cmd.MarkFlagRequired("client-id")

	return cmd
}

func extenderShow(options extenderOptions) error {
	api := userstore.New(curl())

	client, err := api.TrustedClient(options.trustedClientID)
	if err != nil {
		return err
	}

	return stdout(client)
}

//
//
func caListCmd() *cobra.Command {
	options := extenderOptions{}

	cmd := &cobra.Command{
		Use:   "list-ca",
		Short: "List CA certificates",
		Long:  `List CA certificates for extender or web-proxy`,
		Example: `
	privx-cli extender [access flags] --type extender | webproxy
	privx-cli extender [access flags] --group-id <ACCESS-GROUP-ID> --type extender | webproxy
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch options.clientType {
			case "extender":
				extenderCAList(options)
			case "webproxy":
				webproxyCAList(options)
			default:
				return fmt.Errorf("client type does not exist: %s", options.clientType)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.clientType, "type", "", "client type")
	flags.StringVar(&options.accessGroupID, "group-id", "", "access group ID filter")
	cmd.MarkFlagRequired("type")

	return cmd
}

func extenderCAList(options extenderOptions) error {
	api := authorizer.New(curl())

	certificates, err := api.ExtenderCACertificates(options.accessGroupID)
	if err != nil {
		return err
	}

	return stdout(certificates)
}

func webproxyCAList(options extenderOptions) error {
	api := authorizer.New(curl())

	certificates, err := api.WebProxyCACertificates(options.accessGroupID)
	if err != nil {
		return err
	}

	return stdout(certificates)
}

//
//
func caShowCmd() *cobra.Command {
	options := extenderOptions{}

	cmd := &cobra.Command{
		Use:   "show-ca",
		Short: "Get CA certificate",
		Long:  `Get CA certificate for extender or web-proxy`,
		Example: `
	privx-cli extender show [access flags] --client-id <EXTENDER-ID> --type extender | webproxy
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch options.clientType {
			case "extender":
				extenderCAShow(options)
			case "webproxy":
				webproxyCAShow(options)
			default:
				return fmt.Errorf("client type does not exist: %s", options.clientType)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.extenderID, "client-id", "", "trusted client ID")
	flags.StringVar(&options.clientType, "type", "", "client type")
	cmd.MarkFlagRequired("client-id")
	cmd.MarkFlagRequired("type")

	return cmd
}

func extenderCAShow(options extenderOptions) error {
	api := authorizer.New(curl())

	certificate, err := api.ExtenderCACertificate(options.trustedClientID)
	if err != nil {
		return err
	}

	return stdout(certificate)
}

func webproxyCAShow(options extenderOptions) error {
	api := authorizer.New(curl())

	certificate, err := api.WebProxyCACertificate(options.trustedClientID)
	if err != nil {
		return err
	}

	return stdout(certificate)
}

//
//
func revocationListCmd() *cobra.Command {
	options := extenderOptions{}

	cmd := &cobra.Command{
		Use:   "show-crl",
		Short: "Get revocation list",
		Long:  `Get revocation list for extender or web-proxy`,
		Example: `
	privx-cli extender revocation-list [access flags] --client-id <TRUSTED-CLIENT-ID> --type extender | webproxy --name <FILE-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch options.clientType {
			case "extender":
				extenderRevocationList(options)
			case "webproxy":
				webproxyRevocationList(options)
			default:
				return fmt.Errorf("client type does not exist: %s", options.clientType)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.extenderID, "client-id", "", "trusted client ID")
	flags.StringVar(&options.fileName, "name", "", "file name")
	flags.StringVar(&options.clientType, "type", "", "client type")
	cmd.MarkFlagRequired("client-id")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("type")

	return cmd
}

func extenderRevocationList(options extenderOptions) error {
	api := authorizer.New(curl())

	err := api.DownloadExtenderCertificateCRL(options.fileName, options.trustedClientID)
	if err != nil {
		return err
	}

	return nil
}

func webproxyRevocationList(options extenderOptions) error {
	api := authorizer.New(curl())

	err := api.DownloadWebProxyCertificateCRL(options.fileName, options.trustedClientID)
	if err != nil {
		return err
	}

	return nil
}

//
//
func preconfigurationDownloadCmd() *cobra.Command {
	options := extenderOptions{}

	cmd := &cobra.Command{
		Use:   "pre-config",
		Short: "Download a pre-configured config file for extender, webproxy or carrier",
		Long:  `Download a pre-configured config file for extender, webproxy or carrier`,
		Example: `
	privx-cli pre-config [access flags] --client-id <ACCESS-GROUP-ID> --type extender | webproxy | carrier --name <FILE-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return preconfigurationDownloadSwitch(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.trustedClientID, "client-id", "", "trusted client ID")
	flags.StringVar(&options.clientType, "type", "", "trusted client type")
	flags.StringVar(&options.fileName, "name", "", "file name")
	cmd.MarkFlagRequired("client-id")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("name")

	return cmd
}

func preconfigurationDownloadSwitch(options extenderOptions) error {
	switch options.clientType {
	case "extender":
		downloadExtenderPreConf(options)
	case "webproxy":
		downloadWebProxyPreConf(options)
	case "carrier":
		downloadCarrierPreConf(options)
	default:
		return fmt.Errorf("client type does not exist: %s", options.clientType)
	}

	return nil
}
func downloadExtenderPreConf(options extenderOptions) error {
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

func downloadWebProxyPreConf(options extenderOptions) error {
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

func downloadCarrierPreConf(options extenderOptions) error {
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
