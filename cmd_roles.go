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

func cmdRoles(client *api.Client) {
	roleID := flag.String("id", "", "Role ID")
	flag.Parse()

	store, err := rolestore.NewClient(client)
	if err != nil {
		log.Fatalf("failed to create role-store client: %s", err)
	}

	if len(flag.Args()) == 0 {
		log.Fatalf("Possible commands are: list, members")
	}

	cmd := flag.Args()[0]

	switch cmd {
	case "list":
		roles, err := store.GetRoles()
		if err != nil {
			log.Fatalf("get failed: %s", err)
		}
		for idx, role := range roles {
			fmt.Printf("Role %d:\n", idx)
			printRole(role, false)
		}

	case "members":
		if len(*roleID) == 0 {
			log.Fatalf("No role ID specified.")
		}
		users, err := store.GetRoleMembers(*roleID)
		if err != nil {
			log.Fatalf("get role members failed: %s", err)
		}
		for idx, user := range users {
			fmt.Printf("User %d:\n", idx)
			printUser(user)
		}

	default:
		log.Fatalf("Unknown command '%s'", cmd)
	}
}

func printRole(role *rolestore.Role, userRoles bool) {
	tab := outputFormat()
	tab.Header("Field").SetAlign(tabulate.MR)
	tab.Header("Value").SetAlign(tabulate.ML)

	err := tabulate.Reflect(tab, tabulate.OmitEmpty, nil, role)
	if err != nil {
		log.Fatalf("Failed to tabulate: %s", err)
	}
	tab.Print(os.Stdout)

	if false {
		fmt.Printf("           ID : %s\n", role.ID)
		fmt.Printf("         Name : %s\n", role.Name)
		if userRoles {
			fmt.Printf("     Explicit : %v\n", role.Explicit)
			fmt.Printf("     Implicit : %v\n", role.Implicit)
			fmt.Printf("   Grant Type : %s\n", role.GrantType)
		}
		fmt.Printf("       System : %v\n", role.System)
		fmt.Printf("      Comment : %s\n", role.Comment)
		fmt.Printf("  Permissions : %s\n", strings.Join(role.Permissions, ", "))
		fmt.Printf(" Member Count : %d\n", role.MemberCount)
	}
}
