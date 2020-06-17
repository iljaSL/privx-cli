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

	"github.com/SSHcom/privx-sdk-go/oauth"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/markkurossi/tabulate"
)

//
// command-line options
type opts struct {
	config  *string
	format  *string
	access  *string
	secret  *string
	verbose *bool
}

//
// Supported commands
var commands = map[string]func(client restapi.Connector){
	"users":   cmdUsers,
	"secrets": cmdSecrets,
	"roles":   cmdRoles,
}

//
// Supported formatting
var outputFormat func() *tabulate.Tabulate

var formats = map[string]func() *tabulate.Tabulate{
	"whitespace": tabulate.NewWS,
	"ascii":      tabulate.NewASCII,
	"unicode":    tabulate.NewUnicode,
	"colon":      tabulate.NewColon,
	"csv":        tabulate.NewCSV,
}

//
func optsParse() *opts {
	fopts := &opts{
		config:  flag.String("config", "", "configuration file"),
		format:  flag.String("format", "unicode", "output format"),
		access:  flag.String("access", "", "access key"),
		secret:  flag.String("secret", "", "secret key"),
		verbose: flag.Bool("v", false, "verbose output"),
	}
	flag.Parse()
	return fopts
}

//
func optsUsage() {
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

//
func auth(opts *opts) restapi.Authorizer {
	curl := restapi.New(
		restapi.UseConfigFile(*opts.config),
		restapi.UseEnvironment(),
	)

	if *opts.access != "" {
		return oauth.WithCredential(
			curl,
			oauth.UseConfigFile(*opts.config),
			oauth.UseEnvironment(),
			oauth.Access(*opts.access),
			oauth.Secret(*opts.secret),
		)
	}

	return oauth.WithClientID(
		curl,
		oauth.UseConfigFile(*opts.config),
		oauth.UseEnvironment(),
	)
}

//
func main() {
	log.SetFlags(0)

	flag.Usage = optsUsage
	opts := optsParse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		return
	}

	outputFormat = formats[*opts.format]
	if outputFormat == nil {
		log.Printf("Invalid output format '%s'", *opts.format)
		log.Printf("Supported formats are:")
		for k := range formats {
			log.Printf(" - %s", k)
		}
		os.Exit(1)
	}

	client := restapi.New(
		restapi.Auth(auth(opts)),
		restapi.UseConfigFile(*opts.config),
		restapi.UseEnvironment(),
	)

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
