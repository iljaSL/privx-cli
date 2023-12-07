//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"errors"
	"os"

	authApi "github.com/SSHcom/privx-sdk-go/api/auth"
	"github.com/SSHcom/privx-sdk-go/api/licensemanager"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mobilegwCmd())
}

func mobilegwCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "mobilegw",
		Short:        "Manage mobilegw registration",
		Long:         `Manage mobilegw registration. Manage user paired devices.`,
		Example:      `privx-cli mobilegw [access flags]`,
		SilenceUsage: true,
	}

	cmd.AddCommand(registerToMobileGwCmd())
	cmd.AddCommand(unregisterToMobileGWCmd())
	cmd.AddCommand(getMobileGwRegistrationCmd())
	cmd.AddCommand(getUserPairedDevicesCmd())
	cmd.AddCommand(unpairUserDeviceCmd())

	return cmd
}

func registerToMobileGwCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "register",
		Short:        "Register",
		Long:         `Register PrivX to the MobileGW`,
		Example:      `privx-cli mobilegw register`,
		SilenceUsage: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return registerToMobileGw()
		},
	}

	return cmd
}

func registerToMobileGw() error {
	client := licensemanager.New(curl())

	err := client.RegisterToMobileGW()
	if err != nil {
		os.Stdout.WriteString("registration failed")
		return err
	}

	os.Stdout.WriteString("registration success")
	return nil
}

func unregisterToMobileGWCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "unregister",
		Short:        "Unregister",
		Long:         `Unregister PrivX from the MobileGW`,
		Example:      `privx-cli mobilegw unregister`,
		SilenceUsage: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return unregisterFromMobileGw()
		},
	}

	return cmd
}

func unregisterFromMobileGw() error {
	client := licensemanager.New(curl())

	err := client.UnregisterToMobileGW()
	if err != nil {
		return err
	}

	return nil
}

func getMobileGwRegistrationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "regstat",
		Short:        "Registration status",
		Long:         `Get PrivX registration status to the MobileGW`,
		Example:      `privx-cli mobilegw regstat`,
		SilenceUsage: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getMobileGwRegistration()
		},
	}

	return cmd
}

func getMobileGwRegistration() error {
	client := licensemanager.New(curl())

	status, err := client.GetMobileGwRegistration()
	if err != nil {
		return err
	}

	return stdout(status)
}

func getUserPairedDevicesCmd() *cobra.Command {
	var userId string

	cmd := &cobra.Command{
		Use:          "paired-devices",
		Short:        "List Paired Devices",
		Long:         `List paired devices of a user`,
		Example:      `privx-cli mobilegw paired-devices --user-id`,
		SilenceUsage: false,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if userId == "" {
				return errors.New("user-id is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return getUserPairedDevices(userId)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&userId, "user-id", "", "User ID")

	return cmd
}

func getUserPairedDevices(userId string) error {
	client := authApi.New(curl())

	devices, err := client.GetUserPairedDevices(userId)
	if err != nil {
		return err
	}

	return stdout(devices)
}

func unpairUserDeviceCmd() *cobra.Command {
	var userId, deviceId string

	cmd := &cobra.Command{
		Use:          "unpair-device",
		Short:        "Unpair Device",
		Long:         `Unaired a user's devices`,
		Example:      `privx-cli mobilegw unpair-device --user-id --device-id`,
		SilenceUsage: false,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if userId == "" || deviceId == "" {
				return errors.New("user-id and device-id are required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return unpairUserDevice(userId, deviceId)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&userId, "user-id", "", "User ID")
	flags.StringVar(&deviceId, "device-id", "", "Device ID")

	return cmd
}

func unpairUserDevice(userId, deviceId string) error {
	client := authApi.New(curl())

	err := client.UnpairUserDevice(userId, deviceId)
	if err != nil {
		return err
	}

	return nil
}
