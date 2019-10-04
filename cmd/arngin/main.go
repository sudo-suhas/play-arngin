package main

import (
	"context"
	"fmt"
	"os"

	arngin "github.com/sudo-suhas/play-arngin"
	"github.com/sudo-suhas/play-arngin/cmd/arngin/internal"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		app     = kingpin.New("arngin", "Addon rule engine playground.")
		args    internal.CommonArgs
		engines = []string{arngin.Opa, arngin.Arango, arngin.Elastic}
	)
	app.HelpFlag.Short('h')
	// Common arguments
	{
		enginesHelp := fmt.Sprintf(
			"Addon rule engines to run command against. Accepted values - %s, %s, %s, all.",
			arngin.Opa, arngin.Arango, arngin.Elastic,
		)
		app.Flag("engines", enginesHelp).Default("elastic").Short('e').EnumsVar(&args.Engines, append(engines, "all")...)

		app.Flag("arango-url", "Endpoint for ArangoDB instance.").Default("http://localhost:8529").Envar("ARANGO_URL").StringVar(&args.ArangoCfg.Endpoint)
		app.Flag("arango-db-name", "Database name for ArangoDB.").Default("arngin").Envar("ARANGO_DB_NAME").StringVar(&args.ArangoCfg.DB)
		app.Flag("arango-username", "Username for ArangoDB instance.").Envar("ARANGO_USERNAME").StringVar(&args.ArangoCfg.Username)
		app.Flag("arango-password", "Password for ArangoDB instance.").Envar("ARANGO_PASSWORD").StringVar(&args.ArangoCfg.Password)

		app.Flag("elastic-url", "Elasticsearch server URL.").Default("http://localhost:9200").Envar("ELASTIC_URL").StringVar(&args.ElasticURL)
	}

	// Sub commands and their respective flags
	var (
		load      = app.Command("load", "Load addon rules into the rule engine(s). Addon rules are generated randomly.")
		ruleCount = load.Flag("count", "Number of rules to generate.").Default("10000").Short('n').Int()

		query       = app.Command("query", "Run queries against the rule engine(s). Queries are generated randomly.")
		queryCount  = query.Flag("count", "Number of queres to generate.").Default("1000").Short('n').Int()
		concurrency = query.Flag("concurrency", "Number of goroutines to spawn to run queries in parallel.").Default("4").Short('c').Int()
	)

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	// Resolve 'all' to the 3 engines.
	for _, e := range args.Engines {
		if e == "all" {
			args.Engines = engines
			break
		}
	}

	switch cmd {
	// Load rules into engines.
	case load.FullCommand():
		if *ruleCount < 1 {
			app.FatalUsage("Invalid rule count: %d", *ruleCount)
		}

		app.FatalIfError(internal.LoadCommand(ctx, args, *ruleCount), cmd)

	// Run queries against engines.
	case query.FullCommand():
		if *queryCount < 1 {
			app.FatalUsage("Invalid query count: %d", *queryCount)
		}
		if *concurrency < 1 {
			app.FatalUsage("Invalid concurrency: %d", *concurrency)
		}

		app.FatalIfError(internal.QueryCommand(ctx, args, *queryCount, *concurrency), cmd)
	}
}
