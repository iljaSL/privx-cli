//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/monitor"
	"github.com/spf13/cobra"
)

type auditeventOptions struct {
	sortkey    string
	sortdir    string
	fuzzyCount bool
	limit      int
	offset     int
}

func init() {
	rootCmd.AddCommand(auditEventListCmd())
}

//
//
func auditEventListCmd() *cobra.Command {
	options := auditeventOptions{}

	cmd := &cobra.Command{
		Use:   "auditevents",
		Short: "List and manage audit events",
		Long:  `List and manage audit events`,
		Example: `
	privx-cli auditevents [access flags] --limit <LIMIT> --fuzzycount=true
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return auditEventsList(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC (default ASC)")
	flags.BoolVarP(&options.fuzzyCount, "fuzzycount", "", false, "return a fuzzy total count instead of exact total count")

	cmd.AddCommand(auditEventSearchCmd())
	cmd.AddCommand(auditEventCodeListCmd())

	return cmd
}

func auditEventsList(options auditeventOptions) error {
	api := monitor.New(curl())

	events, err := api.AuditEvents(options.offset, options.limit, options.sortkey,
		strings.ToUpper(options.sortdir), options.fuzzyCount)
	if err != nil {
		return err
	}

	return stdout(events)
}

//
//
func auditEventSearchCmd() *cobra.Command {
	options := auditeventOptions{}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search audit events",
		Long:  `Search audit events`,
		Example: `
	privx-cli auditevents search [access flags] --offset <OFFSET> --limit <LIMIT>
	privx-cli auditevents search [access flags] JSON-FILE
		`,
		SilenceUsage: true,
		Args:         cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return auditEventSearch(options, args)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC")
	flags.BoolVarP(&options.fuzzyCount, "fuzzycount", "", false, "return a fuzzy total count instead of exact total count")

	return cmd
}

func auditEventSearch(options auditeventOptions, args []string) error {
	var searchObject monitor.AuditEventSearchObject
	api := monitor.New(curl())

	if len(args) == 1 {
		err := decodeJSON(args[0], &searchObject)
		if err != nil {
			return err
		}
	}

	events, err := api.SearchAuditEvents(options.offset, options.limit, options.sortkey,
		strings.ToUpper(options.sortdir), options.fuzzyCount, &searchObject)
	if err != nil {
		return err
	}

	return stdout(events)
}

//
//
func auditEventCodeListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "codes",
		Short: "Get audit event codes",
		Long:  `Get audit event codes`,
		Example: `
	privx-cli auditevents codes [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return auditEventCodeList()
		},
	}

	return cmd
}

func auditEventCodeList() error {
	api := monitor.New(curl())

	codes, err := api.AuditEventCodes()
	if err != nil {
		return err
	}

	return stdout(codes)
}
