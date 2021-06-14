//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/vault"
	"github.com/spf13/cobra"
)

var (
	secretName   string
	keywords     string
	sortkey      string
	vaultReadTo  []string
	vaultWriteTo []string
)

func init() {
	rootCmd.AddCommand(secretListCmd)
	secretListCmd.Flags().IntVar(&offset, "offset", 0, "where to start fetching the items")
	secretListCmd.Flags().IntVar(&limit, "limit", 50, "number of items to return")

	secretListCmd.AddCommand(secretShowCmd)
	secretShowCmd.Flags().StringVar(&secretName, "name", "", "secret name")
	secretShowCmd.MarkFlagRequired("name")

	secretListCmd.AddCommand(secretCreateCmd)
	secretCreateCmd.Flags().StringVar(&secretName, "name", "", "secret name")
	secretCreateCmd.Flags().StringArrayVar(&vaultReadTo, "allow-read-to", []string{}, "read by role ID")
	secretCreateCmd.Flags().StringArrayVar(&vaultWriteTo, "allow-write-to", []string{}, "write by role ID")
	secretCreateCmd.MarkFlagRequired("name")
	secretCreateCmd.MarkFlagRequired("read-by")
	secretCreateCmd.MarkFlagRequired("write-by")

	secretListCmd.AddCommand(vaultUpdateCmd)
	vaultUpdateCmd.Flags().StringVar(&secretName, "name", "", "secret name")
	vaultUpdateCmd.Flags().StringArrayVar(&vaultReadTo, "allow-read-to", []string{}, "read by role ID")
	vaultUpdateCmd.Flags().StringArrayVar(&vaultWriteTo, "allow-write-to", []string{}, "write by role ID")
	vaultUpdateCmd.MarkFlagRequired("name")

	secretListCmd.AddCommand(secretDeleteCmd)
	secretDeleteCmd.Flags().StringVar(&secretName, "name", "", "secret name")
	secretDeleteCmd.MarkFlagRequired("name")

	secretListCmd.AddCommand(secretMetadataShowCmd)
	secretMetadataShowCmd.Flags().StringVar(&secretName, "name", "", "secret name")
	secretMetadataShowCmd.MarkFlagRequired("name")

	secretListCmd.AddCommand(secretSearchCmd)
	secretSearchCmd.Flags().IntVar(&offset, "offset", 0, "where to start fetching the items")
	secretSearchCmd.Flags().IntVar(&limit, "limit", 50, "number of items to return")
	secretSearchCmd.Flags().StringVar(&sortdir, "sortdir", "", "sort direction, ASC or DESC (default ASC)")
	secretSearchCmd.Flags().StringVar(&sortkey, "sortkey", "", "sort object by name, updated, or created.")
	secretSearchCmd.Flags().StringVar(&keywords, "keywords", "", "comma or space-separated string to search in secret's names")

	secretListCmd.AddCommand(secretSchemasShowCmd)
}

//
//
var secretListCmd = &cobra.Command{
	Use:   "secrets",
	Short: "PrivX secrets",
	Long:  `List and manage PrivX secrets`,
	Example: `
privx-cli secrets [access flags] --offset <OFFSET> --limit <LIMIT>
	`,
	SilenceUsage: true,
	RunE:         secretList,
}

func secretList(cmd *cobra.Command, args []string) error {
	api := vault.New(curl())

	secrets, err := api.Secrets(offset, limit)
	if err != nil {
		return err
	}

	return stdout(secrets)
}

//
//
var secretShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Get a secret",
	Long:  `Get a secret`,
	Example: `
privx-cli secrets show [access flags] --name <SECRET-NAME>
	`,
	SilenceUsage: true,
	RunE:         secretShow,
}

func secretShow(cmd *cobra.Command, args []string) error {
	api := vault.New(curl())

	secret, err := api.Secret(secretName)
	if err != nil {
		return err
	}

	return stdout(secret)
}

