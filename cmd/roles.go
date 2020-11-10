package cmd

import (
	"errors"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rolesCmd)

	rolesCmd.AddCommand(rolesMembersCmd)
}

//
//
var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "PrivX roles",
	Long:  `List and manage PrivX roles`,
	Example: `
privx-cli roles [access flags]
	`,
	SilenceUsage: true,
	RunE:         roles,
}

func roles(cmd *cobra.Command, args []string) error {
	store := rolestore.New(curl())
	roles, err := store.Roles()
	if err != nil {
		return err
	}

	return stdout(roles)
}

//
//
var rolesMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "Get members of PrivX role",
	Long:  `Get members of PrivX role`,
	Example: `
privx-cli roles members [access flags] UID ... 
	`,
	SilenceUsage: true,
	RunE:         members,
}

func members(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires at least one roles id as argument")
	}

	store := rolestore.New(curl())
	users := []rolestore.User{}

	for _, role := range args {
		seq, err := store.GetRoleMembers(role)
		if err != nil {
			return err
		}
		users = append(users, seq...)
	}

	return stdout(users)
}
