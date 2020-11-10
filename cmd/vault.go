package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/SSHcom/privx-sdk-go/api/vault"
	"github.com/spf13/cobra"
)

var (
	vaultID      string
	vaultReadTo  []string
	vaultWriteTo []string
)

func init() {
	rootCmd.AddCommand(vaultCmd)

	vaultCmd.AddCommand(vaultCreateCmd)
	vaultCreateCmd.Flags().StringVar(&vaultID, "id", "", "secret identity")
	vaultCreateCmd.Flags().StringArrayVar(&vaultReadTo, "allow-read-to", []string{}, "read by role ID")
	vaultCreateCmd.Flags().StringArrayVar(&vaultWriteTo, "allow-write-to", []string{}, "write by role ID")
	vaultCreateCmd.MarkFlagRequired("id")
	vaultCreateCmd.MarkFlagRequired("read-by")
	vaultCreateCmd.MarkFlagRequired("write-by")

	vaultCmd.AddCommand(vaultUpdateCmd)
	vaultUpdateCmd.Flags().StringVar(&vaultID, "id", "", "secret identity")
	vaultUpdateCmd.Flags().StringArrayVar(&vaultReadTo, "allow-read-to", []string{}, "read by role ID")
	vaultUpdateCmd.Flags().StringArrayVar(&vaultWriteTo, "allow-write-to", []string{}, "write by role ID")
	vaultUpdateCmd.MarkFlagRequired("id")

	vaultCmd.AddCommand(vaultRemoveCmd)
	vaultRemoveCmd.Flags().StringVar(&vaultID, "id", "", "secret identity")
	vaultRemoveCmd.MarkFlagRequired("id")
}

//
//
var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "PrivX secrets",
	Long:  `List and manage PrivX secrets from Vault`,
	Example: `
privx-cli vault [access flags] ID ...
	`,
	SilenceUsage: true,
	RunE:         secrets,
}

func secrets(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires at least one secret name as argument")
	}

	api := vault.New(curl())
	secrets := []vault.Secret{}

	for _, uid := range args {
		bag, err := api.Get(uid)
		if err != nil {
			return err
		}
		secrets = append(secrets, *bag)
	}

	return stdout(secrets)
}

//
//
var vaultCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create new secret",
	Long:  `create new secret to PrivX Vault`,
	Example: `
privx-cli vault create [access flags] --id secret 
	--allow-read-to role-id
	--allow-write-to role-id
	...
	JSON-FILE
	`,
	SilenceUsage: true,
	RunE:         secretCreate,
}

func secretCreate(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires json file as argument")
	}

	secret, err := readJSON(args[0])
	if err != nil {
		return err
	}

	api := vault.New(curl())
	if err := api.Create(vaultID, vaultReadTo, vaultWriteTo, secret); err != nil {
		return err
	}

	return stdout(secret)
}

//
//
var vaultUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update secret",
	Long:  `update existing secret at PrivX Vault`,
	Example: `
privx-cli vault update [access flags] --id secret JSON-FILE

privx-cli vault update [access flags] --id secret
	--allow-read-to role-id
	--allow-write-to role-id
	...
	JSON-FILE
	`,
	SilenceUsage: true,
	RunE:         secretUpdate,
}

func secretUpdate(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires json file as argument")
	}

	secret, err := readJSON(args[0])
	if err != nil {
		return err
	}

	api := vault.New(curl())
	bag, err := api.Get(vaultID)

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

	if err := api.Update(vaultID, vaultReadTo, vaultWriteTo, secret); err != nil {
		return err
	}

	return stdout(secret)
}

//
//
var vaultRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove secret",
	Long:  `remove secret from PrivX Vault`,
	Example: `
privx-cli vault remove [access flags] --id secret
	`,
	SilenceUsage: true,
	RunE:         secretsRemove,
}

func secretsRemove(cmd *cobra.Command, args []string) error {
	api := vault.New(curl())

	if err := api.Remove(vaultID); err != nil {
		return err
	}

	return stdout(vaultID)
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
