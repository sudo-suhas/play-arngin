package elasticsearch

const mapping = `
{
	"settings": {
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings": {
		"properties": {
			"id": { "type":"keyword" },
			"sources": { "type":"integer" },
			"destinations": { "type":"integer" },
			"boardingPts": { "type":"integer" },
			"droppingPts": { "type":"integer" },
			"boardingTime": { "type":"float_range" },
			"droppingTime": { "type":"float_range" },
			"seatCount": {
				"properties": {
					"op": { "type": "keyword" },
					"value": { "type": "integer" }
				}
			},
			"busOperators": { "type":"keyword" },
			"duration": {
				"properties": {
					"op": { "type": "keyword" },
					"value": { "type": "integer" }
				}
			},
			"appversion": {
				"properties": {
					"op": { "type": "keyword" },
					"value": { "type": "integer" }
				}
			},
			"channels": { "type":"keyword" },
			"addons": {
				"type":"keyword",
				"store": true
			}
		}
	}
}`
