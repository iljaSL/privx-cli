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
	"strings"

	"github.com/SSHcom/privx-sdk-go/api"
	"github.com/SSHcom/privx-sdk-go/api/rolestore"
)

func cmdRoles(client *api.Client) {
	flag.Parse()

	store, err := rolestore.NewClient(client)
	if err != nil {
		log.Fatalf("failed to create role-store client: %s", err)
	}

	if len(flag.Args()) == 0 {
		log.Fatalf("Possible commands are: list")
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
			printRole(role)
		}

	default:
		log.Fatalf("Unknown command '%s'", cmd)
	}
}

func printRole(role *rolestore.Role) {
	fmt.Printf("           ID : %s\n", role.ID)
	fmt.Printf("         Name : %s\n", role.Name)
	fmt.Printf("     Explicit : %v\n", role.Explicit)
	fmt.Printf("     Implicit : %v\n", role.Implicit)
	fmt.Printf("       System : %v\n", role.System)
	fmt.Printf("    GrantType : %s\n", role.GrantType)
	fmt.Printf("      Comment : %s\n", role.Comment)
	fmt.Printf("  Permissions : %s\n", strings.Join(role.Permissions, ", "))
	fmt.Printf(" Member Count : %d\n", role.MemberCount)
}
