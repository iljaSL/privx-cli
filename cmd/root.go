package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/SSHcom/privx-sdk-go/oauth"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/spf13/cobra"
)

// Execute is entry point to application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		e := err.Error()
		fmt.Println(strings.ToUpper(e[:1]) + e[1:])
		os.Exit(1)
	}
}

var (
	config string
	access string
	secret string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&access, "config", "c",
		"",
		"path to config file")

	rootCmd.PersistentFlags().StringVarP(
		&access, "access", "a",
		"",
		"either access key of api client or username.")

	rootCmd.PersistentFlags().StringVarP(
		&access, "secret", "s",
		"",
		"either secret key of api client or password.")
}

//
//
var rootCmd = &cobra.Command{
	Use:   "privx-cli",
	Short: "PrivX command line client",
	Long: `PrivX command line client

`,
	Run:     root,
	Version: "v0",
}

func root(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func auth() restapi.Authorizer {
	curl := restapi.New(
		restapi.UseConfigFile(config),
		restapi.UseEnvironment(),
		restapi.Verbose(),
	)

	return oauth.WithCredential(
		curl,
		oauth.UseConfigFile(config),
		oauth.Access(access),
		oauth.Secret(secret),
	)
}

func curl() restapi.Connector {
	return restapi.New(
		restapi.Auth(auth()),
		restapi.UseConfigFile(config),
		restapi.UseEnvironment(),
	)
}
