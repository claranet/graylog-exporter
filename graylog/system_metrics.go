
package graylog

import "encoding/json"
import "log"
import "io/ioutil"

type Metrics struct {
	Version		string			`json:"version"`
//	Gauges		map[string]MetricGauge	`json:"gauges"`
//	Counters	*[]MetricCounter
//	Histograms	*[]MetricHistogram
}

type MetricGauge struct {
	Value		float64
}

func (m *Metrics) GetGauge(key string) float64 {
	return 0 //float64(m.Gauges[key].Value)
}

func (g *Graylog) GetSystemMetrics() Metrics {
	resp, _ := g.makeRequest("GET", "/system/metrics")
	data := json.NewDecoder(resp.Body)

	x, _ :=ioutil.ReadAll(resp.Body)

	var d Metrics
	data.Decode(&d)

	log.Printf("%s", x)

	return d
}
