//
// Copyright (c) 2020 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/spf13/cobra"
)

var (
	userID          string
	trustedClientID string
	userName        string
	password        string
	clientID        string
	query           string
	apiClientRoles  string
	offset          int
	limit           int
	userQuery       []string
	userRoleGrant   []string
	userRoleRevoke  []string
)

func init() {
	rootCmd.AddCommand(usersCmd)
	usersCmd.Flags().StringArrayVarP(&userQuery, "query", "q", []string{}, "query PrivX users with keyword")

	//
	// local user commands
	usersCmd.AddCommand(localUsersCmd)
	localUsersCmd.Flags().IntVar(&offset, "offset", 0, "where to start fetching the items")
	localUsersCmd.Flags().IntVar(&limit, "limit", 50, "number of items to return")
	localUsersCmd.Flags().StringVar(&userID, "uid", "", "unique user id")
	localUsersCmd.Flags().StringVar(&userName, "username", "", "unique user name")

	localUsersCmd.AddCommand(localUserCmd)
	localUserCmd.Flags().StringVar(&userID, "uid", "", "unique user id")

	localUserCmd.AddCommand(createLocalUserCmd)

	localUserCmd.AddCommand(updateLocalUserCmd)
	updateLocalUserCmd.Flags().StringVar(&userID, "uid", "", "unique user id")
	updateLocalUserCmd.MarkFlagRequired("uid")

	localUserCmd.AddCommand(deleteLocalUserCmd)
	deleteLocalUserCmd.Flags().StringVar(&userID, "uid", "", "unique user id")
	deleteLocalUserCmd.MarkFlagRequired("uid")

	localUserCmd.AddCommand(updateLocalUserPasswordCmd)
	updateLocalUserPasswordCmd.Flags().StringVar(&userID, "uid", "", "unique user id")
	updateLocalUserPasswordCmd.Flags().StringVar(&password, "password", "", "new password for local user")
	updateLocalUserPasswordCmd.MarkFlagRequired("uid")
	updateLocalUserPasswordCmd.MarkFlagRequired("password")

	localUserCmd.AddCommand(localUserTagsCmd)
	localUserTagsCmd.Flags().IntVar(&offset, "offset", 0, "where to start fetching the items")
	localUserTagsCmd.Flags().IntVar(&limit, "limit", 50, "number of items to return")
	localUserTagsCmd.Flags().StringVar(&query, "query", "", "query string matches the tags")
	localUserTagsCmd.Flags().StringVar(&sortdir, "sortdir", "", "sort direction, ASC or DESC (default ASC)")

	//
	// trusted clients commands
	localUsersCmd.AddCommand(trustedClientsCmd)

	trustedClientsCmd.AddCommand(createTrustedClientCmd)

	trustedClientsCmd.AddCommand(trustedClientCmd)
	trustedClientCmd.Flags().StringVar(&trustedClientID, "id", "", "unique trusted client id")
	trustedClientCmd.MarkFlagRequired("id")

	trustedClientsCmd.AddCommand(deleteTrustedClientCmd)
	deleteTrustedClientCmd.Flags().StringVar(&trustedClientID, "id", "", "unique trusted client id")
	deleteTrustedClientCmd.MarkFlagRequired("id")

	trustedClientsCmd.AddCommand(updateTrustedClientCmd)
	updateTrustedClientCmd.Flags().StringVar(&trustedClientID, "id", "", "unique trusted client id")
	updateTrustedClientCmd.MarkFlagRequired("id")

	//
	// extender clients commands
	localUsersCmd.AddCommand(extenderClientsCmd)

	//
	// API clients commands
	localUsersCmd.AddCommand(apiClientsCmd)

	apiClientsCmd.AddCommand(createAPIClientCmd)
	createAPIClientCmd.Flags().StringVar(&name, "name", "", "API client name")
	createAPIClientCmd.Flags().StringVar(&apiClientRoles, "roles", "", "list of roles possessed by the API client")

	apiClientsCmd.AddCommand(apiClientCmd)
	apiClientCmd.Flags().StringVar(&clientID, "id", "", "unique API client id")
	apiClientCmd.MarkFlagRequired("id")

	apiClientsCmd.AddCommand(deleteAPIClientCmd)
	deleteAPIClientCmd.Flags().StringVar(&clientID, "id", "", "unique API client id")
	deleteAPIClientCmd.MarkFlagRequired("id")

	apiClientsCmd.AddCommand(updateAPIClientCmd)
	updateAPIClientCmd.Flags().StringVar(&clientID, "id", "", "unique API client id")
	updateAPIClientCmd.MarkFlagRequired("id")

	//
	// users commands
	usersCmd.AddCommand(usersInfoCmd)

	usersCmd.AddCommand(usersRolesCmd)
	usersRolesCmd.Flags().StringVar(&userID, "uid", "", "user unique id")
	usersRolesCmd.Flags().StringArrayVar(&userRoleGrant, "grant", []string{}, "grant role to user, requires role unique id.")
	usersRolesCmd.Flags().StringArrayVar(&userRoleRevoke, "revoke", []string{}, "revoke role from user, requires role unique id.")
	usersRolesCmd.MarkFlagRequired("uid")
}

