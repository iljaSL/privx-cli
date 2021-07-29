//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/connectionmanager"
	"github.com/spf13/cobra"
)

type connectionOptions struct {
	channID   string
	fileID    string
	connID    string
	sessionID string
	roleID    string
	userID    string
	hostID    string
	fileName  string
	sortkey   string
	sortdir   string
	format    string
	filter    string
	offset    int
	limit     int
}

func init() {
	rootCmd.AddCommand(connectionListCmd())
}

//
//
func connectionListCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "connections",
		Short: "List and manage connections",
		Long:  `List and manage connections`,
		Example: `
	privx-cli connections [access flags] --offset <OFFSET> --sortkey <SORTKEY>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionList(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC (default ASC)")

	cmd.AddCommand(connectionSearchCmd())
	cmd.AddCommand(connectionShowCmd())
	cmd.AddCommand(sessionIDFileDownloadCreateCmd())
	cmd.AddCommand(storedFileDownloadCmd())
	cmd.AddCommand(sessionIDTrailLogDownloadCreateCmd())
	cmd.AddCommand(trailLogDownloadCmd())
	cmd.AddCommand(accessRoleShowCmd())
	cmd.AddCommand(connectionAccessRoleGrantCmd())
	cmd.AddCommand(connectionAccessRoleRevokeCmd())
	cmd.AddCommand(allConnectionAccessRoleRevokeCmd())
	cmd.AddCommand(connectionTerminateCmd())
	cmd.AddCommand(connectionByTargetHostTerminateCmd())
	cmd.AddCommand(connectionByUserTerminateCmd())

	return cmd
}

func connectionList(options connectionOptions) error {
	api := connectionmanager.New(curl())

	conn, err := api.Connections(options.offset, options.limit,
		options.sortkey, options.sortdir)
	if err != nil {
		return err
	}

	return stdout(conn)
}

//
//
func connectionSearchCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for connections",
		Long:  `Search for connections`,
		Example: `
	privx-cli connections search [access flags] --offset <OFFSET> --sortkey <SORTKEY>
	privx-cli connections search [access flags] --limit <LIMIT> JSON-FILE
		`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionSearch(options, args)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC (default ASC)")

	return cmd
}

func connectionSearch(options connectionOptions, args []string) error {
	var searchObject connectionmanager.ConnectionSearch
	api := connectionmanager.New(curl())

	if len(args) == 1 {
		err := decodeJSON(args[0], &searchObject)
		if err != nil {
			return err
		}
	}

	conn, err := api.SearchConnections(options.offset, options.limit, options.sortdir,
		options.sortkey, searchObject)
	if err != nil {
		return err
	}

	return stdout(conn)
}

//
//
func connectionShowCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get connection by ID",
		Long:  `Get connection by ID. Connection ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli connections show [access flags] --id <CONN-ID>,<CONN-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "id", "", "connection ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func connectionShow(options connectionOptions) error {
	api := connectionmanager.New(curl())
	conns := []connectionmanager.Connection{}

	for _, id := range strings.Split(options.connID, ",") {
		conn, err := api.Connection(id)
		if err != nil {
			return err
		}
		conns = append(conns, *conn)
	}

	return stdout(conns)
}

//
//
func sessionIDFileDownloadCreateCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "create-session-id-file",
		Short: "Create session ID for trail stored file download",
		Long:  `Create session ID for trail stored file download`,
		Example: `
	privx-cli connections session-id-file [access flags] --id <CONN-ID> --fid <FILE-ID> --chid <CHANNEL-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sessionIDFileDownloadCreate(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "id", "", "connection ID")
	flags.StringVar(&options.channID, "chid", "", "channel ID")
	flags.StringVar(&options.fileID, "fid", "", "file ID")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("chid")
	cmd.MarkFlagRequired("fid")

	return cmd
}

func sessionIDFileDownloadCreate(options connectionOptions) error {
	api := connectionmanager.New(curl())

	sessionID, err := api.CreateSessionIDFileDownload(options.connID, options.channID, options.fileID)
	if err != nil {
		return err
	}

	return stdout(sessionID)
}

//
//
func storedFileDownloadCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "download-file",
		Short: "Download trail stored file",
		Long:  `Download trail stored file`,
		Example: `
	privx-cli connections download-file [access flags] --id <CONN-ID> --fid <FILE-ID> --chid <CHANNEL-ID> --name <FILE-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return storedFileDownload(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "id", "", "connection ID")
	flags.StringVar(&options.channID, "chid", "", "channel ID")
	flags.StringVar(&options.fileID, "fid", "", "file ID")
	flags.StringVar(&options.sessionID, "sid", "", "session ID")
	flags.StringVar(&options.fileName, "name", "", "file name")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("chid")
	cmd.MarkFlagRequired("fid")
	cmd.MarkFlagRequired("sid")
	cmd.MarkFlagRequired("name")

	return cmd
}

func storedFileDownload(options connectionOptions) error {
	store := connectionmanager.New(curl())

	err := store.DownloadStoredFile(options.connID, options.channID, options.fileID,
		options.sessionID, options.fileName)
	if err != nil {
		return err
	}

	return nil
}

