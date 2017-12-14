
package graylog

import "encoding/json"

type SystemThroughput struct {
	Throughput	int
}

func (g *Graylog) GetSystemThroughput() SystemThroughput {
	res, _ := g.makeRequest("GET", "/system/throughput")
	data := json.NewDecoder(res.Body)

	var d SystemThroughput
	data.Decode(&d)

	return d
}
