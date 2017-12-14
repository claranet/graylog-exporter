
package graylog

import "encoding/json"

type SystemClusterNodes struct {
	Nodes	[]SystemCluserNode
	Total	int
}

type SystemCluserNode  struct {
	ClusterId		string	`json:"cluster_id"`
	NodeId			string	`json:"node_id"`
	Type			string	`json:"type"`
	TransportAddress	string	`json:"transport_address"`
	LastSeen		string	`json:"last_seen"`
	ShortNodeId		string	`json:"short_node_id"`
	Hostname		string	`json:"hostname"`
	IsMaster		bool	`json:"is_master"`
}

func (g *Graylog) GetSystemClusterNodes() SystemClusterNodes {
	res, _ := g.makeRequest("GET", "/system/cluster/nodes")
	data := json.NewDecoder(res.Body)

	var d SystemClusterNodes
	data.Decode(&d)

	return d
}
