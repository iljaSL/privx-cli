//
// Copyright (c) 2020 SSH Communications Security Inc.
//
// All rights reserved.
//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api"
	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/markkurossi/tabulate"
)

func cmdUsers(client *api.Client) {
	userID := flag.String("id", "", "User ID")
	flag.Parse()

	store, err := rolestore.NewClient(client)
	if err != nil {
		log.Fatalf("failed to create role-store client: %s", err)
	}

	if len(flag.Args()) == 0 {
		log.Fatalf("Possible commands are: add-role, info, remove-role, roles, search")
	}

	cmd := flag.Args()[0]
	args := flag.Args()[1:]

	switch cmd {
	case "add-role":
		if len(*userID) == 0 {
			log.Fatalf("No user ID specified.")
		}
		if len(args) == 0 {
			log.Fatalf("No role IDs specified.")
		}
		for _, roleID := range args {
			err = store.AddUserRole(*userID, roleID)
			if err != nil {
				log.Fatalf("Failed to add role '%s': %s", roleID, err)
			}
		}

	case "info":
		if len(*userID) == 0 {
			log.Fatalf("No user ID specified.")
		}
		user, err := store.GetUser(*userID)
		if err != nil {
			log.Fatalf("get info failed: %s", err)
		}
		printUser(user)

	case "remove-role":
		if len(*userID) == 0 {
			log.Fatalf("No user ID specified.")
		}
		if len(args) == 0 {
			log.Fatalf("No role IDs specified.")
		}
		for _, roleID := range args {
			err = store.RemoveUserRole(*userID, roleID)
			if err != nil {
				log.Fatalf("Failed to remove role '%s': %s", roleID, err)
			}
		}

	case "roles":
		if len(*userID) == 0 {
			log.Fatalf("No user ID specified.")
		}
		roles, err := store.GetUserRoles(*userID)
		if err != nil {
			log.Fatalf("get roles failed: %s", err)
		}
		for idx, role := range roles {
			fmt.Printf("Role %d:\n", idx)
			printRole(role, true)
		}

	case "search":
		users, err := store.SearchUsers(strings.Join(args, ","), "")
		if err != nil {
			log.Fatalf("search failed: %s", err)
		}
		for idx, user := range users {
			fmt.Printf("Result %d:\n", idx)
			printUser(user)
		}

	default:
		log.Fatalf("Unknown command '%s'", cmd)
	}
}

func printUser(user *rolestore.User) {
	tab := format()
	tab.Header("Field").SetAlign(tabulate.ML)
	tab.Header("Value").SetAlign(tabulate.ML)

	err := tabulate.Reflect(tab, tabulate.OmitEmpty, []string{"userCtx"}, user)
	if err != nil {
		log.Fatalf("Failed to tabulate: %s", err)
	}
	tab.Print(os.Stdout)
}
