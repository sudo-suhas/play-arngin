package internal

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	arngin "github.com/sudo-suhas/play-arngin"
	"github.com/sudo-suhas/play-arngin/gen"
	"golang.org/x/sync/errgroup"
)

// QueryCommand runs n queries against the engines specified in args.
func QueryCommand(ctx context.Context, args CommonArgs, n, concurrency int) error {
	if n < 10 {
		concurrency = 1
	}

	var errs arngin.Errors
	for _, name := range args.Engines {
		e, err := makeEngine(ctx, name, args)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		var (
			g      errgroup.Group
			avgDur int64
			qcnt   = n / concurrency // we could lose some accuracy here but that's okay
		)
		for i := 0; i < concurrency; i++ {
			g.Go(func() error {
				qs := gen.Qs(qcnt)
				start := time.Now()
				for _, q := range qs {
					if _, err := e.RunQuery(ctx, q); err != nil {
						return fmt.Errorf("run query[%#v]: %w", q, err)
					}
				}

				atomic.AddInt64(
					&avgDur,
					(int64)(time.Since(start))/(int64)(qcnt*concurrency),
				)
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			errs = append(errs, fmt.Errorf("run queries[engine: %s]: %w", name, err))
			continue
		}

		fmt.Printf(
			"Ran %d queries with concurrency of %d against %s [avg: %s]\n",
			n, concurrency, name, (time.Duration)(avgDur),
		)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
