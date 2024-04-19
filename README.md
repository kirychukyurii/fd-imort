# FreshDesk Import

Easy import .json files, exported from API

## Table of Contents

- [Getting Started](#getting-started)


## Getting Started

```
Usage:
  fd-import [flags]
  fd-import [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  migrate     Apply database schema migrations

Flags:
  -a, --attachment string      directory to store attachment files (default "./attachments")
      --domain string          domain name
  -d, --dsn string             database connection string
  -h, --help                   help for fd-import
      --log-file string        log file (default "./fd-import.log")
  -l, --log-level string       log level (default "debug")
  -p, --path string            base path to exported files (default "./export-data")
      --s3.access-key string   S3 access key ID
      --s3.bucket string       S3 bucket
      --s3.region string       S3 region
      --s3.secret-key string   S3 secret access key
  -v, --version                version for fd-import
  -w, --workers-count int      number of concurrent workers (default 100)

Use "fd-import [command] --help" for more information about a command.
```