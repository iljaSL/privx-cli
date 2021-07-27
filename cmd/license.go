//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"github.com/SSHcom/privx-sdk-go/api/licensemanager"
	"github.com/spf13/cobra"
)

type licenseOptions struct {
	licenseKey string
	optin      bool
}

func init() {
	rootCmd.AddCommand(licenseListCmd())
}

//
//
func licenseListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "license",
		Short: "List and manage PrivX license",
		Long:  `List and manage PrivX license keys`,
		Example: `
	privx-cli license [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return licenseList()
		},
	}

	cmd.AddCommand(licenseSetCmd())
	cmd.AddCommand(licenseRefreshCmd())
	cmd.AddCommand(licenseStatisticsSetCmd())
	cmd.AddCommand(licenseDeactivateCmd())

	return cmd
}

func licenseList() error {
	store := licensemanager.New(curl())

	license, err := store.License()
	if err != nil {
		return err
	}

	return stdout(license)
}

//
//
func licenseSetCmd() *cobra.Command {
	options := licenseOptions{}

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set new license",
		Long:  `Set new license`,
		Example: `
	privx-cli license set [access flags] --key <LICENSE-KEY>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return licenseSet(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.licenseKey, "key", "", "PrivX license key")

	return cmd
}

func licenseSet(options licenseOptions) error {
	api := licensemanager.New(curl())

	err := api.SetLicense(options.licenseKey)
	if err != nil {
		return err
	}

	return err
}

//
//
func licenseRefreshCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "Refresh license info",
		Long:  `Refresh for PrivX license infos and apply any new changes found in your subscription`,
		Example: `
	privx-cli license refresh [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return licenseRefresh()
		},
	}

	return cmd
}

func licenseRefresh() error {
	api := licensemanager.New(curl())

	license, err := api.RefreshLicense()
	if err != nil {
		return err
	}

	return stdout(license)
}

//
//
func licenseStatisticsSetCmd() *cobra.Command {
	options := licenseOptions{}

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Update license statistics",
		Long:  `Update PrivX license statistics`,
		Example: `
	privx-cli license stats [access flags]
	privx-cli license stats [access flags] --optin=false 
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return licenseStatisticsSet(options)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&options.optin, "optin", "", true, "enable or disable license statistics")

	return cmd
}

func licenseStatisticsSet(options licenseOptions) error {
	api := licensemanager.New(curl())

	err := api.SetLicenseStatistics(options.optin)

	return err
}

//
//
func licenseDeactivateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deactivate",
		Short: "Deactivate license",
		Long:  `Deactivate PrivX license`,
		Example: `
	privx-cli license deactivate [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return licenseDeactivate()
		},
	}

	return cmd
}

func licenseDeactivate() error {
	api := licensemanager.New(curl())

	err := api.DeactivateLicense()

	return err
}
