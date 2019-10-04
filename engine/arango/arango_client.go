package arango

import (
	"context"
	"fmt"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/vst"
	arngin "github.com/sudo-suhas/play-arngin"
)

// Client implements the arngin.Engine interface using ArangoDB.
type Client struct {
	db driver.Database
}

// ClientConfig is the parameters required for connecting to an Arango
// database.
type ClientConfig struct {
	Endpoint string
	DB       string
	Username string
	Password string
}

const (
	dbName         = "arngin"
	collectionName = "addonRules"
)

// NewClient creates an instance of the arango Client which satisfies
// the arngin.Engine interface.
func NewClient(ctx context.Context, cfg ClientConfig) (arngin.Engine, error) {
	// Create a VST connection to the database
	conn, err := vst.NewConnection(vst.ConnectionConfig{
		Endpoints: []string{cfg.Endpoint},
	})
	if err != nil {
		return nil, fmt.Errorf("create connection to DB: %w", err)
	}

	// Create a client
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(cfg.Username, cfg.Password),
	})
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	// Open the database
	db, err := client.Database(ctx, cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("open DB: %w", err)
	}

	// Check if our collection exists, create if it does not.
	exists, err := db.CollectionExists(ctx, collectionName)
	if err != nil {
		return nil, fmt.Errorf("check collection exists: %w", err)
	}

	if !exists {
		opts := driver.CreateCollectionOptions{
			KeyOptions: &driver.CollectionKeyOptions{
				AllowUserKeys: true, // supply own key values in the _key attribute of a document
			},
		}
		_, err = db.CreateCollection(ctx, collectionName, &opts)
		if err != nil {
			return nil, fmt.Errorf("create collection: %w", err)
		}
	}

	return &Client{db: db}, nil
}

func (c *Client) LoadRules(ctx context.Context, rules []arngin.AddonRule) error {
	col, err := c.db.Collection(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("open collection: %w", err)
	}

	if err := col.Truncate(ctx); err != nil {
		return fmt.Errorf("truncate collection: %w", err)
	}

	// Precompute and store ignore clauses as it is slightly faster to check a
	// boolean than checking if a field is null.
	optimRules := make([]addonRuleOptim, 0, len(rules))
	for _, r := range rules {
		optim := addonRuleOptim{
			Key:       r.ID,
			AddonRule: r,
			Ignore: ignoreClause{
				Sources:      len(r.Sources) == 0,
				Destinations: len(r.Destinations) == 0,
				BoardingPts:  len(r.BoardingPts) == 0,
				DroppingPts:  len(r.DroppingPts) == 0,
				BoardingTime: r.BoardingTime == nil,
				DroppingTime: r.DroppingTime == nil,
				SeatCount:    r.SeatCount == nil,
				BusOperators: len(r.BusOperators) == 0,
				Duration:     r.Duration == nil,
				Appversion:   r.Appversion == nil,
				Channels:     len(r.Channels) == 0,
			},
		}
		optimRules = append(optimRules, optim)
	}

	_, errs, err := col.CreateDocuments(ctx, optimRules)
	if err != nil {
		return fmt.Errorf("create documents: %w", err)
	} else if err := errs.FirstNonNil(); err != nil {
		return fmt.Errorf("create documents: first error: %w", err)
	}

	return nil
}

const qry = `RETURN UNIQUE(FLATTEN(FOR r IN addonRules
FILTER r.ignore.sources OR @source IN r.sources
FILTER r.ignore.destinations OR @destination IN r.destinations
FILTER r.ignore.sources OR @source IN r.sources
FILTER r.ignore.destinations OR @destination IN r.destinations
FILTER r.ignore.sources OR @source IN r.sources
FILTER r.ignore.destinations OR @destination IN r.destinations
FILTER r.ignore.boardingTime OR (@boardingTime >= r.boardingTime.gte AND @boardingTime <= r.boardingTime.lte)
FILTER r.ignore.droppingTime OR (@droppingTime >= r.droppingTime.gte AND @droppingTime <= r.droppingTime.lte)
FILTER r.ignore.seatCount OR
	(r.seatCount.op == "eq" AND @seats == r.seatCount.value) OR
	(r.seatCount.op == "lte" AND @seats <= r.seatCount.value) OR
	(r.seatCount.op == "gte" AND @seats >= r.seatCount.value)
FILTER r.ignore.busOperators OR @busOperator IN r.busOperators
FILTER r.ignore.duration OR
	(r.duration.op == "eq" AND @duration == r.duration.value) OR
	(r.duration.op == "lte" AND @duration <= r.duration.value) OR
	(r.duration.op == "gte" AND @duration >= r.duration.value)
FILTER r.ignore.appversion OR @appversion == 0 OR
	(r.appversion.op == "eq" AND @appversion == r.appversion.value) OR
	(r.appversion.op == "lte" AND @appversion <= r.appversion.value) OR
	(r.appversion.op == "gte" AND @appversion >= r.appversion.value)
FILTER r.ignore.channels OR @channel IN r.channels
RETURN r.addons))`

func (c *Client) RunQuery(ctx context.Context, q arngin.AddonsQ) ([]string, error) {
	bindVars := map[string]interface{}{
		"source":       q.Source,
		"destination":  q.Destination,
		"boardingTime": q.BoardingTime,
		"droppingTime": q.DroppingTime,
		"seats":        q.Seats,
		"busOperator":  q.BusOperator,
		"duration":     q.Duration,
		"appversion":   q.Appversion,
		"channel":      q.Channel,
	}
	cursor, err := c.db.Query(ctx, qry, bindVars)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer cursor.Close()

	var addons []string
	if _, err := cursor.ReadDocument(ctx, &addons); err != nil {
		return nil, fmt.Errorf("read doc from cursor: %w", err)
	}
	return addons, nil
}
