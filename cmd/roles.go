package cmd

import (
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rolesCmd)
}

var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "PrivX roles",
	Long:  `PrivX roles`,
	Example: `
privx-cli roles -a ... -s ... 
	`,
	SilenceUsage: true,
	RunE:         roles,
}

func roles(cmd *cobra.Command, args []string) error {
	store := rolestore.New(curl())
	roles, err := store.GetRoles()
	if err != nil {
		return err
	}

	fmt.Println(roles)
	return nil
}
