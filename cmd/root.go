//
// Copyright (c) 2020 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"encoding/json"
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
	config  string
	baseURL string
	access  string
	secret  string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "path to config file")
	rootCmd.PersistentFlags().StringVar(&access, "url", "", "PrivX absolute URL (e.g. https://your-instance.privx.io)")
	rootCmd.PersistentFlags().StringVarP(&access, "access", "a", "", "either access key of api client or username.")
	rootCmd.PersistentFlags().StringVarP(&secret, "secret", "s", "", "either secret key of api client or password.")
}

//
//
var rootCmd = &cobra.Command{
	Use:   "privx-cli",
	Short: "PrivX command line client",
	Long:  `PrivX command line client`,
	Example: `
See https://github.com/SSHcom/privx-cli about client configurations

Configure client with environment variables
export PRIVX_API_BASE_URL=https://your-instance.privx.io
export PRIVX_API_ACCESS_KEY=your-username
export PRIVX_API_SECRET_KEY=your-password

Configure client with cli flags
privx-cli --url https://your-instance.privx.io \
	--access your-username \
	--secret your-password
`,
	Run:     root,
	Version: "v1",
}

func root(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func auth() restapi.Authorizer {
	curl := restapi.New(
		restapi.UseConfigFile(config),
		restapi.UseEnvironment(),
	)

	return oauth.With(
		curl,
		oauth.UseConfigFile(config),
		oauth.UseEnvironment(),
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

func stdout(data interface{}) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(encoded)
	return err
}