//
//
var localUsersCmd = &cobra.Command{
	Use:   "local",
	Short: "local users",
	Long:  `get information about privx local users`,
	Example: `
privx-cli users local [access flags] COMMANDS
privx-cli users local [access flags] --uid UID
privx-cli users local [access flags] --username USERNAME
privx-cli users local [access flags] --offset OFFSET --limit LIMIT
	`,
	SilenceUsage: true,
	RunE:         localUsers,
}

func localUsers(cmd *cobra.Command, args []string) error {
	store := userstore.New(curl())

	users, err := store.LocalUsers(offset, limit, userID, userName)
	if err != nil {
		return err
	}

	return stdout(users)
}

//
//
var createLocalUserCmd = &cobra.Command{
	Use:   "create",
	Short: "create new local user",
	Long:  `create new local user to privx local user store`,
	Example: `
privx-cli users local user create [access flags] JSON-FILE
	`,
	SilenceUsage: true,
	RunE:         createLocalUser,
}

func createLocalUser(cmd *cobra.Command, args []string) error {
	var newUser userstore.LocalUser
	api := userstore.New(curl())

	if len(args) != 1 {
		return errors.New("requires json file as argument")
	}

	file, err := openJSON(args[0])
	if err != nil {
		return err
	}

	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&newUser)
	if err != nil {
		return errors.New("json file decoding failed")
	}

	uid, err := api.CreateLocalUser(newUser)
	if err != nil {
		return err
	}

	return stdout(uid)
}

//
//
var localUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Get a local user",
	Long:  `Get a local user by user ID`,
	Example: `
privx-cli users local user [access flags] --uid UID
	`,
	SilenceUsage: true,
	RunE:         localUser,
}

func localUser(cmd *cobra.Command, args []string) error {
	store := userstore.New(curl())

	user, err := store.LocalUser(userID)
	if err != nil {
		return err
	}

	return stdout(user)
}

//
//
var updateLocalUserCmd = &cobra.Command{
	Use:   "update",
	Short: "update local user",
	Long:  `update a local user inside the privx local user store`,
	Example: `
privx-cli users local user update [access flags] JSON-FILE --uid UID
	`,
	SilenceUsage: true,
	RunE:         updateLocalUser,
}

func updateLocalUser(cmd *cobra.Command, args []string) error {
	var updateUser userstore.LocalUser
	api := userstore.New(curl())

	if len(args) != 1 {
		return errors.New("requires json file as argument")
	}

	file, err := openJSON(args[0])
	if err != nil {
		return err
	}

	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&updateUser)
	if err != nil {
		return errors.New("json file decoding failed")
	}

	err = api.UpdateLocalUser(userID, &updateUser)
	if err != nil {
		return err
	}

	return err
}

//
//
var deleteLocalUserCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete local user",
	Long:  `delete a local user from the privx local user store`,
	Example: `
privx-cli users local user delete [access flags] --uid UID
	`,
	SilenceUsage: true,
	RunE:         deleteLocalUser,
}

func deleteLocalUser(cmd *cobra.Command, args []string) error {
	store := userstore.New(curl())

	err := store.DeleteLocalUser(userID)

	return err
}

