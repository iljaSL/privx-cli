//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"fmt"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/hoststore"
	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/spf13/cobra"
)

type tagOptions struct {
	tagType string
	sortdir string
	query   string
	limit   int
	offset  int
}

func init() {
	rootCmd.AddCommand(tagListCmd())
}

//
//
func tagListCmd() *cobra.Command {
	options := tagOptions{}

	cmd := &cobra.Command{
		Use:   "tags",
		Short: "User | Host tags",
		Long:  `Get privx user or host tags`,
		Example: `
	privx-cli tags [access flags] --type user
	privx-cli tags [access flags] --type host --sortdir DESC
	privx-cli tags [access flags] --type host --query TAG
	privx-cli tags [access flags] --type user --offset OFFSET --limit LIMIT
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return tagList(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.query, "query", "", "query string matches the tags")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC (default ASC)")
	flags.StringVar(&options.tagType, "type", "", "choose the tag type, user or host")
	cmd.MarkFlagRequired("type")

	return cmd
}

func tagList(options tagOptions) error {
	switch options.tagType {
	case "user":
		userTags(options)
	case "host":
		hostTags(options)
	default:
		return fmt.Errorf("tag type does not exist: %s", options.tagType)
	}

	return nil
}

func userTags(options tagOptions) error {
	api := userstore.New(curl())

	tags, err := api.LocalUserTags(options.offset, options.limit,
		strings.ToUpper(options.sortdir), options.query)
	if err != nil {
		return err
	}

	return stdout(tags)
}

func hostTags(options tagOptions) error {
	api := hoststore.New(curl())

	tags, err := api.HostTags(options.offset, options.limit,
		strings.ToUpper(options.sortdir), options.query)
	if err != nil {
		return err
	}

	return stdout(tags)
}