//
//
func sessionIDTrailLogDownloadCreateCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "create-session-id-log",
		Short: "Create session ID for trail log download",
		Long:  `Create session ID for trail log download`,
		Example: `
	privx-cli connections create-session-id-log [access flags] --id <CONN-ID> --chid <CHANNEL-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sessionIDTrailLogDownloadCreate(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "id", "", "connection ID")
	flags.StringVar(&options.channID, "chid", "", "channel ID")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("chid")

	return cmd
}

func sessionIDTrailLogDownloadCreate(options connectionOptions) error {
	api := connectionmanager.New(curl())

	sessionID, err := api.CreateSessionIDTrailLog(options.connID, options.channID)
	if err != nil {
		return err
	}

	return stdout(sessionID)
}

//
//
func trailLogDownloadCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "download-log",
		Short: "Download trail log",
		Long:  `Download trail log`,
		Example: `
	privx-cli connections download-log [access flags] --id <CONN-ID> --chid <CHANNEL-ID> --sid <SESSION-ID> --name <FILE-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return trailLogDownload(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "id", "", "connection ID")
	flags.StringVar(&options.channID, "chid", "", "channel ID")
	flags.StringVar(&options.sessionID, "sid", "", "session ID")
	flags.StringVar(&options.fileName, "name", "", "file name")
	flags.StringVar(&options.format, "format", "", "trail log format, json or hex")
	flags.StringVar(&options.filter, "filter", "", "trail log event filter")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("chid")
	cmd.MarkFlagRequired("sid")
	cmd.MarkFlagRequired("name")

	return cmd
}

func trailLogDownload(options connectionOptions) error {
	api := connectionmanager.New(curl())

	err := api.DownloadTrailLog(options.connID, options.channID, options.sessionID,
		options.format, options.filter, options.fileName)
	if err != nil {
		return err
	}

	return nil
}

//
//
func accessRoleShowCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "roles",
		Short: "Get saved access roles for a connection",
		Long:  `Get saved access roles for a connection`,
		Example: `
	privx-cli connections roles [access flags] --id <CONN-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return accessRoleShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "id", "", "connection ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func accessRoleShow(options connectionOptions) error {
	api := connectionmanager.New(curl())

	roles, err := api.AccessRoles(options.connID)
	if err != nil {
		return err
	}

	return stdout(roles)
}

//
//
func connectionAccessRoleGrantCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "grant-access-role-permission",
		Short: "Grant a permission for a role for a connection",
		Long:  `Grant a permission for a role for a connection`,
		Example: `
	privx-cli connections grant-role-permission [access flags] --id <CONN-ID> --rid <ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionAccessRoleGrant(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "id", "", "connection ID")
	flags.StringVar(&options.roleID, "rid", "", "role ID")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("rid")

	return cmd
}

func connectionAccessRoleGrant(options connectionOptions) error {
	api := connectionmanager.New(curl())

	err := api.GrantAccessRoleToConnection(options.connID, options.roleID)
	if err != nil {
		return err
	}

	return nil
}

//
//
func connectionAccessRoleRevokeCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "revoke-role-from-connection",
		Short: "Revoke a permission for a role from a connection",
		Long:  `Revoke a permission for a role from a connection`,
		Example: `
	privx-cli connections revoke-role-from-connection [access flags] --id <CONN-ID> --rid <ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionAccessRoleRevoke(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "id", "", "connection ID")
	flags.StringVar(&options.roleID, "rid", "", "role ID")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("rid")

	return cmd
}

func connectionAccessRoleRevoke(options connectionOptions) error {
	api := connectionmanager.New(curl())

	err := api.RevokeAccessRoleFromConnection(options.connID, options.roleID)
	if err != nil {
		return err
	}

	return nil
}

//
//
func allConnectionAccessRoleRevokeCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "revoke-role-from-all-connections",
		Short: "Revoke permissions for a role from connections",
		Long:  `Revoke permissions for a role from connections`,
		Example: `
	privx-cli connections revoke-role-from-all-connections [access flags] --rid <ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return allConnectionAccessRoleRevoke(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.roleID, "rid", "", "role ID")
	cmd.MarkFlagRequired("rid")

	return cmd
}

func allConnectionAccessRoleRevoke(options connectionOptions) error {
	api := connectionmanager.New(curl())

	err := api.RevokeAccessRoleFromAllConnections(options.roleID)
	if err != nil {
		return err
	}

	return nil
}

//
//
func connectionTerminateCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "terminate-connection",
		Short: "Terminate connection by ID",
		Long:  `Terminate connection by ID`,
		Example: `
	privx-cli connections terminate-connection [access flags] --id <CONN-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionTerminate(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "id", "", "connection ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func connectionTerminate(options connectionOptions) error {
	api := connectionmanager.New(curl())

	err := api.TerminateConnection(options.connID)
	if err != nil {
		return err
	}

	return nil
}

//
//
func connectionByTargetHostTerminateCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "terminate-host-connection",
		Short: "Terminate connection from host",
		Long:  `Terminate connection from host`,
		Example: `
	privx-cli terminate-host-connection [access flags] -hid <HOST-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionByTargetHostTerminate(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.hostID, "hid", "", "host ID")
	cmd.MarkFlagRequired("hid")

	return cmd
}

func connectionByTargetHostTerminate(options connectionOptions) error {
	api := connectionmanager.New(curl())

	err := api.TerminateConnectionsByTargetHost(options.hostID)
	if err != nil {
		return err
	}

	return nil
}

//
//
func connectionByUserTerminateCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "terminate-user-connection",
		Short: "Terminate connection of a user",
		Long:  `Terminate connection of a user`,
		Example: `
	privx-cli connections terminate-user-connection [access flags] --uid <USER-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionByUserTerminate(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "uid", "", "user ID")
	cmd.MarkFlagRequired("uid")

	return cmd
}

func connectionByUserTerminate(options connectionOptions) error {
	api := connectionmanager.New(curl())

	err := api.TerminateConnectionsByUser(options.userID)
	if err != nil {
		return err
	}

	return nil
}
