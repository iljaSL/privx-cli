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
	requestID string
)

func init() {
	rootCmd.AddCommand(requestListCmd)
	requestListCmd.Flags().IntVar(&offset, "offset", 0, "where to start fetching the items")
	requestListCmd.Flags().IntVar(&limit, "limit", 50, "number of items to return")
	requestListCmd.Flags().StringVar(&filter, "filter", "", "filter request items")

	requestListCmd.AddCommand(requestCreateCmd)

	requestListCmd.AddCommand(requestShowCmd)
	requestShowCmd.Flags().StringVar(&requestID, "id", "", "request ID")
	requestShowCmd.MarkFlagRequired("id")

	requestListCmd.AddCommand(requestDeleteCmd)
	requestDeleteCmd.Flags().StringVar(&requestID, "id", "", "request ID")
	requestDeleteCmd.MarkFlagRequired("id")

	requestListCmd.AddCommand(decisionOnRequestCmd)
	decisionOnRequestCmd.Flags().StringVar(&requestID, "id", "", "request ID")
	decisionOnRequestCmd.MarkFlagRequired("id")

	requestListCmd.AddCommand(requestSearchCmd)
	requestSearchCmd.Flags().IntVar(&offset, "offset", 0, "where to start fetching the items")
	requestSearchCmd.Flags().IntVar(&limit, "limit", 50, "number of items to return")
	requestSearchCmd.Flags().StringVar(&filter, "filter", "", "filter request items(requests, active_requests, approvals, etc.)")
	requestSearchCmd.Flags().StringVar(&sortkey, "sortkey", "", "sort by specific object property")
	requestSearchCmd.Flags().StringVar(&sortdir, "sortdir", "", "sort direction, ASC or DESC")
}

//
//
var requestListCmd = &cobra.Command{
	Use:   "requests",
	Short: "List and manage request queues",
	Long:  `List and manage the request queue for the user`,
	Example: `
privx-cli requests [access flags] --offset <OFFSET> --limit <LIMIT> --filter <FILTER>
	`,
	SilenceUsage: true,
	RunE:         requestList,
}

func requestList(cmd *cobra.Command, args []string) error {
	api := workflow.New(curl())

	requests, err := api.Requests(offset, limit, filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	return stdout(requests)
}

//
//
var requestCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create request",
	Long:  `Add a workflow to the request queue`,
	Example: `
privx-cli requests create [access flags] JSON-FILE
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         requestCreate,
}

func requestCreate(cmd *cobra.Command, args []string) error {
	var newRequest workflow.Request
	api := workflow.New(curl())

	err := decodeJSON(args[0], &newRequest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	id, err := api.CreateRequest(&newRequest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	return stdout(id)
}

//
//
var requestShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Get request by ID",
	Long:  `Get a request object by ID. Request ID's are separated by commas when using multiple values, see example.`,
	Example: `
privx-cli requests show [access flags] --id <REQUEST-ID>,<REQUEST-ID>
	`,
	SilenceUsage: true,
	RunE:         requestShow,
}

func requestShow(cmd *cobra.Command, args []string) error {
	api := workflow.New(curl())
	requests := []workflow.Request{}

	for _, id := range strings.Split(requestID, ",") {
		request, err := api.Request(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		requests = append(requests, *request)
	}

	return stdout(requests)
}

//
//
var requestDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete request",
	Long:  `Delete request item by ID. Request ID's are separated by commas when using multiple values, see example.`,
	Example: `
privx-cli requests delete [access flags] --id <REQUEST-ID>,<REQUEST-ID>
	`,
	SilenceUsage: true,
	RunE:         requestDelete,
}

func requestDelete(cmd *cobra.Command, args []string) error {
	api := workflow.New(curl())

	for _, id := range strings.Split(requestID, ",") {
		err := api.DeleteRequest(id)
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
var decisionOnRequestCmd = &cobra.Command{
	Use:   "decision-request",
	Short: "Update a request in queue",
	Long:  `Update a request in queue. Only users with matching role are permitted to change the status of a step requiring such role.`,
	Example: `
privx-cli requests request-decision [access flags] JSON-FILE
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         decisionOnRequest,
}

func decisionOnRequest(cmd *cobra.Command, args []string) error {
	var request workflow.Decision
	api := workflow.New(curl())

	err := decodeJSON(args[0], &request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	err = api.MakeDecisionOnRequest(hostID, request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	return nil
}

//
//
var requestSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search access requests",
	Long:  `Search access requests`,
	Example: `
privx-cli requests search [access flags] --offset <OFFSET> --limit <LIMIT> --filter <FILTER>
privx-cli requests search [access flags] JSON-FILE
	`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE:         requestSearch,
}

func requestSearch(cmd *cobra.Command, args []string) error {
	var searchObject workflow.Search
	api := workflow.New(curl())

	if len(args) == 1 {
		err := decodeJSON(args[0], &searchObject)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}

	requests, err := api.SearchRequests(offset, limit, strings.ToUpper(sortdir), sortkey, strings.ToUpper(filter), &searchObject)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	return stdout(requests)
}
