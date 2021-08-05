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

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/spf13/cobra"
)

var (
	userID         string
	password       string
	query          string
	offset         int
	limit          int
	userQuery      []string
	userRoleGrant  []string
	userRoleRevoke []string
)

func init() {
	rootCmd.AddCommand(userListCmd)
	userListCmd.Flags().StringArrayVarP(&userQuery, "query", "q", []string{}, "query PrivX users with keyword")

	userListCmd.AddCommand(userShowCmd)

	userListCmd.AddCommand(usersRolesCmd)
	usersRolesCmd.Flags().StringVar(&userID, "uid", "", "user unique id")
	usersRolesCmd.Flags().StringArrayVar(&userRoleGrant, "grant", []string{}, "grant role to user, requires role unique id.")
	usersRolesCmd.Flags().StringArrayVar(&userRoleRevoke, "revoke", []string{}, "revoke role from user, requires role unique id.")
	usersRolesCmd.MarkFlagRequired("uid")
}

//
//
var userListCmd = &cobra.Command{
	Use:   "users",
	Short: "PrivX users",
	Long:  `List and manage PrivX users`,
	Example: `
privx-cli users [access flags]
	`,
	SilenceUsage: true,
	RunE:         userList,
}

func userList(cmd *cobra.Command, args []string) error {
	api := rolestore.New(curl())
	users, err := api.SearchUsers(strings.Join(userQuery, ","), "")
	if err != nil {
		return err
	}

	return stdout(users)
}

//
//
var userShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Description about PrivX user",
	Long:  `Description about PrivX user`,
	Example: `
privx-cli users show [access flags] UID ...
	`,
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE:         userShow,
}

func userShow(cmd *cobra.Command, args []string) error {
	api := rolestore.New(curl())
	users := []rolestore.User{}

	for _, uid := range args {
		user, err := api.User(uid)
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
	api := rolestore.New(curl())

	for _, role := range userRoleGrant {
		err := api.GrantUserRole(userID, role)
		if err != nil {
			return err
		}
	}

	for _, role := range userRoleRevoke {
		err := api.RevokeUserRole(userID, role)
		if err != nil {
			return err
		}
	}

	roles, err := api.UserRoles(userID)
	if err != nil {
		return err
	}
	return stdout(roles)
}

func decodeJSON(name string, object interface{}) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &object)
	if err != nil {
		return err
	}

	return nil
}
