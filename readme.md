# play-arngin

Addon rule engine playground.

## Setup

First install [Go][go-install] `v1.13` or above. We use [Go modules][go-modules]
for building the binary without relying on `GOPATH`. And errors are wrapped
using the new `%w` verb(just cause).

The dependencies can be installed by running `go get` but this is optional since
`go build` is module aware and will fetch dependencies if required.

## Build

```
$ go build ./cmd/arngin
```

## Usage

```
$ ./arngin
usage: arngin [<flags>] <command> [<args> ...]

Addon rule engine playground.

Flags:
  -h, --help                     Show context-sensitive help (also try --help-long and --help-man).
  -e, --engines=elastic ...      Addon rule engines to run command against. Accepted values - opa, arango, elastic, all.
      --arango-url="http://localhost:8529"
                                 Endpoint for ArangoDB instance.
      --arango-db-name="arngin"  Database name for ArangoDB.
      --arango-username=ARANGO-USERNAME
                                 Username for ArangoDB instance.
      --arango-password=ARANGO-PASSWORD
                                 Password for ArangoDB instance.
      --elastic-url="http://localhost:9200"
                                 Elasticsearch server URL.

Commands:
  help [<command>...]
    Show help.

  load [<flags>]
    Load addon rules into the rule engine(s). Addon rules are generated randomly.

  query [<flags>]
    Run queries against the rule engine(s). Queries are generated randomly.

```

### Examples

```
# Load 1000 rules into all engines, use default endpoints
$ ./arngin load --engines all --count 10000 --arango-username testuser --arango-password testpass --arango-db-name arnginTest
Loaded 1000 rules into opa [took: 9.034712154s]
Loaded 1000 rules into arango [took: 109.521924ms]
Loaded 1000 rules into elastic [took: 848.229276ms]

# Run 100 queries with a concurrency of 4 against opa
$ ./arngin query --engines opa --count 100 --concurrency 4
Ran 100 queries with concurrency of 4 against opa [avg: 362.49385ms]

# Run 1000 queries with a concurrency of 10 against arango and elastic
$ ./arngin query -e arango -e elastic --count 1000 --concurrency 10 --arango-username testuser --arango-password testpass --arango-db-name arnginTest
Ran 1000 queries with concurrency of 10 against arango [avg: 82.168104ms]
Ran 1000 queries with concurrency of 10 against elastic [avg: 7.499732ms]
```

[go-install]: https://golang.org/doc/install
[go-modules]: https://blog.golang.org/modules2019
