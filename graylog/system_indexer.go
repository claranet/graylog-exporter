
package graylog

import "encoding/json"

type SystemIndexOverview struct {
	Deflector		*IndexDeflector
	IndexerCluster		*IndexCluster	`json:"indexer_cluser"`
	Counts			*IndexCount
	Indices			map[string]SystemIndex
}

type IndexDeflector struct {
	CurrentTarget		string	`json:"current_target"`
	IsUp			bool	`json:"is_up"`
}

type IndexCount struct {
	Events			int
}

type IndexCluster struct {
	Name			string
	Health			*IndexClusterHealth
}

type IndexClusterHealth struct {
	Status			string
	Shards			*IndexClusterHealthShards
}

type IndexClusterHealthShards struct {
	Shards			int
	Initializing		int
	Relocating		int
	Unassigned		int
}

type SystemIndex struct {

}

func (g *Graylog) GetSystemIndexOverview() SystemIndex {

	res, _ := g.makeRequest("GET", "/system/indexer/overview")
	data := json.NewDecoder(res.Body)

	var d SystemIndex
	data.Decode(&d)

	return d

}