//
//
var updateLocalUserPasswordCmd = &cobra.Command{
	Use:   "update-password",
	Short: "update local user password",
	Long:  `update a local users password inside the privx local user store`,
	Example: `
privx-cli users local user update-password [access flags] --uid UID --password NEW-PASSWORD
	`,
	SilenceUsage: true,
	RunE:         updateLocalUserPassword,
}

func updateLocalUserPassword(cmd *cobra.Command, args []string) error {
	newPassword := userstore.Password{
		Password: password,
	}
	api := userstore.New(curl())

	err := api.UpdateLocalUserPassword(userID, &newPassword)

	return err
}

//
//
var localUserTagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "local user tags",
	Long:  `get privx local user tags`,
	Example: `
privx-cli users local tags [access flags]
privx-cli users local tags [access flags] --sortdir DESC
privx-cli users local tags [access flags] --query TAG
privx-cli users local tags [access flags] --offset OFFSET --limit LIMIT
	`,
	SilenceUsage: true,
	RunE:         localUserTags,
}

func localUserTags(cmd *cobra.Command, args []string) error {
	store := userstore.New(curl())

	tags, err := store.LocalUserTags(offset, limit, strings.ToUpper(sortdir), query)
	if err != nil {
		return err
	}

	return stdout(tags)
}

//
//
var trustedClientsCmd = &cobra.Command{
	Use:   "trusted-clients",
	Short: "get trusted clients",
	Long:  `get trusted clients from the privx local user store`,
	Example: `
privx-cli users trusted-clients [access flags]
	`,
	SilenceUsage: true,
	RunE:         trustedClients,
}

func trustedClients(cmd *cobra.Command, args []string) error {
	store := userstore.New(curl())

	trustedClients, err := store.TrustedClients()
	if err != nil {
		return err
	}

	return stdout(trustedClients)
}

//
//
var createTrustedClientCmd = &cobra.Command{
	Use:   "create",
	Short: "create new trusted-client",
	Long:  `create new trusted client to privX local user store`,
	Example: `
privx-cli users local trusted-clients create [access flags] JSON-FILE
	`,
	SilenceUsage: true,
	RunE:         createTrustedClient,
}

func createTrustedClient(cmd *cobra.Command, args []string) error {
	var trustedClient userstore.TrustedClient
	api := userstore.New(curl())

	if len(args) != 1 {
		return errors.New("requires json file as argument")
	}

	file, err := openJSON(args[0])
	if err != nil {
		return err
	}

	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&trustedClient)
	if err != nil {
		return errors.New("json file decoding failed")
	}

	id, err := api.CreateTrustedClient(trustedClient)
	if err != nil {
		return err
	}

	return stdout(id)
}

//
//
var trustedClientCmd = &cobra.Command{
	Use:   "client",
	Short: "get trusted client by ID",
	Long:  `get trusted client by ID from the privx local user store`,
	Example: `
privx-cli users trusted-clients client [access flags] --id TRUSTED-CLIENT-ID
	`,
	SilenceUsage: true,
	RunE:         trustedClient,
}

func trustedClient(cmd *cobra.Command, args []string) error {
	store := userstore.New(curl())

	trustedClient, err := store.TrustedClient(trustedClientID)
	if err != nil {
		return err
	}

	return stdout(trustedClient)
}

//
//
var deleteTrustedClientCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete trusted client",
	Long:  `delete a trusted client from the privx local user store`,
	Example: `
privx-cli users local trusted-clients delete [access flags] --id TRUSTED-CLIENT-ID
	`,
	SilenceUsage: true,
	RunE:         deleteTrustedClient,
}

func deleteTrustedClient(cmd *cobra.Command, args []string) error {
	store := userstore.New(curl())

	err := store.DeleteTrustedClient(trustedClientID)

	return err
}

//
//
var updateTrustedClientCmd = &cobra.Command{
	Use:   "update",
	Short: "update trusted client",
	Long:  `update an existing trusted client inside the privx local user store`,
	Example: `
privx-cli users local trusted-clients update [access flags] --id TRUSTED-CLIENT-ID JSON-FILE
	`,
	SilenceUsage: true,
	RunE:         updateTrustedClient,
}

