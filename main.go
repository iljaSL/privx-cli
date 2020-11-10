//
// Copyright (c) 2020 SSH Communications Security Inc.
//
// All rights reserved.
//

package main

import (
	"github.com/SSHcom/privx-cli/cmd"
)

//
/*
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
*/

//
func main() { cmd.Execute() }
