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
	rootCmd.AddCommand(indexingCmd())
}

//
//
func indexingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "List and manage SSH trail files",
		Long: `List and manage SSH trail files. Indexing commands are operation based and not a resource itself.
The index operation commands can only be applied to SSH trails.`,
		SilenceUsage: true,
	}

	cmd.AddCommand(indexingStatusCmd())
	cmd.AddCommand(contentSearchCmd())
	cmd.AddCommand(indexingStartCmd())

	return cmd
}

//
//
func indexingStatusCmd() *cobra.Command {
	options := trailindexOptions{}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "List trail-index status",
		Long:  `List trail-index status. Connection ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli index status [access flags] --conn-id <CONNECTION-ID>,<CONNECTION-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return indexingStatus(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	cmd.MarkFlagRequired("conn-id")

	return cmd
}

func indexingStatus(options trailindexOptions) error {
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
	privx-cli index search [access flags] --offset <OFFSET> --limit <LIMIT> --sortdir <SORTDIR>
	privx-cli index search [access flags] JSON-FILE
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
func indexingStartCmd() *cobra.Command {
	options := trailindexOptions{}

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start indexing connections",
		Long:  `Start indexing connections. Connection ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli index start [access flags] --conn-id <CONNECTION-ID>,<CONNECTION-ID>
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