func updateTrustedClient(cmd *cobra.Command, args []string) error {
	var trustedClient userstore.TrustedClient
	api := userstore.New(curl())

	if len(args) != 1 {
		return errors.New("requires json file as argument")
	}

	file, err := openJSON(args[0])
	if err != nil {
		return err
	}

	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&trustedClient)
	if err != nil {
		return errors.New("json file decoding failed")
	}

	err = api.UpdateTrustedClient(trustedClientID, &trustedClient)

	return err
}

//
//
var extenderClientsCmd = &cobra.Command{
	Use:   "extender-clients",
	Short: "get extender clients",
	Long:  `get extender clients from the privx local user store`,
	Example: `
privx-cli users local extender-clients [access flags]
	`,
	SilenceUsage: true,
	RunE:         extenderClients,
}

func extenderClients(cmd *cobra.Command, args []string) error {
	store := userstore.New(curl())

	extenderClients, err := store.TrustedClients()
	if err != nil {
		return err
	}

	return stdout(extenderClients)
}

//
//
var apiClientsCmd = &cobra.Command{
	Use:   "api-clients",
	Short: "get API clients",
	Long:  `get all API clients from the privx local user store`,
	Example: `
privx-cli users local api-clients [access flags]
	`,
	SilenceUsage: true,
	RunE:         apiClients,
}

func apiClients(cmd *cobra.Command, args []string) error {
	store := userstore.New(curl())
	result, err := store.APIClients()
	if err != nil {
		return err
	}

	return stdout(result)
}

//
//
var createAPIClientCmd = &cobra.Command{
	Use:   "create",
	Short: "create new API client",
	Long:  `create new API client to privX local user store`,
	Example: `
privx-cli users local api-clients create [access flags] --name NAME --roles ROLE-ID,ROLE-ID
	`,
	SilenceUsage: true,
	RunE:         createAPIClient,
}

func createAPIClient(cmd *cobra.Command, args []string) error {
	api := userstore.New(curl())

	id, err := api.CreateAPIClient(name, strings.Split(apiClientRoles, ","))
	if err != nil {
		return err
	}

	return stdout(id)
}

//
//
var apiClientCmd = &cobra.Command{
	Use:   "client",
	Short: "get API client",
	Long:  `get API client by ID from the privx local user store`,
	Example: `
privx-cli users local api-clients client [access flags] --id API-CLIENT-ID
	`,
	SilenceUsage: true,
	RunE:         apiClient,
}

func apiClient(cmd *cobra.Command, args []string) error {
	store := userstore.New(curl())

	result, err := store.APIClient(clientID)
	if err != nil {
		return err
	}

	return stdout(result)
}

//
//
var deleteAPIClientCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete API client",
	Long:  `delete a API client from the privx local user store`,
	Example: `
privx-cli users local api-clients delete [access flags] --id API-CLIENT-ID
	`,
	SilenceUsage: true,
	RunE:         deleteAPIClient,
}

func deleteAPIClient(cmd *cobra.Command, args []string) error {
	store := userstore.New(curl())

	err := store.DeleteAPIClient(clientID)

	return err
}

//
//
var updateAPIClientCmd = &cobra.Command{
	Use:   "update",
	Short: "update API client",
	Long:  `update an existing API client inside the privx local user store`,
	Example: `
privx-cli users local api-clients update [access flags] --id API-CLIENT-ID JSON-FILE
	`,
	SilenceUsage: true,
	RunE:         updateAPIClient,
}

func updateAPIClient(cmd *cobra.Command, args []string) error {
	var apiClient userstore.APIClient
	api := userstore.New(curl())

	if len(args) != 1 {
		return errors.New("requires json file as argument")
	}

	file, err := openJSON(args[0])
	if err != nil {
		return err
	}

	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&apiClient)
	if err != nil {
		return errors.New("json file decoding failed")
	}

	err = api.UpdateAPIClient(clientID, &apiClient)

	return err
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
		err := store.GrantUserRole(userID, role)
		if err != nil {
			return err
		}
	}

	for _, role := range userRoleRevoke {
		err := store.RevokeUserRole(userID, role)
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

func openJSON(name string) (*os.File, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return file, err
}
