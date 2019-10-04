package internal

import (
	"context"
	"fmt"
	"time"

	arngin "github.com/sudo-suhas/play-arngin"
	"github.com/sudo-suhas/play-arngin/gen"
)

// LoadCommand loads rules into the engines specified in the args.
func LoadCommand(ctx context.Context, args CommonArgs, n int) error {
	rules := gen.Rules(n)

	var errs arngin.Errors
	for _, name := range args.Engines {
		e, err := makeEngine(ctx, name, args)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		start := time.Now()
		if err := e.LoadRules(ctx, rules); err != nil {
			errs = append(errs, fmt.Errorf("load rules[engine: %s]: %w", name, err))
			continue
		}

		fmt.Printf("Loaded %d rules into %s [took: %s]\n", n, name, time.Since(start))
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
