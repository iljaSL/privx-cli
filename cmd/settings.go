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
	merge   string
}

func init() {
	rootCmd.AddCommand(scopeSettingListCmd())
}

//
//
func scopeSettingListCmd() *cobra.Command {
	options := settingsOptions{}

	cmd := &cobra.Command{
		Use:   "settings",
		Short: "List and manage settings",
		Long:  `List and manage settings`,
		Example: `
	privx-cli settings [access flags] --scope <SCOPE>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return scopeSettingList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.scope, "scope", "", "scope setting name")
	flags.StringVar(&options.merge, "merge", "", "signal whether service specific settings should be merged with shared settings")
	cmd.MarkFlagRequired("scope")

	cmd.AddCommand(scopeSettingUpdateCmd())
	cmd.AddCommand(scopeSectionSettingShowCmd())
	cmd.AddCommand(scopeSectionSettingUpdateCmd())

	return cmd
}

func scopeSettingList(options settingsOptions) error {
	api := settings.New(curl())

	res, err := api.ScopeSettings(strings.ToUpper(options.scope), options.merge)
	if err != nil {
		return err
	}

	return stdout(res)
}

//
//
func scopeSettingUpdateCmd() *cobra.Command {
	options := settingsOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update scope settings",
		Long:  `Update settings for a scope`,
		Example: `
	privx-cli settings update-scope [access flags] --scope <SCOPE> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return scopeSettingUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.scope, "scope", "", "scope setting name")
	cmd.MarkFlagRequired("scope")

	return cmd
}

func scopeSettingUpdate(options settingsOptions, args []string) error {
	var scopeSettings json.RawMessage
	api := settings.New(curl())

	err := decodeJSON(args[0], &scopeSettings)
	if err != nil {
		return err
	}

	err = api.UpdateScopeSettings(&scopeSettings, strings.ToUpper(options.scope))
	if err != nil {
		return err
	}

	return nil
}

//
//
func scopeSectionSettingShowCmd() *cobra.Command {
	options := settingsOptions{}

	cmd := &cobra.Command{
		Use:   "scope-section",
		Short: "Get scope section settings",
		Long:  `Get scope section settings`,
		Example: `
	privx-cli settings scope-section [access flags] --scope <SCOPE> --section <SECTION>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return scopeSectionSettingShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.scope, "scope", "", "scope setting name")
	flags.StringVar(&options.section, "section", "", "section setting name")
	cmd.MarkFlagRequired("scope")
	cmd.MarkFlagRequired("section")

	return cmd
}

func scopeSectionSettingShow(options settingsOptions) error {
	api := settings.New(curl())

	res, err := api.ScopeSectionSettings(strings.ToUpper(options.scope),
		strings.ToLower(options.section))
	if err != nil {
		return err
	}

	return stdout(res)
}

//
//
func scopeSectionSettingUpdateCmd() *cobra.Command {
	options := settingsOptions{}

	cmd := &cobra.Command{
		Use:   "update-scope-section",
		Short: "Update scope section settings",
		Long:  `Update settings for a scope and section combination`,
		Example: `
	privx-cli settings update-scope-section [access flags] --scope <SCOPE> --section <SECTION> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return scopeSectionSettingUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.scope, "scope", "", "scope setting name")
	flags.StringVar(&options.section, "section", "", "section setting name")
	cmd.MarkFlagRequired("scope")
	cmd.MarkFlagRequired("section")

	return cmd
}

func scopeSectionSettingUpdate(options settingsOptions, args []string) error {
	var scopeSectionSettings json.RawMessage
	api := settings.New(curl())

	err := decodeJSON(args[0], &scopeSectionSettings)
	if err != nil {
		return err
	}

	err = api.UpdateScopeSectionSettings(&scopeSectionSettings, strings.ToUpper(options.scope),
		strings.ToLower(options.section))
	if err != nil {
		return err
	}

	return nil
}