//
//
var secretCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new secret",
	Long:  `Create new secret`,
	Example: `
privx-cli secrets create [access flags] --name <SECRET-NAME> 
	--allow-read-to <ROLE-ID>
	--allow-write-to <ROLE-ID>
	...
	JSON-FILE
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         secretCreate,
}

func secretCreate(cmd *cobra.Command, args []string) error {
	secret, err := readJSON(args[0])
	if err != nil {
		return err
	}

	api := vault.New(curl())
	if err := api.CreateSecret(secretName, vaultReadTo, vaultWriteTo, secret); err != nil {
		return err
	}

	return stdout(secret)
}

//
//
var vaultUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update secret",
	Long:  `Update existing secret`,
	Example: `
privx-cli secrets update [access flags] --name <SECRET-NAME> JSON-FILE

privx-cli secrets update [access flags] --name <SECRET-NAME>
	--allow-read-to <ROLE-ID>
	--allow-write-to <ROLE-ID>
	...
	JSON-FILE
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         secretUpdate,
}

func secretUpdate(cmd *cobra.Command, args []string) error {
	secret, err := readJSON(args[0])
	if err != nil {
		return err
	}

	api := vault.New(curl())
	bag, err := api.Secret(secretName)
	if err != nil {
		return err
	}

	if len(vaultReadTo) == 0 {
		for _, ref := range bag.AllowRead {
			vaultReadTo = append(vaultReadTo, ref.ID)
		}
	}

	if len(vaultWriteTo) == 0 {
		for _, ref := range bag.AllowWrite {
			vaultWriteTo = append(vaultWriteTo, ref.ID)
		}
	}

	if err := api.UpdateSecret(secretName, vaultReadTo, vaultWriteTo, secret); err != nil {
		return err
	}

	return stdout(secret)
}

//
//
var secretDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete secret",
	Long:  `Delete secret from PrivX Vault`,
	Example: `
privx-cli vault delete [access flags] --name <SECRET-NAME>
	`,
	SilenceUsage: true,
	RunE:         secretDelete,
}

func secretDelete(cmd *cobra.Command, args []string) error {
	api := vault.New(curl())

	err := api.DeleteSecret(secretName)

	return err
}

//
//
var secretMetadataShowCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Get a secrets metadata",
	Long:  `Get a secrets metadata`,
	Example: `
privx-cli secrets metadata [access flags] --name <SECRET-NAME>
	`,
	SilenceUsage: true,
	RunE:         secretMetadataShow,
}

func secretMetadataShow(cmd *cobra.Command, args []string) error {
	api := vault.New(curl())

	meta, err := api.SecretMetadata(secretName)
	if err != nil {
		return err
	}

	return stdout(meta)
}

//
//
var secretSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for secrets",
	Long:  `Search for secrets`,
	Example: `
privx-cli secrets search [access flags] --offset <OFFSET> --keywords "<KEYWORD>,<KEYWORD>"
privx-cli secrets search [access flags] --limit <LIMIT> --sortkey <SORTKEY>
	`,
	SilenceUsage: true,
	RunE:         secretSearch,
}

func secretSearch(cmd *cobra.Command, args []string) error {
	api := vault.New(curl())

	secret, err := api.SearchSecrets(offset, limit, keywords, sortkey, strings.ToUpper(sortdir))
	if err != nil {
		return err
	}

	return stdout(secret)
}

//
//
var secretSchemasShowCmd = &cobra.Command{
	Use:   "schemas",
	Short: "Returns the defined schemas",
	Long:  `Returns the defined schemas`,
	Example: `
privx-cli secrets schemas [access flags]
	`,
	SilenceUsage: true,
	RunE:         secretSchemasShow,
}

func secretSchemasShow(cmd *cobra.Command, args []string) error {
	api := vault.New(curl())

	schemas, err := api.VaultSchemas()
	if err != nil {
		return err
	}

	return stdout(schemas)
}

func readJSON(name string) (secret interface{}, err error) {
	file, err := os.Open(name)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &secret)
	return
}
