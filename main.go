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

	"github.com/SSHcom/privx-sdk-go/api"
	"github.com/SSHcom/privx-sdk-go/config"
	"github.com/SSHcom/privx-sdk-go/oauth"
	"github.com/markkurossi/tabulate"
)

var commands = map[string]func(client *api.Client){
	"users":   cmdUsers,
	"secrets": cmdSecrets,
	"roles":   cmdRoles,
}

var outputFormat func() *tabulate.Tabulate

var formats = map[string]func() *tabulate.Tabulate{
	"whitespace": tabulate.NewWS,
	"ascii":      tabulate.NewASCII,
	"unicode":    tabulate.NewUnicode,
	"colon":      tabulate.NewColon,
	"csv":        tabulate.NewCSV,
}

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s [options] COMMAND [command options] [ARG]...\n",
			os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\nCommands:\n")
		for key := range commands {
			fmt.Fprintf(os.Stderr, "  - %s\n", key)
		}
		fmt.Fprintf(os.Stderr,
			"\nType %s COMMAND -h for help about COMMAND\n",
			os.Args[0])
	}

	apiEndpoint := flag.String("api", "", "API endpoint URL")
	configFile := flag.String("config", config.Default(), "configuration file")
	verbose := flag.Bool("v", false, "verbose output")
	formatFlag := flag.String("format", "unicode", "output format")
	flag.Parse()

	outputFormat = formats[*formatFlag]
	if outputFormat == nil {
		log.Printf("Invalid output format '%s'", *formatFlag)
		log.Printf("Supported formats are:")
		for k := range formats {
			log.Printf(" - %s", k)
		}
		os.Exit(1)
	}

	config, err := config.Read(*configFile)
	if err != nil {
		log.Fatalf("Failed to read config file '%s': %s", *configFile, err)
	}

	// Command line overrides.
	if len(*apiEndpoint) > 0 {
		config.API.Endpoint = *apiEndpoint
	}

	if *verbose {
		tab := outputFormat()
		tab.Header("Field").SetAlign(tabulate.MR)
		tab.Header("Value").SetAlign(tabulate.ML)

		err = tabulate.Reflect(tab, tabulate.OmitEmpty, nil, config)
		if err != nil {
			log.Fatalf("Failed to tabulate: %s", err)
		}
		tab.Print(os.Stdout)
	}

	// Construct API client.
	auth, err := oauth.New(config.Auth, config.API.Endpoint,
		config.API.Certificate.X509, *verbose)
	if err != nil {
		log.Fatal(err)
	}
	client := api.New(api.AccessToken(auth),
		api.Endpoint(config.API.Endpoint),
		api.X509(config.API.Certificate.X509))
	if err != nil {
		log.Fatal(err)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		return
	}
	os.Args = flag.Args()
	fn, ok := commands[flag.Arg(0)]
	if !ok {
		fmt.Printf("Unknown command: %s\n", flag.Arg(0))
		os.Exit(1)
	}
	flag.CommandLine = flag.NewFlagSet(
		fmt.Sprintf("privx-cli %s", os.Args[0]),
		flag.ExitOnError)
	fn(client)
}