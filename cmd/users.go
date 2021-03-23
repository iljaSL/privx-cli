//
// Copyright (c) 2020 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"errors"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/spf13/cobra"
)

var (
	userQuery      []string
	userID         string
	userRoleGrant  []string
	userRoleRevoke []string
)

func init() {
	rootCmd.AddCommand(usersCmd)
	usersCmd.Flags().StringArrayVarP(&userQuery, "query", "q", []string{}, "query PrivX users with keyword")

	usersCmd.AddCommand(usersInfoCmd)

	usersCmd.AddCommand(usersRolesCmd)
	usersRolesCmd.Flags().StringVar(&userID, "uid", "", "user unique id")
	usersRolesCmd.Flags().StringArrayVar(&userRoleGrant, "grant", []string{}, "grant role to user, requires role unique id.")
	usersRolesCmd.Flags().StringArrayVar(&userRoleRevoke, "revoke", []string{}, "revoke role from user, requires role unique id.")
	usersRolesCmd.MarkFlagRequired("uid")
}

//
//
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "PrivX users",
	Long:  `List and manage PrivX users`,
	Example: `
privx-cli users [access flags]
	`,
	SilenceUsage: true,
	RunE:         users,
}

func users(cmd *cobra.Command, args []string) error {
	store := rolestore.New(curl())
	users, err := store.SearchUsers(strings.Join(userQuery, ","), "")
	if err != nil {
		return err
	}

	return stdout(users)
}

//
//
var usersInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Description about PrivX user",
	Long:  `Description about PrivX user`,
	Example: `
privx-cli users info [access flags] UID ...
	`,
	SilenceUsage: true,
	RunE:         info,
}

func info(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires at least one user id as argument")
	}

	store := rolestore.New(curl())
	users := []rolestore.User{}

	for _, uid := range args {
		user, err := store.User(uid)
		if err != nil {
			return err
		}
		users = append(users, *user)
	}

	return stdout(users)
}

//
//
var usersRolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "user roles",
	Long:  `list and manage user roles`,
	Example: `
privx-cli users roles [access flags] --uid UID
privx-cli users roles [access flags] --uid UID --grant role-uid
privx-cli users roles [access flags] --uid UID --revoke role-uid
	`,
	SilenceUsage: true,
	RunE:         userRoles,
}

func userRoles(cmd *cobra.Command, args []string) error {
	store := rolestore.New(curl())

	for _, role := range userRoleGrant {
		err := store.AddUserRole(userID, role)
		if err != nil {
			return err
		}
	}

	for _, role := range userRoleRevoke {
		err := store.RemoveUserRole(userID, role)
		if err != nil {
			return err
		}
	}

	roles, err := store.UserRoles(userID)
	if err != nil {
		return err
	}
	return stdout(roles)
}
