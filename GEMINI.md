# Project Title: Aurora

A database migration tool for [Aurora
DSQL](https://docs.aws.amazon.com/aurora-dsql/latest/userguide/getting-started.html),
compatible with [Atlas](https://github.com/ariga/atlas), focusing on `migrate
apply` and `migrate status commands.

## Overview

Aurora is a command-line interface (CLI) tool designed to manage database
schema migrations for Aurora DSQL databases.

The aurora command-line tool must be compatible with
[atlas](https://github.com/ariga/atlas) command-line tool below:

```
'atlas migrate' wraps several sub-commands for migration management.

Usage:
  atlas migrate [command]

Available Commands:
  apply       Applies pending migration files on the connected database.
  status      Get information about the current migration status.

Flags:
  -c, --config string        select config (project) file using URL format (default "file://atlas.hcl")
      --env string           set which env from the config file to use
  -h, --help                 help for migrate

Use "atlas migrate [command] --help" for more information about a command.
```

The configuration file is in `hcl` format. The aurora command-line tool should
read the configuration file that atlas command-line tool provides:

```
env "aws" {
  migration {
    dir          = "file://database/migration"
    format       = atlas
    lock_timeout = "20m"
  }

  url = "postgres://${var.aws_dsql_username}:${urlescape(data.aws_dsql_token.this)}@${var.aws_dsql_host}/pharmacy-api"
}

data "aws_dsql_token" "this" {
  username = var.aws_dsql_username
  endpoint = var.aws_dsql_host
  region   = var.aws_region
}

variable "aws_dsql_username" {
  type    = string
  default = "pharmacy-api"
}

variable "aws_dsql_host" {
  type    = string
  default = getenv("PHARMACY_API_DATABASE_HOST")
}

variable "aws_region" {
  type    = string
  default = getenv("AWS_REGION")
}
```

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go: Version 1.18 or higher. You can check your version with go version.
- Aurora DSQL Database: Access to a running Aurora DSQL compatible database instance.

## Build

If you have to build the project, you should use the following command:

```bash
go build -o ./bin/aurora github.com/aws-contrib/atlas/cmd/atlas
```

## Test

You can run all tests in the project:

```bash
go tool ginkgo -r
```

You can run all tests for a given file:

```bash
go tool ginkgo <PATH_TO_TEST_FILE>
```

If you have to re-generate the auto-generate code in the project, you should use the following command:

```bash
go generate ./...
```

## Core Tools and Packages

The following packages and tools are used by the project:

- [CLI](github.com/urfave/cli) a package for implementing a command-line applications
- [Ginkgo](https://github.com/onsi/ginkgo) a testing package and command-line tool for running tests
- [Counterfeiter](https://github.com/maxbrunsfeld/counterfeiter) a command-line tool for stub/mock generation

If you have to install other Golang based tools, you should use

```bash
go get -tool <PACKAGE_URL>
```
