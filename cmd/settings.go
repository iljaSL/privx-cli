//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"encoding/json"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/settings"
	"github.com/spf13/cobra"
)

type settingsOptions struct {
	scope   string
	section string
}

func (m settingsOptions) normalize_scope() string {
	return strings.ToUpper(m.scope)
}

func (m settingsOptions) normalize_section() string {
	return strings.ToLower(m.section)
}

func init() {
	rootCmd.AddCommand(settingsCmd())
}

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
	cmd.AddCommand(settingRestartRequiredCmd())

	return cmd
}

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

	return cmd
}

func settingShow(options settingsOptions) error {
	api := settings.New(curl())

	if options.section != "" {
		res, err := api.ScopeSectionSettings(options.normalize_scope(),
			options.normalize_section())
		if err != nil {
			return err
		}

		return stdout(res)
	} else {
		res, err := api.ScopeSettings(options.normalize_scope(), "")
		if err != nil {
			return err
		}

		return stdout(res)
	}
}

func settingUpdateCmd() *cobra.Command {
	options := settingsOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update scope/section settings",
		Long:  `Update scope/section settings.`,
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
	flags.StringVar(&options.scope, "scope", "", "scope setting name")
	flags.StringVar(&options.section, "section", "", "section setting name")
	cmd.MarkFlagRequired("scope")

	return cmd
}

func settingUpdate(options settingsOptions, args []string) error {
	var updateSettings json.RawMessage
	api := settings.New(curl())

	err := decodeJSON(args[0], &updateSettings)
	if err != nil {
		return err
	}

	switch options.section {
	case "":
		err = api.UpdateScopeSettings(&updateSettings, options.normalize_scope())
	default:
		err = api.UpdateScopeSectionSettings(&updateSettings, options.normalize_scope(),
			options.normalize_section())
	}

	return err
}

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

	res, err := api.ScopeSchema(options.normalize_scope())
	if err != nil {
		return err
	}

	return stdout(res)
}

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

	res, err := api.SectionSchema(options.normalize_scope(),
		options.normalize_section())
	if err != nil {
		return err
	}

	return stdout(res)
}

func settingRestartRequiredCmd() *cobra.Command {
	options := settingsOptions{}

	cmd := &cobra.Command{
		Use:   "restart-required",
		Short: "Verify if restart is required.",
		Long:  `Verify if restart is required for given settings scope.`,
		Example: `
	privx-cli settings restart-required [access flags] --scope <SCOPE>
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return settingRestartRequired(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.scope, "scope", "", "scope setting name")
	cmd.MarkFlagRequired("scope")

	return cmd
}

func settingRestartRequired(options settingsOptions, args []string) error {
	var s json.RawMessage
	api := settings.New(curl())

	err := decodeJSON(args[0], &s)
	if err != nil {
		return err
	}

	res, err := api.RestartRequired(&s, options.normalize_scope())
	if err != nil {
		return err
	}

	return stdout(res)
}
