package elasticsearch

import (
	"context"
	"errors"
	"fmt"

	"github.com/olivere/elastic/v7"
	arngin "github.com/sudo-suhas/play-arngin"
)

// Client implements the arngin.Engine interface using Elasticsearch.
type Client struct {
	*elastic.Client
}

const indexName = "addon_rules"

// NewClient creates an instance of the Elasticsearch Client which
// satisfies the arngin.Engine interface.
func NewClient(ctx context.Context, url string) (arngin.Engine, error) {
	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		// elastic.SetTraceLog(log.New(os.Stdout, "", 0)),
	)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	// See https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-exists.html
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("check index exists: %w", err)
	}

	if !exists {
		// See https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-create-index.html
		createIndex, err := client.CreateIndex(indexName).
			Body(mapping).
			// Pretty(true).
			Do(ctx)
		if err != nil {
			return nil, fmt.Errorf("create index: %w", err)
		}
		if !createIndex.Acknowledged {
			return nil, fmt.Errorf("create index not acknowledged: %w", err)
		}
	}

	return &Client{Client: client}, nil
}

func (c *Client) LoadRules(ctx context.Context, rules []arngin.AddonRule) error {
	// See https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-delete-by-query.html
	_, err := c.DeleteByQuery(indexName).
		Query(elastic.NewMatchAllQuery()).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("delete all docs in index: %w", err)
	}

	// See https://github.com/olivere/elastic/blob/v7.0.6/recipes/bulk_insert/bulk_insert.go
	bulk := c.Bulk().Index(indexName).Refresh("true")
	for _, r := range rules {
		bulk.Add(elastic.NewBulkIndexRequest().Id(r.ID.String()).Doc(r))
	}

	res, err := bulk.Do(ctx)
	if err != nil {
		return fmt.Errorf("execute bulk request: %w", err)
	}
	if res.Errors {
		errDtls := res.Failed()[0].Error
		if errDtls == nil {
			return errors.New("bulk commit failed: unknown reason")
		}
		// Look up the failed documents with res.Failed(), and e.g. recommit
		return fmt.Errorf("bulk commit failed: first error: %s [type=%s]", errDtls.Reason, errDtls.Reason)
	}

	return nil
}

func (c *Client) RunQuery(ctx context.Context, q arngin.AddonsQ) ([]string, error) {
	// See https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-bool-query.html
	qry := elastic.NewBoolQuery().Filter(
		termQueryIfExists("sources", q.Source),
		termQueryIfExists("destinations", q.Destination),
		termQueryIfExists("boardingPts", q.BoardingPt),
		termQueryIfExists("droppingPts", q.DroppingPt),
		termQueryIfExists("boardingTime", q.BoardingTime),
		termQueryIfExists("droppingTime", q.DroppingTime),
		numRuleQueryIfExists("seatCount", q.Seats),
		termQueryIfExists("busOperators", q.BusOperator),
		numRuleQueryIfExists("duration", q.Duration),
		termQueryIfExists("channels", q.Channel),
	)
	if q.Channel == "MOBILE_APP" {
		qry.Filter(numRuleQueryIfExists("appversion", q.Appversion))
	}
	// See https://www.elastic.co/guide/en/elasticsearch/reference/current/search-aggregations-bucket-terms-aggregation.html
	termsAggName := "addon_terms_agg"
	termsAgg := elastic.NewTermsAggregation().Field("addons").Size(1000)
	// See https://www.elastic.co/guide/en/elasticsearch/reference/current/search-aggregations-bucket-filter-aggregation.html
	filterAggName := "filter_agg"
	filterAgg := elastic.NewFilterAggregation().
		Filter(qry).
		SubAggregation(termsAggName, termsAgg)

	// See https://www.elastic.co/guide/en/elasticsearch/reference/current/search-request-body.html
	searchResult, err := c.Search().
		Index(indexName).                      // search in index "addon_test"
		Aggregation(filterAggName, filterAgg). // specify the aggregation with filter query
		FilterPath("took", "aggregations").    // Reduce the response returned by Elasticsearch
		Size(0).                               // Return only aggregation results (See https://www.elastic.co/guide/en/elasticsearch/reference/current/returning-only-agg-results.html)
		Do(ctx)                                // execute
	if err != nil {
		return nil, fmt.Errorf("execute search: %w", err)
	}

	filterAggResult, ok := searchResult.Aggregations.Filter(filterAggName)
	if !ok {
		return nil, errors.New("failed to parse filter agg")
	}

	termsAggResult, ok := filterAggResult.Terms(termsAggName)
	if !ok {
		return nil, errors.New("failed to parse terms agg")
	}

	addons := make([]string, len(termsAggResult.Buckets))
	for i, b := range termsAggResult.Buckets {
		id, ok := b.Key.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected key type: %T", b.Key)
		}

		addons[i] = id
	}

	return addons, nil
}

func termQueryIfExists(field string, v interface{}) *elastic.BoolQuery {
	return withIfExists(field, elastic.NewTermQuery(field, v))
}

func numRuleQueryIfExists(field string, n int) *elastic.BoolQuery {
	return withIfExists(
		field,
		elastic.NewBoolQuery().Should(
			elastic.NewBoolQuery().Must(
				elastic.NewTermQuery(field+".op", arngin.GTEOp),
				elastic.NewRangeQuery(field+".value").Gte(n),
			),
			elastic.NewBoolQuery().Must(
				elastic.NewTermQuery(field+".op", arngin.EqOp),
				elastic.NewTermQuery(field+".value", n),
			),
			elastic.NewBoolQuery().Must(
				elastic.NewTermQuery(field+".op", arngin.LTEOp),
				elastic.NewRangeQuery(field+".value").Lte(n),
			),
		),
	)
}

func withIfExists(field string, qry elastic.Query) *elastic.BoolQuery {
	return elastic.NewBoolQuery().
		Should(elastic.NewBoolQuery().MustNot(elastic.NewExistsQuery(field))).
		Should(qry)
}
