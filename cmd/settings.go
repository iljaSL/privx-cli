//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/settings"
	"github.com/spf13/cobra"
)

type settingsOptions struct {
	scope   string
	section string
	merge   string
}

func init() {
	rootCmd.AddCommand(settingsCmd())
}

//
//
func settingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "settings",
		Short:        "List and manage settings",
		Long:         `List and manage settings`,
		SilenceUsage: true,
	}

	cmd.AddCommand(settingShowCmd())
	cmd.AddCommand(settingUpdateCmd())
	cmd.AddCommand(schemaListCmd())
	cmd.AddCommand(schemaShowCmd())

	return cmd
}

//
//
func settingShowCmd() *cobra.Command {
	options := settingsOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show settings for a specific scope/section",
		Long:  `Show settings for a specific scope/section. Scope is by default GLOBAL.`,
		Example: `
	privx-cli settings show [access flags] --scope <SCOPE>
	privx-cli settings show [access flags] --scope <SCOPE> --section <SECTION>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return settingShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&options.scope, "scope", "", "GLOBAL", "scope setting name")
	flags.StringVar(&options.section, "section", "", "section setting name")
	flags.StringVar(&options.merge, "merge", "", "signal whether service specific settings should be merged with shared settings. Compatible with scope settings only")

	return cmd
}

func settingShow(options settingsOptions) error {
	api := settings.New(curl())

	if options.section != "" {
		if options.merge != "" {
			fmt.Fprintln(os.Stderr, "Error: --merge flag is compatible with scope settings only")
			os.Exit(1)
		}

		res, err := api.ScopeSectionSettings(strings.ToUpper(options.scope),
			strings.ToLower(options.section))
		if err != nil {
			return err
		}

		return stdout(res)
	} else {
		res, err := api.ScopeSettings(strings.ToUpper(options.scope), options.merge)
		if err != nil {
			return err
		}

		return stdout(res)
	}
}

//
//
func settingUpdateCmd() *cobra.Command {
	options := settingsOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update scope/section settings",
		Long:  `Update scope/section settings. Scope is by default GLOBAL.`,
		Example: `
	privx-cli settings update [access flags] --scope <SCOPE> JSON-FILE
	privx-cli settings update [access flags] --scope <SCOPE> --section <SECTION> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return settingUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&options.scope, "scope", "", "GLOBAL", "scope setting name")
	flags.StringVar(&options.section, "section", "", "section setting name")

	return cmd
}

func settingUpdate(options settingsOptions, args []string) error {
	var updateSettings json.RawMessage
	api := settings.New(curl())

	err := decodeJSON(args[0], &updateSettings)
	if err != nil {
		return err
	}

	if options.section != "" {
		err = api.UpdateScopeSectionSettings(&updateSettings, strings.ToUpper(options.scope),
			strings.ToLower(options.section))
		if err != nil {
			return err
		}
	} else {
		err = api.UpdateScopeSettings(&updateSettings, strings.ToUpper(options.scope))
		if err != nil {
			return err
		}
	}

	return nil
}

//
//
func schemaListCmd() *cobra.Command {
	options := settingsOptions{}

	cmd := &cobra.Command{
		Use:   "list-schema",
		Short: "Get schema for the scope",
		Long:  `Get schema for the scope. Scope is by default GLOBAL.`,
		Example: `
	privx-cli settings list-schema [access flags] --scope <SCOPE>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return schemaList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&options.scope, "scope", "", "GLOBAL", "scope setting name")

	return cmd
}

func schemaList(options settingsOptions) error {
	api := settings.New(curl())

	res, err := api.ScopeSchema(strings.ToUpper(options.scope))
	if err != nil {
		return err
	}

	return stdout(res)
}

//
//
func schemaShowCmd() *cobra.Command {
	options := settingsOptions{}

	cmd := &cobra.Command{
		Use:   "show-schema",
		Short: "Get scope section schema",
		Long:  `Get scope section schema. Scope is by default GLOBAL.`,
		Example: `
	privx-cli settings show-schema [access flags] --scope <SCOPE> --section <SECTION>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return schemaShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&options.scope, "scope", "", "GLOBAL", "scope setting name")
	flags.StringVar(&options.section, "section", "", "section setting name")
	cmd.MarkFlagRequired("section")

	return cmd
}

func schemaShow(options settingsOptions) error {
	api := settings.New(curl())

	res, err := api.SectionSchema(strings.ToUpper(options.scope),
		strings.ToLower(options.section))
	if err != nil {
		return err
	}

	return stdout(res)
}
