//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/trailindex"
	"github.com/spf13/cobra"
)

type trailindexOptions struct {
	connID  string
	sortdir string
	limit   int
	offset  int
}

func init() {
	rootCmd.AddCommand(indexingStatusListCmd())
}

//
//
func indexingStatusListCmd() *cobra.Command {
	options := trailindexOptions{}

	cmd := &cobra.Command{
		Use:   "trailindex",
		Short: "List and manage trail files",
		Long:  `List and manage trail files. Connection ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli trailindex [access flags] --conn-id <CONNECTION-ID>,<CONNECTION-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return indexingStatusList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	cmd.MarkFlagRequired("conn-id")

	cmd.AddCommand(contentSearchCmd())
	cmd.AddCommand(indexingStatusShowCmd())
	cmd.AddCommand(indexingStartCmd())

	return cmd
}

func indexingStatusList(options trailindexOptions) error {
	api := trailindex.New(curl())

	status, err := api.IndexingStatuses(strings.Split(options.connID, ","))
	if err != nil {
		return err
	}

	return stdout(status)
}

//
//
func contentSearchCmd() *cobra.Command {
	options := trailindexOptions{}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for the content based on the search parameters defined",
		Long:  `Search for the content based on the search parameters defined`,
		Example: `
	privx-cli trailindex [access flags] --offset <OFFSET> --limit <LIMIT> --sortdir <SORTDIR>
	privx-cli trailindex [access flags] JSON-FILE
		`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return contentSearch(options, args)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC")

	cmd.AddCommand(indexingStatusShowCmd())

	return cmd
}

func contentSearch(options trailindexOptions, args []string) error {
	var searchObject trailindex.SearchRequestObject
	api := trailindex.New(curl())

	if len(args) == 1 {
		err := decodeJSON(args[0], &searchObject)
		if err != nil {
			return err
		}
	}

	content, err := api.SearchContent(options.offset, options.limit,
		strings.ToUpper(options.sortdir), searchObject)
	if err != nil {
		return err
	}

	return stdout(content)
}

//
//
func indexingStatusShowCmd() *cobra.Command {
	options := trailindexOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get indexing status of a connection",
		Long:  `Get indexing status of a connection`,
		Example: `
	privx-cli trailindex show [access flags] --conn-id <CONNECTION-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return indexingStatusShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	cmd.MarkFlagRequired("conn-id")

	return cmd
}

func indexingStatusShow(options trailindexOptions) error {
	api := trailindex.New(curl())

	status, err := api.IndexingStatus(options.connID)
	if err != nil {
		return err
	}

	return stdout(status)
}

//
//
func indexingStartCmd() *cobra.Command {
	options := trailindexOptions{}

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start indexing connections",
		Long:  `Start indexing connections. Connection ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli trailindex start [access flags] --conn-id <CONNECTION-ID>,<CONNECTION-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return indexingStart(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	cmd.MarkFlagRequired("conn-id")

	return cmd
}

func indexingStart(options trailindexOptions) error {
	api := trailindex.New(curl())

	conn, err := api.StartIndexing(strings.Split(options.connID, ","))
	if err != nil {
		return err
	}

	return stdout(conn)
}
