package internal

import (
	"context"
	"fmt"

	arngin "github.com/sudo-suhas/play-arngin"
	"github.com/sudo-suhas/play-arngin/engine/arango"
	"github.com/sudo-suhas/play-arngin/engine/elasticsearch"
	"github.com/sudo-suhas/play-arngin/engine/opa"
)

// CommonArgs is the common set of arguments applicable for all
// commands.
type CommonArgs struct {
	Engines    []string
	ArangoCfg  arango.ClientConfig
	ElasticURL string
}

func makeEngine(ctx context.Context, name string, args CommonArgs) (arngin.Engine, error) {
	var (
		e   arngin.Engine
		err error
	)
	switch name {
	case arngin.Opa:
		e, err = opa.NewRegoEngine(ctx)

	case arngin.Arango:
		e, err = arango.NewClient(ctx, args.ArangoCfg)

	case arngin.Elastic:
		e, err = elasticsearch.NewClient(ctx, args.ElasticURL)
	}

	if err != nil {
		return nil, fmt.Errorf("create engine[%s]: %w", name, err)
	}

	return e, nil
}
