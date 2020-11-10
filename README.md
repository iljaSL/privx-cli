# privx-cli

This is a command line application to use with PrivX. The application supports On-Prem PrivX deployments, [PrivX on AWS](https://github.com/SSHcom/privx-on-aws) and PrivX SaaS.

## Getting Started

```bash
go get github.com/sshcom/privx-cli
```

Obtain access/secret keys from PrivX Instance so that the command line application is able to access PrivX Rest API
* Login as superuser
* Go to: Settings > Deployment > Integrate with PrivX Using API clients
* Create new API client (or use existing one)

Alternatively, you can use api client on behalf of existing user using its credentials.

```
usage: privx-cli [options] <command> <subcommand> [<subcommand> ...] [parameters]
To see help text, you can run:

  privx-cli help
  privx-cli <command> help
  privx-cli <command> <subcommand> help
```

## Configure client

It is mandatory to define HTTPS address to PrivX API (e.g. https://example.privx.io)

```bash
# with command line flag
privx-cli -c https://example.privx.io

# with environment variable(s)
export PRIVX_API_BASE_URL=https://example.privx.io
privx-cli
```

### Use user credentials

```bash
# with command line flag
privx-cli --access superuser --secret xhaSgasAU...As

# with environment variable(s)
export PRIVX_API_ACCESS_KEY=superuser
export PRIVX_API_SECRET_KEY=xhaSgasAU...As
privx-cli
```

### Use api client credentials

```bash
# with environment variable(s)
export PRIVX_API_ACCESS_KEY=00000000-0000-0000-0000-000000000000
export PRIVX_API_SECRET_KEY=xhaSgasAU...As
export PRIVX_API_OAUTH_CLIENT_ID=privx-external
export PRIVX_API_OAUTH_CLIENT_SECRET=another-random-base64
privx-cli
```

### Use config file

```bash
privx-cli -c ~/.privx.toml
```

the application support following config file format

```conf
[api]

# PRIVX_API_BASE_URL
base_url="https://your-instance.privx.io"

# restapi.X509(...)
api_ca_crt=""" PEM certificate chain """

[auth]

# PRIVX_API_ACCESS_KEY
api_client_id="00000000-0000-0000-0000-000000000000"

# PRIVX_API_SECRET_KEY
api_client_secret="some-random-base64"

# PRIVX_API_OAUTH_CLIENT_ID
oauth_client_id="privx-external"
# PRIVX_API_OAUTH_CLIENT_SECRET
oauth_client_secret="another-random-base64"
```

## Workflows

Note, each invocation of privx-cli causes a new authentication request using supplied credentials. It becomes inefficient if sequence of commands needs to be executed. It is possible to login once using any of supported methods and then reuse same access token

```bash
#!/bin/bash
export PRIVX_API_SECRET_KEY=$(privx-cli --access superuser --secret xhaSgasAU...As)

privx-cli roles
...
```


## Bugs

If you experience any issues with the application, please let us know via [GitHub issues](https://github.com/SSHcom/privx-cli/issues). We appreciate detailed and accurate reports that help us to identity and replicate the issue.

* **Specify** the configuration of your environment. Include which operating system you use and the versions of runtime environments.

* **Attach** logs, screenshots and exceptions, in possible.

* **Reveal** the steps you took to reproduce the problem, include code snippet or links to your project.


## How To Contribute

The project is [Apache 2.0](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request


## License

[![See LICENSE](https://img.shields.io/github/license/SSHcom/privx-cli.svg?style=for-the-badge)](LICENSE)
