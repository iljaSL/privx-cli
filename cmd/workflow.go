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

	"github.com/SSHcom/privx-sdk-go/api/workflow"
	"github.com/spf13/cobra"
)

var (
	workflowID string
)

func init() {
	rootCmd.AddCommand(workflowListCmd)
	workflowListCmd.Flags().IntVar(&offset, "offset", 0, "where to start fetching the items")
	workflowListCmd.Flags().IntVar(&limit, "limit", 50, "number of items to return")

	workflowListCmd.AddCommand(workflowCreateCmd)

	workflowListCmd.AddCommand(workflowShowCmd)
	workflowShowCmd.Flags().StringVar(&workflowID, "id", "", "unique workflow ID")
	workflowShowCmd.MarkFlagRequired("id")

	workflowListCmd.AddCommand(workflowDeleteCmd)
	workflowDeleteCmd.Flags().StringVar(&workflowID, "id", "", "unique workflow ID")
	workflowDeleteCmd.MarkFlagRequired("id")

	workflowListCmd.AddCommand(workflowUpdateCmd)
	workflowUpdateCmd.Flags().StringVar(&workflowID, "id", "", "unique workflow ID")
	workflowUpdateCmd.MarkFlagRequired("id")

	workflowListCmd.AddCommand(workflowSettingListCmd)

	workflowListCmd.AddCommand(workflowSettingsUpdateCmd)

	workflowListCmd.AddCommand(testEmailNotificationCmd)
}

//
//
var workflowListCmd = &cobra.Command{
	Use:   "workflows",
	Short: "List and manage workflows",
	Long:  `List and manage PrivX workflows`,
	Example: `
privx-cli workflows [access flags] --offset <OFFSET> --limit <LIMIT>
	`,
	SilenceUsage: true,
	RunE:         workflowList,
}

func workflowList(cmd *cobra.Command, args []string) error {
	api := workflow.New(curl())

	workflows, err := api.Workflows(offset, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	return stdout(workflows)
}

//
//
var workflowCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new workflow",
	Long:  `Create a new workflow`,
	Example: `
privx-cli workflows create [access flags] JSON-FILE
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         workflowCreate,
}

func workflowCreate(cmd *cobra.Command, args []string) error {
	var newWorkflow workflow.Workflow
	api := workflow.New(curl())

	err := decodeJSON(args[0], &newWorkflow)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	id, err := api.CreateWorkflow(&newWorkflow)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	return stdout(id)
}

//
//
var workflowShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Get workflow by ID",
	Long:  `Get workflow object by ID. Workflow ID's are separated by commas when using multiple values, see example.`,
	Example: `
privx-cli workflows show [access flags] --id <WORKFLOW-ID>,<WORKFLOW-ID>
	`,
	SilenceUsage: true,
	RunE:         workflowShow,
}

func workflowShow(cmd *cobra.Command, args []string) error {
	api := workflow.New(curl())
	workflows := []workflow.Workflow{}

	for _, id := range strings.Split(workflowID, ",") {
		result, err := api.Workflow(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		workflows = append(workflows, *result)
	}

	return stdout(workflows)
}

//
//
var workflowDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete workflow by ID",
	Long:  `Delete workflow by ID. Workflow ID's are separated by commas when using multiple values, see example.`,
	Example: `
privx-cli workflows delete [access flags] --id <WORKFLOW-ID>,<WORKFLOW-ID>
	`,
	SilenceUsage: true,
	RunE:         workflowDelete,
}

func workflowDelete(cmd *cobra.Command, args []string) error {
	api := workflow.New(curl())

	for _, id := range strings.Split(workflowID, ",") {
		err := api.DeleteWorkflow(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		} else {
			fmt.Println(id)
		}
	}

	return nil
}

//
//
var workflowUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a workflow",
	Long:  `Update a workflow`,
	Example: `
privx-cli workflows update [access flags] JSON-FILE --id <WORKFLOW-ID>
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         workflowUpdate,
}

func workflowUpdate(cmd *cobra.Command, args []string) error {
	var updateWorkflow workflow.Workflow
	api := workflow.New(curl())

	err := decodeJSON(args[0], &updateWorkflow)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	err = api.UpdateWorkflow(workflowID, &updateWorkflow)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	return nil
}

//
//
var workflowSettingListCmd = &cobra.Command{
	Use:   "settings",
	Short: "Get workflow settings",
	Long:  `Get workflow settings`,
	Example: `
privx-cli workflows settings [access flags]
	`,
	SilenceUsage: true,
	RunE:         workflowSettingList,
}

func workflowSettingList(cmd *cobra.Command, args []string) error {
	api := workflow.New(curl())

	settings, err := api.Settings()
	if err != nil {
		return err
	}

	return stdout(settings)
}

//
//
var workflowSettingsUpdateCmd = &cobra.Command{
	Use:   "update-settings",
	Short: "Update workflow settings",
	Long:  `Update workflow settings`,
	Example: `
privx-cli workflows update-settings [access flags] JSON-FILE
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         workflowSettingsUpdate,
}

func workflowSettingsUpdate(cmd *cobra.Command, args []string) error {
	var updateSettings workflow.Settings
	api := workflow.New(curl())

	err := decodeJSON(args[0], &updateSettings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	err = api.UpdateSettings(&updateSettings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	return nil
}

//
//
var testEmailNotificationCmd = &cobra.Command{
	Use:   "testsmtp",
	Short: "Test the email settings",
	Long:  `Test the email settings`,
	Example: `
privx-cli workflows testsmtp [access flags] JSON-FILE
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         testEmailNotification,
}

func testEmailNotification(cmd *cobra.Command, args []string) error {
	var smtp workflow.Settings
	api := workflow.New(curl())

	err := decodeJSON(args[0], &smtp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	testResult, err := api.TestEmailNotification(&smtp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	return stdout(testResult)
}
