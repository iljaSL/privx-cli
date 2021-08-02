//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/connectionmanager"
	"github.com/spf13/cobra"
)

type connectionOptions struct {
	channID  string
	fileID   string
	connID   string
	roleID   string
	userID   string
	hostID   string
	fileName string
	sortkey  string
	sortdir  string
	format   string
	filter   string
	offset   int
	limit    int
	force    bool
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
	cmd.AddCommand(storedFileDownloadCmd())
	cmd.AddCommand(trailLogDownloadCmd())
	cmd.AddCommand(accessRoleListCmd())
	cmd.AddCommand(connectionAccessRoleGrantCmd())
	cmd.AddCommand(connectionAccessRoleRevokeCmd())
	cmd.AddCommand(connectionTerminateCmd())

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
	privx-cli connections show [access flags] --conn-id <CONN-ID>,<CONN-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	cmd.MarkFlagRequired("conn-id")

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
func storedFileDownloadCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "download-file",
		Short: "Download trail stored file",
		Long:  `Download trail stored file`,
		Example: `
	privx-cli connections download-file [access flags] --conn-id <CONN-ID> --file-id <FILE-ID> --channel-id <CHANNEL-ID> --name <FILE-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return storedFileDownload(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	flags.StringVar(&options.channID, "channel-id", "", "channel ID")
	flags.StringVar(&options.fileID, "file-id", "", "file ID")
	flags.StringVar(&options.fileName, "name", "", "file name")
	cmd.MarkFlagRequired("conn-id")
	cmd.MarkFlagRequired("channel-id")
	cmd.MarkFlagRequired("file-id")
	cmd.MarkFlagRequired("name")

	return cmd
}

func storedFileDownload(options connectionOptions) error {
	api := connectionmanager.New(curl())

	sessionID, err := api.CreateSessionIDFileDownload(options.connID, options.channID, options.fileID)
	if err != nil {
		return err
	}

	err = api.DownloadStoredFile(options.connID, options.channID, options.fileID,
		sessionID, options.fileName)
	if err != nil {
		return err
	}

	return nil
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
	privx-cli connections download-log [access flags] --conn-id <CONN-ID> --channel-id <CHANNEL-ID> --sid <SESSION-ID> --name <FILE-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return trailLogDownload(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	flags.StringVar(&options.channID, "channel-id", "", "channel ID")
	flags.StringVar(&options.fileName, "name", "", "file name")
	flags.StringVar(&options.format, "format", "", "trail log format, json or hex")
	flags.StringVar(&options.filter, "filter", "", "trail log event filter")
	cmd.MarkFlagRequired("conn-id")
	cmd.MarkFlagRequired("channel-id")
	cmd.MarkFlagRequired("name")

	return cmd
}

func trailLogDownload(options connectionOptions) error {
	api := connectionmanager.New(curl())

	sessionID, err := api.CreateSessionIDTrailLog(options.connID, options.channID)
	if err != nil {
		return err
	}

	err = api.DownloadTrailLog(options.connID, options.channID, sessionID,
		options.format, options.filter, options.fileName)
	if err != nil {
		return err
	}

	return nil
}

//
//
func accessRoleListCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "access-roles",
		Short: "List access roles for a connection",
		Long:  `List access roles for a connection`,
		Example: `
	privx-cli connections access-roles [access flags] --conn-id <CONN-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return accessRoleList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	cmd.MarkFlagRequired("conn-id")

	return cmd
}

func accessRoleList(options connectionOptions) error {
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
		Use:   "grant-access-role",
		Short: "Add the role to a list of roles that can explicitly access this connection data",
		Long:  `Add the role to a list of roles that can explicitly access this connection data`,
		Example: `
	privx-cli connections grant-role [access flags] --conn-id <CONN-ID> --role-id <ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionAccessRoleGrant(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	flags.StringVar(&options.roleID, "role-id", "", "role ID")
	cmd.MarkFlagRequired("conn-id")
	cmd.MarkFlagRequired("role-id")

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
		Use:   "revoke-access-role",
		Short: "Remove the role from a list of roles that can explicitly access this connection data",
		Long:  `Remove the role from a list of roles that can explicitly access this connection data`,
		Example: `
	privx-cli connections revoke-access-role [access flags] --conn-id <CONN-ID> --role-id <ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionAccessRoleRevoke(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	flags.StringVar(&options.roleID, "role-id", "", "role ID")
	flags.BoolVarP(&options.force, "force", "f", false, "force command")
	cmd.MarkFlagRequired("role-id")

	return cmd
}

func connectionAccessRoleRevoke(options connectionOptions) error {
	api := connectionmanager.New(curl())

	if options.connID != "" {
		err := api.RevokeAccessRoleFromConnection(options.connID, options.roleID)
		if err != nil {
			return err
		}
	}

	if options.connID == "" {
		if !options.force {
			fmt.Println("You are about to delete the roles from ALL the connections. Please use the --force | -f flag to proceed with the command")
			os.Exit(1)
		} else {
			err := api.RevokeAccessRoleFromAllConnections(options.roleID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//
//
func connectionTerminateCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "terminate",
		Short: "Terminate connection by ID",
		Long:  `Terminate connection by ID`,
		Example: `
	privx-cli connections terminate [access flags] --conn-id <CONN-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionTerminate(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "terminate connection by ID")
	flags.StringVar(&options.hostID, "by-target", "", "terminate connection by host ID")
	flags.StringVar(&options.userID, "by-user", "", "terminate connection by user ID")

	return cmd
}

func connectionTerminate(options connectionOptions) error {
	if (options == connectionOptions{}) {
		fmt.Println("Specify at least one flag for the termination type of the connection: --conn-id, --by-target or --by-user")
	} else if options.hostID != "" {
		terminateConnectionByTargerHost(options)
	} else if options.userID != "" {
		terminateConnectionByUser(options)
	} else if options.connID != "" {
		terminateConnectionByConnection(options)
	}

	return nil
}

func terminateConnectionByConnection(options connectionOptions) error {
	api := connectionmanager.New(curl())

	err := api.TerminateConnection(options.connID)
	if err != nil {
		return err
	}

	return nil
}

func terminateConnectionByTargerHost(options connectionOptions) error {
	api := connectionmanager.New(curl())

	err := api.TerminateConnectionsByTargetHost(options.hostID)
	if err != nil {
		return err
	}

	return nil
}

func terminateConnectionByUser(options connectionOptions) error {
	api := connectionmanager.New(curl())

	err := api.TerminateConnectionsByUser(options.userID)
	if err != nil {
		return err
	}

	return nil
}
