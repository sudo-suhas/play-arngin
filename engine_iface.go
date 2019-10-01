package arngin

import "context"

// Engine describes the functionality of an addon rule engine.
type Engine interface {
	// LoadRules loads the rules into the engine in preparation for
	// running the queries.
	LoadRules(ctx context.Context, rules []AddonRule) error

	// RunQuery matches the query against the addon rules loaded
	// into the rule engine.
	RunQuery(ctx context.Context, q AddonsQ) ([]string, error)
}
