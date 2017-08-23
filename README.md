# smtpoo

Fake SMTP server, caching outbound emails on Redis. Probably useless, useful for testing.

## Usage

```
NAME:
   smtpoo - a fake SMTP server caching outbound emails on Redis

USAGE:
   smtpoo [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value, -p value    SMTP port (default: 25)
   --redis-addr value        Redis address (default: "localhost")
   --redis-port value        Redis port (default: 6379)
   --redis-db value          Redis db number (default: 0)
   --redis-pass value        Redis password
   --redis-expiration value  Redis keys expiration time (seconds) (default: 0)
   --help, -h                show help
   --version, -v             print the version
```

Emails are stored as JSON objects on Redis, with keys being `mail:<timestamp>`.