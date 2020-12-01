package cmd

import (
	"errors"
	"os"

	apiConfig "github.com/SSHcom/privx-sdk-go/api/config"
	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(hostsCmd)
	hostsCmd.AddCommand(hostsDeployCmd)
}

//
//
var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "PrivX hosts",
	Long:  `List and manage PrivX hosts`,
	Example: `
privx-cli hosts [access flags]
	`,
	SilenceUsage: true,
	RunE:         hosts,
}

func hosts(cmd *cobra.Command, args []string) error {
	return nil
}

//
//
var hostsDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Creates target hosts deployment config",
	Long:  `Creates target hosts deployment config`,
	Example: `
privx-cli hosts deploy [access flags] Name ...
	`,
	SilenceUsage: true,
	RunE:         deploy,
}

func deploy(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires name of deployment configuration as an argument")
	}
	name := args[0]

	curl := curl()
	store := userstore.New(curl)

	seq, err := store.TrustedClients()
	if err != nil {
		return err
	}

	cli := findClientID(seq, name)
	if cli == "" {
		cli, err = store.CreateTrustedClient(
			userstore.HostProvisioning(name),
		)
		if err != nil {
			return err
		}
	}

	conf := apiConfig.New(curl)
	file, err := conf.ConfigDeploy(cli)
	if err != nil {
		return err
	}

	os.Stdout.Write(file)
	return nil
}

func findClientID(seq []userstore.TrustedClient, name string) string {
	for _, cli := range seq {
		if cli.Name == name {
			return cli.ID
		}
	}
	return ""
}
