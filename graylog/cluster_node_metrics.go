
package graylog

import (
	"encoding/json"
	"fmt"
)

type ClusterNodeMetricBase struct {
	Total			int  `json:"total"`
	Metrics			[]ClusterNodeMetric  `json:"metrics"`
}

type ClusterNodeMetric  struct {
	FullName		string	`json:"full_name"`
	Name			string  `json:"name"`
	Type			string	`json:"type"`
	Metric			*ClusterNodeMetricValue  `json:"metric"`
}

type ClusterNodeMetricValue struct {
	Count			float64	`json:"count"`
	Value			float64	`json:"value"`
}

func (g *Graylog) GetClusterNodeMetrics(nodeId string, metrics []string) ClusterNodeMetricBase {
	bodyMetrics, _ := json.Marshal(metrics)
	params := &RequestParams{
		body: fmt.Sprintf(`{"metrics": %s}`, bodyMetrics),
	}
	action := fmt.Sprintf("/cluster/%s/metrics/multiple", nodeId)
	res, _ := g.makeRequestWithParams("POST", action, *params)

	data := json.NewDecoder(res.Body)
	var d ClusterNodeMetricBase
	data.Decode(&d)

	return d
}
