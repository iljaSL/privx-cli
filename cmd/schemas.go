//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/settings"
	"github.com/spf13/cobra"
)

type schemaOptions struct {
	scope   string
	section string
}

func init() {
	rootCmd.AddCommand(scopeSchemaListCmd())
}

//
//
func scopeSchemaListCmd() *cobra.Command {
	options := schemaOptions{}

	cmd := &cobra.Command{
		Use:   "schemas",
		Short: "List settings scope and section schemas",
		Long:  `List settings scope and section schemas`,
		Example: `
	privx-cli schemas [access flags] --scope <SCOPE>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return scopeSchemaList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.scope, "scope", "", "scope setting name")
	cmd.MarkFlagRequired("scope")

	cmd.AddCommand(scopeSectionSchemaListCmd())

	return cmd
}

func scopeSchemaList(options schemaOptions) error {
	api := settings.New(curl())

	res, err := api.ScopeSchema(strings.ToUpper(options.scope))
	if err != nil {
		return err
	}

	return stdout(res)
}

//
//
func scopeSectionSchemaListCmd() *cobra.Command {
	options := schemaOptions{}

	cmd := &cobra.Command{
		Use:   "scope-section",
		Short: "Get scope section schema",
		Long:  `Get scope section schema`,
		Example: `
	privx-cli schemas scope-section [access flags] --scope <SCOPE> --section <SECTION>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return scopeSectionSchemaList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.scope, "scope", "", "scope setting name")
	flags.StringVar(&options.section, "section", "", "section setting name")
	cmd.MarkFlagRequired("scope")
	cmd.MarkFlagRequired("section")

	return cmd
}

func scopeSectionSchemaList(options schemaOptions) error {
	api := settings.New(curl())

	res, err := api.SectionSchema(strings.ToUpper(options.scope),
		strings.ToLower(options.section))
	if err != nil {
		return err
	}

	return stdout(res)
}
