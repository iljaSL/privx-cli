//
// Copyright (c) 2023 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	authApi "github.com/SSHcom/privx-sdk-go/api/auth"
	"github.com/spf13/cobra"
)

type sessionStorageOptions struct {
	userID    string
	sourceID  string
	sessionID string
	sortkey   string
	sortdir   string
	offset    int
	limit     int
}

func init() {
	rootCmd.AddCommand(sessionStorageCmd())
}

//
//
func sessionStorageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sessions",
		Short: "List and manage sessions",
		Long:  `List and manage sessions`,
		Example: `
	privx-cli sessions [access flags]
		`,
		SilenceUsage: true,
	}

	cmd.AddCommand(userSessionsShowCmd())
	cmd.AddCommand(sourceSessionsShowCmd())
	cmd.AddCommand(sessionsSearchCmd())
	cmd.AddCommand(terminateSessionCmd())
	cmd.AddCommand(terminateUserSessionsCmd())

	return cmd
}

//
//
func userSessionsShowCmd() *cobra.Command {
	options := sessionStorageOptions{}

	cmd := &cobra.Command{
		Use:   "show-by-user",
		Short: "Get sessions by userID",
		Long:  `Get sessions by userID. Fetch valid sessions for specified user`,
		Example: `
	privx-cli sessions show-by-user [access flags] --id <USER-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return userSessionsShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "id", "", "User ID")
	flags.IntVar(&options.offset, "offset", 1, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "expires", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "ASC", "sort direction, ASC or DESC")
	cmd.MarkFlagRequired("id")

	return cmd
}

func userSessionsShow(options sessionStorageOptions) error {
	api := authApi.New(curl())

	userSessions, err := api.UserSessions(options.offset, options.limit,
		options.sortkey, options.sortdir, options.userID)
	if err != nil {
		return err
	}

	return stdout(userSessions)
}

//
//
func sourceSessionsShowCmd() *cobra.Command {
	options := sessionStorageOptions{}

	cmd := &cobra.Command{
		Use:   "show-by-source",
		Short: "Get sessions by sourceID",
		Long:  `Get sessions by sourceID. Fetch valid sessions for specified source`,
		Example: `
	privx-cli sessions show-by-source [access flags] --id <SOURCE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sourceSessionsShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.sourceID, "id", "", "Source ID")
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "expires", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "ASC", "sort direction, ASC or DESC")
	cmd.MarkFlagRequired("id")

	return cmd
}

func sourceSessionsShow(options sessionStorageOptions) error {
	api := authApi.New(curl())

	sourceSessions, err := api.SourceSessions(options.offset, options.limit,
		options.sortkey, options.sortdir, options.sourceID)
	if err != nil {
		return err
	}

	return stdout(sourceSessions)
}

//
//
func sessionsSearchCmd() *cobra.Command {
	options := sessionStorageOptions{}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for sessions",
		Long:  `Search for sessions. Search with keywords, userid or type`,
		Example: `
	privx-cli sessions search [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sessionsSearch(options, args)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "expires", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "ASC", "sort direction, ASC or DESC")

	return cmd
}

func sessionsSearch(options sessionStorageOptions, args []string) error {
	var searchParams authApi.SearchParams
	api := authApi.New(curl())

	err := decodeJSON(args[0], &searchParams)
	if err != nil {
		return err
	}

	sessions, err := api.SearchSessions(options.offset, options.limit,
		options.sortkey, options.sortdir, &searchParams)
	if err != nil {
		return err
	}

	return stdout(sessions)
}

//
//
func terminateSessionCmd() *cobra.Command {
	options := sessionStorageOptions{}

	cmd := &cobra.Command{
		Use:   "terminate",
		Short: "Terminate single session by ID",
		Long:  `Terminate single session by sessionID`,
		Example: `
	privx-cli sessions terminate [access flags] --id <SESSION-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return terminateSession(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.sessionID, "id", "", "Session ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func terminateSession(options sessionStorageOptions) error {
	api := authApi.New(curl())

	return api.TerminateSession(options.sessionID)
}

//
//
func terminateUserSessionsCmd() *cobra.Command {
	options := sessionStorageOptions{}

	cmd := &cobra.Command{
		Use:   "terminate-user",
		Short: "Terminate all sessions for a user",
		Long:  `Terminate all sessions for a user by userID`,
		Example: `
	privx-cli sessions terminate-user [access flags] --id <USER-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return terminateUserSessions(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "id", "", "User ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func terminateUserSessions(options sessionStorageOptions) error {
	api := authApi.New(curl())

	return api.TerminateUserSessions(options.userID)
}
