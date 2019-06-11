package graylog

import "encoding/json"

type Metrics struct {
	Version string                 `json:"version"`
	Gauges  map[string]MetricGauge `json:"gauges"`
	//	Counters	*[]MetricCounter
	//	Histograms	*[]MetricHistogram
}

type MetricGauge struct {
	Value float64
}

func (m *Metrics) GetGauge(key string) float64 {
	return float64(m.Gauges[key].Value)
}

func (g *Graylog) GetSystemMetrics() Metrics {
	resp, _ := g.makeRequest("GET", "/system/metrics")
	data := json.NewDecoder(resp.Body)

	var d Metrics
	data.Decode(&d)

	return d
}
